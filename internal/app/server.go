package app

import (
	"log"

	"booking-service/internal/generated"
	"booking-service/internal/handlers"
	"booking-service/internal/repositories/booking"
	"booking-service/internal/repositories/employee"
	"booking-service/internal/repositories/guest"
	"booking-service/internal/repositories/house_keeping_request"
	"booking-service/internal/repositories/notifications"
	"booking-service/internal/repositories/payments"
	"booking-service/internal/repositories/review"
	"booking-service/internal/repositories/room"
	bookingService "booking-service/internal/services/booking"
	employeeService "booking-service/internal/services/employee"
	guestService "booking-service/internal/services/guest"
	houseKeepingService "booking-service/internal/services/house_keeping_request"
	reviewService "booking-service/internal/services/review"
	roomService "booking-service/internal/services/room"
	"booking-service/internal/storage"

	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Storage struct {
	db *sqlx.DB
}

type Handlers struct {
	booking *handlers.Handler
}

type Services struct {
	booking             *bookingService.Service
	employee            *employeeService.Service
	guest               *guestService.Service
	houseKeepingRequest *houseKeepingService.Service
	review              *reviewService.Service
	room                *roomService.Service
}

type Repositories struct {
	booking             *booking.BookingRepository
	employee            *employee.Repository
	guest               *guest.GuestRepository
	houseKeepingRequest *house_keeping_request.HouseKeepingRequestRepository
	review              *review.ReviewRepository
	room                *room.RoomRepository
}

type Clients struct {
	notificationQueue   amqp.Queue
	notificationChannel *amqp.Channel

	notifications *notifications.Wrapper
	payments      *payments.Wrapper
}

func (a *Application) initDB() error {
	db, err := storage.NewDB(a.config.Db.Host, a.config.Db.User, a.config.Db.Password, a.config.Db.Name, a.config.Db.Port)
	if err != nil {
		return err
	}
	a.storage = &Storage{
		db: db,
	}
	return nil
}

func (a *Application) initRepositories() {
	a.repositories = &Repositories{
		booking:             booking.NewBookingRepository(a.storage.db),
		employee:            employee.NewRepository(a.storage.db),
		guest:               guest.NewGuestRepository(a.storage.db),
		houseKeepingRequest: house_keeping_request.NewHouseKeepingRequestRepository(a.storage.db),
		review:              review.NewReviewRepository(a.storage.db),
		room:                room.NewRoomRepository(a.storage.db),
	}
}

func (a *Application) initServices() {
	a.services = &Services{}
	a.services.booking = bookingService.NewBookingService(a.repositories.booking, a.repositories.guest)
	a.services.employee = employeeService.NewService(a.repositories.employee)
	a.services.guest = guestService.NewService(a.repositories.guest)
	a.services.houseKeepingRequest = houseKeepingService.New(a.repositories.houseKeepingRequest)
	a.services.review = reviewService.New(a.repositories.booking, a.repositories.review)
	a.services.room = roomService.New(a.repositories.room, a.clients.notifications)
}

func (a *Application) initHandlers() {
	a.handlers = &Handlers{}
	a.handlers.booking = handlers.New(
		a.services.booking,
		a.services.employee,
		a.services.guest,
		a.services.houseKeepingRequest,
		a.services.review,
		a.services.room,
	)
}

func (a *Application) initClients() {
	if a.clients == nil {
		a.clients = &Clients{}
	}
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
	a.clients.notifications = notifications.New(
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
