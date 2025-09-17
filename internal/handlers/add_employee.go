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

func (h *Handler) AddEmployee(ctx context.Context, in *generated.AddEmployeeRequest) (*generated.AddEmployeeResponse, error) {
	log.Printf("[handlers.AddEmployee] received request: %s", in)
	employee, err := h.EmployeeService.AddEmployee(ctx, entities.EmployeeDTO{
		Name: in.GetName(),
		Role: in.GetRole(),
	})
	if err != nil {
		switch {
		case errors.Is(err, entities.ErrInvalidName):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "internal error: %v", err)
		}
	}

	return &generated.AddEmployeeResponse{
		Employee: &generated.Employee{
			Id:        employee.ID,
			CreatedAt: timestamppb.New(employee.CreatedAt),
			UpdatedAt: timestamppb.New(employee.UpdatedAt),
			Name:      employee.Name,
			Role:      employee.Role,
		},
	}, nil
}
