package booking

import (
	"context"
	"errors"
	"time"

	"booking-service/internal/entities"
	"booking-service/internal/repositories/booking"
	"booking-service/internal/repositories/guest"
)

//go:generate mockery --disable-version-string --case=underscore --name=BookingService --structname=BookingServiceMock
var (
	BookingCancelled                                   = errors.New("booking cancelled")
	CantRescheduleLaterThanSevenDaysBeforeStartBooking = errors.New("can't reschedule later than seven days before start booking")
)

const (
	MESSAGE_HAS_BEEN_SENT     = "has_been_sent"
	MESSAGE_HAS_BEEN_RECEIVED = "has_not_been_cent"

	BOOKING_CREATE_MESSAGE = "Бронирование успешно создано"

	PAYMENT_PAID      = "paid"
	PAYMENT_UNPAID    = "unpaid"
	PAYMENT_CANCELLED = "cancelled"

	BOOKING_PENDING   = "pending"
	BOOKING_CONFIRM   = "confirmed"
	BOOKING_CANCELLED = "cancelled"
)

type BookingService interface {
	CreateBooking(ctx context.Context, input entities.CreateBookingDTO) (entities.Booking, error)
	CancelBooking(ctx context.Context, bookingID uint64) error
	ModifyBooking(ctx context.Context, bookingID uint64, startDate, endDate time.Time) (entities.Booking, error)
}

type Service struct {
	bookingRepo booking.BookingRepo
	guestRepo   guest.GuestRepo
}

func NewBookingService(
	bookingRepo booking.BookingRepo,
	guestRepo guest.GuestRepo,
) *Service {
	return &Service{
		bookingRepo: bookingRepo,
		guestRepo:   guestRepo,
	}
}

func (s *Service) CreateBooking(ctx context.Context, input entities.CreateBookingDTO) (entities.Booking, error) {
	available, err := s.bookingRepo.IsRoomAvailable(ctx, input.RoomID, input.StartDate, input.EndDate)
	if err != nil {
		return entities.Booking{}, err
	}
	if !available {
		return entities.Booking{}, entities.ErrRoomNotAvailable
	}
	if input.StartDate.After(input.EndDate) {
		return entities.Booking{}, entities.ErrStartDateIsAfterEndDate
	}

	var guests []entities.Guest
	for _, guest := range input.Guests {
		guests = append(guests, entities.Guest{
			Name: guest.Name,
		})
	}

	for i := range guests {
		guest, err := s.guestRepo.SaveAndReturnIt(ctx, &guests[i])
		if err != nil {
			return entities.Booking{}, err
		}
		guests[i] = guest
	}

	booking := entities.Booking{
		RoomID:        input.RoomID,
		StartDate:     input.StartDate,
		EndDate:       input.EndDate,
		Comment:       input.Comment,
		Guests:        guests,
		Status:        BOOKING_PENDING,
		PaymentStatus: PAYMENT_UNPAID,
	}

	err = s.bookingRepo.Save(ctx, &booking)
	if err != nil {
		return entities.Booking{}, err
	}

	return booking, nil
}

func (s *Service) CancelBooking(ctx context.Context, bookingID uint64) error {
	booking, err := s.bookingRepo.FindById(ctx, bookingID)
	if err != nil {
		return err
	}

	booking.Status = BOOKING_CANCELLED

	err = s.bookingRepo.Update(ctx, booking)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) ModifyBooking(ctx context.Context, bookingID uint64, startDate, endDate time.Time) (entities.Booking, error) {
	canReschedule, err := s.bookingRepo.CanReschedule(ctx, bookingID, startDate, endDate)
	if err != nil {
		return entities.Booking{}, err
	}
	if !canReschedule {
		return entities.Booking{}, entities.ErrRoomNotAvailable
	}

	booking, err := s.bookingRepo.FindById(ctx, bookingID)
	if err != nil {
		return entities.Booking{}, err
	}
	if booking.Status == BOOKING_CANCELLED {
		return entities.Booking{}, BookingCancelled
	}

	if time.Now().AddDate(0, 0, 7).After(booking.StartDate) {
		return entities.Booking{}, CantRescheduleLaterThanSevenDaysBeforeStartBooking
	}

	booking.StartDate = startDate
	booking.EndDate = endDate

	err = s.bookingRepo.Update(ctx, booking)
	if err != nil {
		return entities.Booking{}, err
	}

	return *booking, err
}
