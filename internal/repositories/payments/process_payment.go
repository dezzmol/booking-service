package payments

import (
	"context"
	"errors"

	"booking-service/internal/generated"
)

func (w *Wrapper) ProcessPayment(ctx context.Context, bookingID uint64, amount float32) error {
	resp, err := w.client.ProcessPayment(ctx, &generated.ProcessRequest{
		BookingId: bookingID,
		Amount:    amount,
	})
	if err != nil {
		return err
	}
	if !resp.Status {
		return errors.New(resp.GetError())
	}

	return nil
}
