package tests

import (
	"context"
	"testing"
	"time"

	"booking-service/internal/entities"
	"booking-service/internal/generated"
	"booking-service/internal/handlers"
	"booking-service/internal/services/house_keeping_request/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestHandler_AssignEmployee(t *testing.T) {
	// Фиксируем время для тестов
	now := time.Now()
	requestTime := now.Add(1 * time.Hour)

	tests := []struct {
		name          string
		input         *generated.AssignEmployeeRequest
		mockSetup     func(serviceMock *mocks.HouseKeepingRequestServiceMock)
		expected      *generated.AssignEmployeeResponse
		expectedError bool
	}{
		{
			name: "successful employee assignment",
			input: &generated.AssignEmployeeRequest{
				RequestId:  uint64(123),
				EmployeeId: uint64(456),
			},
			mockSetup: func(mockService *mocks.HouseKeepingRequestServiceMock) {
				mockService.On("AssignEmployee", mock.Anything, uint64(123), uint64(456)).
					Return(entities.HousekeepingRequest{
						ID:          uint64(123),
						CreatedAt:   now,
						UpdatedAt:   now,
						RoomID:      uint64(789),
						RequestTime: requestTime,
						Status:      "ASSIGNED",
					}, nil)
			},
			expected: &generated.AssignEmployeeResponse{
				Request: &generated.HousekeepingRequest{
					Id:          uint64(123),
					CreatedAt:   timestamppb.New(now),
					UpdatedAt:   timestamppb.New(now),
					RoomId:      uint64(789),
					RequestTime: timestamppb.New(requestTime),
					Status:      "ASSIGNED",
				},
				Links: nil,
			},
			expectedError: false,
		},
		{
			name: "service returns error",
			input: &generated.AssignEmployeeRequest{
				RequestId:  uint64(123),
				EmployeeId: uint64(456),
			},
			mockSetup: func(mockService *mocks.HouseKeepingRequestServiceMock) {
				mockService.On("AssignEmployee", mock.Anything, uint64(123), uint64(456)).
					Return(entities.HousekeepingRequest{}, assert.AnError)
			},
			expected:      nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок сервиса
			mockService := new(mocks.HouseKeepingRequestServiceMock)
			tt.mockSetup(mockService)

			// Создаем хендлер с моком сервиса
			h := &handlers.Handler{
				HousekeepingService: mockService,
			}

			// Вызываем тестируемую функцию
			res, err := h.AssignEmployee(context.Background(), tt.input)

			// Проверяем ошибки
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, res)
			}

			// Проверяем, что все ожидания по моку выполнены
			mockService.AssertExpectations(t)
		})
	}
}
