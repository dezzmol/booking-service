package storage

import (
	"context"
	"database/sql"
	"time"

	"booking-service/internal/entities"
)

func (s *Storage) SaveHotel(ctx context.Context, tx *sql.Tx, hotel entities.Hotel) (entities.Hotel, error) {
	query := `INSERT INTO hotels (name, created_at, updated_at)
				VALUES ($1, $2, $3)
				RETURNING id`

	err := tx.QueryRowContext(ctx,
		query,
		hotel.Name,
		time.Now().UTC(),
		time.Now().UTC(),
	).Scan(&hotel.ID)
	if err != nil {
		return hotel, err
	}

	return hotel, nil
}

func (s *Storage) FindHotelByID(ctx context.Context, tx *sql.Tx, id uint64) (entities.Hotel, error) {
	var hotel entities.Hotel
	query := `SELECT (id, name, created_at, updated_at) FROM hotels WHERE id = $1`

	err := tx.QueryRowContext(ctx, query, id).Scan(
		&hotel.ID,
		&hotel.Name,
		&hotel.CreatedAt,
		&hotel.UpdatedAt,
	)
	if err != nil {
		return entities.Hotel{}, err
	}

	return hotel, nil
}
