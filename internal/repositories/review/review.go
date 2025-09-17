package review

import (
	"context"

	"booking-service/internal/entities"

	"github.com/jmoiron/sqlx"
)

//go:generate mockery --disable-version-string --case=underscore --name=ReviewRepo --structname=ReviewRepoMock
type ReviewRepo interface {
	// Save сохраняет новый отзыв в базу данных.
	// После успешного сохранения в объект review будут установлены сгенерированные поля id, created_at и updated_at.
	Save(ctx context.Context, review *entities.Review) error

	// FindByGuest возвращает список отзывов, оставленных гостем с указанным guestID.
	FindByGuest(ctx context.Context, guestID uint64) ([]entities.Review, error)
}

type ReviewRepository struct {
	db *sqlx.DB
}

func NewReviewRepository(db *sqlx.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Save(ctx context.Context, review *entities.Review) error {
	query := `
        INSERT INTO reviews (booking_id, guest_id, rating, comment)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at
    `
	return r.db.QueryRowxContext(ctx, query, review.BookingID, review.GuestID, review.Rating, review.Comment).
		Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)
}

func (r *ReviewRepository) FindByGuest(ctx context.Context, guestID uint64) ([]entities.Review, error) {
	var reviews []entities.Review
	query := `
        SELECT id, booking_id, guest_id, rating, comment, date, created_at, updated_at
        FROM reviews
        WHERE guest_id = $1
        ORDER BY date DESC
    `
	err := r.db.SelectContext(ctx, &reviews, query, guestID)
	return reviews, err
}
