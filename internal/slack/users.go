package slack

import (
	"context"

	slackapi "github.com/slack-go/slack"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	RealName string `json:"real_name"`
	Email    string `json:"email,omitempty"`
	IsAdmin  bool   `json:"is_admin"`
	IsBot    bool   `json:"is_bot"`
	Deleted  bool   `json:"deleted"`
	TZ       string `json:"tz,omitempty"`
	Presence string `json:"presence,omitempty"`
}

func userFromAPI(u slackapi.User) User {
	return User{
		ID:       u.ID,
		Name:     u.Name,
		RealName: u.RealName,
		Email:    u.Profile.Email,
		IsAdmin:  u.IsAdmin,
		IsBot:    u.IsBot,
		Deleted:  u.Deleted,
		TZ:       u.TZ,
	}
}

func (c *Client) ListUsers(params PaginationParams) (*PaginatedResult[User], error) {
	// slack-go uses GetUsersPaginated for paginated user lists
	var allUsers []User
	pager := c.api.GetUsersPaginated(slackapi.GetUsersOptionLimit(params.EffectiveLimit()))

	for {
		var err error
		pager, err = pager.Next(context.Background())
		if err != nil {
			if pager.Done(err) {
				break
			}
			return nil, classifyError(err)
		}
		for _, u := range pager.Users {
			allUsers = append(allUsers, userFromAPI(u))
		}
		if !params.All {
			// Return first page only
			return &PaginatedResult[User]{
				Items:   allUsers,
				HasMore: !pager.Done(nil),
			}, nil
		}
	}

	return &PaginatedResult[User]{
		Items:   allUsers,
		HasMore: false,
	}, nil
}

func (c *Client) GetUserInfo(userID string) (*User, error) {
	u, err := retry(func() (*slackapi.User, error) {
		return c.api.GetUserInfo(userID)
	})
	if err != nil {
		return nil, classifyError(err)
	}
	result := userFromAPI(*u)
	return &result, nil
}

func (c *Client) GetUserPresence(userID string) (string, error) {
	p, err := retry(func() (*slackapi.UserPresence, error) {
		return c.api.GetUserPresence(userID)
	})
	if err != nil {
		return "", classifyError(err)
	}
	return p.Presence, nil
}
