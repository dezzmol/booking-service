package booking

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"booking-service/internal/entities"

	"github.com/jmoiron/sqlx"
)

//go:generate mockery --disable-version-string --case=underscore --name=BookingRepo --structname=BookingRepoMock
type BookingRepo interface {
	// Save сохраняет новое бронирование и связанные с ним гостей.
	// Вставка выполняется в рамках транзакции: сначала добавляется запись в bookings, затем для каждого гостя создаётся связь в bookings_guests.
	Save(ctx context.Context, booking *entities.Booking) error

	// FindById возвращает бронирование по идентификатору, включая связанных гостей.
	FindById(ctx context.Context, bookingID uint64) (*entities.Booking, error)

	// FindByDate возвращает список бронирований, активных на заданную дату.
	// Для каждого бронирования дополнительно загружаются связанные гости.
	FindByDate(ctx context.Context, date time.Time) ([]entities.Booking, error)

	// Update обновляет данные бронирования, а при наличии обновляемого списка гостей – сначала удаляет старые связи, затем добавляет новые.
	Update(ctx context.Context, booking *entities.Booking) error

	// Delete удаляет бронирование по идентификатору, включая связанные записи в таблице bookings_guests.
	Delete(ctx context.Context, bookingID uint64) error

	// FindByRoomID возвращает список бронирований для заданной комнаты.
	// Для каждого бронирования дополнительно загружаются связанные гости.
	FindByRoomID(ctx context.Context, roomID uint64) ([]entities.Booking, error)

	IsRoomAvailable(ctx context.Context, roomID uint64, startDate, endDate time.Time) (bool, error)

	CanReschedule(ctx context.Context, bookingID uint64, newStartDate, newEndDate time.Time) (bool, error)
}

type BookingRepository struct {
	db *sqlx.DB
}

