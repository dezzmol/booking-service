package employee

import (
	"context"

	"booking-service/internal/entities"
	"booking-service/internal/repositories/employee"
)

//go:generate mockery --disable-version-string --case=underscore --name=EmployeeService --structname=EmployeeServiceMock
type EmployeeService interface {
	GetEmployee(ctx context.Context, employeeID uint64) (entities.Employee, error)
	AddEmployee(ctx context.Context, employee entities.EmployeeDTO) (entities.Employee, error)
	UpdateEmployee(ctx context.Context, employeeID uint64, employee entities.EmployeeDTO) (entities.Employee, error)
}

type Service struct {
	repo employee.EmployeeRepository
}

func NewService(repo employee.EmployeeRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetEmployee(ctx context.Context, employeeID uint64) (entities.Employee, error) {
	return s.repo.FindByID(ctx, employeeID)
}

func (s *Service) AddEmployee(ctx context.Context, employee entities.EmployeeDTO) (entities.Employee, error) {
	baseEmployee := entities.Employee{
		Name: employee.Name,
		Role: employee.Role,
	}
	if baseEmployee.Name == "" || baseEmployee.Role == "" {
		return entities.Employee{}, entities.ErrInvalidName
	}

	err := s.repo.Save(ctx, &baseEmployee)
	if err != nil {
		return entities.Employee{}, err
	}
	return baseEmployee, nil
}

func (s *Service) UpdateEmployee(ctx context.Context, employeeID uint64, employee entities.EmployeeDTO) (entities.Employee, error) {
	baseEmployee, err := s.repo.FindByID(ctx, employeeID)
	if err != nil {
		return baseEmployee, entities.ErrNotFound
	}

	baseEmployee.Name = employee.Name
	baseEmployee.Role = employee.Role

	err = s.repo.Update(ctx, &baseEmployee)
	if err != nil {
		return baseEmployee, err
	}
	return baseEmployee, nil
}
