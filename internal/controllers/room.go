package controllers

import (
	"context"
	"database/sql"

	"booking-service/internal/entities"
	"booking-service/internal/storage"
)

//go:generate mockery --disable-version-string --case=underscore --name=RoomService --structname=RoomServiceMock

const (
	ChunkSize           = 100
	NotificationMessage = "Информация о номере была обновлена, пожалуйста, просмотрите ее"
)

func (c *Controller) CreateRooms(ctx context.Context, rooms []entities.RoomDTO) error {
	baseRooms := make([]entities.Room, len(rooms))

	for i, room := range rooms {
		baseRooms[i] = entities.Room{
			Number:  room.Number,
			Type:    room.Type,
			HotelID: room.HotelID,
		}
	}

	err := storage.WithWriteTransaction(ctx, c.sql, func(ctx context.Context, tx *sql.Tx) error {
		txErr := c.ds.SaveAllRooms(ctx, tx, baseRooms)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
