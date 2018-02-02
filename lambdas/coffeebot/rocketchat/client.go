package rocketchat

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrNotSuccessful is returned when the "status" property of a Response is not "success"
	ErrNotSuccessful = errors.New("API returned non-success status")

	// ErrInvalidAPIResponse is returned when the data structures do not match (this package might be outdated then!)
	ErrInvalidAPIResponse = errors.New("Invalid API Response")
)

// UserCredentials are a combination of Username and Password used to authenticate against Rocket.Chat
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Authentication represents Rocket.Chat authentication parameters
type Authentication struct {
	AuthToken string `json:"authToken"`
	UserID    string `json:"userId"`
}

// Client Rocket.Chat API Client struct
type Client struct {
	auth Authentication
	URL  string

	Users *Users
	Chat  *Chat
}

// NewClient creates a new instance of RocketChat
func NewClient(url string) *Client {
	c := &Client{
		URL: url,
	}

	c.Users = &Users{c}
	c.Chat = &Chat{c}

	return c
}

// Login authenticates the user with Rocket.Chat and requests temporary auth credentials
func (c *Client) Login(username, password string) error {
	type Response struct {
		Data   Authentication `json:"data"`
		Status string         `json:"status"`
	}

	res, err := resty.R().
		SetBody(UserCredentials{
			username,
			password,
		}).
		SetResult(Response{}).
		Post(fmt.Sprintf("%s/api/v1/login", c.URL))
	if err != nil {
		return err
	}

	d, ok := res.Result().(*Response)
	if !ok {
		log.Debugf("%v", string(res.Body()))
		return ErrInvalidAPIResponse
	}

	if d.Status != "success" {
		return ErrNotSuccessful
	}

	c.auth = d.Data
	return nil
}

// NewRequest creates a new request including the current auth
func (c *Client) NewRequest() *resty.Request {
	return resty.R().
		SetHeader("X-Auth-Token", c.auth.AuthToken).
		SetHeader("X-User-Id", c.auth.UserID)
}
