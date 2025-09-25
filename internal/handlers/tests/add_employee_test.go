package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"booking-service/internal/controllers/employee/mocks"
	"booking-service/internal/entities"
	"booking-service/internal/generated"
	"booking-service/internal/handlers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAddEmployee_Success(t *testing.T) {
	// Создаем мок сервиса
	employeeServiceMock := &mocks.EmployeeServiceMock{}

	// Подготавливаем тестовые данные
	now := time.Now()
	expectedEmployee := entities.Employee{
		ID:        uint64(123),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "John Doe",
		Role:      "Developer",
	}

	// Настраиваем ожидания
	employeeServiceMock.On("AddEmployee", mock.Anything, entities.EmployeeDTO{
		Name: "John Doe",
		Role: "Developer",
	}).Return(expectedEmployee, nil)

	// Создаем хендлер с моком
	h := &handlers.Handler{
		EmployeeService: employeeServiceMock,
	}

	// Вызываем тестируемую функцию
	resp, err := h.AddEmployee(context.Background(), &generated.AddEmployeeRequest{
		Name: "John Doe",
		Role: "Developer",
	})

	// Проверяем результаты
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Employee)
	assert.Equal(t, uint64(123), resp.Employee.Id)
	assert.Equal(t, "John Doe", resp.Employee.Name)
	assert.Equal(t, "Developer", resp.Employee.Role)
	assert.True(t, timestamppb.New(now).AsTime().Equal(resp.Employee.CreatedAt.AsTime()))
	assert.True(t, timestamppb.New(now).AsTime().Equal(resp.Employee.UpdatedAt.AsTime()))

	// Проверяем, что все ожидания по моку выполнены
	employeeServiceMock.AssertExpectations(t)
}

func TestAddEmployee_InternalError(t *testing.T) {
	// Создаем мок сервиса
	employeeServiceMock := &mocks.EmployeeServiceMock{}

	// Настраиваем ожидания с возвратом ошибки
	expectedErr := errors.New("database error")
	employeeServiceMock.On("AddEmployee", mock.Anything, entities.EmployeeDTO{
		Name: "John Doe",
		Role: "Developer",
	}).Return(entities.Employee{}, expectedErr)

	// Создаем хендлер с моком
	h := &handlers.Handler{
		EmployeeService: employeeServiceMock,
	}

	// Вызываем тестируемую функцию
	resp, err := h.AddEmployee(context.Background(), &generated.AddEmployeeRequest{
		Name: "John Doe",
		Role: "Developer",
	})

	// Проверяем результаты
	assert.Nil(t, resp)
	assert.Error(t, err)

	// Проверяем, что это именно gRPC ошибка с кодом Internal
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "internal error")

	// Проверяем, что все ожидания по моку выполнены
	employeeServiceMock.AssertExpectations(t)
}

func TestAddEmployee_EmptyName(t *testing.T) {
	// Создаем мок сервиса (не должен вызываться в этом случае)
	employeeServiceMock := &mocks.EmployeeServiceMock{}

	// Создаем хендлер с моком
	h := &handlers.Handler{
		EmployeeService: employeeServiceMock,
	}

	employeeServiceMock.On("AddEmployee", mock.Anything, entities.EmployeeDTO{
		Name: "",
		Role: "Developer",
	}).Return(entities.Employee{}, entities.ErrInvalidName)

	// Вызываем тестируемую функцию с пустым именем
	resp, err := h.AddEmployee(context.Background(), &generated.AddEmployeeRequest{
		Name: "",
		Role: "Developer",
	})

	// Проверяем результаты
	assert.Nil(t, resp)
	assert.Error(t, err)

	// Проверяем, что это именно gRPC ошибка с кодом InvalidArgument
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}
