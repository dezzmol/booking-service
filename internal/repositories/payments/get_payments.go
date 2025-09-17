package payments

import (
	"context"

	"booking-service/internal/entities"
	"booking-service/internal/generated"
)

func (w *Wrapper) GetPayments(ctx context.Context, bookingID uint64) ([]entities.Payment, error) {
	resp, err := w.client.GetPaymentsInfo(ctx, &generated.BookingInfo{
		BookingId: bookingID,
	})
	if err != nil {
		return nil, err
	}

	return w.makeGetPaymentsResponse(resp.Payments), nil
}

func (w *Wrapper) makeGetPaymentsResponse(genPayments []*generated.Payment) []entities.Payment {
	payments := make([]entities.Payment, 0, len(genPayments))
	for _, genPayment := range genPayments {
		payments = append(payments, entities.Payment{
			ID:          genPayment.Id,
			CreatedAt:   genPayment.CreatedAt.AsTime(),
			UpdatedAt:   genPayment.UpdatedAt.AsTime(),
			BookingID:   genPayment.BookingId,
			Amount:      genPayment.Amount,
			PaymentDate: genPayment.PaymentDate.AsTime(),
			Status:      genPayment.Status,
		})
	}
	return payments
}
