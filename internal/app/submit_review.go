package app

import (
	"context"

	"booking-service/internal/entities"
	"booking-service/internal/generated"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *Handler) SubmitReview(ctx context.Context, in *generated.SubmitReviewRequest) (
	*generated.SubmitReviewResponse, error,
) {
	review, err := h.bookingController.SubmitReview(ctx, entities.ReviewDTO{
		BookingID: in.BookingId,
		GuestID:   in.GuestId,
		Rating:    int(in.Rating),
		Comment:   in.Comment,
	})
	if err != nil {
		return nil, err
	}

	return &generated.SubmitReviewResponse{
		Review: &generated.Review{
			Id:        review.ID,
			CreatedAt: timestamppb.New(review.CreatedAt),
			UpdatedAt: timestamppb.New(review.UpdatedAt),
			BookingId: review.BookingID,
			Rating:    int32(review.Rating),
			Comment:   review.Comment,
		},
	}, nil
}
