package app

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

func (h *Handler) CreateBooking(ctx context.Context, in *generated.CreateBookingRequest) (
	*generated.CreateBookingResponse, error,
) {
	log.Printf("[handlers.CreateBooking] received request: %v", in)

	booking, err := h.bookingController.CreateBooking(ctx, h.makeBookingDTO(in))
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrStartDateIsAfterEndDate):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, entities.ErrRoomNotAvailable):
			return nil, status.Error(codes.InvalidArgument, "room is not available")
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &generated.CreateBookingResponse{
		Booking: h.makeBookingToResponse(booking),
	}, nil
}

func (h *Handler) makeBookingDTO(in *generated.CreateBookingRequest) entities.CreateBookingDTO {
	guests := make([]entities.GuestDTO, 0, len(in.Guests))
	for _, g := range in.Guests {
		guests = append(guests, entities.GuestDTO{
			Name: g.Name,
		})
	}

	return entities.CreateBookingDTO{
		RoomID:    in.RoomId,
		StartDate: in.StartDate.AsTime(),
		EndDate:   in.EndDate.AsTime(),
		Comment:   in.Comment,
		Guests:    guests,
	}
}

func (h *Handler) makeBookingToResponse(in entities.Booking) *generated.Booking {

	return &generated.Booking{
		Id:        in.ID,
		CreatedAt: timestamppb.New(in.CreatedAt),
		UpdatedAt: timestamppb.New(in.UpdatedAt),
		RoomId:    in.RoomID,
		StartDate: timestamppb.New(in.StartDate),
		EndDate:   timestamppb.New(in.EndDate),
		Comment:   in.Comment,
		Status:    generated.BookingStatus(in.Status),
	}
}

func (h *Handler) makeGuestToResponse(in entities.Guest) *generated.Guest {
	return &generated.Guest{
		Id:        in.ID,
		CreatedAt: timestamppb.New(in.CreatedAt),
		UpdatedAt: timestamppb.New(in.UpdatedAt),
		Name:      in.Name,
	}
}
