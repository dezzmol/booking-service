package payments

import (
	"context"
	"errors"

	"booking-service/internal/generated"
)

func (w *Wrapper) CancelPayment(ctx context.Context, bookingID uint64) error {
	resp, err := w.client.CancelPayment(ctx, &generated.BookingInfo{
		BookingId: bookingID,
	})
	if err != nil {
		return err
	}
	if !resp.Status {
		return errors.New(resp.GetError())
	}
	return nil
}
