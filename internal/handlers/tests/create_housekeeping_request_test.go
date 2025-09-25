package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"booking-service/internal/controllers/house_keeping_request/mocks"
	"booking-service/internal/entities"
	"booking-service/internal/generated"
	"booking-service/internal/handlers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateHouseKeepingRequest_Success(t *testing.T) {
	// Arrange
	serviceMock := &mocks.HouseKeepingRequestServiceMock{}
	h := &handlers.Handler{HousekeepingService: serviceMock}

	userID := uint64(123)
	roomID := uint64(456)
	requestTime := time.Now().UTC()
	expectedRequest := entities.HousekeepingRequest{
		ID:          userID,
		RoomID:      roomID,
		RequestTime: requestTime,
		Status:      "pending",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	serviceMock.On("CreateHouseKeepingRequest", mock.Anything, entities.HouseKeepingRequestDTO{
		RoomID:      roomID,
		RequestTime: requestTime,
	}).Return(expectedRequest, nil)

	// Act
	resp, err := h.CreateHouseKeepingRequest(context.Background(), &generated.CreateHouseKeepingRequestRequest{
		RoomId:      roomID,
		RequestTime: timestamppb.New(requestTime),
	})

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Request)
	assert.Equal(t, userID, resp.Request.Id)
	assert.Equal(t, roomID, resp.Request.RoomId)
	assert.Equal(t, "pending", resp.Request.Status)
	assert.True(t, requestTime.Equal(resp.Request.RequestTime.AsTime()))
	assert.Len(t, resp.Links, 1)
	assert.Equal(t, "self", resp.Links[0].Rel)

	serviceMock.AssertExpectations(t)
}

func TestCreateHouseKeepingRequest_RoomNotFound(t *testing.T) {
	// Arrange
	serviceMock := &mocks.HouseKeepingRequestServiceMock{}
	h := &handlers.Handler{HousekeepingService: serviceMock}

	roomID := uint64(456)
	requestTime := time.Now().UTC()
	serviceMock.On("CreateHouseKeepingRequest", mock.Anything, entities.HouseKeepingRequestDTO{
		RoomID:      roomID,
		RequestTime: requestTime,
	}).Return(entities.HousekeepingRequest{}, entities.ErrNotFound)

	// Act
	resp, err := h.CreateHouseKeepingRequest(context.Background(), &generated.CreateHouseKeepingRequestRequest{
		RoomId:      roomID,
		RequestTime: timestamppb.New(requestTime),
	})

	// Assert
	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Equal(t, "room not found", st.Message())

	serviceMock.AssertExpectations(t)
}

func TestCreateHouseKeepingRequest_InternalError(t *testing.T) {
	// Arrange
	serviceMock := &mocks.HouseKeepingRequestServiceMock{}
	h := &handlers.Handler{HousekeepingService: serviceMock}

	roomID := uint64(456)
	requestTime := time.Now().UTC()
	expectedErr := errors.New("database error")
	serviceMock.On("CreateHouseKeepingRequest", mock.Anything, entities.HouseKeepingRequestDTO{
		RoomID:      roomID,
		RequestTime: requestTime,
	}).Return(entities.HousekeepingRequest{}, expectedErr)

	// Act
	resp, err := h.CreateHouseKeepingRequest(context.Background(), &generated.CreateHouseKeepingRequestRequest{
		RoomId:      roomID,
		RequestTime: timestamppb.New(requestTime),
	})

	// Assert
	assert.Nil(t, resp)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
	assert.Contains(t, st.Message(), "failed to create the housekeeping request")

	serviceMock.AssertExpectations(t)
}
