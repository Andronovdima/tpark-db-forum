package models

type Post struct {
	Author   string `json:"author" `
	Created  string `json:"created" `
	Forum    string `json:"forum" `
	Id       int64  `json:"id" `
	IsEdited bool   `json:"isEdited" `
	Message  string `json:"message" `
	Parent   int64  `json:"parent" `
	Thread   int32  `json:"thread" `
	Path          []int64 `json:"-"`

}
