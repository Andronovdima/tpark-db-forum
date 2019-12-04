package models

type PostFull struct {
	Author User
	Forum Forum
	Post Post
	Thread Thread
}