package repository

import (
	"context"
	stderrors "errors"

	"github.com/jackc/pgx/v5"

	"github.com/wylu1037/go-micro-boilerplate/pkg/db"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/errors"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/model"
)

type ShowRepository interface {
	Create(ctx context.Context, show *model.Show) error
	GetByID(ctx context.Context, id string) (*model.Show, error)
	List(ctx context.Context, category, status, city *string, offset, limit int) ([]*model.Show, int64, error)
	Update(ctx context.Context, show *model.Show) error
	Delete(ctx context.Context, id string) error
}

type showRepository struct {
	db *db.Pool
}

func NewShowRepository(db *db.Pool) ShowRepository {
	return &showRepository{db: db}
}

func (repo *showRepository) Create(ctx context.Context, show *model.Show) error {
	query := `
		INSERT INTO catalog.shows (title, description, artist, category, poster_url, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := repo.db.QueryRow(ctx, query,
		show.Title,
		show.Description,
		show.Artist,
		show.Category,
		show.PosterURL,
		show.Status,
	).Scan(&show.ID, &show.CreatedAt, &show.UpdatedAt)

	return err
}

func (repo *showRepository) GetByID(ctx context.Context, id string) (*model.Show, error) {
	query := `
		SELECT id, title, description, artist, category, poster_url, status, created_at, updated_at
		FROM catalog.shows
		WHERE id = $1
	`

	show := &model.Show{}
	err := repo.db.QueryRow(ctx, query, id).Scan(
		&show.ID,
		&show.Title,
		&show.Description,
		&show.Artist,
		&show.Category,
		&show.PosterURL,
		&show.Status,
		&show.CreatedAt,
		&show.UpdatedAt,
	)

	if stderrors.Is(err, pgx.ErrNoRows) {
		return nil, errors.ErrShowNotFound
	}
	if err != nil {
		return nil, err
	}

	return show, nil
}

func (repo *showRepository) List(ctx context.Context, category, status, city *string, offset, limit int) ([]*model.Show, int64, error) {
	// Build dynamic query
	query := `
		SELECT DISTINCT s.id, s.title, s.description, s.artist, s.category, s.poster_url, s.status, s.created_at, s.updated_at
		FROM catalog.shows s
		LEFT JOIN catalog.sessions se ON s.id = se.show_id
		LEFT JOIN catalog.venues v ON se.venue_id = v.id
		WHERE 1=1
	`
	countQuery := `
		SELECT COUNT(DISTINCT s.id)
		FROM catalog.shows s
		LEFT JOIN catalog.sessions se ON s.id = se.show_id
		LEFT JOIN catalog.venues v ON se.venue_id = v.id
		WHERE 1=1
	`

	args := []any{}
	argIndex := 1

	if category != nil {
		query += ` AND s.category = $` + string(rune('0'+argIndex))
		countQuery += ` AND s.category = $` + string(rune('0'+argIndex))
		args = append(args, *category)
		argIndex++
	}

	if status != nil {
		query += ` AND s.status = $` + string(rune('0'+argIndex))
		countQuery += ` AND s.status = $` + string(rune('0'+argIndex))
		args = append(args, *status)
		argIndex++
	}

	if city != nil {
		query += ` AND v.city = $` + string(rune('0'+argIndex))
		countQuery += ` AND v.city = $` + string(rune('0'+argIndex))
		args = append(args, *city)
		argIndex++
	}

	// Get total count
	var total int64
	countArgs := make([]any, len(args))
	copy(countArgs, args)
	if err := repo.db.QueryRow(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Add pagination
	query += ` ORDER BY s.created_at DESC LIMIT $` + string(rune('0'+argIndex)) + ` OFFSET $` + string(rune('0'+argIndex+1))
	args = append(args, limit, offset)

	rows, err := repo.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var shows []*model.Show
	for rows.Next() {
		show := &model.Show{}
		if err := rows.Scan(
			&show.ID,
			&show.Title,
			&show.Description,
			&show.Artist,
			&show.Category,
			&show.PosterURL,
			&show.Status,
			&show.CreatedAt,
			&show.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}
		shows = append(shows, show)
	}

	return shows, total, nil
}

func (repo *showRepository) Update(ctx context.Context, show *model.Show) error {
	query := `
		UPDATE catalog.shows
		SET title = $1, description = $2, artist = $3, category = $4, poster_url = $5, status = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at
	`

	err := repo.db.QueryRow(ctx, query,
		show.Title,
		show.Description,
		show.Artist,
		show.Category,
		show.PosterURL,
		show.Status,
		show.ID,
	).Scan(&show.UpdatedAt)

	if stderrors.Is(err, pgx.ErrNoRows) {
		return errors.ErrShowNotFound
	}

	return err
}

func (repo *showRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM catalog.shows WHERE id = $1`

	result, err := repo.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.ErrShowNotFound
	}

	return nil
}
