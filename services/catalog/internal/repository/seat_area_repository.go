package repository

import (
	"context"
	stderrors "errors"

	"github.com/jackc/pgx/v5"

	"github.com/wylu1037/go-micro-boilerplate/pkg/db"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/errors"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/model"
)

type SeatAreaRepository interface {
	Create(ctx context.Context, seatArea *model.SeatArea) error
	GetByID(ctx context.Context, id string) (*model.SeatArea, error)
	ListBySessionID(ctx context.Context, sessionID string) ([]*model.SeatArea, error)
	UpdateAvailableSeats(ctx context.Context, id string, delta int32) error
	CheckAndReserve(ctx context.Context, id string, quantity int32) (bool, error)
	ReleaseSeats(ctx context.Context, id string, quantity int32) error
}

type seatAreaRepository struct {
	db *db.Pool
}

func NewSeatAreaRepository(db *db.Pool) SeatAreaRepository {
	return &seatAreaRepository{db: db}
}

func (repo *seatAreaRepository) Create(ctx context.Context, seatArea *model.SeatArea) error {
	query := `
		INSERT INTO catalog.seat_areas (session_id, name, price, total_seats, available_seats)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id, created_at
	`

	return repo.db.QueryRow(ctx, query,
		seatArea.SessionID,
		seatArea.Name,
		seatArea.Price,
		seatArea.TotalSeats,
	).Scan(&seatArea.ID, &seatArea.CreatedAt)
}

func (repo *seatAreaRepository) GetByID(ctx context.Context, id string) (*model.SeatArea, error) {
	query := `
		SELECT id, session_id, name, price, total_seats, available_seats, created_at
		FROM catalog.seat_areas
		WHERE id = $1
	`

	seatArea := &model.SeatArea{}
	err := repo.db.QueryRow(ctx, query, id).Scan(
		&seatArea.ID,
		&seatArea.SessionID,
		&seatArea.Name,
		&seatArea.Price,
		&seatArea.TotalSeats,
		&seatArea.AvailableSeats,
		&seatArea.CreatedAt,
	)

	if stderrors.Is(err, pgx.ErrNoRows) {
		return nil, errors.ErrSeatAreaNotFound
	}
	if err != nil {
		return nil, err
	}

	return seatArea, nil
}

func (repo *seatAreaRepository) ListBySessionID(ctx context.Context, sessionID string) ([]*model.SeatArea, error) {
	query := `
		SELECT id, session_id, name, price, total_seats, available_seats, created_at
		FROM catalog.seat_areas
		WHERE session_id = $1
		ORDER BY price DESC
	`

	rows, err := repo.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seatAreas []*model.SeatArea
	for rows.Next() {
		seatArea := &model.SeatArea{}
		if err := rows.Scan(
			&seatArea.ID,
			&seatArea.SessionID,
			&seatArea.Name,
			&seatArea.Price,
			&seatArea.TotalSeats,
			&seatArea.AvailableSeats,
			&seatArea.CreatedAt,
		); err != nil {
			return nil, err
		}
		seatAreas = append(seatAreas, seatArea)
	}

	return seatAreas, nil
}

func (repo *seatAreaRepository) UpdateAvailableSeats(ctx context.Context, id string, delta int32) error {
	query := `
		UPDATE catalog.seat_areas
		SET available_seats = available_seats + $1
		WHERE id = $2 AND available_seats + $1 >= 0 AND available_seats + $1 <= total_seats
	`

	result, err := repo.db.Exec(ctx, query, delta, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.ErrInsufficientSeats
	}

	return nil
}

func (repo *seatAreaRepository) CheckAndReserve(ctx context.Context, id string, quantity int32) (bool, error) {
	// Use optimistic locking with CAS operation
	query := `
		UPDATE catalog.seat_areas
		SET available_seats = available_seats - $1
		WHERE id = $2 AND available_seats >= $1
		RETURNING available_seats
	`

	var newAvailable int32
	err := repo.db.QueryRow(ctx, query, quantity, id).Scan(&newAvailable)

	if stderrors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo *seatAreaRepository) ReleaseSeats(ctx context.Context, id string, quantity int32) error {
	query := `
		UPDATE catalog.seat_areas
		SET available_seats = available_seats + $1
		WHERE id = $2 AND available_seats + $1 <= total_seats
	`

	result, err := repo.db.Exec(ctx, query, quantity, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.ErrSeatAreaNotFound
	}

	return nil
}
