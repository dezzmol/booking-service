package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"booking-service/internal/entities"
	"booking-service/internal/generated"
	"booking-service/internal/handlers"
	"booking-service/internal/services/guest/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetGuestsByBookingID_Success(t *testing.T) {
	// Arrange
	guestServiceMock := &mocks.GuestServiceMock{}
	h := &handlers.Handler{GuestService: guestServiceMock}

	guestID1 := uint64(1)
	guestID2 := uint64(2)
	now := time.Now()
	mockGuests := []entities.Guest{
		{
			ID:        guestID1,
			Name:      "John Doe",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        guestID2,
			Name:      "Jane Smith",
			CreatedAt: now.Add(-24 * time.Hour),
			UpdatedAt: now.Add(-12 * time.Hour),
		},
	}

	guestServiceMock.On("FindGuestsByBookingID", mock.Anything, uint(123)).
		Return(mockGuests, nil)

	// Act
	resp, err := h.GetGuestsByBookingID(context.Background(), &generated.GetGuestsByBookingIDRequest{
		BookingId: 123,
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Guests, 2)

	// Verify first guest
	assert.Equal(t, guestID1, resp.Guests[0].Id)
	assert.Equal(t, "John Doe", resp.Guests[0].Name)
	assert.True(t, timestamppb.New(now).AsTime().Equal(resp.Guests[0].CreatedAt.AsTime()))

	// Verify second guest
	assert.Equal(t, guestID2, resp.Guests[1].Id)
	assert.Equal(t, "Jane Smith", resp.Guests[1].Name)

	// Verify links
	assert.Len(t, resp.Links, 1)
	assert.Equal(t, "self", resp.Links[0].Rel)
	assert.Equal(t, "/v1/booking/{booking_id}/guests", resp.Links[0].Href)

	guestServiceMock.AssertExpectations(t)
}

func TestGetGuestsByBookingID_NotFound(t *testing.T) {
	// Arrange
	guestServiceMock := &mocks.GuestServiceMock{}
	h := &handlers.Handler{GuestService: guestServiceMock}

	guestServiceMock.On("FindGuestsByBookingID", mock.Anything, uint(404)).
		Return(nil, entities.ErrNotFound)

	// Act
	resp, err := h.GetGuestsByBookingID(context.Background(), &generated.GetGuestsByBookingIDRequest{
		BookingId: 404,
	})

	// Assert
	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, "booking not found", st.Message())

	guestServiceMock.AssertExpectations(t)
}

func TestGetGuestsByBookingID_InternalError(t *testing.T) {
	// Arrange
	guestServiceMock := &mocks.GuestServiceMock{}
	h := &handlers.Handler{GuestService: guestServiceMock}

	expectedErr := errors.New("database connection failed")
	guestServiceMock.On("FindGuestsByBookingID", mock.Anything, uint(500)).
		Return(nil, expectedErr)

	// Act
	resp, err := h.GetGuestsByBookingID(context.Background(), &generated.GetGuestsByBookingIDRequest{
		BookingId: 500,
	})

	// Assert
	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "database connection failed")

	guestServiceMock.AssertExpectations(t)
}

func TestGetGuestsByBookingID_EmptyResult(t *testing.T) {
	// Arrange
	guestServiceMock := &mocks.GuestServiceMock{}
	h := &handlers.Handler{GuestService: guestServiceMock}

	guestServiceMock.On("FindGuestsByBookingID", mock.Anything, uint(123)).
		Return([]entities.Guest{}, nil)

	// Act
	resp, err := h.GetGuestsByBookingID(context.Background(), &generated.GetGuestsByBookingIDRequest{
		BookingId: 123,
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Empty(t, resp.Guests)
	assert.Len(t, resp.Links, 1)

	guestServiceMock.AssertExpectations(t)
}
