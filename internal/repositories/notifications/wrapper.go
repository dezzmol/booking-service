package notifications

import (
	"context"

	"booking-service/internal/generated"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

//go:generate mockery --disable-version-string --case=underscore --name=client --structname=ClientMock
type (
	client interface {
		SendNotification(ctx context.Context, in *generated.NotificationDTO, opts ...grpc.CallOption) (*generated.Response, error)
	}

	Wrapper struct {
		client     client
		clientHost string
		clientPort string
		channel    *amqp.Channel
		queueName  string
	}
)

func New(client client, clientHost string, clientPort string, channel *amqp.Channel, queueName string) *Wrapper {
	return &Wrapper{
		client:     client,
		clientHost: clientHost,
		clientPort: clientPort,
		channel:    channel,
		queueName:  queueName,
	}
}
