package tests

import (
	"context"
	"errors"
	"testing"

	"booking-service/internal/entities"
	"booking-service/internal/generated"
	"booking-service/internal/handlers"
	"booking-service/internal/services/booking/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCancelBooking_Success(t *testing.T) {
	bookingServiceMock := &mocks.BookingServiceMock{}

	// Настраиваем ожидания
	bookingServiceMock.On("CancelBooking", mock.Anything, uint64(123)).Return(nil)

	// Создаем хендлер с моком
	h := &handlers.Handler{
		BookingService: bookingServiceMock,
	}

	// Вызываем тестируемую функцию
	resp, err := h.CancelBooking(context.Background(), &generated.CancelBookingRequest{
		BookingId: 123,
	})

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Links, 1)
	assert.Equal(t, "self", resp.Links[0].Rel)

	bookingServiceMock.AssertExpectations(t)
}

func TestCancelBooking_NotFoundError(t *testing.T) {
	bookingServiceMock := &mocks.BookingServiceMock{}

	expectedErr := entities.ErrNotFound
	bookingServiceMock.On("CancelBooking", mock.Anything, uint64(123)).Return(expectedErr)

	h := &handlers.Handler{
		BookingService: bookingServiceMock,
	}

	resp, err := h.CancelBooking(context.Background(), &generated.CancelBookingRequest{
		BookingId: 123,
	})

	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "booking not found")

	bookingServiceMock.AssertExpectations(t)
}

func TestCancelBooking_InternalError(t *testing.T) {
	// Создаем мок сервиса
	bookingServiceMock := &mocks.BookingServiceMock{}

	// Настраиваем ожидания с возвратом неожиданной ошибки
	unexpectedErr := errors.New("some unexpected error")
	bookingServiceMock.On("CancelBooking", mock.Anything, uint64(123)).Return(unexpectedErr)

	// Создаем хендлер с моком
	h := &handlers.Handler{
		BookingService: bookingServiceMock,
	}

	// Вызываем тестируемую функцию
	resp, err := h.CancelBooking(context.Background(), &generated.CancelBookingRequest{
		BookingId: 123,
	})

	// Проверяем результаты
	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "internal error")

	// Проверяем, что все ожидания по моку выполнены
	bookingServiceMock.AssertExpectations(t)
}
