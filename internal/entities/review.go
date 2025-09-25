package entities

import "time"

type Review struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	BookingID uint64    `db:"booking_id"`
	Rating    int       `db:"rating"`
	Comment   string    `db:"comment"`
}

type ReviewDTO struct {
	BookingID uint64
	GuestID   uint64
	Rating    int
	Comment   string
}
