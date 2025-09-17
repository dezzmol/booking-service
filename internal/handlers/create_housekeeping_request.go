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

func (h *Handler) CreateHouseKeepingRequest(ctx context.Context, in *generated.CreateHouseKeepingRequestRequest) (
	*generated.CreateHouseKeepingRequestResponse, error,
) {
	request, err := h.HousekeepingService.CreateHouseKeepingRequest(ctx, entities.HouseKeepingRequestDTO{
		RoomID:      in.RoomId,
		RequestTime: in.RequestTime.AsTime(),
	})
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrNotFound):
			return nil, status.Error(codes.NotFound, "room not found")
		default:
			return nil, status.Error(codes.Internal, "failed to create the housekeeping request "+err.Error())
		}
	}

	return &generated.CreateHouseKeepingRequestResponse{
		Request: &generated.HousekeepingRequest{
			Id:          request.ID,
			CreatedAt:   timestamppb.New(request.CreatedAt),
			UpdatedAt:   timestamppb.New(request.UpdatedAt),
			RoomId:      request.RoomID,
			RequestTime: timestamppb.New(request.RequestTime),
			Status:      request.Status,
		},
		Links: []*generated.BookingServiceLink{
			{
				Rel:    "self",
				Href:   "/v1/housekeeping",
				Method: "post",
			},
		},
	}, nil
}
