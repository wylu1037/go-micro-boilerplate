package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/wylu1037/go-micro-boilerplate/pkg/db"
	"github.com/wylu1037/go-micro-boilerplate/services/booking/internal/model"
)

type BookingRepository interface {
	Create(ctx context.Context, booking *model.Booking) error
	GetByID(ctx context.Context, id string) (*model.Booking, error)
	UpdateStatus(ctx context.Context, id string, status model.BookingStatus) error
	List(ctx context.Context, page, pageSize int, userID string, status *model.BookingStatus) ([]*model.Booking, int64, error)
}

type bookingRepository struct {
	db *db.Pool
}

func NewBookingRepository(db *db.Pool) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *model.Booking) error {
	query := `
		INSERT INTO booking.bookings (user_id, show_id, session_id, seat_area_id, quantity, total_price, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		booking.UserID,
		booking.ShowID,
		booking.SessionID,
		booking.SeatAreaID,
		booking.Quantity,
		booking.TotalPrice,
		booking.Status,
	).Scan(&booking.ID, &booking.CreatedAt, &booking.UpdatedAt)

	return err
}

func (r *bookingRepository) GetByID(ctx context.Context, id string) (*model.Booking, error) {
	query := `
		SELECT id, user_id, show_id, session_id, seat_area_id, quantity, total_price, status, created_at, updated_at
		FROM booking.bookings
		WHERE id = $1
	`

	booking := &model.Booking{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.ShowID,
		&booking.SessionID,
		&booking.SeatAreaID,
		&booking.Quantity,
		&booking.TotalPrice,
		&booking.Status,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil // Return nil, nil consistent with previous gorm implementation or define ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id string, status model.BookingStatus) error {
	query := `
		UPDATE booking.bookings
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

func (r *bookingRepository) List(ctx context.Context, page, pageSize int, userID string, status *model.BookingStatus) ([]*model.Booking, int64, error) {
	offset := (page - 1) * pageSize

	baseQuery := `FROM booking.bookings WHERE 1=1`
	args := []interface{}{}
	argCount := 0

	if userID != "" {
		argCount++
		baseQuery += ` AND user_id = $` + string(rune('0'+argCount))
		args = append(args, userID)
	}

	if status != nil {
		argCount++
		baseQuery += ` AND status = $` + string(rune('0'+argCount))
		args = append(args, *status)
	}

	// Count total
	var total int64
	countQuery := `SELECT COUNT(*) ` + baseQuery
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// List items
	listQuery := `
		SELECT id, user_id, show_id, session_id, seat_area_id, quantity, total_price, status, created_at, updated_at
		` + baseQuery + `
		ORDER BY created_at DESC
		LIMIT $` + string(rune('0'+argCount+1)) + ` OFFSET $` + string(rune('0'+argCount+2))

	args = append(args, pageSize, offset)

	rows, err := r.db.Query(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	bookings := []*model.Booking{}
	for rows.Next() {
		b := &model.Booking{}
		err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.ShowID,
			&b.SessionID,
			&b.SeatAreaID,
			&b.Quantity,
			&b.TotalPrice,
			&b.Status,
			&b.CreatedAt,
			&b.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		bookings = append(bookings, b)
	}

	return bookings, total, nil
}
