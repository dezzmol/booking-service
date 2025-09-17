package entities

import "time"

type Employee struct {
	ID        uint64    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Name      string    `db:"name"`
	Role      string    `db:"role"`
}

type EmployeeDTO struct {
	Name string
	Role string
}

func (e *Employee) GetID() uint64 {
	if e == nil {
		return 0
	}
	return e.ID
}
