package entities

import "time"

type Guest struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Name      string    `db:"name"`
}

type GuestDTO struct {
	Name string
}
