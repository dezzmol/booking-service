package controllers

import (
	"context"
	"database/sql"

	"booking-service/internal/entities"
	"booking-service/internal/storage"
)

func (c *Controller) CreateHotel(ctx context.Context, hotelName string) (res entities.Hotel, err error) {
	if err = storage.WithWriteTransaction(ctx, c.sql, func(ctx context.Context, tx *sql.Tx) (errTx error) {
		if res, errTx = c.ds.SaveHotel(ctx, tx, entities.Hotel{
			Name: hotelName,
		}); errTx != nil {
			return errTx
		}
		return nil
	}); err != nil {
		return res, err
	}

	return res, nil
}
