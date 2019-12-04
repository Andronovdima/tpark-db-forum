package models



type Post struct {
	Author string
	Created string // is nullable true
	Forum string
	Id int64
	IsEdited bool
	Message string
	Parent int64
	Thread int32
}

