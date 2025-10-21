package app

import (
	"booking-service/internal/controllers"
	"booking-service/internal/generated"
)

type (
	Handler struct {
		generated.UnimplementedBookingServiceServer

		bookingController *controllers.Controller
	}
)

func New(
	bookingController *controllers.Controller,
) *Handler {
	return &Handler{
		bookingController: bookingController,
	}
}
