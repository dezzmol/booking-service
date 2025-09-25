package controllers

import (
	"context"
	"database/sql"

	"booking-service/internal/entities"
	"booking-service/internal/storage"
)

func (c *Controller) FindGuestsByBookingID(ctx context.Context, bookingID uint) ([]entities.Guest, error) {
	return s.repo.FindByBookingID(ctx, bookingID)
}

func (c *Controller) CreateGuest(ctx context.Context, input entities.GuestDTO) (entities.Guest, error) {
	if input.Name == "" {
		return entities.Guest{}, entities.ErrNameIsRequired
	}
	if len([]rune(input.Name)) > 40 {
		return entities.Guest{}, entities.ErrNameIsTooLong
	}
	guest := entities.Guest{
		Name: input.Name,
	}

	err := storage.WithWriteTransaction(ctx, c.sql, func(ctx context.Context, tx *sql.Tx) error {
		var errTx error
		guest, errTx = c.ds.SaveGuestAndReturnIt(ctx, tx, guest)
		if errTx != nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		return entities.Guest{}, err
	}

	return guest, nil
}
