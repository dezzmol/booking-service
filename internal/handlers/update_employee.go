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

func (h *Handler) UpdateEmployee(ctx context.Context, in *generated.UpdateEmployeeRequest) (
	*generated.UpdateEmployeeResponse, error,
) {
	employee, err := h.EmployeeService.UpdateEmployee(ctx, in.EmployeeId, entities.EmployeeDTO{
		Name: in.Name,
		Role: in.Role,
	})
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrNotFound):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &generated.UpdateEmployeeResponse{
		Employee: &generated.Employee{
			Id:        employee.ID,
			CreatedAt: timestamppb.New(employee.CreatedAt),
			UpdatedAt: timestamppb.New(employee.UpdatedAt),
			Name:      employee.Name,
			Role:      employee.Role,
		},
		Links: []*generated.BookingServiceLink{
			{
				Rel:    "self",
				Href:   "/v1/employee/{employee_id}",
				Method: "put",
			},
		},
	}, nil
}
