package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"booking-service/internal/controllers/booking/mocks"
	"booking-service/internal/entities"
	"booking-service/internal/generated"
	"booking-service/internal/handlers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateBooking_Success(t *testing.T) {
	// Arrange
	bookingServiceMock := &mocks.BookingServiceMock{}
	h := &handlers.Handler{BookingService: bookingServiceMock}

	startDate := time.Now()
	endDate := startDate.Add(24 * time.Hour)

	expectedBooking := entities.Booking{
		ID:        uint64(123),
		RoomID:    uint64(123),
		StartDate: startDate,
		EndDate:   endDate,
		Status:    "confirmed",
	}

	bookingServiceMock.On("CreateBooking", mock.Anything, mock.AnythingOfType("entities.CreateBookingDTO")).
		Return(expectedBooking, nil)

	// Act
	resp, err := h.CreateBooking(context.Background(), &generated.CreateBookingRequest{
		RoomId:    uint64(123),
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
		Guests: []*generated.GuestDTO{
			{Name: "John Doe"},
		},
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint64(123), resp.Booking.Id)
	assert.Equal(t, uint64(123), resp.Booking.RoomId)
	assert.Len(t, resp.Links, 1)
	assert.Equal(t, "self", resp.Links[0].Rel)

	bookingServiceMock.AssertExpectations(t)
}

func TestCreateBooking_RoomNotAvailable(t *testing.T) {
	// Arrange
	bookingServiceMock := &mocks.BookingServiceMock{}
	h := &handlers.Handler{BookingService: bookingServiceMock}

	startDate := time.Now()
	endDate := startDate.Add(24 * time.Hour)

	bookingServiceMock.On("CreateBooking", mock.Anything, mock.AnythingOfType("entities.CreateBookingDTO")).
		Return(entities.Booking{}, entities.ErrRoomNotAvailable)

	// Act
	resp, err := h.CreateBooking(context.Background(), &generated.CreateBookingRequest{
		RoomId:    uint64(123),
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	})

	// Assert
	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Equal(t, "room is not available", st.Message())

	bookingServiceMock.AssertExpectations(t)
}

func TestCreateBooking_InternalError(t *testing.T) {
	// Arrange
	bookingServiceMock := &mocks.BookingServiceMock{}
	h := &handlers.Handler{BookingService: bookingServiceMock}

	startDate := time.Now()
	endDate := startDate.Add(24 * time.Hour)

	expectedErr := errors.New("database connection failed")
	bookingServiceMock.On("CreateBooking", mock.Anything, mock.AnythingOfType("entities.CreateBookingDTO")).
		Return(entities.Booking{}, expectedErr)

	// Act
	resp, err := h.CreateBooking(context.Background(), &generated.CreateBookingRequest{
		RoomId:    uint64(123),
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	})

	// Assert
	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "database connection failed")

	bookingServiceMock.AssertExpectations(t)
}

func TestCreateBooking_InvalidDates(t *testing.T) {
	// Arrange
	bookingServiceMock := &mocks.BookingServiceMock{}
	h := &handlers.Handler{BookingService: bookingServiceMock}

	startDate := time.Now()
	endDate := startDate.Add(-24 * time.Hour) // End date before start date

	bookingServiceMock.On("CreateBooking", mock.Anything, mock.AnythingOfType("entities.CreateBookingDTO")).
		Return(entities.Booking{}, entities.ErrStartDateIsAfterEndDate)
	// Act
	resp, err := h.CreateBooking(context.Background(), &generated.CreateBookingRequest{
		RoomId:    uint64(123),
		StartDate: timestamppb.New(startDate),
		EndDate:   timestamppb.New(endDate),
	})

	// Assert
	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	// Сервис не должен вызываться при невалидных датах
	bookingServiceMock.AssertNotCalled(t, "CreateBooking")
}
