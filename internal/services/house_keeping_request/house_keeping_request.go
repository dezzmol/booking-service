package house_keeping_request

import (
	"context"
	"time"

	"booking-service/internal/entities"
	"booking-service/internal/repositories/house_keeping_request"
)

//go:generate mockery --disable-version-string --case=underscore --name=HouseKeepingRequestService --structname=HouseKeepingRequestServiceMock

const (
	REQUESTED = "requested"
	ACCEPTED  = "accepted"
)

type HouseKeepingRequestService interface {
	CreateHouseKeepingRequest(ctx context.Context, request entities.HouseKeepingRequestDTO) (entities.HousekeepingRequest, error)
	AssignEmployee(ctx context.Context, requestID uint64, EmployeeID uint64) (entities.HousekeepingRequest, error)
}

type Service struct {
	repo house_keeping_request.HouseKeepingRequestRepo
}

func New(repo house_keeping_request.HouseKeepingRequestRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateHouseKeepingRequest(ctx context.Context, request entities.HouseKeepingRequestDTO) (entities.HousekeepingRequest, error) {
	req := entities.HousekeepingRequest{
		RoomID:      request.RoomID,
		RequestTime: time.Time{},
		Status:      REQUESTED,
	}

	err := s.repo.Save(ctx, &req)
	if err != nil {
		return entities.HousekeepingRequest{}, err
	}
	return req, nil
}

func (s *Service) AssignEmployee(ctx context.Context, requestID uint64, EmployeeID uint64) (entities.HousekeepingRequest, error) {
	req, err := s.repo.FindByID(ctx, requestID)
	if err != nil {
		return entities.HousekeepingRequest{}, err
	}
	req.Status = ACCEPTED
	err = s.repo.Update(ctx, *req)
	if err != nil {
		return entities.HousekeepingRequest{}, err
	}

	return *req, nil
}
