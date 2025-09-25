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

func TestGetEmployee_Success(t *testing.T) {
	// Arrange
	employeeServiceMock := &mocks.EmployeeServiceMock{}
	h := &handlers.Handler{EmployeeService: employeeServiceMock}

	id := uint64(123)
	now := time.Now()
	expectedEmployee := entities.Employee{
		ID:        id,
		Name:      "John Doe",
		Role:      "Manager",
		CreatedAt: now,
		UpdatedAt: now,
	}

	employeeServiceMock.On("GetEmployee", mock.Anything, id).
		Return(expectedEmployee, nil)

	// Act
	resp, err := h.GetEmployee(context.Background(), &generated.GetEmployeeRequest{
		EmployeeId: id,
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Employee)
	assert.Equal(t, id, resp.Employee.Id)
	assert.Equal(t, "John Doe", resp.Employee.Name)
	assert.Equal(t, "Manager", resp.Employee.Role)
	assert.True(t, timestamppb.New(now).AsTime().Equal(resp.Employee.CreatedAt.AsTime()))
	assert.True(t, timestamppb.New(now).AsTime().Equal(resp.Employee.UpdatedAt.AsTime()))
	assert.Len(t, resp.Links, 1)
	assert.Equal(t, "self", resp.Links[0].Rel)
	assert.Equal(t, "/v1/employee/{employee_id}", resp.Links[0].Href)

	employeeServiceMock.AssertExpectations(t)
}

func TestGetEmployee_NotFound(t *testing.T) {
	// Arrange
	employeeServiceMock := &mocks.EmployeeServiceMock{}
	h := &handlers.Handler{EmployeeService: employeeServiceMock}
	id := uint64(404)

	employeeServiceMock.On("GetEmployee", mock.Anything, id).
		Return(entities.Employee{}, entities.ErrNotFound)

	// Act
	resp, err := h.GetEmployee(context.Background(), &generated.GetEmployeeRequest{
		EmployeeId: id,
	})

	// Assert
	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, "Employee not found", st.Message())

	employeeServiceMock.AssertExpectations(t)
}

func TestGetEmployee_InternalError(t *testing.T) {
	// Arrange
	employeeServiceMock := &mocks.EmployeeServiceMock{}
	h := &handlers.Handler{EmployeeService: employeeServiceMock}
	id := uint64(500)

	employeeServiceMock.On("GetEmployee", mock.Anything, id).
		Return(entities.Employee{}, errors.New("database error"))

	// Act
	resp, err := h.GetEmployee(context.Background(), &generated.GetEmployeeRequest{
		EmployeeId: id,
	})

	// Assert
	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Equal(t, "Failed to get employee", st.Message())

	employeeServiceMock.AssertExpectations(t)
}
