package guest

import (
	"context"

	"booking-service/internal/entities"
	"booking-service/internal/repositories/guest"
)

//go:generate mockery --disable-version-string --case=underscore --name=GuestService --structname=GuestServiceMock
type GuestService interface {
	FindGuestsByBookingID(ctx context.Context, bookingID uint) ([]entities.Guest, error)
	CreateGuest(ctx context.Context, guest entities.GuestDTO) (entities.Guest, error)
}

type Service struct {
	repo guest.GuestRepo
}

func NewService(repo guest.GuestRepo) *Service {
	return &Service{repo: repo}
}

func (s *Service) FindGuestsByBookingID(ctx context.Context, bookingID uint) ([]entities.Guest, error) {
	return s.repo.FindByBookingID(ctx, bookingID)
}

func (s *Service) CreateGuest(ctx context.Context, input entities.GuestDTO) (entities.Guest, error) {
	if input.Name == "" {
		return entities.Guest{}, entities.ErrNameIsRequired
	}
	if len([]rune(input.Name)) > 40 {
		return entities.Guest{}, entities.ErrNameIsTooLong
	}
	guest := entities.Guest{
		Name: input.Name,
	}

	return s.repo.SaveAndReturnIt(ctx, &guest)
}
