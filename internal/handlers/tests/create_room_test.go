package tests

import (
	"context"
	"errors"
	"testing"

	"booking-service/internal/controllers/room/mocks"
	"booking-service/internal/generated"
	"booking-service/internal/handlers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateRoom_ServiceError(t *testing.T) {
	// Arrange
	roomServiceMock := &mocks.RoomServiceMock{}
	h := &handlers.Handler{RoomService: roomServiceMock}

	expectedErr := errors.New("database error")
	roomServiceMock.On("CreateRooms", mock.Anything, mock.AnythingOfType("[]entities.RoomDTO")).
		Return(expectedErr)

	// Act
	resp, err := h.CreateRoom(context.Background(), &generated.CreateRoomRequest{
		Rooms: []*generated.RoomDTO{
			{
				Number:  "101",
				Type:    "standard",
				HotelId: 1,
			},
		},
	})

	// Assert
	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	roomServiceMock.AssertExpectations(t)
}
