package models

type Thread struct {
	Author  string `json:"author" `
	Created string	`json:"created,omitempty" `
	Forum   string	`json:"forum" `
	Id      int32	`json:"id" `
	Message string  `json:"message" `
	Slug    string	`json:"slug,omitempty"`
	Title   string	`json:"title" `
	Votes   int32	`json:"votes" `
}
