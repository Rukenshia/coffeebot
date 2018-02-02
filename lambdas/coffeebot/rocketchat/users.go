package rocketchat

import (
	"fmt"
)

// Users client API
type Users struct {
	c *Client
}

// User of Rocket.Chat
type User struct {
	ID        string  `json:"_id"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
	Active    bool    `json:"active"`
	Name      string  `json:"name"`
	Username  string  `json:"username"`
	UtcOffset float32 `json:"utcOffset"`
}

// List returns a List of all Users, "all" being restricted by your permission of the client
func (u *Users) List() ([]User, error) {
	type Response struct {
		Users   []User `json:"users"`
		Success bool   `json:"success"`
	}

	res, err := u.c.NewRequest().
		SetQueryParam("count", "0").
		SetResult(Response{}).
		Get(fmt.Sprintf("%s/api/v1/users.list", u.c.URL))

	if err != nil {
		return []User{}, err
	}

	r, ok := res.Result().(*Response)
	if !ok {
		return []User{}, ErrInvalidAPIResponse
	}

	if !r.Success {
		return []User{}, ErrNotSuccessful
	}

	return r.Users, nil
}
