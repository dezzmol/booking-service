package guest

import (
	"context"
	"database/sql"
	"errors"

	"booking-service/internal/entities"

	"github.com/jmoiron/sqlx"
)

//go:generate mockery --disable-version-string --case=underscore --name=GuestRepo --structname=GuestRepoMock
type GuestRepo interface {
	// FindByBookingID возвращает список гостей, связанных с заданным bookingID.
	FindByBookingID(ctx context.Context, bookingID uint) ([]entities.Guest, error)
	SaveAndReturnIt(ctx context.Context, guest *entities.Guest) (entities.Guest, error)
}

type GuestRepository struct {
	db *sqlx.DB
}

func NewGuestRepository(db *sqlx.DB) *GuestRepository {
	return &GuestRepository{db: db}
}

func (r *GuestRepository) FindByBookingID(ctx context.Context, bookingID uint) ([]entities.Guest, error) {
	var guests []entities.Guest
	query := `
		SELECT g.id, g.name, g.created_at, g.updated_at
		FROM guests g
		JOIN bookings_guests bg ON g.id = bg.guests_id
		WHERE bg.booking_id = $1
	`
	err := r.db.SelectContext(ctx, &guests, query, bookingID)
	return guests, err
}

func (r *GuestRepository) SaveAndReturnIt(ctx context.Context, input *entities.Guest) (entities.Guest, error) {
	var guest entities.Guest
	query := `SELECT g.id, g.name, g.created_at, g.updated_at FROM guests g WHERE g.name = $1`

	err := r.db.SelectContext(ctx, &guest, query, input.Name)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return entities.Guest{}, err
	}

	if err == nil {
		return guest, nil
	}

	insertQuery := `INSERT INTO guests (name) VALUES ($1) RETURNING id, name, created_at, updated_at`
	err = r.db.QueryRowContext(ctx, insertQuery, input.Name).Scan(&guest.ID, &guest.Name, &guest.CreatedAt, &guest.UpdatedAt)
	if err != nil {
		return entities.Guest{}, err
	}

	return guest, nil
}
