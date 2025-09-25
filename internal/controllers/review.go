package controllers

import (
	"context"

	"booking-service/internal/entities"
)

func (c *Controller) SubmitReview(ctx context.Context, reviewDTO entities.ReviewDTO) (entities.Review, error) {
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
