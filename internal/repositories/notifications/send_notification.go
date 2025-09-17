package notifications

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"booking-service/internal/generated"
)

func (w *Wrapper) SendNotification(ctx context.Context, message string, receiver uint64) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	body := &generated.NotificationDTO{Message: message, RecipientId: receiver}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	_, err = http.Post("http://"+w.clientHost+":"+w.clientPort+"/v1/notification",
		"application/json", strings.NewReader(string(jsonBody)))
	if err != nil {
		return err
	}

	return nil
}
