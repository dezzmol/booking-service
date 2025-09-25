package employee_test

import (
	"context"
	"errors"
	"testing"

	employeeService "booking-service/internal/controllers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"booking-service/internal/entities"
	"booking-service/internal/repositories/employee/mocks"
)

// TestGetEmployee_Success проверяет успешное получение сотрудника.
func TestGetEmployee_Success(t *testing.T) {
	mockRepo := new(mocks.EmployeeRepositoryMock)
	service := employeeService.NewService(mockRepo)

	ctx := context.Background()
	expectedEmployee := entities.Employee{
		ID:   1,
		Name: "Алиса",
		Role: "Разработчик",
	}

	mockRepo.
		On("FindByID", mock.Anything, uint64(1)).
		Return(expectedEmployee, nil).
		Once()

	employee, err := service.GetEmployee(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, expectedEmployee, employee)

	mockRepo.AssertExpectations(t)
}

// TestGetEmployee_Error проверяет сценарий, когда сотрудник не найден.
func TestGetEmployee_Error(t *testing.T) {
	mockRepo := new(mocks.EmployeeRepositoryMock)
	service := employeeService.NewService(mockRepo)

	ctx := context.Background()
	expErr := errors.New("сотрудник не найден")

	mockRepo.
		On("FindByID", mock.Anything, uint64(2)).
		Return(entities.Employee{}, expErr).
		Once()

	employee, err := service.GetEmployee(ctx, 2)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
	assert.Equal(t, entities.Employee{}, employee)

	mockRepo.AssertExpectations(t)
}

// TestAddEmployee_Success проверяет успешное добавление сотрудника.
func TestAddEmployee_Success(t *testing.T) {
	mockRepo := new(mocks.EmployeeRepositoryMock)
	service := employeeService.NewService(mockRepo)

	ctx := context.Background()
	dto := entities.EmployeeDTO{
		Name: "Боб",
		Role: "Тестировщик",
	}

	// Ожидаем, что метод Save будет вызван с объектом, поля которого соответствуют значениями DTO.
	mockRepo.
		On("Save", ctx, mock.MatchedBy(func(emp *entities.Employee) bool {
			return emp.Name == dto.Name && emp.Role == dto.Role
		})).
		Return(nil).
		Once()

	employee, err := service.AddEmployee(ctx, dto)
	assert.NoError(t, err)
	assert.Equal(t, dto.Name, employee.Name)
	assert.Equal(t, dto.Role, employee.Role)

	mockRepo.AssertExpectations(t)
}

// TestAddEmployee_SaveError проверяет сценарий, когда Save возвращает ошибку.
func TestAddEmployee_SaveError(t *testing.T) {
	mockRepo := new(mocks.EmployeeRepositoryMock)
	service := employeeService.NewService(mockRepo)

	ctx := context.Background()
	dto := entities.EmployeeDTO{
		Name: "Чарли",
		Role: "Поддержка",
	}

	expErr := errors.New("ошибка сохранения")
	mockRepo.
		On("Save", ctx, mock.AnythingOfType("*entities.Employee")).
		Return(expErr).
		Once()

	employee, err := service.AddEmployee(ctx, dto)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
	assert.Equal(t, entities.Employee{}, employee)

	mockRepo.AssertExpectations(t)
}

// TestUpdateEmployee_Success проверяет успешное обновление данных сотрудника.
func TestUpdateEmployee_Success(t *testing.T) {
	mockRepo := new(mocks.EmployeeRepositoryMock)
	service := employeeService.NewService(mockRepo)

	ctx := context.Background()
	employeeID := uint64(1)
	// Исходный сотрудник с устаревшими данными.
	existingEmployee := entities.Employee{
		ID:   employeeID,
		Name: "Дэвид",
		Role: "Старый Роль",
	}
	mockRepo.
		On("FindByID", ctx, employeeID).
		Return(existingEmployee, nil).
		Once()

	dto := entities.EmployeeDTO{
		Name: "Дэвид Обновлённый",
		Role: "Новая Роль",
	}

	// Ожидаем, что Update будет вызван с объектом, в котором обновлены Name и Role.
	mockRepo.
		On("Update", ctx, mock.MatchedBy(func(emp *entities.Employee) bool {
			return emp.ID == employeeID && emp.Name == dto.Name && emp.Role == dto.Role
		})).
		Return(nil).
		Once()

	_, err := service.UpdateEmployee(ctx, employeeID, dto)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestUpdateEmployee_FindError проверяет ситуацию, когда сотрудник не найден.
func TestUpdateEmployee_FindError(t *testing.T) {
	mockRepo := new(mocks.EmployeeRepositoryMock)
	service := employeeService.NewService(mockRepo)

	ctx := context.Background()
	employeeID := uint64(2)
	expErr := errors.New("entity not found")

	mockRepo.
		On("FindByID", ctx, employeeID).
		Return(entities.Employee{}, expErr).
		Once()

	dto := entities.EmployeeDTO{
		Name: "Любой",
		Role: "Любой",
	}

	_, err := service.UpdateEmployee(ctx, employeeID, dto)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)

	mockRepo.AssertExpectations(t)
}

// TestUpdateEmployee_UpdateError проверяет сценарий, когда метод Update возвращает ошибку.
func TestUpdateEmployee_UpdateError(t *testing.T) {
	mockRepo := new(mocks.EmployeeRepositoryMock)
	service := employeeService.NewService(mockRepo)

	ctx := context.Background()
	employeeID := uint64(3)
	existingEmployee := entities.Employee{
		ID:   employeeID,
		Name: "Эдвард",
		Role: "Старый Роль",
	}
	mockRepo.
		On("FindByID", ctx, employeeID).
		Return(existingEmployee, nil).
		Once()

	dto := entities.EmployeeDTO{
		Name: "Эдвард Обновлённый",
		Role: "Новая Роль",
	}
	expErr := errors.New("ошибка обновления")
	mockRepo.
		On("Update", ctx, mock.AnythingOfType("*entities.Employee")).
		Return(expErr).
		Once()

	_, err := service.UpdateEmployee(ctx, employeeID, dto)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)

	mockRepo.AssertExpectations(t)
}
