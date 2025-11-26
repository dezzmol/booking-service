package app

import (
	"context"
	"log"

	"booking-service/internal/generated"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateHotel(ctx context.Context, req *generated.CreateHotelRequest) (*generated.CreateHotelResponse, error) {
	log.Printf("[handlers.CreateHotel] received request: %v", req)

	hotel, err := h.bookingController.CreateHotel(ctx, req.GetName())
	if err != nil {
		return nil, err
	}

	return &generated.CreateHotelResponse{
		Hotel: &generated.Hotel{
			Id:        hotel.ID,
			CreatedAt: timestamppb.New(hotel.CreatedAt),
			UpdatedAt: timestamppb.New(hotel.UpdatedAt),
			Name:      hotel.Name,
		},
	}, nil
}
