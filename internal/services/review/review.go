package review

import (
	"context"

	"booking-service/internal/entities"
	"booking-service/internal/repositories/booking"
	"booking-service/internal/repositories/review"
)

//go:generate mockery --disable-version-string --case=underscore --name=ReviewService --structname=ReviewServiceMock

type ReviewService interface {
	SubmitReview(ctx context.Context, reviewDTO entities.ReviewDTO) (entities.Review, error)
}

type Service struct {
	bookingRepo booking.BookingRepo
	reviewRepo  review.ReviewRepo
}

func New(
	bookingRepo booking.BookingRepo,
	reviewRepo review.ReviewRepo,
) *Service {
	return &Service{
		bookingRepo: bookingRepo,
		reviewRepo:  reviewRepo,
	}
}

func (s *Service) SubmitReview(ctx context.Context, reviewDTO entities.ReviewDTO) (entities.Review, error) {
	booking, err := s.bookingRepo.FindById(ctx, reviewDTO.BookingID)
	if err != nil {
		return entities.Review{}, err
	}

	review := entities.Review{
		BookingID: booking.ID,
		GuestID:   reviewDTO.GuestID,
		Rating:    reviewDTO.Rating,
		Comment:   reviewDTO.Comment,
	}

	err = s.reviewRepo.Save(ctx, &review)
	if err != nil {
		return entities.Review{}, err
	}

	return review, nil
}
