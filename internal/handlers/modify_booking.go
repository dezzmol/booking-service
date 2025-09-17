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

func (h *Handler) ModifyBooking(ctx context.Context, in *generated.ModifyBookingRequest) (
	*generated.ModifyBookingResponse, error,
) {
	booking, err := h.BookingService.ModifyBooking(ctx, in.BookingId, in.StartDate.AsTime(), in.EndDate.AsTime())
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrNotFound):
			return nil, status.Error(codes.NotFound, "booking not found")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	guests := make([]*generated.Guest, 0, len(booking.Guests))
	for _, g := range booking.Guests {
		guests = append(guests, &generated.Guest{
			Id:        g.ID,
			CreatedAt: timestamppb.New(g.CreatedAt),
			UpdatedAt: timestamppb.New(g.UpdatedAt),
			Name:      g.Name,
		})
	}

	return &generated.ModifyBookingResponse{
		Booking: &generated.Booking{
			Id:            booking.ID,
			CreatedAt:     timestamppb.New(booking.CreatedAt),
			UpdatedAt:     timestamppb.New(booking.UpdatedAt),
			RoomId:        booking.RoomID,
			StartDate:     timestamppb.New(booking.StartDate),
			EndDate:       timestamppb.New(booking.EndDate),
			Comment:       booking.Comment,
			Status:        booking.Status,
			PaymentStatus: booking.PaymentStatus,
			Guests:        guests,
		},
		Links: []*generated.BookingServiceLink{
			{
				Rel:    "self",
				Href:   "/v1/booking",
				Method: "put",
			},
		},
	}, nil
}
