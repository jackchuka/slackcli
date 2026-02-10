package slack

import (
	"strconv"
	"time"

	slackapi "github.com/slack-go/slack"
)

type Message struct {
	Timestamp string `json:"timestamp"`
	User      string `json:"user"`
	Text      string `json:"text"`
	ThreadTS  string `json:"thread_ts,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Type      string `json:"type"`
}

func messageFromAPI(msg slackapi.Message) Message {
	return Message{
		Timestamp: msg.Timestamp,
		User:      msg.User,
		Text:      msg.Text,
		ThreadTS:  msg.ThreadTimestamp,
		Type:      msg.Type,
	}
}

type ListMessagesParams struct {
	ChannelID  string
	Pagination PaginationParams
	Oldest     time.Time
	Latest     time.Time
}

func (c *Client) ListMessages(params ListMessagesParams) (*PaginatedResult[Message], error) {
	if params.Pagination.All {
		return c.listAllMessages(params)
	}
	return c.listMessagesPage(params)
}

func (c *Client) listMessagesPage(params ListMessagesParams) (*PaginatedResult[Message], error) {
	type result struct {
		messages []slackapi.Message
		hasMore  bool
		cursor   string
	}

	histParams := &slackapi.GetConversationHistoryParameters{
		ChannelID: params.ChannelID,
		Cursor:    params.Pagination.Cursor,
		Limit:     params.Pagination.EffectiveLimit(),
	}
	if !params.Oldest.IsZero() {
		histParams.Oldest = formatTimestamp(params.Oldest)
	}
	if !params.Latest.IsZero() {
		histParams.Latest = formatTimestamp(params.Latest)
	}

	r, err := retry(func() (result, error) {
		resp, err := c.api.GetConversationHistory(histParams)
		if err != nil {
			return result{}, err
		}
		cursor := ""
		if resp.ResponseMetaData.NextCursor != "" {
			cursor = resp.ResponseMetaData.NextCursor
		}
		return result{messages: resp.Messages, hasMore: resp.HasMore, cursor: cursor}, nil
	})
	if err != nil {
		return nil, classifyError(err)
	}

	items := make([]Message, len(r.messages))
	for i, msg := range r.messages {
		items[i] = messageFromAPI(msg)
		items[i].Channel = params.ChannelID
	}
	return &PaginatedResult[Message]{
		Items:      items,
		NextCursor: r.cursor,
		HasMore:    r.hasMore,
	}, nil
}

func (c *Client) listAllMessages(params ListMessagesParams) (*PaginatedResult[Message], error) {
	var allItems []Message
	cursor := params.Pagination.Cursor
	for {
		p := params
		p.Pagination = PaginationParams{Cursor: cursor, Limit: params.Pagination.EffectiveLimit()}
		page, err := c.listMessagesPage(p)
		if err != nil {
			return nil, err
		}
		allItems = append(allItems, page.Items...)
		if !page.HasMore {
			break
		}
		cursor = page.NextCursor
	}
	return &PaginatedResult[Message]{Items: allItems, HasMore: false}, nil
}

type SendMessageParams struct {
	ChannelID string
	Text      string
	ThreadTS  string
}

func (c *Client) SendMessage(params SendMessageParams) (*Message, error) {
	opts := []slackapi.MsgOption{
		slackapi.MsgOptionText(params.Text, false),
	}
	if params.ThreadTS != "" {
		opts = append(opts, slackapi.MsgOptionTS(params.ThreadTS))
	}

	type result struct {
		channel   string
		timestamp string
		text      string
	}

	r, err := retry(func() (result, error) {
		ch, ts, txt, err := c.api.SendMessage(params.ChannelID, opts...)
		return result{ch, ts, txt}, err
	})
	if err != nil {
		return nil, classifyError(err)
	}
	return &Message{
		Channel:   r.channel,
		Timestamp: r.timestamp,
		Text:      r.text,
		ThreadTS:  params.ThreadTS,
	}, nil
}

func (c *Client) EditMessage(channelID, timestamp, text string) (*Message, error) {
	type result struct {
		channel   string
		timestamp string
		text      string
	}

	r, err := retry(func() (result, error) {
		ch, ts, txt, err := c.api.UpdateMessage(channelID,
			timestamp,
			slackapi.MsgOptionText(text, false),
		)
		return result{ch, ts, txt}, err
	})
	if err != nil {
		return nil, classifyError(err)
	}
	return &Message{
		Channel:   r.channel,
		Timestamp: r.timestamp,
		Text:      r.text,
	}, nil
}

func (c *Client) DeleteMessage(channelID, timestamp string) error {
	_, err := retry(func() (struct{}, error) {
		_, _, err := c.api.DeleteMessage(channelID, timestamp)
		return struct{}{}, err
	})
	if err != nil {
		return classifyError(err)
	}
	return nil
}

type SearchParams struct {
	Query      string
	Sort       string
	SortDir    string
	Pagination PaginationParams
}

type SearchResult struct {
	Matches []Message `json:"matches"`
	Total   int       `json:"total"`
}

func (c *Client) SearchMessages(params SearchParams) (*SearchResult, error) {
	searchParams := slackapi.SearchParameters{
		Sort:          params.Sort,
		SortDirection: params.SortDir,
		Count:         params.Pagination.EffectiveLimit(),
		Page:          1,
	}

	r, err := retry(func() (*slackapi.SearchMessages, error) {
		msgs, err := c.api.SearchMessages(params.Query, searchParams)
		return msgs, err
	})
	if err != nil {
		return nil, classifyError(err)
	}

	matches := make([]Message, len(r.Matches))
	for i, m := range r.Matches {
		matches[i] = Message{
			Timestamp: m.Timestamp,
			User:      m.User,
			Text:      m.Text,
			Channel:   m.Channel.ID,
			Type:      "message",
		}
	}
	return &SearchResult{
		Matches: matches,
		Total:   r.Total,
	}, nil
}

func formatTimestamp(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10) + ".000000"
}
