package booking_test

import (
	"context"
	"errors"
	"testing"
	"time"

	bookingService "booking-service/internal/controllers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"booking-service/internal/entities"
	"booking-service/internal/repositories/booking/mocks"
	guestMocks "booking-service/internal/repositories/guest/mocks"
)

//
// Тесты для метода CreateBooking
//

// TestCreateBooking_Success проверяет успешное создание бронирования.
func TestCreateBooking_Success(t *testing.T) {
	mockBookingRepo := new(mocks.BookingRepoMock)
	mockGuestRepo := new(guestMocks.GuestRepoMock)
	service := bookingService.NewBookingService(mockBookingRepo, mockGuestRepo)

	ctx := context.Background()
	// Задаем входные данные: бронирование с двумя гостями.
	startDate := time.Now().Add(48 * time.Hour)
	endDate := startDate.Add(24 * time.Hour)
	input := entities.CreateBookingDTO{
		RoomID:    101,
		StartDate: startDate,
		EndDate:   endDate,
		Comment:   "Тестовое бронирование",
		Guests: []entities.GuestDTO{
			{Name: "Иван Иванов"},
			{Name: "Петр Петров"},
		},
	}

	// Ожидаем, что комната доступна.
	mockBookingRepo.
		On("IsRoomAvailable", ctx, input.RoomID, input.StartDate, input.EndDate).
		Return(true, nil).
		Once()

	// Для каждого гостя вызывается SaveAndReturnIt, возвращаем обновленные данные (например, с присвоенным ID).
	returnedGuest1 := entities.Guest{ID: 1, Name: "Иван Иванов"}
	returnedGuest2 := entities.Guest{ID: 2, Name: "Петр Петров"}

	// Первый вызов для первого гостя.
	mockGuestRepo.
		On("SaveAndReturnIt", ctx, mock.MatchedBy(func(g *entities.Guest) bool {
			return g.Name == "Иван Иванов"
		})).
		Return(returnedGuest1, nil).
		Once()

	// Второй вызов для второго гостя.
	mockGuestRepo.
		On("SaveAndReturnIt", ctx, mock.MatchedBy(func(g *entities.Guest) bool {
			return g.Name == "Петр Петров"
		})).
		Return(returnedGuest2, nil).
		Once()

	// Ожидаем вызов метода Save для бронирования.
	mockBookingRepo.
		On("Save", ctx, mock.AnythingOfType("*entities.Booking")).
		Return(nil).
		Once()

	booking, err := service.CreateBooking(ctx, input)
	assert.NoError(t, err)
	assert.Equal(t, input.RoomID, booking.RoomID)
	assert.Equal(t, input.Comment, booking.Comment)
	// Проверяем, что статус бронирования выставлен как pending и платеж не оплачен.
	assert.Equal(t, bookingService.BOOKING_PENDING, booking.Status)
	assert.Equal(t, bookingService.PAYMENT_UNPAID, booking.PaymentStatus)
	// Также проверяем, что гости обновились (возвращены с ID)
	assert.Len(t, booking.Guests, 2)
	assert.Equal(t, returnedGuest1, booking.Guests[0])
	assert.Equal(t, returnedGuest2, booking.Guests[1])

	mockBookingRepo.AssertExpectations(t)
	mockGuestRepo.AssertExpectations(t)
}

// TestCreateBooking_RoomNotAvailable проверяет ситуацию, когда комната недоступна.
func TestCreateBooking_RoomNotAvailable(t *testing.T) {
	mockBookingRepo := new(mocks.BookingRepoMock)
	mockGuestRepo := new(guestMocks.GuestRepoMock)
	service := bookingService.NewBookingService(mockBookingRepo, mockGuestRepo)

	ctx := context.Background()
	input := entities.CreateBookingDTO{
		RoomID:    101,
		StartDate: time.Now().Add(48 * time.Hour),
		EndDate:   time.Now().Add(72 * time.Hour),
		Comment:   "Тест бронирования",
		Guests: []entities.GuestDTO{
			{Name: "Иван Иванов"},
		},
	}

	mockBookingRepo.
		On("IsRoomAvailable", ctx, input.RoomID, input.StartDate, input.EndDate).
		Return(false, nil).
		Once()

	booking, err := service.CreateBooking(ctx, input)
	assert.Error(t, err)
	assert.Equal(t, entities.ErrRoomNotAvailable, err)
	assert.Equal(t, entities.Booking{}, booking)

	mockBookingRepo.AssertExpectations(t)
	// guestRepo не должен вызываться, поэтому можно проверить отсутствие вызовов.
	mockGuestRepo.AssertExpectations(t)
}

