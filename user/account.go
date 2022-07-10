package user

import "time"

type Account struct {
	Id        int
	Firstname string
	Lastname  string
	Nickname  string
	Password  string
	Email     string
	Country   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
