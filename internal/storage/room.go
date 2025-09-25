package storage

import (
	"context"
	"database/sql"
	"fmt"

	"booking-service/internal/entities"
)

func (s *Storage) FindRoomById(ctx context.Context, tx *sql.Tx, roomId int64) (entities.Room, error) {
	var room entities.Room
	query := `
		SELECT id, number, type, hotel_id, created_at, updated_at
		FROM rooms
		WHERE id = $1
	`
	if err := tx.QueryRowContext(ctx, query, roomId).Scan(&room); err != nil {
		return entities.Room{}, fmt.Errorf("[RoomRepository]: FindById: %w ", err)
	}
	return room, nil
}

func (s *Storage) SaveRoom(ctx context.Context, tx *sql.Tx, room *entities.Room) error {
	query := `
		INSERT INTO rooms (number, type, hotel_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	if err := tx.QueryRowContext(ctx, query, room.Number, room.Type, room.HotelID).
		Scan(&room.ID, &room.CreatedAt, &room.UpdatedAt); err != nil {
		return fmt.Errorf("[RoomRepository]: Save: %w ", err)
	}

	return nil
}

func (s *Storage) SaveAllRooms(ctx context.Context, tx *sql.Tx, rooms []entities.Room) error {
	query := `
		INSERT INTO rooms (number, type, hotel_id)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING
	`
	// Итерация по всем комнатам для сохранения.
	for i := range rooms {
		if _, err := tx.ExecContext(ctx, query, rooms[i].Number, rooms[i].Type, rooms[i].HotelID); err != nil {
			return err
		}
	}
	return nil
}
