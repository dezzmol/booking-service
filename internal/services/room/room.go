package room

import (
	"context"
	"log"
	"slices"

	"booking-service/internal/entities"
	"booking-service/internal/repositories/room"
)

//go:generate mockery --disable-version-string --case=underscore --name=RoomService --structname=RoomServiceMock

const (
	ChunkSize           = 100
	NotificationMessage = "Информация о номере была обновлена, пожалуйста, просмотрите ее"
)

type (
	RoomService interface {
		CreateRooms(ctx context.Context, rooms []entities.RoomDTO) error
		Update(ctx context.Context, room entities.RoomDTO) error
	}

	Notifications interface {
		SendNotification(ctx context.Context, message string, receiver uint64) error
		SendNotificationWithQueue(ctx context.Context, message string, receiver uint64) error
	}

	Service struct {
		roomRepo      room.RoomRepo
		notifications Notifications
	}
)

func New(room room.RoomRepo, notifications Notifications) *Service {
	return &Service{roomRepo: room, notifications: notifications}
}

func (s *Service) CreateRooms(ctx context.Context, rooms []entities.RoomDTO) error {
	baseRooms := make([]entities.Room, len(rooms))

	for i, room := range rooms {
		baseRooms[i] = entities.Room{
			Number:  room.Number,
			Type:    room.Type,
			HotelID: room.HotelID,
		}
	}

	for room := range slices.Chunk(baseRooms, ChunkSize) {
		err := s.roomRepo.SaveAll(ctx, room)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Update(ctx context.Context, room entities.RoomDTO) error {
	baseRoom := &entities.Room{
		Number:  room.Number,
		Type:    room.Type,
		HotelID: room.HotelID,
	}

	err := s.roomRepo.Update(ctx, baseRoom)
	if err != nil {
		return err
	}
	go func() {
		err := s.notifications.SendNotificationWithQueue(ctx, NotificationMessage, 1)
		if err != nil {
			log.Printf("[services.RoomService.Update] notification send error: %s", err)
		}
	}()

	return err
}
