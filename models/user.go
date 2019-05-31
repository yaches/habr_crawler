package models

import "time"

type User struct {
	Username           string
	Name               string
	Spec               string
	About              string
	Birthday           time.Time
	Badges             []string
	Hubs               []string
	Works              []string
	SubscribeCompanies []string
	Invites            []string
	InvitedBy          string
	Karma              float32
	Rating             float32
	Subscribers        int
	From               []string
	RegDate            time.Time
	PostsCount         int
	CommentsCount      int
}
