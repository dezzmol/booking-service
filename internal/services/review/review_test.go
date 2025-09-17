package review_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"booking-service/internal/entities"
	bookingMocks "booking-service/internal/repositories/booking/mocks"
	"booking-service/internal/repositories/review/mocks"
	reviewService "booking-service/internal/services/review"
)

func TestSubmitReview_Success(t *testing.T) {
	mockBookingRepo := new(bookingMocks.BookingRepoMock)
	mockReviewRepo := new(mocks.ReviewRepoMock)
	service := reviewService.New(mockBookingRepo, mockReviewRepo)

	ctx := context.Background()
	reviewDTO := entities.ReviewDTO{
		BookingID: 1,
		GuestID:   2,
		Rating:    5,
		Comment:   "Отличный сервис!",
	}

	// Мок репозитория бронирования возвращает бронирование с нужным ID.
	booking := &entities.Booking{ID: reviewDTO.BookingID}
	mockBookingRepo.
		On("FindById", ctx, reviewDTO.BookingID).
		Return(booking, nil).
		Once()

	// Мок репозитория отзывов ожидает вызов Save с отзывом, содержащим данные из reviewDTO.
	mockReviewRepo.
		On("Save", ctx, mock.MatchedBy(func(r *entities.Review) bool {
			return r.BookingID == booking.ID &&
				r.GuestID == reviewDTO.GuestID &&
				r.Rating == reviewDTO.Rating &&
				r.Comment == reviewDTO.Comment
		})).
		Return(nil).
		Once()

	_, err := service.SubmitReview(ctx, reviewDTO)
	assert.NoError(t, err)

	mockBookingRepo.AssertExpectations(t)
	mockReviewRepo.AssertExpectations(t)
}

// TestSubmitReview_BookingNotFound проверяет случай, когда бронирование не найдено.
func TestSubmitReview_BookingNotFound(t *testing.T) {
	mockBookingRepo := new(bookingMocks.BookingRepoMock)
	mockReviewRepo := new(mocks.ReviewRepoMock)
	service := reviewService.New(mockBookingRepo, mockReviewRepo)

	ctx := context.Background()
	reviewDTO := entities.ReviewDTO{
		BookingID: 1,
		GuestID:   2,
		Rating:    3,
		Comment:   "Было неплохо",
	}

	mockBookingRepo.
		On("FindById", ctx, reviewDTO.BookingID).
		Return(&entities.Booking{}, errors.New("бронирование не найдено")).
		Once()

	_, err := service.SubmitReview(ctx, reviewDTO)
	assert.Error(t, err)
	assert.Equal(t, "бронирование не найдено", err.Error())

	mockBookingRepo.AssertExpectations(t)
	mockReviewRepo.AssertNotCalled(t, "Save")
}

// TestSubmitReview_SaveError проверяет случай, когда при сохранении отзыва происходит ошибка.
func TestSubmitReview_SaveError(t *testing.T) {
	mockBookingRepo := new(bookingMocks.BookingRepoMock)
	mockReviewRepo := new(mocks.ReviewRepoMock)
	service := reviewService.New(mockBookingRepo, mockReviewRepo)

	ctx := context.Background()
	reviewDTO := entities.ReviewDTO{
		BookingID: 1,
		GuestID:   2,
		Rating:    3,
		Comment:   "Можно лучше",
	}

	booking := &entities.Booking{ID: reviewDTO.BookingID}
	mockBookingRepo.
		On("FindById", ctx, reviewDTO.BookingID).
		Return(booking, nil).
		Once()

	expErr := errors.New("save error")
	mockReviewRepo.
		On("Save", ctx, mock.AnythingOfType("*entities.Review")).
		Return(expErr).
		Once()

	_, err := service.SubmitReview(ctx, reviewDTO)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)

	mockBookingRepo.AssertExpectations(t)
	mockReviewRepo.AssertExpectations(t)
}
