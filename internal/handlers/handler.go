package handlers

import (
	"booking-service/internal/generated"
	"booking-service/internal/services/booking"
	"booking-service/internal/services/employee"
	"booking-service/internal/services/guest"
	houseKeepingService "booking-service/internal/services/house_keeping_request"
	"booking-service/internal/services/review"
	"booking-service/internal/services/room"
)

type (
	Handler struct {
		generated.UnimplementedBookingServiceServer

		BookingService      booking.BookingService
		EmployeeService     employee.EmployeeService
		GuestService        guest.GuestService
		HousekeepingService houseKeepingService.HouseKeepingRequestService
		ReviewService       review.ReviewService
		RoomService         room.RoomService
	}
)

func New(
	bookingService booking.BookingService,
	employeeService employee.EmployeeService,
	guestService guest.GuestService,
	housekeepingService houseKeepingService.HouseKeepingRequestService,
	reviewService review.ReviewService,
	roomService room.RoomService,
) *Handler {
	return &Handler{
		BookingService:      bookingService,
		EmployeeService:     employeeService,
		GuestService:        guestService,
		HousekeepingService: housekeepingService,
		ReviewService:       reviewService,
		RoomService:         roomService,
	}
}
