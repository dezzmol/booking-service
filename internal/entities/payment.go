package entities

import "time"

type PaymentStatus int8

var (
	PaymentStatusUnknown  PaymentStatus = 0
	PaymentStatusSuccess  PaymentStatus = 1
	PaymentStatusFailed   PaymentStatus = 2
	PaymentStatusCanceled PaymentStatus = 3
)

type Payment struct {
	ID          uint64        `db:"id"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
	BookingID   uint64        `db:"booking_id"`
	Amount      float64       `db:"amount"`
	PaymentDate time.Time     `db:"payment_date"`
	Status      PaymentStatus `db:"status"`
}
