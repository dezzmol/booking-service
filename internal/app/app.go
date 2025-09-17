package app

import (
	"log"
)

type (
	App interface {
		Run()
		Stop()
	}

	Application struct {
		config  *Config
		storage *Storage

		repositories *Repositories
		clients      *Clients
		services     *Services
		handlers     *Handlers
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
	a.initRepositories()
	a.initServices()
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
	a.storage.db.Close()
}
