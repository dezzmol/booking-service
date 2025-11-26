package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"booking-service/internal/app"
	"booking-service/internal/controllers"
	"booking-service/internal/generated"
	"booking-service/internal/storage"

	"github.com/jmoiron/sqlx"
)

type (
	App struct {
		config *Config

		Handlers struct {
			booking *app.Handler
		}

		Controllers struct {
			BookingController *controllers.Controller
		}

		PostgreSQL *sqlx.DB
		Storage    *storage.Storage

		Clients struct {
			payment generated.PaymentServiceClient
		}
	}
)

func New() *App {
	return &App{}
}

func (a *App) Run() {
	err := a.InitConfig()
	if err != nil {
		log.Printf("failed to initialize config: %s\n", err)
		return
	}

	err = a.initDB()
	if err != nil {
		log.Printf("failed to initialize database: %s\n", err)
		return
	}
	defer a.Stop()

	a.initClients()
	a.initControllers()
	a.initHandlers()

	go a.initGRPC()
	go a.initHTTP()
	//err = a.initConsul()
	//if err != nil {
	//	log.Printf("failed to initialize consul: %s\n", err)
	//	return
	//}
	a.initSwagger()
	log.Println("application started")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Received shutdown signal, initiating graceful shutdown...")
}

func (a *App) Stop() {
	a.PostgreSQL.Close()
}
