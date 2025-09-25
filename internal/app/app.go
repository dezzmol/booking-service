package app

import (
	"log"

	"booking-service/internal/controllers"
	"booking-service/internal/handlers"
	"booking-service/internal/storage"

	"github.com/jmoiron/sqlx"
	"github.com/streadway/amqp"
)

type (
	App interface {
		Run()
		Stop()
	}

	Application struct {
		config *Config

		Handlers struct {
			booking *handlers.Handler
		}

		Controllers struct {
			BookingController *controllers.Controller
		}

		PostgreSQL *sqlx.DB
		Storage    *storage.Storage

		Clients struct {
			notificationQueue   amqp.Queue
			notificationChannel *amqp.Channel
		}
	}
)

func New() *Application {
	return &Application{}
}

func (a *Application) Run() {
	err := a.InitConfig()
	if err != nil {
		log.Printf("failed to initialize config: %s\n", err)
		return
	}

	err = a.initDB()
	defer a.Stop()
	if err != nil {
		log.Printf("failed to initialize database: %s\n", err)
		return
	}

	err = a.initProducer()
	if err != nil {
		log.Printf("failed to initialize producer: %s\n", err)
		return
	}

	a.initClients()
	a.initControllers()
	a.initHandlers()

	go a.initGRPC()
	go a.initHTTP()
	err = a.initConsul()
	if err != nil {
		log.Printf("failed to initialize consul: %s\n", err)
		return
	}
	a.initSwagger()
	log.Println("application started")
	select {}
}

func (a *Application) Stop() {
	a.PostgreSQL.Close()
}
