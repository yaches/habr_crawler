package models

import "time"

type Post struct {
	ID      string
	Author  string
	PubDate time.Time
	Title   string
	Text    string
	Tags    []string
	Hubs    []string
}
