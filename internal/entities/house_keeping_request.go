package entities

import (
	"time"
)

type HousekeepingRequest struct {
	ID          uint64    `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	RoomID      uint64    `db:"room_id"`
	RequestTime time.Time `db:"request_time"`
	Status      string    `db:"status"`
}

type HouseKeepingRequestDTO struct {
	RoomID      uint64
	RequestTime time.Time
}
