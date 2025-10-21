package storage

import (
	"context"
	"database/sql"

	"booking-service/internal/entities"
)

func (s *Storage) SaveReview(ctx context.Context, tx *sql.Tx, review entities.Review) (entities.Review, error) {
	query := `
        INSERT INTO reviews (booking_id, rating, comment)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `
	err := tx.QueryRowContext(ctx, query, review.BookingID, review.Rating, review.Comment).
		Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)
	if err != nil {
		return entities.Review{}, err
	}

	return review, nil
}
