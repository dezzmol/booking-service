package room

import (
	"context"
	"fmt"

	"booking-service/internal/entities"

	"github.com/jmoiron/sqlx"
)

//go:generate mockery --disable-version-string --case=underscore --name=RoomRepo --structname=RoomRepoMock
type RoomRepo interface {
	// FindById ищет комнату по идентификатору и возвращает её.
	FindById(ctx context.Context, roomId int64) (*entities.Room, error)

	// Save сохраняет новую комнату в базу данных.
	// После успешного сохранения в структуру room будут записаны сгенерированные поля (id, created_at, updated_at).
	Save(ctx context.Context, room *entities.Room) error

	// SaveAll сохраняет список комнат в базу данных в рамках одной транзакции.
	SaveAll(ctx context.Context, rooms []entities.Room) error

	// Update обновляет данные комнаты в базе данных.
	// Поле updated_at обновляется на текущее время.
	Update(ctx context.Context, room *entities.Room) error
}

type RoomRepository struct {
	db *sqlx.DB
}

func NewRoomRepository(db *sqlx.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) FindById(ctx context.Context, roomId int64) (*entities.Room, error) {
	var room entities.Room
	query := `
		SELECT id, number, type, hotel_id, created_at, updated_at
		FROM rooms
		WHERE id = $1
	`
	if err := r.db.GetContext(ctx, &room, query, roomId); err != nil {
		return nil, fmt.Errorf("[RoomRepository]: FindById: %w ", err)
	}
	return &room, nil
}

func (r *RoomRepository) Save(ctx context.Context, room *entities.Room) error {
	query := `
		INSERT INTO rooms (number, type, hotel_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowxContext(ctx, query, room.Number, room.Type, room.HotelID).
		Scan(&room.ID, &room.CreatedAt, &room.UpdatedAt)
}

func (r *RoomRepository) SaveAll(ctx context.Context, rooms []entities.Room) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return fmt.Errorf("[RoomRepository]: SaveAll: %w ", err)
	}
	query := `
		INSERT INTO rooms (number, type, hotel_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	// Итерация по всем комнатам для сохранения.
	for i := range rooms {
		if err := tx.QueryRowxContext(ctx, query, rooms[i].Number, rooms[i].Type, rooms[i].HotelID).
			Scan(&rooms[i].ID, &rooms[i].CreatedAt, &rooms[i].UpdatedAt); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *RoomRepository) Update(ctx context.Context, room *entities.Room) error {
	query := `
		UPDATE rooms
		SET number = $1,
		    type = $2,
		    hotel_id = $3,
		    updated_at = NOW()
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, room.Number, room.Type, room.HotelID, room.ID)
	return err
}
