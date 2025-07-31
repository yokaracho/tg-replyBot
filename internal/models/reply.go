package models

type Reply struct {
	Style   string `json:"style"`
	Content string `json:"content"`
	Emoji   string `json:"emoji"`
}

type ReplyRequest struct {
	Message         string   `json:"message"`
	ContextMessages []string `json:"context_messages"`
	Styles          []string `json:"styles"`
}

type ReplyResponse struct {
	Replies []Reply `json:"replies"`
}
