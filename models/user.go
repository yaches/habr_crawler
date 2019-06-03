package models

import "time"

type User struct {
	Username           string    `json:",omitempty"`
	Name               string    `json:",omitempty"`
	Spec               string    `json:",omitempty"`
	About              string    `json:",omitempty"`
	Birthday           time.Time `json:",omitempty"`
	Badges             []string
	Hubs               []string
	Works              []string
	SubscribeCompanies []string
	Invites            []string
	InvitedBy          string `json:",omitempty"`
	Karma              float32
	Rating             float32
	Subscribers        int
	From               []string
	RegDate            time.Time `json:",omitempty"`
	PostsCount         int
	CommentsCount      int
}
