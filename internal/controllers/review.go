package controllers

import (
	"context"
	"database/sql"
	"time"

	"booking-service/internal/entities"
	"booking-service/internal/storage"
)

func (c *Controller) SubmitReview(ctx context.Context, reviewDTO entities.ReviewDTO) (entities.Review, error) {
	var reviewRes entities.Review
	err := storage.WithNoTransaction(ctx, c.sql, func(ctx context.Context, tx *sql.Tx) (errTx error) {
		booking, errTx := c.ds.FindBookingById(ctx, tx, reviewDTO.BookingID)
		if errTx != nil {
			return errTx
		}

		review := entities.Review{
			BookingID: booking.ID,
			Rating:    reviewDTO.Rating,
			Comment:   reviewDTO.Comment,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}

		reviewRes, errTx = c.ds.SaveReview(ctx, tx, review)
		if errTx != nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		return entities.Review{}, err
	}

	return reviewRes, nil
}
