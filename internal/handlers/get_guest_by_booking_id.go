package handlers

import (
	"context"
	"errors"

	"booking-service/internal/entities"
	"booking-service/internal/generated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) GetGuestsByBookingID(ctx context.Context, in *generated.GetGuestsByBookingIDRequest) (
	*generated.GetGuestsByBookingIDResponse, error,
) {
	guests, err := h.GuestService.FindGuestsByBookingID(ctx, uint(in.BookingId))
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrNotFound):
			return nil, status.Error(codes.NotFound, "booking not found")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	guestsResponse := make([]*generated.Guest, 0, len(guests))
	for _, g := range guests {
		guestsResponse = append(guestsResponse, &generated.Guest{
			Id:        g.ID,
			CreatedAt: timestamppb.New(g.CreatedAt),
			UpdatedAt: timestamppb.New(g.UpdatedAt),
			Name:      g.Name,
		})
	}

	return &generated.GetGuestsByBookingIDResponse{
		Guests: guestsResponse,
		Links: []*generated.BookingServiceLink{
			{
				Rel:    "self",
				Href:   "/v1/booking/{booking_id}/guests",
				Method: "get",
			},
		},
	}, nil
}
