package models

import "time"

type Comment struct {
	ID      string
	PostID  string
	Author  string
	Text    string
	PubDate time.Time
}
