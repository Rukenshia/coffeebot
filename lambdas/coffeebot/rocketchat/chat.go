package rocketchat

import (
	"fmt"
)

// Chat client API
type Chat struct {
	c *Client
}

// Message is a Rocket.Chat Message
type Message struct {
	RoomID      string       `json:"roomId,omitempty"`
	Channel     string       `json:"channel"`
	Text        string       `json:"text,omitempty"`
	Alias       string       `json:"alias,omitempty"`
	Emoji       string       `json:"emoji,omitempty"`
	Avatar      string       `json:"avatar,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment is a Message Attachment of a Rocket.Chat Message
type Attachment struct {
	Color             string `json:"color,omitempty"`
	Text              string `json:"text,omitempty"`
	Timestamp         string `json:"ts,omitempty"`
	ThumbURL          string `json:"thumb_url,omitempty"`
	MessageLink       string `json:"message_link,omitempty"`
	Collapsed         bool   `json:"collapsed,omitempty"`
	AuthorName        string `json:"author_name,omitempty"`
	AuthorLink        string `json:"author_link,omitempty"`
	AuthorIcon        string `json:"author_icon,omitempty"`
	Title             string `json:"title,omitempty"`
	TitleLink         string `json:"title_link,omitempty"`
	TitleLinkDownload string `json:"title_link_download,omitempty"`
	ImageURL          string `json:"image_url,omitempty"`
	AudioURL          string `json:"audio_url,omitempty"`
	VideoURL          string `json:"video_url,omitempty"`

	Fields []Field `json:"fields,omitempty"`
}

// Field is an Attachment Field
type Field struct {
	Short bool   `json:"short,omitempty"`
	Title string `json:"title"`
	Value string `json:"value"`
}

// PostMessage posts a Message to Rocket.Chat
func (c *Chat) PostMessage(message Message) error {
	type Response struct {
		Success bool `json:"success"`
	}

	res, err := c.c.NewRequest().
		SetBody(message).
		SetResult(Response{}).
		Post(fmt.Sprintf("%s/api/v1/chat.postMessage", c.c.URL))
	if err != nil {
		return err
	}

	r, ok := res.Result().(*Response)
	if !ok {
		return ErrInvalidAPIResponse
	}

	if !r.Success {
		return ErrNotSuccessful
	}

	return nil
}
