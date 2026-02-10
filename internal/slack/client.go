package slack

import slackapi "github.com/slack-go/slack"

type Client struct {
	api   *slackapi.Client
	token string
}

type Option func(*Client)

func NewClient(token string, opts ...Option) *Client {
	c := &Client{
		api:   slackapi.New(token),
		token: token,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithDebug() Option {
	return func(c *Client) {
		c.api = slackapi.New(c.token, slackapi.OptionDebug(true))
	}
}
