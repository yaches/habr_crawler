package models

import "time"

type User struct {
	Username      string
	Name          string
	About         string
	Contacts      []Contact
	Invites       []UserShortcut
	InvitedBy     string
	InvitedDate   time.Time
	Karma         float32
	Rating        float32
	Subscribers   int
	Country       string
	RegDate       time.Time
	PostsCount    int
	CommentsCount int
}

type UserShortcut struct {
	Username string
	Name     string
}

type Contact struct {
	Name string
	Link string
}
