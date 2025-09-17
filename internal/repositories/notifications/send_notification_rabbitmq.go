package notifications

import (
	"context"
	"encoding/json"
	"log"

	"booking-service/internal/generated"

	"github.com/streadway/amqp"
)

func (w *Wrapper) SendNotificationWithQueue(ctx context.Context, message string, receiver uint64) error {
	body := &generated.NotificationDTO{Message: message, RecipientId: receiver}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	err = w.channel.Publish(
		"",
		w.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		},
	)
	log.Printf("[SendNotificationWithQueue] Sent %s", body)
	if err != nil {
		return err
	}

	return nil
}
