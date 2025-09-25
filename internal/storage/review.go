package storage

import (
	"context"
	"database/sql"

	"booking-service/internal/entities"
)

func (s *Storage) SaveReview(ctx context.Context, tx *sql.Tx, review entities.Review) error {
	query := `
        INSERT INTO reviews (booking_id, rating, comment)
        VALUES ($1, $2, $3)
        RETURNING id, created_at, updated_at
    `
	return tx.QueryRowContext(ctx, query, review.BookingID, review.Rating, review.Comment).
		Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)
}

func (s *Storage) FindByGuest(ctx context.Context, tx *sql.Tx, guestID uint64) ([]entities.Review, error) {
	var reviews []entities.Review
	query := `
        SELECT id, booking_id, rating, comment, created_at, updated_at
        FROM reviews
        WHERE guest_id = $1
        ORDER BY date DESC
    `
	rows, err := tx.QueryContext(ctx, query, guestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var review entities.Review
		err := rows.Scan(
			&review.ID,
			&review.BookingID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
			&review.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return reviews, nil
}
