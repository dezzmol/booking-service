package app

import (
	"context"
	"log"

	"booking-service/internal/entities"
	"booking-service/internal/generated"
)

func (h *Handler) CreateRoom(ctx context.Context, in *generated.CreateRoomRequest) (*generated.CreateRoomResponse, error) {
	log.Printf("[CreateRoom]: Handling CreateRoom request: %+v", in)

	err := h.bookingController.CreateRooms(ctx, h.convertRooms(in.GetDto()))
	if err != nil {
		return nil, err
	}

	return &generated.CreateRoomResponse{}, nil
}

func (h *Handler) convertRooms(in []*generated.CreateRoomRequest_DTO) []entities.RoomDTO {
	rooms := make([]entities.RoomDTO, 0, len(in))
	for _, room := range in {
		rooms = append(rooms, entities.RoomDTO{
			Number:  room.Number,
			Type:    room.Type,
			HotelID: room.HotelId,
		})
	}

	return rooms
}
