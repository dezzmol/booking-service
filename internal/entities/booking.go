package entities

import (
	"time"
)

type Booking struct {
	ID            uint64    `db:"id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	RoomID        uint64    `db:"room_id"`
	StartDate     time.Time `db:"start_date"`
	EndDate       time.Time `db:"end_date"`
	Comment       string    `db:"comment"`
	Status        string    `db:"status"`
	PaymentStatus string    `db:"payment_status"`
	Guests        []Guest   `db:"-"`
}

type CreateBookingDTO struct {
	RoomID    uint64
	StartDate time.Time
	EndDate   time.Time
	Comment   string
	Guests    []GuestDTO
}
