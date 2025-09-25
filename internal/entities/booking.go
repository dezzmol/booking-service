package entities

import (
	"time"
)

type BookingStatus int8

const (
	BookingStatusUnknown   BookingStatus = 0
	BookingStatusSuccess   BookingStatus = 1
	BookingStatusCancelled BookingStatus = 2
	BookingStatusConfirmed BookingStatus = 3
)

type Booking struct {
	ID        uint64        `db:"id"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
	RoomID    uint64        `db:"room_id"`
	StartDate time.Time     `db:"start_date"`
	EndDate   time.Time     `db:"end_date"`
	Comment   string        `db:"comment"`
	Status    BookingStatus `db:"status"`
	IsPaid    bool          `db:"is_paid"`
}

type CreateBookingDTO struct {
	RoomID    uint64
	StartDate time.Time
	EndDate   time.Time
	Comment   string
	Guests    []GuestDTO
}
