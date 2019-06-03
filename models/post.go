package models

import "time"

type Post struct {
	ID            string    `json:",omitempty"`
	Author        string    `json:",omitempty"`
	PubDate       time.Time `json:",omitempty"`
	Title         string    `json:",omitempty"`
	Text          string    `json:",omitempty"`
	Tags          []string
	Hubs          []string
	CommentsCount int
	Rating        int
}
