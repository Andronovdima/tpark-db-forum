package models

type Thread struct {
	Author  string
	Created string
	Forum   string
	Id      int32
	Message string
	Slug    string
	Title   string
	Votes   int32
}
