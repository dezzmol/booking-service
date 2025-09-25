package controllers

import (
	"context"
	"database/sql"
	"errors"

	"booking-service/internal/entities"
	"booking-service/internal/storage"
)

var (
	BookingCancelled                                   = errors.New("booking cancelled")
	CantRescheduleLaterThanSevenDaysBeforeStartBooking = errors.New("can't reschedule later than seven days before start booking")
)

const (
	PAYMENT_PAID      = "paid"
	PAYMENT_UNPAID    = "unpaid"
	PAYMENT_CANCELLED = "cancelled"

	BOOKING_PENDING   = "pending"
	BOOKING_CONFIRM   = "confirmed"
	BOOKING_CANCELLED = "cancelled"
)

func (c *Controller) CreateBooking(ctx context.Context, input entities.CreateBookingDTO) (entities.Booking, error) {
	if input.StartDate.After(input.EndDate) {
		return entities.Booking{}, entities.ErrStartDateIsAfterEndDate
	}

	var available bool
	err := storage.WithNoTransaction(ctx, c.sql, func(ctx context.Context, tx *sql.Tx) error {
		var errTx error
		if available, errTx = c.ds.IsRoomAvailableForBooking(ctx, tx, input.RoomID, input.StartDate, input.EndDate); errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return entities.Booking{}, err
	}
	if !available {
		return entities.Booking{}, entities.ErrRoomNotAvailable
	}

	var guests []entities.Guest
	for _, guest := range input.Guests {
		guests = append(guests, entities.Guest{
			Name: guest.Name,
		})
	}

	var booking entities.Booking
	err = storage.WithWriteTransaction(ctx, c.sql, func(ctx context.Context, tx *sql.Tx) error {
		var errTx error
		for i := range guests {
			var guest entities.Guest
			guest, errTx = c.ds.SaveGuestAndReturnIt(ctx, tx, guests[i])
			if errTx != nil {
				return errTx
			}
			guests[i] = guest
		}

		booking = entities.Booking{
			RoomID:    input.RoomID,
			StartDate: input.StartDate,
			EndDate:   input.EndDate,
			Comment:   input.Comment,
			Status:    entities.BookingStatusConfirmed,
		}

		errTx = c.ds.SaveBooking(ctx, tx, booking)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return entities.Booking{}, err
	}

	return booking, nil
}

func (c *Controller) CancelBooking(ctx context.Context, bookingID uint64) error {
	var booking entities.Booking
	if err := storage.WithWriteTransaction(ctx, c.sql, func(ctx context.Context, tx *sql.Tx) error {
		var errTx error
		booking, errTx = c.ds.FindBookingById(ctx, tx, bookingID)
		if errTx != nil {
			return errTx
		}

		booking.Status = entities.BookingStatusCancelled
		errTx = c.ds.SaveBooking(ctx, tx, booking)
		if errTx != nil {
			return errTx
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}
