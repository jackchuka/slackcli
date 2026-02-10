package slack

type AuthTestResult struct {
	UserID string `json:"user_id"`
	User   string `json:"user"`
	TeamID string `json:"team_id"`
	Team   string `json:"team"`
	URL    string `json:"url"`
}

func (c *Client) AuthTest() (*AuthTestResult, error) {
	resp, err := retry(func() (*AuthTestResult, error) {
		r, err := c.api.AuthTest()
		if err != nil {
			return nil, err
		}
		return &AuthTestResult{
			UserID: r.UserID,
			User:   r.User,
			TeamID: r.TeamID,
			Team:   r.Team,
			URL:    r.URL,
		}, nil
	})
	if err != nil {
		return nil, classifyError(err)
	}
	return resp, nil
}
