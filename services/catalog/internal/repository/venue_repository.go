package repository

import (
	"context"
	stderrors "errors"

	"github.com/jackc/pgx/v5"

	"github.com/wylu1037/go-micro-boilerplate/pkg/db"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/errors"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/model"
)

type VenueRepository interface {
	Create(ctx context.Context, venue *model.Venue) error
	GetByID(ctx context.Context, id string) (*model.Venue, error)
	List(ctx context.Context, city *string, offset, limit int) ([]*model.Venue, int64, error)
}

type venueRepository struct {
	db *db.Pool
}

func NewVenueRepository(db *db.Pool) VenueRepository {
	return &venueRepository{db: db}
}

func (repo *venueRepository) Create(ctx context.Context, venue *model.Venue) error {
	query := `
		INSERT INTO catalog.venues (name, city, address, capacity)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return repo.db.QueryRow(ctx, query,
		venue.Name,
		venue.City,
		venue.Address,
		venue.Capacity,
	).Scan(&venue.ID, &venue.CreatedAt)
}

func (repo *venueRepository) GetByID(ctx context.Context, id string) (*model.Venue, error) {
	query := `
		SELECT id, name, city, address, capacity, created_at
		FROM catalog.venues
		WHERE id = $1
	`

	venue := &model.Venue{}
	err := repo.db.QueryRow(ctx, query, id).Scan(
		&venue.ID,
		&venue.Name,
		&venue.City,
		&venue.Address,
		&venue.Capacity,
		&venue.CreatedAt,
	)

	if stderrors.Is(err, pgx.ErrNoRows) {
		return nil, errors.ErrVenueNotFound
	}
	if err != nil {
		return nil, err
	}

	return venue, nil
}

func (repo *venueRepository) List(ctx context.Context, city *string, offset, limit int) ([]*model.Venue, int64, error) {
	query := `SELECT id, name, city, address, capacity, created_at FROM catalog.venues WHERE 1=1`
	countQuery := `SELECT COUNT(*) FROM catalog.venues WHERE 1=1`

	args := []any{}
	argIndex := 1

	if city != nil {
		query += ` AND city = $1`
		countQuery += ` AND city = $1`
		args = append(args, *city)
		argIndex++
	}

	var total int64
	countArgs := make([]any, len(args))
	copy(countArgs, args)
	if err := repo.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query += ` ORDER BY created_at DESC LIMIT $` + string(rune('0'+argIndex)) + ` OFFSET $` + string(rune('0'+argIndex+1))
	args = append(args, limit, offset)

	rows, err := repo.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var venues []*model.Venue
	for rows.Next() {
		venue := &model.Venue{}
		if err := rows.Scan(
			&venue.ID,
			&venue.Name,
			&venue.City,
			&venue.Address,
			&venue.Capacity,
			&venue.CreatedAt,
		); err != nil {
			return nil, 0, err
		}
		venues = append(venues, venue)
	}

	return venues, total, nil
}
