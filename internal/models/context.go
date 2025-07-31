package models

import "time"

type Context struct {
	UserID      int64     `json:"user_id"`
	Messages    []string  `json:"messages"`
	LastMessage string    `json:"last_message"`
	Replies     Replies   `json:"replies"`
	Updated     time.Time `json:"updated"`
}

type Replies map[string]string

func NewContext(userID int64) *Context {
	return &Context{
		UserID:   userID,
		Messages: make([]string, 0),
		Replies:  make(Replies),
		Updated:  time.Now(),
	}
}

func (c *Context) AddMessage(message string) {
	c.Messages = append(c.Messages, message)
	c.LastMessage = message
	c.Updated = time.Now()

	if len(c.Messages) > 10 {
		c.Messages = c.Messages[len(c.Messages)-10:]
	}
}

func (c *Context) GetRecentMessages(count int) []string {
	if len(c.Messages) <= count {
		return c.Messages
	}
	return c.Messages[len(c.Messages)-count:]
}
