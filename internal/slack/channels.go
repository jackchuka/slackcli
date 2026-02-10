package slack

import (
	slackapi "github.com/slack-go/slack"
)

type Channel struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Topic      string `json:"topic,omitempty"`
	Purpose    string `json:"purpose,omitempty"`
	NumMembers int    `json:"num_members"`
	IsArchived bool   `json:"is_archived"`
	IsPrivate  bool   `json:"is_private"`
	IsMember   bool   `json:"is_member"`
	Created    int    `json:"created"`
}

func channelFromAPI(ch slackapi.Channel) Channel {
	return Channel{
		ID:         ch.ID,
		Name:       ch.Name,
		Topic:      ch.Topic.Value,
		Purpose:    ch.Purpose.Value,
		NumMembers: ch.NumMembers,
		IsArchived: ch.IsArchived,
		IsPrivate:  ch.IsPrivate,
		IsMember:   ch.IsMember,
		Created:    int(ch.Created),
	}
}

func (c *Client) ListChannels(params PaginationParams) (*PaginatedResult[Channel], error) {
	if params.All {
		return c.listAllChannels(params)
	}
	return c.listChannelsPage(params)
}

func (c *Client) listChannelsPage(params PaginationParams) (*PaginatedResult[Channel], error) {
	type result struct {
		channels []slackapi.Channel
		cursor   string
	}
	r, err := retry(func() (result, error) {
		channels, cursor, err := c.api.GetConversations(&slackapi.GetConversationsParameters{
			Cursor:          params.Cursor,
			Limit:           params.EffectiveLimit(),
			ExcludeArchived: false,
			Types:           []string{"public_channel", "private_channel"},
		})
		return result{channels, cursor}, err
	})
	if err != nil {
		return nil, classifyError(err)
	}
	items := make([]Channel, len(r.channels))
	for i, ch := range r.channels {
		items[i] = channelFromAPI(ch)
	}
	return &PaginatedResult[Channel]{
		Items:      items,
		NextCursor: r.cursor,
		HasMore:    r.cursor != "",
	}, nil
}

func (c *Client) listAllChannels(params PaginationParams) (*PaginatedResult[Channel], error) {
	var allItems []Channel
	cursor := params.Cursor
	for {
		page, err := c.listChannelsPage(PaginationParams{
			Cursor: cursor,
			Limit:  params.EffectiveLimit(),
		})
		if err != nil {
			return nil, err
		}
		allItems = append(allItems, page.Items...)
		if !page.HasMore {
			break
		}
		cursor = page.NextCursor
	}
	return &PaginatedResult[Channel]{
		Items:   allItems,
		HasMore: false,
	}, nil
}

func (c *Client) GetChannelInfo(channelID string) (*Channel, error) {
	ch, err := retry(func() (*slackapi.Channel, error) {
		return c.api.GetConversationInfo(&slackapi.GetConversationInfoInput{
			ChannelID: channelID,
		})
	})
	if err != nil {
		return nil, classifyError(err)
	}
	result := channelFromAPI(*ch)
	return &result, nil
}

func (c *Client) CreateChannel(name string, isPrivate bool) (*Channel, error) {
	ch, err := retry(func() (*slackapi.Channel, error) {
		return c.api.CreateConversation(slackapi.CreateConversationParams{
			ChannelName: name,
			IsPrivate:   isPrivate,
		})
	})
	if err != nil {
		return nil, classifyError(err)
	}
	result := channelFromAPI(*ch)
	return &result, nil
}

func (c *Client) ArchiveChannel(channelID string) error {
	_, err := retry(func() (struct{}, error) {
		return struct{}{}, c.api.ArchiveConversation(channelID)
	})
	if err != nil {
		return classifyError(err)
	}
	return nil
}

func (c *Client) InviteToChannel(channelID string, userIDs ...string) error {
	_, err := retry(func() (*slackapi.Channel, error) {
		return c.api.InviteUsersToConversation(channelID, userIDs...)
	})
	if err != nil {
		return classifyError(err)
	}
	return nil
}

func (c *Client) KickFromChannel(channelID, userID string) error {
	_, err := retry(func() (struct{}, error) {
		return struct{}{}, c.api.KickUserFromConversation(channelID, userID)
	})
	if err != nil {
		return classifyError(err)
	}
	return nil
}

func (c *Client) SetChannelTopic(channelID, topic string) error {
	_, err := retry(func() (*slackapi.Channel, error) {
		return c.api.SetTopicOfConversation(channelID, topic)
	})
	if err != nil {
		return classifyError(err)
	}
	return nil
}

func (c *Client) SetChannelPurpose(channelID, purpose string) error {
	_, err := retry(func() (*slackapi.Channel, error) {
		return c.api.SetPurposeOfConversation(channelID, purpose)
	})
	if err != nil {
		return classifyError(err)
	}
	return nil
}
