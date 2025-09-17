package room_test

import (
	"context"
	"errors"
	"testing"

	"booking-service/internal/repositories/notifications"
	notificationMocks "booking-service/internal/repositories/notifications/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"booking-service/internal/entities"
	"booking-service/internal/repositories/room/mocks"
	roomService "booking-service/internal/services/room"
)

// TestCreateRooms_Success проверяет успешное создание комнат.
// При условии, что количество комнат меньше ChunkSize, ожидается один вызов SaveAll.
func TestCreateRooms_Success(t *testing.T) {
	mockRepo := new(mocks.RoomRepoMock)
	notificationClient := new(notificationMocks.ClientMock)
	wrapper := notifications.New(
		notificationClient,
		"",
		"",
		nil,
		"",
	)
	service := roomService.New(mockRepo, wrapper)

	ctx := context.Background()
	// Создаем DTO для комнат.
	roomDTOs := []entities.RoomDTO{
		{Number: "101", Type: "Single", HotelID: 1},
		{Number: "102", Type: "Double", HotelID: 1},
	}

	// Ожидаемый срез после преобразования DTO в entities.Room.
	expectedRooms := []entities.Room{
		{Number: "101", Type: "Single", HotelID: 1},
		{Number: "102", Type: "Double", HotelID: 1},
	}

	// Ожидаем один вызов SaveAll с сформированным срезом комнат.
	mockRepo.
		On("SaveAll", ctx, expectedRooms).
		Return(nil).
		Once()

	err := service.CreateRooms(ctx, roomDTOs)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestCreateRooms_SaveAllError проверяет, что если SaveAll возвращает ошибку, то CreateRooms также возвращает ошибку.
func TestCreateRooms_SaveAllError(t *testing.T) {
	mockRepo := new(mocks.RoomRepoMock)
	notificationClient := new(notificationMocks.ClientMock)
	wrapper := notifications.New(
		notificationClient,
		"",
		"",
		nil,
		"",
	)
	service := roomService.New(mockRepo, wrapper)

	ctx := context.Background()
	roomDTOs := []entities.RoomDTO{
		{Number: "101", Type: "Single", HotelID: 1},
	}

	expectedRooms := []entities.Room{
		{Number: "101", Type: "Single", HotelID: 1},
	}

	expErr := errors.New("ошибка сохранения комнат")
	mockRepo.
		On("SaveAll", ctx, expectedRooms).
		Return(expErr).
		Once()

	err := service.CreateRooms(ctx, roomDTOs)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)

	mockRepo.AssertExpectations(t)
}

// TestUpdate_Success проверяет успешное обновление информации о комнате.
// Проверяем, что метод Update репозитория получает корректно сформированный объект.
func TestUpdate_Success(t *testing.T) {
	mockRepo := new(mocks.RoomRepoMock)
	notificationClient := new(notificationMocks.ClientMock)
	wrapper := notifications.New(
		notificationClient,
		"",
		"",
		nil,
		"",
	)
	service := roomService.New(mockRepo, wrapper)

	ctx := context.Background()
	dto := entities.RoomDTO{
		Number:  "101",
		Type:    "Suite",
		HotelID: 1,
	}

	mockRepo.
		On("Update", ctx, mock.MatchedBy(func(r *entities.Room) bool {
			return r.Number == dto.Number && r.Type == dto.Type && r.HotelID == dto.HotelID
		})).
		Return(nil).
		Once()
	notificationClient.
		On("SendNotificationWithQueue", mock.Anything, roomService.NotificationMessage, 1).
		Return(nil).
		Once()

	err := service.Update(ctx, dto)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestUpdate_Error проверяет сценарий, когда метод Update репозитория возвращает ошибку.
func TestUpdate_Error(t *testing.T) {
	mockRepo := new(mocks.RoomRepoMock)
	notificationClient := new(notificationMocks.ClientMock)
	wrapper := notifications.New(
		notificationClient,
		"",
		"",
		nil,
		"",
	)
	service := roomService.New(mockRepo, wrapper)

	ctx := context.Background()
	dto := entities.RoomDTO{
		Number:  "101",
		Type:    "Suite",
		HotelID: 1,
	}

	expErr := errors.New("ошибка обновления комнаты")
	mockRepo.
		On("Update", mock.Anything, &entities.Room{
			Number:  dto.Number,
			Type:    dto.Type,
			HotelID: dto.HotelID,
		}).
		Return(expErr).
		Once()

	err := service.Update(ctx, dto)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)

	mockRepo.AssertExpectations(t)
}
