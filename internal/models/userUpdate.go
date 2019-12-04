package models

type UserUpdate struct {
	About    string
	Email    string // unique
	Fullname string
}
