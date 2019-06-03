package models

import "time"

type Comment struct {
	ID       string    `json:",omitempty"`
	ParentID string    `json:",omitempty"`
	PostID   string    `json:",omitempty"`
	Author   string    `json:",omitempty"`
	Text     string    `json:",omitempty"`
	PubDate  time.Time `json:",omitempty"`
	Rating   int
}
