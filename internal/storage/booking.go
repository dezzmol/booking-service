package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"booking-service/internal/entities"
)

func (s *Storage) SaveBooking(ctx context.Context, tx *sql.Tx, booking entities.Booking) error {
	queryBooking := `
        INSERT INTO bookings (room_id, guest_id, start_date, end_date, comment, status, is_paid)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT DO UPDATE SET
            room_id = EXCLUDED.room_id,
            guest_id = EXCLUDED.guest_id,
            start_date = EXCLUDED.start_date,
            end_date = EXCLUDED.end_date,
            comment = EXCLUDED.comment,
            status = EXCLUDED.status,
            is_paid = EXCLUDED.is_paid
    
    `

	err := tx.QueryRowContext(ctx,
		queryBooking,
		booking.RoomID,
		booking.StartDate,
		booking.EndDate,
		booking.Comment,
		booking.Status,
		booking.IsPaid,
	).
		Scan(&booking.ID, &booking.CreatedAt, &booking.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// FindBookingById возвращает бронирование по идентификатору, включая связанных гостей.
func (s *Storage) FindBookingById(ctx context.Context, tx *sql.Tx, bookingID uint64) (entities.Booking, error) {
	var booking entities.Booking
	query := `
        SELECT id, room_id, guest_id, start_date, end_date, comment, created_at, updated_at, status, is_paid
        FROM bookings
        WHERE id = $1
    `
	if err := tx.QueryRowContext(ctx, query, bookingID).Scan(&booking); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.Booking{}, entities.ErrNotFound
		}
		return entities.Booking{}, err
	}

	return booking, nil
}

// FindBookingByDate возвращает список бронирований, активных на заданную дату.
// Для каждого бронирования дополнительно загружаются связанные гости.
func (s *Storage) FindBookingByDate(ctx context.Context, tx *sql.Tx, startDate, endDate time.Time) ([]entities.Booking, error) {
	var bookings []entities.Booking
	query := `
        SELECT id, room_id, guest_id, start_date, end_date, comment, created_at, updated_at, status, is_paid
        FROM bookings
        WHERE start_date <= $1 AND end_date >= $2
        ORDER BY start_date
    `
	rows, err := tx.QueryContext(ctx, query, endDate, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]entities.Booking, 0)
	for rows.Next() {
		var booking entities.Booking
		errScan := rows.Scan(&booking)
		if errScan != nil {
			return nil, errScan
		}
		res = append(res, booking)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return bookings, nil
}

// DeleteBooking удаляет бронирование по идентификатору, включая связанные записи в таблице bookings_guests.
func (s *Storage) DeleteBooking(ctx context.Context, tx *sql.Tx, bookingID uint64) error {
	deleteGuestsQuery := `
        DELETE FROM bookings
        WHERE id = $1
    `
	if _, err := tx.ExecContext(ctx, deleteGuestsQuery, bookingID); err != nil {
		return err
	}

	return nil
}

// FindBookingByRoomIDAndDate возвращает список бронирований для заданной комнаты.
// Для каждого бронирования дополнительно загружаются связанные гости.
func (s *Storage) FindBookingByRoomIDAndDate(
	ctx context.Context, tx *sql.Tx, roomID uint64, startDate, endDate time.Time,
) ([]entities.Booking, error) {
	var bookings []entities.Booking
	query := `
        SELECT id, room_id, guest_id, start_date, end_date, comment, created_at, updated_at, status, is_paid
        FROM bookings
        WHERE room_id = $1 AND start_date <= $2 AND end_date >= $3
        ORDER BY start_date
    `
	rows, err := tx.QueryContext(ctx, query, roomID, endDate, startDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]entities.Booking, 0)
	for rows.Next() {
		var booking entities.Booking
		errScan := rows.Scan(&booking)
		if errScan != nil {
			return nil, errScan
		}
		res = append(res, booking)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return bookings, nil
}

func (s *Storage) IsRoomAvailableForBooking(ctx context.Context, tx *sql.Tx, roomID uint64, startDate, endDate time.Time) (bool, error) {
	query := `SELECT 
    NOT EXISTS (
        SELECT 1 
        FROM bookings 
        WHERE room_id = $1  -- ID конкретной комнаты
          AND status IN (1, 3)  -- активные статусы бронирований (уточните по вашей системе)
          AND start_date < $2 -- конечная дата желаемого бронирования
          AND end_date > $3    -- начальная дата желаемого бронирования
    ) as is_available;`

	var exist bool
	if err := tx.QueryRowContext(ctx, query, roomID, startDate, endDate).Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}
