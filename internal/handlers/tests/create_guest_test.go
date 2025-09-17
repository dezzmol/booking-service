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

func TestCreateGuest_Success(t *testing.T) {
	// Arrange
	guestServiceMock := &mocks.GuestServiceMock{}
	h := &handlers.Handler{GuestService: guestServiceMock}

	now := time.Now()
	expectedGuest := entities.Guest{
		ID:        uint64(123),
		Name:      "John Doe",
		CreatedAt: now,
		UpdatedAt: now,
	}

	guestServiceMock.On("CreateGuest", mock.Anything, entities.GuestDTO{
		Name: "John Doe",
	}).Return(expectedGuest, nil)

	// Act
	resp, err := h.CreateGuest(context.Background(), &generated.CreateGuestRequest{
		Name: "John Doe",
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Guest)
	assert.Equal(t, uint64(123), resp.Guest.Id)
	assert.Equal(t, "John Doe", resp.Guest.Name)
	assert.True(t, timestamppb.New(now).AsTime().Equal(resp.Guest.CreatedAt.AsTime()))
	assert.True(t, timestamppb.New(now).AsTime().Equal(resp.Guest.UpdatedAt.AsTime()))
	assert.Len(t, resp.Links, 1)
	assert.Equal(t, "self", resp.Links[0].Rel)
	assert.Equal(t, "/v1/guests", resp.Links[0].Href)

	guestServiceMock.AssertExpectations(t)
}

func TestCreateGuest_InternalError(t *testing.T) {
	// Arrange
	guestServiceMock := &mocks.GuestServiceMock{}
	h := &handlers.Handler{GuestService: guestServiceMock}

	expectedErr := errors.New("database error")
	guestServiceMock.On("CreateGuest", mock.Anything, entities.GuestDTO{
		Name: "John Doe",
	}).Return(entities.Guest{}, expectedErr)

	// Act
	resp, err := h.CreateGuest(context.Background(), &generated.CreateGuestRequest{
		Name: "John Doe",
	})

	// Assert
	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "database error")

	guestServiceMock.AssertExpectations(t)
}

func TestCreateGuest_EmptyName(t *testing.T) {
	// Arrange
	guestServiceMock := &mocks.GuestServiceMock{}
	h := &handlers.Handler{GuestService: guestServiceMock}
	guestServiceMock.On("CreateGuest", mock.Anything, entities.GuestDTO{
		Name: "",
	}).Return(entities.Guest{}, entities.ErrNameIsRequired)

	// Act
	resp, err := h.CreateGuest(context.Background(), &generated.CreateGuestRequest{
		Name: "",
	})

	// Assert
	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	// Сервис не должен вызываться при пустом имени
	guestServiceMock.AssertNotCalled(t, "CreateGuest")
}

func TestCreateGuest_ValidationError(t *testing.T) {
	// Arrange
	guestServiceMock := &mocks.GuestServiceMock{}
	longName := "This is a very long name that exceeds the maximum allowed length of characters for a guest name in our system. This is a very long name that exceeds the maximum allowed length of characters for a guest name in our system. This is a very long name that exceeds the maximum allowed length of characters for a guest name in our system"
	h := &handlers.Handler{GuestService: guestServiceMock}
	guestServiceMock.On("CreateGuest", mock.Anything, entities.GuestDTO{
		Name: longName,
	}).Return(entities.Guest{}, entities.ErrNameIsTooLong)
	// Act (имя слишком длинное)
	_, err := h.CreateGuest(context.Background(), &generated.CreateGuestRequest{
		Name: longName,
	})

	// Assert
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())

	// Сервис не должен вызываться при невалидных данных
	guestServiceMock.AssertNotCalled(t, "CreateGuest")
}
