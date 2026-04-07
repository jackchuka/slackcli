package slack

import (
	slackapi "github.com/slack-go/slack"
)

type Reaction struct {
	Name  string   `json:"name"`
	Count int      `json:"count"`
	Users []string `json:"users"`
}

type ReactedItem struct {
	Type      string     `json:"type"`
	Channel   string     `json:"channel,omitempty"`
	Timestamp string     `json:"timestamp,omitempty"`
	Reactions []Reaction `json:"reactions"`
}

func (c *Client) AddReaction(channelID, timestamp, name string) error {
	ref := slackapi.ItemRef{
		Channel:   channelID,
		Timestamp: timestamp,
	}
	_, err := retry(func() (struct{}, error) {
		return struct{}{}, c.api.AddReaction(name, ref)
	})
	if err != nil {
		return classifyError(err)
	}
	return nil
}

func (c *Client) RemoveReaction(channelID, timestamp, name string) error {
	ref := slackapi.ItemRef{
		Channel:   channelID,
		Timestamp: timestamp,
	}
	_, err := retry(func() (struct{}, error) {
		return struct{}{}, c.api.RemoveReaction(name, ref)
	})
	if err != nil {
		return classifyError(err)
	}
	return nil
}

func (c *Client) ListReactions(userID string, params PaginationParams) (*PaginatedResult[ReactedItem], error) {
	if params.All {
		return c.listAllReactions(userID, params)
	}
	return c.listReactionsPage(userID, params, "")
}

func (c *Client) listReactionsPage(userID string, params PaginationParams, cursor string) (*PaginatedResult[ReactedItem], error) {
	listParams := slackapi.ListReactionsParameters{
		User:   userID,
		Limit:  params.EffectiveLimit(),
		Cursor: cursor,
		Full:   true,
	}

	type listResult struct {
		items      []slackapi.ReactedItem
		nextCursor string
	}

	r, err := retry(func() (listResult, error) {
		items, nextCursor, err := c.api.ListReactions(listParams)
		return listResult{items, nextCursor}, err
	})
	if err != nil {
		return nil, classifyError(err)
	}

	result := convertReactedItems(r.items)

	return &PaginatedResult[ReactedItem]{
		Items:      result,
		NextCursor: r.nextCursor,
		HasMore:    r.nextCursor != "",
	}, nil
}

func (c *Client) listAllReactions(userID string, params PaginationParams) (*PaginatedResult[ReactedItem], error) {
	var allItems []ReactedItem
	cursor := ""
	for {
		result, err := c.listReactionsPage(userID, params, cursor)
		if err != nil {
			return nil, err
		}
		allItems = append(allItems, result.Items...)
		if !result.HasMore {
			break
		}
		cursor = result.NextCursor
	}
	return &PaginatedResult[ReactedItem]{
		Items:   allItems,
		HasMore: false,
	}, nil
}

func convertReactedItems(items []slackapi.ReactedItem) []ReactedItem {
	var result []ReactedItem
	for _, item := range items {
		ri := ReactedItem{
			Type: item.Type,
		}
		if item.Message != nil {
			ri.Channel = item.Message.Channel
			ri.Timestamp = item.Message.Timestamp
		}
		for _, rx := range item.Reactions {
			ri.Reactions = append(ri.Reactions, Reaction{
				Name:  rx.Name,
				Count: rx.Count,
				Users: rx.Users,
			})
		}
		result = append(result, ri)
	}
	return result
}