// TestCreateBooking_ErrorOnGuestSave проверяет ошибку при сохранении гостя.
func TestCreateBooking_ErrorOnGuestSave(t *testing.T) {
	mockBookingRepo := new(mocks.BookingRepoMock)
	mockGuestRepo := new(guestMocks.GuestRepoMock)
	service := bookingService.NewBookingService(mockBookingRepo, mockGuestRepo)

	ctx := context.Background()
	input := entities.CreateBookingDTO{
		RoomID:    101,
		StartDate: time.Now().Add(48 * time.Hour),
		EndDate:   time.Now().Add(72 * time.Hour),
		Comment:   "Тест бронирования",
		Guests: []entities.GuestDTO{
			{Name: "Иван Иванов"},
		},
	}

	// Комната доступна.
	mockBookingRepo.
		On("IsRoomAvailable", ctx, input.RoomID, input.StartDate, input.EndDate).
		Return(true, nil).
		Once()

	// guestRepo возвращает ошибку при сохранении гостя.
	expErr := errors.New("ошибка сохранения гостя")
	mockGuestRepo.
		On("SaveAndReturnIt", ctx, mock.AnythingOfType("*entities.Guest")).
		Return(entities.Guest{}, expErr).
		Once()

	booking, err := service.CreateBooking(ctx, input)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
	assert.Equal(t, entities.Booking{}, booking)

	mockBookingRepo.AssertExpectations(t)
	mockGuestRepo.AssertExpectations(t)
}

// TestCreateBooking_ErrorOnBookingSave проверяет ошибку при сохранении бронирования.
func TestCreateBooking_ErrorOnBookingSave(t *testing.T) {
	mockBookingRepo := new(mocks.BookingRepoMock)
	mockGuestRepo := new(guestMocks.GuestRepoMock)
	service := bookingService.NewBookingService(mockBookingRepo, mockGuestRepo)

	ctx := context.Background()
	input := entities.CreateBookingDTO{
		RoomID:    101,
		StartDate: time.Now().Add(48 * time.Hour),
		EndDate:   time.Now().Add(72 * time.Hour),
		Comment:   "Тест бронирования",
		Guests: []entities.GuestDTO{
			{Name: "Иван Иванов"},
		},
	}

	// Комната доступна.
	mockBookingRepo.
		On("IsRoomAvailable", ctx, input.RoomID, input.StartDate, input.EndDate).
		Return(true, nil).
		Once()

	// guestRepo успешно сохраняет гостя.
	returnedGuest := entities.Guest{ID: 1, Name: "Иван Иванов"}
	mockGuestRepo.
		On("SaveAndReturnIt", ctx, mock.AnythingOfType("*entities.Guest")).
		Return(returnedGuest, nil).
		Once()

	// bookingRepo.Save возвращает ошибку.
	expErr := errors.New("ошибка сохранения бронирования")
	mockBookingRepo.
		On("Save", ctx, mock.AnythingOfType("*entities.Booking")).
		Return(expErr).
		Once()

	booking, err := service.CreateBooking(ctx, input)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
	assert.Equal(t, entities.Booking{}, booking)

	mockBookingRepo.AssertExpectations(t)
	mockGuestRepo.AssertExpectations(t)
}

//
// Тесты для метода CancelBooking
//

// TestCancelBooking_Success проверяет успешную отмену бронирования.
func TestCancelBooking_Success(t *testing.T) {
	mockBookingRepo := new(mocks.BookingRepoMock)

	// guestRepo здесь не используется
	service := bookingService.NewBookingService(mockBookingRepo, nil)

	ctx := context.Background()
	bookingID := uint64(1)
	// Исходное бронирование со статусом pending.
	existingBooking := &entities.Booking{
		ID:     bookingID,
		Status: bookingService.BOOKING_PENDING,
	}
	mockBookingRepo.
		On("FindById", ctx, bookingID).
		Return(existingBooking, nil).
		Once()

	// При обновлении ожидается, что статус изменится на cancelled.
	mockBookingRepo.
		On("Update", ctx, existingBooking).
		Return(nil).
		Once()

	err := service.CancelBooking(ctx, bookingID)
	assert.NoError(t, err)
	assert.Equal(t, bookingService.BOOKING_CANCELLED, existingBooking.Status)

	mockBookingRepo.AssertExpectations(t)
}

