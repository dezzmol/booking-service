package payments

import (
	"context"

	"booking-service/internal/generated"

	"google.golang.org/grpc"
)

//go:generate mockery --disable-version-string --case=underscore --name=client --structname=ClientMock
type (
	client interface {
		ProcessPayment(ctx context.Context, in *generated.ProcessRequest, opts ...grpc.CallOption) (*generated.ProcessResponse, error)
		CancelPayment(ctx context.Context, in *generated.BookingInfo, opts ...grpc.CallOption) (*generated.ProcessResponse, error)
		GetPaymentsInfo(ctx context.Context, in *generated.BookingInfo, opts ...grpc.CallOption) (*generated.PaymentsResponse, error)
	}

	Wrapper struct {
		client client
	}
)

func New(client client) *Wrapper {
	return &Wrapper{
		client: client,
	}
}
