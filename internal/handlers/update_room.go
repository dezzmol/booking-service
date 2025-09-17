package handlers

import (
	"context"
	"errors"
	"log"

	"booking-service/internal/entities"
	"booking-service/internal/generated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) UpdateRoom(ctx context.Context, in *generated.UpdateRoomRequest) (*generated.UpdateRoomResponse, error) {
	log.Printf("[UpdateRoom]: Handling UpdateRoom request: %+v", in)

	err := h.RoomService.Update(ctx, entities.RoomDTO{
		Number:  in.Number,
		Type:    in.Type,
		HotelID: in.HotelId,
	})
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return nil, nil
}
