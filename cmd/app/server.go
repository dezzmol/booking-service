package app

import (
	"booking-service/internal/app"
	"booking-service/internal/controllers"
	"booking-service/internal/storage"
)

func (a *App) initDB() error {
	db, err := storage.NewDB(a.config.Db.Host, a.config.Db.User, a.config.Db.Password, a.config.Db.Name, a.config.Db.Port)
	if err != nil {
		return err
	}
	a.PostgreSQL = db
	return nil
}

func (a *App) initControllers() {
	a.Controllers.BookingController = controllers.New(a.PostgreSQL, a.Storage)
}

func (a *App) initHandlers() {
	a.Handlers.booking = app.New(a.Controllers.BookingController)
}
