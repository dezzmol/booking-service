package app

import (
	"log"

	"booking-service/internal/controllers"
	"booking-service/internal/generated"
	"booking-service/internal/handlers"
	"booking-service/internal/storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (a *Application) initDB() error {
	db, err := storage.NewDB(a.config.Db.Host, a.config.Db.User, a.config.Db.Password, a.config.Db.Name, a.config.Db.Port)
	if err != nil {
		return err
	}
	a.PostgreSQL = db
	return nil
}

func (a *Application) initControllers() {
	a.Controllers.BookingController = controllers.New(a.PostgreSQL, a.Storage)
}

func (a *Application) initHandlers() {
	a.Handlers.booking = handlers.New(a.Controllers.BookingController)
}

func (a *Application) initClients() {
	a.initNotificationClient()
	a.initPaymentClient()
}

func (a *Application) initNotificationClient() {
	conn, err := grpc.NewClient("dns:///"+a.config.NotificationClient.Host+":"+a.config.NotificationClient.Grpc.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := generated.NewNotificationServiceClient(conn)
	a.Clients.notificationQueue = notifications.New(
		client,
		a.config.NotificationClient.Host,
		a.config.NotificationClient.Port,
		a.clients.notificationChannel,
		a.clients.notificationQueue.Name,
	)
	log.Printf("Notification client initialized on %v\n",
		a.config.NotificationClient.Host+":"+a.config.NotificationClient.Grpc.Port)
}

func (a *Application) initPaymentClient() {
	conn, err := grpc.NewClient("dns:///"+a.config.PaymentClient.Host+":"+a.config.PaymentClient.Grpc.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client := generated.NewPaymentServiceClient(conn)
	a.clients.payments = payments.New(client)
	log.Printf("Payment client initialized on %v\n",
		a.config.PaymentClient.Host+":"+a.config.PaymentClient.Grpc.Port)
}
