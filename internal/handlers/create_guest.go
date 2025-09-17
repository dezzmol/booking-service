package handlers

import (
	"context"
	"errors"
	"log"

	"booking-service/internal/entities"
	"booking-service/internal/generated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) CreateGuest(ctx context.Context, in *generated.CreateGuestRequest) (
	*generated.CreateGuestResponse, error,
) {
	log.Printf("[handlers.CreateGuest]: received request: %v", in)

	guest, err := h.GuestService.CreateGuest(ctx, entities.GuestDTO{
		Name: in.GetName(),
	})
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrNameIsRequired) ||
			errors.Is(err, entities.ErrNameIsTooLong):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &generated.CreateGuestResponse{
		Guest: &generated.Guest{
			Id:        guest.ID,
			CreatedAt: timestamppb.New(guest.CreatedAt),
			UpdatedAt: timestamppb.New(guest.UpdatedAt),
			Name:      guest.Name,
		},
		Links: []*generated.BookingServiceLink{
			{
				Rel:    "self",
				Href:   "/v1/guests",
				Method: "post",
			},
		},
	}, nil
}
