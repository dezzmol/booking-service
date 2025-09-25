package house_keeping_request_test

import (
	"context"
	"errors"
	"testing"

	houseKeepingService "booking-service/internal/controllers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"booking-service/internal/entities"
	"booking-service/internal/repositories/house_keeping_request/mocks"
)

// TestCreateHouseKeepingRequest_Success проверяет успешное создание запроса на уборку.
func TestCreateHouseKeepingRequest_Success(t *testing.T) {
	mockRepo := new(mocks.HouseKeepingRequestRepoMock)
	service := houseKeepingService.New(mockRepo)

	ctx := context.Background()

	// Создаем DTO для запроса уборки.
	dto := entities.HouseKeepingRequestDTO{
		RoomID: 101,
	}

	// Ожидаем вызов метода Save с любым указателем на entities.HousekeepingRequest и возврат nil.
	mockRepo.
		On("Save", ctx, mock.AnythingOfType("*entities.HousekeepingRequest")).
		Return(nil).Once()

	result, err := service.CreateHouseKeepingRequest(ctx, dto)
	assert.NoError(t, err)
	// Проверяем, что статус запроса соответствует REQUESTED
	assert.Equal(t, "requested", result.Status)
	assert.Equal(t, dto.RoomID, result.RoomID)

	mockRepo.AssertExpectations(t)
}

// TestCreateHouseKeepingRequest_SaveError проверяет обработку ошибки при сохранении запроса.
func TestCreateHouseKeepingRequest_SaveError(t *testing.T) {
	mockRepo := new(mocks.HouseKeepingRequestRepoMock)
	service := houseKeepingService.New(mockRepo)

	ctx := context.Background()

	dto := entities.HouseKeepingRequestDTO{
		RoomID: 101,
	}

	expErr := errors.New("ошибка сохранения запроса")
	mockRepo.
		On("Save", ctx, mock.AnythingOfType("*entities.HousekeepingRequest")).
		Return(expErr).Once()

	result, err := service.CreateHouseKeepingRequest(ctx, dto)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
	// Результат должен быть пустым
	assert.Equal(t, entities.HousekeepingRequest{}, result)

	mockRepo.AssertExpectations(t)
}

// TestAssignEmployee_Success проверяет успешное назначение сотрудника для запроса на уборку.
func TestAssignEmployee_Success(t *testing.T) {
	mockRepo := new(mocks.HouseKeepingRequestRepoMock)
	service := houseKeepingService.New(mockRepo)

	ctx := context.Background()
	requestID := uint64(1)
	employeeID := uint64(1001)

	// Создаем существующий запрос с первоначальным статусом "0" (REQUESTED)
	existingRequest := &entities.HousekeepingRequest{
		RoomID: 101,
		Status: "0",
		// Можно добавить и другие поля, если требуется
	}

	// Ожидаем, что метод FindByID вернет существующий запрос.
	mockRepo.
		On("FindByID", ctx, requestID).
		Return(existingRequest, nil).Once()

	// Ожидаем вызов метода Update, где статус должен измениться на "1" (ACCEPTED)
	mockRepo.
		On("Update", ctx, mock.MatchedBy(func(r entities.HousekeepingRequest) bool {
			return r.Status == houseKeepingService.ACCEPTED
		})).
		Return(nil).Once()

	result, err := service.AssignEmployee(ctx, requestID, employeeID)
	assert.NoError(t, err)
	assert.Equal(t, houseKeepingService.ACCEPTED, result.Status)

	mockRepo.AssertExpectations(t)
}

// TestAssignEmployee_FindError проверяет ситуацию, когда запрос не найден.
func TestAssignEmployee_FindError(t *testing.T) {
	mockRepo := new(mocks.HouseKeepingRequestRepoMock)
	service := houseKeepingService.New(mockRepo)

	ctx := context.Background()
	requestID := uint64(1)
	employeeID := uint64(1001)

	expErr := errors.New("запрос не найден")
	mockRepo.
		On("FindByID", ctx, requestID).
		Return((*entities.HousekeepingRequest)(nil), expErr).Once()

	result, err := service.AssignEmployee(ctx, requestID, employeeID)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
	// Результат должен быть пустым
	assert.Equal(t, entities.HousekeepingRequest{}, result)

	mockRepo.AssertExpectations(t)
}

// TestAssignEmployee_UpdateError проверяет обработку ошибки при обновлении запроса.
func TestAssignEmployee_UpdateError(t *testing.T) {
	mockRepo := new(mocks.HouseKeepingRequestRepoMock)
	service := houseKeepingService.New(mockRepo)

	ctx := context.Background()
	requestID := uint64(1)
	employeeID := uint64(1001)

	existingRequest := &entities.HousekeepingRequest{
		RoomID: 101,
		Status: "0", // REQUESTED
	}

	expErr := errors.New("ошибка обновления запроса")

	mockRepo.
		On("FindByID", ctx, requestID).
		Return(existingRequest, nil).Once()

	mockRepo.
		On("Update", ctx, mock.AnythingOfType("entities.HousekeepingRequest")).
		Return(expErr).Once()

	result, err := service.AssignEmployee(ctx, requestID, employeeID)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
	assert.Equal(t, entities.HousekeepingRequest{}, result)

	mockRepo.AssertExpectations(t)
}
