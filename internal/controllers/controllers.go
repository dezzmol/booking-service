package controllers

import (
	"context"
	"database/sql"
	"time"

	"booking-service/internal/entities"

	"github.com/jmoiron/sqlx"
)

type (
	ds interface {
		FindRoomById(ctx context.Context, tx *sql.Tx, roomId int64) (entities.Room, error)
		SaveRoom(ctx context.Context, tx *sql.Tx, room *entities.Room) error
		SaveAllRooms(ctx context.Context, tx *sql.Tx, rooms []entities.Room) error
		SaveGuestAndReturnIt(ctx context.Context, tx *sql.Tx, input entities.Guest) (entities.Guest, error)
		SaveReview(ctx context.Context, tx *sql.Tx, review entities.Review) (entities.Review, error)
		SaveBooking(ctx context.Context, tx *sql.Tx, booking entities.Booking) error
		FindBookingById(ctx context.Context, tx *sql.Tx, bookingID uint64) (entities.Booking, error)
		FindBookingByDate(ctx context.Context, tx *sql.Tx, startDate time.Time, endDate time.Time) ([]entities.Booking, error)
		DeleteBooking(ctx context.Context, tx *sql.Tx, bookingID uint64) error
		FindBookingByRoomIDAndDate(
			ctx context.Context, tx *sql.Tx, roomID uint64, startDate time.Time, endDate time.Time,
		) ([]entities.Booking, error)
		IsRoomAvailableForBooking(ctx context.Context, tx *sql.Tx, roomID uint64, startDate, endDate time.Time) (bool, error)
		SaveHotel(ctx context.Context, tx *sql.Tx, hotel entities.Hotel) (entities.Hotel, error)
		FindHotelByID(ctx context.Context, tx *sql.Tx, id uint64) (entities.Hotel, error)
	}

	Controller struct {
		sql *sqlx.DB
		ds  ds
	}
)

func New(
	db *sqlx.DB,
	ds ds,
) *Controller {
	return &Controller{
		sql: db,
		ds:  ds,
	}
}
