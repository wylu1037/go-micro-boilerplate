package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/wylu1037/go-micro-boilerplate/pkg/db"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/model"
)

type SessionRepository interface {
	Create(ctx context.Context, session *model.Session) error
	GetByID(ctx context.Context, id string) (*model.Session, error)
	ListByShowID(ctx context.Context, showID string) ([]*model.Session, error)
}

type sessionRepository struct {
	db *db.Pool
}

func NewSessionRepository(db *db.Pool) SessionRepository {
	return &sessionRepository{db: db}
}

func (repo *sessionRepository) Create(ctx context.Context, session *model.Session) error {
	query := `
		INSERT INTO catalog.sessions (show_id, venue_id, start_time, end_time, sale_start_time, sale_end_time, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at
	`

	return repo.db.QueryRow(ctx, query,
		session.ShowID,
		session.VenueID,
		session.StartTime,
		session.EndTime,
		session.SaleStartTime,
		session.SaleEndTime,
		session.Status,
	).Scan(&session.ID, &session.CreatedAt)
}

func (repo *sessionRepository) GetByID(ctx context.Context, id string) (*model.Session, error) {
	query := `
		SELECT s.id, s.show_id, s.venue_id, s.start_time, s.end_time, s.sale_start_time, s.sale_end_time, s.status, s.created_at,
		       v.id, v.name, v.city, v.address, v.capacity, v.created_at
		FROM catalog.sessions s
		JOIN catalog.venues v ON s.venue_id = v.id
		WHERE s.id = $1
	`

	session := &model.Session{Venue: &model.Venue{}}
	err := repo.db.QueryRow(ctx, query, id).Scan(
		&session.ID,
		&session.ShowID,
		&session.VenueID,
		&session.StartTime,
		&session.EndTime,
		&session.SaleStartTime,
		&session.SaleEndTime,
		&session.Status,
		&session.CreatedAt,
		&session.Venue.ID,
		&session.Venue.Name,
		&session.Venue.City,
		&session.Venue.Address,
		&session.Venue.Capacity,
		&session.Venue.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (repo *sessionRepository) ListByShowID(ctx context.Context, showID string) ([]*model.Session, error) {
	query := `
		SELECT s.id, s.show_id, s.venue_id, s.start_time, s.end_time, s.sale_start_time, s.sale_end_time, s.status, s.created_at,
		       v.id, v.name, v.city, v.address, v.capacity, v.created_at
		FROM catalog.sessions s
		JOIN catalog.venues v ON s.venue_id = v.id
		WHERE s.show_id = $1
		ORDER BY s.start_time
	`

	rows, err := repo.db.Query(ctx, query, showID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*model.Session
	for rows.Next() {
		session := &model.Session{Venue: &model.Venue{}}
		if err := rows.Scan(
			&session.ID,
			&session.ShowID,
			&session.VenueID,
			&session.StartTime,
			&session.EndTime,
			&session.SaleStartTime,
			&session.SaleEndTime,
			&session.Status,
			&session.CreatedAt,
			&session.Venue.ID,
			&session.Venue.Name,
			&session.Venue.City,
			&session.Venue.Address,
			&session.Venue.Capacity,
			&session.Venue.CreatedAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, nil
}
