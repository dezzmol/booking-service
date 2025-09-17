package handlers

import (
	"context"
	"log"

	"booking-service/internal/entities"
	"booking-service/internal/generated"
)

func (h *Handler) CreateRoom(ctx context.Context, in *generated.CreateRoomRequest) (*generated.VoidResponse, error) {
	log.Printf("[CreateRoom]: Handling CreateRoom request: %+v", in)

	err := h.RoomService.CreateRooms(ctx, h.convertRooms(in.GetRooms()))
	if err != nil {
		return nil, err
	}

	return &generated.VoidResponse{
		Links: []*generated.BookingServiceLink{
			{
				Rel:    "self",
				Href:   "/v1/room",
				Method: "post",
			},
		},
	}, nil
}

func (h *Handler) convertRooms(in []*generated.RoomDTO) []entities.RoomDTO {
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
