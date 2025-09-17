package guest_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"booking-service/internal/entities"
	"booking-service/internal/repositories/guest/mocks"
	guestService "booking-service/internal/services/guest"
)

// TestFindGuestsByBookingID_Success проверяет успешное получение списка гостей по ID бронирования.
func TestFindGuestsByBookingID_Success(t *testing.T) {
	mockRepo := new(mocks.GuestRepoMock)
	service := guestService.NewService(mockRepo)

	ctx := context.Background()
	bookingID := uint(1)
	expectedGuests := []entities.Guest{
		{ID: 1, Name: "Алиса"},
		{ID: 2, Name: "Боб"},
	}

	mockRepo.
		On("FindByBookingID", ctx, bookingID).
		Return(expectedGuests, nil).
		Once()

	guests, err := service.FindGuestsByBookingID(ctx, bookingID)
	assert.NoError(t, err)
	assert.Equal(t, expectedGuests, guests)

	mockRepo.AssertExpectations(t)
}

// TestFindGuestsByBookingID_Empty проверяет случай, когда по ID бронирования нет гостей.
func TestFindGuestsByBookingID_Empty(t *testing.T) {
	mockRepo := new(mocks.GuestRepoMock)
	service := guestService.NewService(mockRepo)

	ctx := context.Background()
	bookingID := uint(2)

	mockRepo.
		On("FindByBookingID", ctx, bookingID).
		Return([]entities.Guest{}, nil).
		Once()

	guests, err := service.FindGuestsByBookingID(ctx, bookingID)
	assert.NoError(t, err)
	assert.Empty(t, guests)

	mockRepo.AssertExpectations(t)
}

// TestFindGuestsByBookingID_Error проверяет ситуацию, когда метод репозитория возвращает ошибку.
func TestFindGuestsByBookingID_Error(t *testing.T) {
	mockRepo := new(mocks.GuestRepoMock)
	service := guestService.NewService(mockRepo)

	ctx := context.Background()
	bookingID := uint(3)
	expErr := errors.New("ошибка запроса в базу данных")

	mockRepo.
		On("FindByBookingID", ctx, bookingID).
		Return([]entities.Guest{}, expErr).
		Once()

	guests, err := service.FindGuestsByBookingID(ctx, bookingID)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
	assert.Empty(t, guests)

	mockRepo.AssertExpectations(t)
}

// TestCreateGuest_Success проверяет успешное создание гостя.
func TestCreateGuest_Success(t *testing.T) {
	mockRepo := new(mocks.GuestRepoMock)
	service := guestService.NewService(mockRepo)

	ctx := context.Background()
	dto := entities.GuestDTO{Name: "Чарли"}
	expectedGuest := entities.Guest{ID: 10, Name: "Чарли"}

	mockRepo.
		On("SaveAndReturnIt", ctx, mock.MatchedBy(func(g *entities.Guest) bool {
			return g.Name == dto.Name
		})).
		Return(expectedGuest, nil).
		Once()

	guest, err := service.CreateGuest(ctx, dto)
	assert.NoError(t, err)
	assert.Equal(t, expectedGuest, guest)

	mockRepo.AssertExpectations(t)
}

// TestCreateGuest_Error проверяет случай, когда создание гостя завершается ошибкой.
func TestCreateGuest_Error(t *testing.T) {
	mockRepo := new(mocks.GuestRepoMock)
	service := guestService.NewService(mockRepo)

	ctx := context.Background()
	dto := entities.GuestDTO{Name: "Давид"}
	expErr := errors.New("ошибка сохранения")

	mockRepo.
		On("SaveAndReturnIt", ctx, mock.AnythingOfType("*entities.Guest")).
		Return(entities.Guest{}, expErr).
		Once()

	guest, err := service.CreateGuest(ctx, dto)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
	assert.Equal(t, entities.Guest{}, guest)

	mockRepo.AssertExpectations(t)
}
