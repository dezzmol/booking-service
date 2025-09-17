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

func (h *Handler) GetEmployee(ctx context.Context, in *generated.GetEmployeeRequest) (*generated.GetEmployeeResponse, error) {
	employee, err := h.EmployeeService.GetEmployee(ctx, in.EmployeeId)
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrNotFound):
			return nil, status.Error(codes.NotFound, "Employee not found")
		default:
			return nil, status.Error(codes.Internal, "Failed to get employee")
		}
	}

	return &generated.GetEmployeeResponse{
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
				Method: "get",
			},
		},
	}, nil
}
