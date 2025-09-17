package handlers

import (
	"context"
	"log"

	"booking-service/internal/generated"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) AssignEmployee(ctx context.Context, in *generated.AssignEmployeeRequest) (
	*generated.AssignEmployeeResponse, error,
) {
	log.Printf("[handlers.AssignEmployee] received request for assign employee: %v", in)

	request, err := h.HousekeepingService.AssignEmployee(ctx, in.GetRequestId(), in.GetEmployeeId())
	if err != nil {
		return nil, err
	}

	return &generated.AssignEmployeeResponse{
		Request: &generated.HousekeepingRequest{
			Id:          request.ID,
			CreatedAt:   timestamppb.New(request.CreatedAt),
			UpdatedAt:   timestamppb.New(request.UpdatedAt),
			RoomId:      request.RoomID,
			RequestTime: timestamppb.New(request.RequestTime),
			Status:      request.Status,
		},
		Links: nil,
	}, nil
}
