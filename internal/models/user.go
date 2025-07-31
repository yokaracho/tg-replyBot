package models

import "time"

type User struct {
	ID       int64     `json:"id"`
	Username string    `json:"username"`
	Context  *Context  `json:"context"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}
