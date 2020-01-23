package models

type Status struct {
	Forum int64	`json:"forum" `
	Post int64	`json:"post" `
	Thread int32 `json:"thread" `
	User int32 	`json:"user" `
}
