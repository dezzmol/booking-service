package house_keeping_request

import (
	"context"

	"booking-service/internal/entities"

	"github.com/jmoiron/sqlx"
)

//go:generate mockery --disable-version-string --case=underscore --name=HouseKeepingRequestRepo --structname=HouseKeepingRequestRepoMock
type HouseKeepingRequestRepo interface {
	Save(ctx context.Context, request *entities.HousekeepingRequest) error
	FindByID(ctx context.Context, requestId uint64) (*entities.HousekeepingRequest, error)
	Update(ctx context.Context, request entities.HousekeepingRequest) error
}

type HouseKeepingRequestRepository struct {
	db *sqlx.DB
}

func NewHouseKeepingRequestRepository(db *sqlx.DB) *HouseKeepingRequestRepository {
	return &HouseKeepingRequestRepository{db: db}
}

func (r *HouseKeepingRequestRepository) Save(ctx context.Context, request *entities.HousekeepingRequest) error {
	query := `
		INSERT INTO house_keeping_requests (room_id, request_time, status)
		VALUES ($1, $2, $3)
		RETURNING created_at, updated_at
	`

	return r.db.QueryRowContext(ctx, query, request.RoomID, request.RequestTime, request.Status).
		Scan(&request.ID, &request.CreatedAt, &request.UpdatedAt)
}

func (r *HouseKeepingRequestRepository) FindByID(ctx context.Context, requestId uint64) (*entities.HousekeepingRequest, error) {
	var request entities.HousekeepingRequest
	query := `
		SELECT id, room_id, request_time, status, created_at, updated_at 
		FROM house_keeping_requests 
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &request, query, requestId)
	return &request, err
}

func (r *HouseKeepingRequestRepository) Update(ctx context.Context, request entities.HousekeepingRequest) error {
	query := `
		UPDATE house_keeping_requests
		SET room_id = $1, status = $2, request_time = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query, request.RoomID, request.Status, request.RequestTime, request.ID)
	return err
}
