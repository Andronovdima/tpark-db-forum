package models

type Vote struct {
	Nickname string
	Voice    int // enum [-1 , 1]
}