func NewBookingRepository(db *sqlx.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Save(ctx context.Context, booking *entities.Booking) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	queryBooking := `
        INSERT INTO bookings (room_id, start_date, end_date, comment)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at
    `
	err = tx.QueryRowxContext(ctx, queryBooking, booking.RoomID, booking.StartDate, booking.EndDate, booking.Comment).
		Scan(&booking.ID, &booking.CreatedAt, &booking.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Если переданы связанные гости, добавляем записи в таблицу bookings_guests.
	if len(booking.Guests) > 0 {
		queryGuest := `
            INSERT INTO bookings_guests (booking_id, guests_id)
            VALUES ($1, $2)
        `
		for _, guest := range booking.Guests {
			if _, err := tx.Exec(queryGuest, booking.ID, guest.ID); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

// FindById возвращает бронирование по идентификатору, включая связанных гостей.
func (r *BookingRepository) FindById(ctx context.Context, bookingID uint64) (*entities.Booking, error) {
	var booking entities.Booking
	query := `
        SELECT id, room_id, start_date, end_date, comment, created_at, updated_at
        FROM bookings
        WHERE id = $1
    `
	if err := r.db.GetContext(ctx, &booking, query, bookingID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entities.ErrNotFound
		}
		return nil, err
	}

	// Получаем связанных гостей через соединение с таблицей bookings_guests.
	guestQuery := `
        SELECT g.id, g.name, g.created_at, g.updated_at
        FROM guests g
        JOIN bookings_guests bg ON g.id = bg.guests_id
        WHERE bg.booking_id = $1
    `
	if err := r.db.SelectContext(ctx, &booking.Guests, guestQuery, booking.ID); err != nil {
		return nil, err
	}

	return &booking, nil
}

// FindByDate возвращает список бронирований, активных на заданную дату.
// Для каждого бронирования дополнительно загружаются связанные гости.
func (r *BookingRepository) FindByDate(ctx context.Context, date time.Time) ([]entities.Booking, error) {
	var bookings []entities.Booking
	query := `
        SELECT id, room_id, start_date, end_date, comment, created_at, updated_at
        FROM bookings
        WHERE start_date <= $1 AND end_date >= $1
        ORDER BY start_date
    `
	if err := r.db.SelectContext(ctx, &bookings, query, date); err != nil {
		return nil, err
	}

	guestQuery := `
        SELECT g.id, g.name, g.created_at, g.updated_at
        FROM guests g
        JOIN bookings_guests bg ON g.id = bg.guests_id
        WHERE bg.booking_id = $1
    `
	for i := range bookings {
		if err := r.db.SelectContext(ctx, &bookings[i].Guests, guestQuery, bookings[i].ID); err != nil {
			return nil, err
		}
	}

	return bookings, nil
}

// Update обновляет данные бронирования, а при наличии обновляемого списка гостей – сначала удаляет старые связи, затем добавляет новые.
func (r *BookingRepository) Update(ctx context.Context, booking *entities.Booking) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	queryUpdate := `
        UPDATE bookings
        SET room_id = $1,
            start_date = $2,
            end_date = $3,
            comment = $4
        WHERE id = $5
    `
	_, err = tx.ExecContext(ctx, queryUpdate, booking.RoomID, booking.StartDate, booking.EndDate, booking.Comment, booking.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Обновляем связи с гостями, если список гостей не nil.
	if booking.Guests != nil {
		deleteQuery := `
            DELETE FROM bookings_guests
            WHERE booking_id = $1
        `
		if _, err = tx.ExecContext(ctx, deleteQuery, booking.ID); err != nil {
			tx.Rollback()
			return err
		}

		insertQuery := `
            INSERT INTO bookings_guests (booking_id, guests_id)
            VALUES ($1, $2)
        `
		for _, guest := range booking.Guests {
			if _, err = tx.ExecContext(ctx, insertQuery, booking.ID, guest.ID); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

// Delete удаляет бронирование по идентификатору, включая связанные записи в таблице bookings_guests.
func (r *BookingRepository) Delete(ctx context.Context, bookingID uint64) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	deleteGuestsQuery := `
        DELETE FROM bookings_guests
        WHERE booking_id = $1
    `
	if _, err = tx.ExecContext(ctx, deleteGuestsQuery, bookingID); err != nil {
		tx.Rollback()
		return err
	}

	deleteBookingQuery := `
        DELETE FROM bookings
        WHERE id = $1
    `
	if _, err = tx.ExecContext(ctx, deleteBookingQuery, bookingID); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// FindByRoomID возвращает список бронирований для заданной комнаты.
// Для каждого бронирования дополнительно загружаются связанные гости.
func (r *BookingRepository) FindByRoomID(ctx context.Context, roomID uint64) ([]entities.Booking, error) {
	var bookings []entities.Booking
	query := `
        SELECT id, room_id, start_date, end_date, comment, created_at, updated_at
        FROM bookings
        WHERE room_id = $1
        ORDER BY start_date
    `
	if err := r.db.SelectContext(ctx, &bookings, query, roomID); err != nil {
		return nil, err
	}

	guestQuery := `
        SELECT g.id, g.name, g.created_at, g.updated_at
        FROM guests g
        JOIN bookings_guests bg ON g.id = bg.guests_id
        WHERE bg.booking_id = $1
    `
	for i := range bookings {
		if err := r.db.SelectContext(ctx, &bookings[i].Guests, guestQuery, bookings[i].ID); err != nil {
			return nil, err
		}
	}

	return bookings, nil
}

func (r *BookingRepository) IsRoomAvailable(ctx context.Context, roomID uint64, startDate, endDate time.Time) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM bookings 
		WHERE room_id = $1 AND status NOT IN ('pending', 'confirmed', 'checked-in')
		  AND NOT (end_date <= $2 OR start_date >= $3)
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, roomID, startDate, endDate).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *BookingRepository) CanReschedule(ctx context.Context, bookingID uint64, newStartDate, newEndDate time.Time) (bool, error) {
	var canReschedule bool
	query := `
        SELECT COUNT(*) = 0 FROM bookings b
        WHERE b.room_id = (SELECT room_id FROM bookings WHERE id = $1)
          AND b.id <> $1
          AND b.status IN ('pending', 'confirmed', 'checked-in')
          AND ($2, $3) OVERLAPS (b.start_date, b.end_date)
    `
	err := r.db.QueryRowxContext(ctx, query, bookingID, newStartDate, newEndDate).Scan(&canReschedule)
	if err != nil {
		return false, err
	}
	return canReschedule, nil
}
