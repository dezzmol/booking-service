package entities

import "time"

type Payment struct {
	ID          uint64    `db:"id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	BookingID   uint64    `db:"booking_id"`
	Amount      float64   `db:"amount"`
	PaymentDate time.Time `db:"payment_date"`
	Status      string    `db:"status"`
}
