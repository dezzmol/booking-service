package app

import (
	"context"
	"errors"
	"log"

	"booking-service/internal/entities"
	"booking-service/internal/generated"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Handler) CancelBooking(ctx context.Context, in *generated.CancelBookingRequest) (
	*generated.CancelBookingResponse, error,
) {
	log.Printf("[handlers.CancelBooking] received request with: %+v", in)

	err := h.bookingController.CancelBooking(ctx, in.BookingId)
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrNotFound):
			return nil, status.Errorf(codes.NotFound, "booking not found: %v", err)
		default:
			return nil, status.Errorf(codes.Internal, "internal error: %v", err)
		}
	}

	return &generated.CancelBookingResponse{}, nil
}