// TestCancelBooking_FindError проверяет, что ошибка при поиске бронирования передается дальше.
func TestCancelBooking_FindError(t *testing.T) {
	mockBookingRepo := new(mocks.BookingRepoMock)
	service := bookingService.NewBookingService(mockBookingRepo, nil)

	ctx := context.Background()
	bookingID := uint64(1)
	expErr := errors.New("бронирование не найдено")
	mockBookingRepo.
		On("FindById", ctx, bookingID).
		Return((*entities.Booking)(nil), expErr).
		Once()

	err := service.CancelBooking(ctx, bookingID)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)

	mockBookingRepo.AssertExpectations(t)
}

// TestCancelBooking_UpdateError проверяет ситуацию, когда обновление бронирования завершается ошибкой.
func TestCancelBooking_UpdateError(t *testing.T) {
	mockBookingRepo := new(mocks.BookingRepoMock)
	service := bookingService.NewBookingService(mockBookingRepo, nil)

	ctx := context.Background()
	bookingID := uint64(1)
	existingBooking := &entities.Booking{
		ID:     bookingID,
		Status: bookingService.BOOKING_PENDING,
	}
	expErr := errors.New("ошибка обновления")

	mockBookingRepo.
		On("FindById", ctx, bookingID).
		Return(existingBooking, nil).
		Once()

	mockBookingRepo.
		On("Update", ctx, existingBooking).
		Return(expErr).
		Once()

	err := service.CancelBooking(ctx, bookingID)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)

	mockBookingRepo.AssertExpectations(t)
}

//
// Тесты для метода ModifyBooking
//

// TestModifyBooking_Success проверяет успешное изменение бронирования.
func TestModifyBooking_Success(t *testing.T) {
	mockBookingRepo := new(mocks.BookingRepoMock)
	service := bookingService.NewBookingService(mockBookingRepo, nil)

	ctx := context.Background()
	bookingID := uint64(1)
	// Новые даты переноса
	newStart := time.Now().Add(10 * 24 * time.Hour) // более чем через 7 дней
	newEnd := newStart.Add(24 * time.Hour)

	// Исходное бронирование с датой начала достаточно в будущем.
	existingBooking := &entities.Booking{
		ID:        bookingID,
		Status:    bookingService.BOOKING_PENDING,
		StartDate: time.Now().Add(15 * 24 * time.Hour),
		EndDate:   time.Now().Add(16 * 24 * time.Hour),
	}

	// Ожидается, что перенести бронирование можно.
	mockBookingRepo.
		On("CanReschedule", ctx, bookingID, newStart, newEnd).
		Return(true, nil).
		Once()

	mockBookingRepo.
		On("FindById", ctx, bookingID).
		Return(existingBooking, nil).
		Once()

	// Ожидаем вызов обновления с измененными датами.
	mockBookingRepo.
		On("Update", ctx, existingBooking).
		Return(nil).
		Once()

	updatedBooking, err := service.ModifyBooking(ctx, bookingID, newStart, newEnd)
	assert.NoError(t, err)
	assert.True(t, updatedBooking.StartDate.Equal(newStart))
	assert.True(t, updatedBooking.EndDate.Equal(newEnd))

	mockBookingRepo.AssertExpectations(t)
}

// TestModifyBooking_CanRescheduleFalse проверяет, если перенос невозможен.
func TestModifyBooking_CanRescheduleFalse(t *testing.T) {
	mockBookingRepo := new(mocks.BookingRepoMock)
	service := bookingService.NewBookingService(mockBookingRepo, nil)

	ctx := context.Background()
	bookingID := uint64(1)
	newStart := time.Now().Add(10 * 24 * time.Hour)
	newEnd := newStart.Add(24 * time.Hour)

	mockBookingRepo.
		On("CanReschedule", ctx, bookingID, newStart, newEnd).
		Return(false, nil).
		Once()

	updatedBooking, err := service.ModifyBooking(ctx, bookingID, newStart, newEnd)
	assert.Error(t, err)
	assert.Equal(t, entities.ErrRoomNotAvailable, err)
	assert.Equal(t, entities.Booking{}, updatedBooking)

	mockBookingRepo.AssertExpectations(t)
}

