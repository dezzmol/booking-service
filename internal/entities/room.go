package entities

import "time"

type Room struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Number    string    `db:"number"`
	Type      string    `db:"type"`
	HotelID   uint64    `db:"hotel_id"`
}

type RoomDTO struct {
	Number  string
	Type    string
	HotelID uint64
}
