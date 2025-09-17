package employee

import (
	"context"

	"booking-service/internal/entities"

	"github.com/jmoiron/sqlx"
)

//go:generate mockery --disable-version-string --case=underscore --name=EmployeeRepository --structname=EmployeeRepositoryMock
type EmployeeRepository interface {
	Save(ctx context.Context, employee *entities.Employee) error
	FindByID(ctx context.Context, employee uint64) (entities.Employee, error)
	Update(ctx context.Context, employee *entities.Employee) error
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, employee *entities.Employee) error {
	query := `INSERT INTO employees (name, role) 
			VALUES ($1, $2)
			RETURNING id, created_at, updated_at`

	return r.db.QueryRowxContext(ctx, query, employee.Name, employee.Role).Scan(&employee.ID, &employee.CreatedAt, &employee.UpdatedAt)
}

func (r *Repository) FindByID(ctx context.Context, id uint64) (entities.Employee, error) {
	query := `SELECT id, name, role, created_at, updated_at FROM employees WHERE id = $1`
	var employee entities.Employee
	err := r.db.QueryRowxContext(ctx, query, id).
		Scan(&employee.ID, &employee.Name, &employee.Role, &employee.CreatedAt, &employee.UpdatedAt)
	if err != nil {
		return entities.Employee{}, err
	}
	return employee, nil
}

func (r *Repository) Update(ctx context.Context, employee *entities.Employee) error {
	query := `UPDATE employees SET name = $1, role = $2 WHERE id = $4 RETURNING id, created_at, updated_at`

	return r.db.QueryRowxContext(ctx, query, employee.Name, employee.Role, employee.ID).Scan(&employee.ID, &employee.CreatedAt, &employee.UpdatedAt)
}
