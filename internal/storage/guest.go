package storage

import (
	"context"
	"database/sql"
	"errors"

	"booking-service/internal/entities"
)

func (s *Storage) SaveGuestAndReturnIt(ctx context.Context, tx *sql.Tx, input entities.Guest) (entities.Guest, error) {
	var guest entities.Guest
	query := `SELECT g.id, g.name, g.created_at, g.updated_at FROM guests g WHERE g.name = $1`

	err := tx.QueryRowContext(ctx, query, input.Name).Scan(&guest)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return entities.Guest{}, err
	}

	if err == nil {
		return guest, nil
	}

	insertQuery := `INSERT INTO guests (name) VALUES ($1) RETURNING id, name, created_at, updated_at`
	err = tx.QueryRowContext(ctx, insertQuery, input.Name).Scan(&guest.ID, &guest.Name, &guest.CreatedAt, &guest.UpdatedAt)
	if err != nil {
		return entities.Guest{}, err
	}

	return guest, nil
}
