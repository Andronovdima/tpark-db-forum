package models

type User struct {
	About    string
	Email    string // unique
	Fullname string
	Nickname string // unique
}
