package entity

import "time"

type Book struct {
	ID        string
	Name      string
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