// TestModifyBooking_BookingCancelled проверяет попытку изменения отмененного бронирования.
func TestModifyBooking_BookingCancelled(t *testing.T) {
	mockBookingRepo := new(mocks.BookingRepoMock)
	service := bookingService.NewBookingService(mockBookingRepo, nil)

	ctx := context.Background()
	bookingID := uint64(1)
	newStart := time.Now().Add(10 * 24 * time.Hour)
	newEnd := newStart.Add(24 * time.Hour)

	// Бронирование уже отменено.
	existingBooking := &entities.Booking{
		ID:     bookingID,
		Status: bookingService.BOOKING_CANCELLED,
	}

	mockBookingRepo.
		On("CanReschedule", ctx, bookingID, newStart, newEnd).
		Return(true, nil).
		Once()

	mockBookingRepo.
		On("FindById", ctx, bookingID).
		Return(existingBooking, nil).
		Once()

	updatedBooking, err := service.ModifyBooking(ctx, bookingID, newStart, newEnd)
	assert.Error(t, err)
	assert.Equal(t, bookingService.BookingCancelled, err)
	assert.Equal(t, entities.Booking{}, updatedBooking)

	mockBookingRepo.AssertExpectations(t)
}

// TestModifyBooking_TooLateReschedule проверяет ошибку, если попытка переноса бронирования производится слишком поздно.
func TestModifyBooking_TooLateReschedule(t *testing.T) {
	mockBookingRepo := new(mocks.BookingRepoMock)
	service := bookingService.NewBookingService(mockBookingRepo, nil)

	ctx := context.Background()
	bookingID := uint64(1)
	newStart := time.Now().Add(10 * 24 * time.Hour)
	newEnd := newStart.Add(24 * time.Hour)
	// Бронирование с датой начала менее чем через 7 дней.
	existingBooking := &entities.Booking{
		ID:        bookingID,
		Status:    bookingService.BOOKING_PENDING,
		StartDate: time.Now().Add(5 * 24 * time.Hour),
		EndDate:   time.Now().Add(6 * 24 * time.Hour),
	}

	mockBookingRepo.
		On("CanReschedule", ctx, bookingID, newStart, newEnd).
		Return(true, nil).
		Once()

	mockBookingRepo.
		On("FindById", ctx, bookingID).
		Return(existingBooking, nil).
		Once()

	updatedBooking, err := service.ModifyBooking(ctx, bookingID, newStart, newEnd)
	assert.Error(t, err)
	assert.Equal(t, bookingService.CantRescheduleLaterThanSevenDaysBeforeStartBooking, err)
	assert.Equal(t, entities.Booking{}, updatedBooking)

	mockBookingRepo.AssertExpectations(t)
}

// TestModifyBooking_UpdateError проверяет ситуацию, когда обновление бронирования завершается ошибкой.
func TestModifyBooking_UpdateError(t *testing.T) {
	mockBookingRepo := new(mocks.BookingRepoMock)
	service := bookingService.NewBookingService(mockBookingRepo, nil)

	ctx := context.Background()
	bookingID := uint64(1)
	newStart := time.Now().Add(10 * 24 * time.Hour)
	newEnd := newStart.Add(24 * time.Hour)

	existingBooking := &entities.Booking{
		ID:        bookingID,
		Status:    bookingService.BOOKING_PENDING,
		StartDate: time.Now().Add(15 * 24 * time.Hour),
		EndDate:   time.Now().Add(16 * 24 * time.Hour),
	}

	mockBookingRepo.
		On("CanReschedule", ctx, bookingID, newStart, newEnd).
		Return(true, nil).
		Once()

	mockBookingRepo.
		On("FindById", ctx, bookingID).
		Return(existingBooking, nil).
		Once()

	expErr := errors.New("ошибка обновления")
	mockBookingRepo.
		On("Update", ctx, existingBooking).
		Return(expErr).
		Once()

	updatedBooking, err := service.ModifyBooking(ctx, bookingID, newStart, newEnd)
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
	assert.Equal(t, entities.Booking{}, updatedBooking)

	mockBookingRepo.AssertExpectations(t)
}
