package service

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/model"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/repository"
)

type CatalogService interface {
	// Show
	CreateShow(ctx context.Context, show *model.Show) error
	GetShow(ctx context.Context, id string) (*model.Show, error)
	ListShows(ctx context.Context, category, status, city *string, offset, limit int) ([]*model.Show, int64, error)
	UpdateShow(ctx context.Context, show *model.Show) error
	DeleteShow(ctx context.Context, id string) error

	// Venue
	CreateVenue(ctx context.Context, venue *model.Venue) error
	GetVenue(ctx context.Context, id string) (*model.Venue, error)
	ListVenues(ctx context.Context, city *string, offset, limit int) ([]*model.Venue, int64, error)

	// Session
	CreateSession(ctx context.Context, session *model.Session) error
	GetSession(ctx context.Context, id string) (*model.Session, error)
	ListSessions(ctx context.Context, showID string) ([]*model.Session, error)

	// SeatArea
	CreateSeatArea(ctx context.Context, seatArea *model.SeatArea) error
	ListSeatAreas(ctx context.Context, sessionID string) ([]*model.SeatArea, error)

	// Inventory
	CheckAvailability(ctx context.Context, sessionID, seatAreaID string, quantity int32) (bool, int32, decimal.Decimal, error)
	ReserveSeats(ctx context.Context, sessionID, seatAreaID string, quantity int32, orderID string) error
	ReleaseSeats(ctx context.Context, sessionID, seatAreaID string, quantity int32, orderID string) error
}

type catalogService struct {
	showRepo     repository.ShowRepository
	venueRepo    repository.VenueRepository
	sessionRepo  repository.SessionRepository
	seatAreaRepo repository.SeatAreaRepository
	logger       *zerolog.Logger
}

func NewCatalogService(
	showRepo repository.ShowRepository,
	venueRepo repository.VenueRepository,
	sessionRepo repository.SessionRepository,
	seatAreaRepo repository.SeatAreaRepository,
	logger *zerolog.Logger,
) CatalogService {
	return &catalogService{
		showRepo:     showRepo,
		venueRepo:    venueRepo,
		sessionRepo:  sessionRepo,
		seatAreaRepo: seatAreaRepo,
		logger:       logger,
	}
}

func (svc *catalogService) CreateShow(ctx context.Context, show *model.Show) error {
	if err := svc.showRepo.Create(ctx, show); err != nil {
		return err
	}
	svc.logger.Info().Str("show_id", show.ID).Msg("Show created")
	return nil
}

func (svc *catalogService) GetShow(ctx context.Context, id string) (*model.Show, error) {
	return svc.showRepo.GetByID(ctx, id)
}

func (svc *catalogService) ListShows(ctx context.Context, category, status, city *string, offset, limit int) ([]*model.Show, int64, error) {
	return svc.showRepo.List(ctx, category, status, city, offset, limit)
}

func (svc *catalogService) UpdateShow(ctx context.Context, show *model.Show) error {
	return svc.showRepo.Update(ctx, show)
}

func (svc *catalogService) DeleteShow(ctx context.Context, id string) error {
	return svc.showRepo.Delete(ctx, id)
}

// Venue

func (svc *catalogService) CreateVenue(ctx context.Context, venue *model.Venue) error {
	if err := svc.venueRepo.Create(ctx, venue); err != nil {
		return err
	}
	svc.logger.Info().Str("venue_id", venue.ID).Msg("Venue created")
	return nil
}

func (svc *catalogService) GetVenue(ctx context.Context, id string) (*model.Venue, error) {
	return svc.venueRepo.GetByID(ctx, id)
}

func (svc *catalogService) ListVenues(ctx context.Context, city *string, offset, limit int) ([]*model.Venue, int64, error) {
	return svc.venueRepo.List(ctx, city, offset, limit)
}

// Session

func (svc *catalogService) CreateSession(ctx context.Context, session *model.Session) error {
	// Verify Show and Venue exist
	if _, err := svc.showRepo.GetByID(ctx, session.ShowID); err != nil {
		return err
	}
	if _, err := svc.venueRepo.GetByID(ctx, session.VenueID); err != nil {
		return err
	}

	if err := svc.sessionRepo.Create(ctx, session); err != nil {
		return err
	}
	svc.logger.Info().Str("session_id", session.ID).Msg("Session created")
	return nil
}

func (svc *catalogService) GetSession(ctx context.Context, id string) (*model.Session, error) {
	return svc.sessionRepo.GetByID(ctx, id)
}

func (svc *catalogService) ListSessions(ctx context.Context, showID string) ([]*model.Session, error) {
	return svc.sessionRepo.ListByShowID(ctx, showID)
}

func (svc *catalogService) CreateSeatArea(ctx context.Context, seatArea *model.SeatArea) error {
	// Verify Session exists
	if _, err := svc.sessionRepo.GetByID(ctx, seatArea.SessionID); err != nil {
		return err
	}

	if err := svc.seatAreaRepo.Create(ctx, seatArea); err != nil {
		return err
	}
	svc.logger.Info().Str("seat_area_id", seatArea.ID).Msg("Seat area created")
	return nil
}

func (svc *catalogService) ListSeatAreas(ctx context.Context, sessionID string) ([]*model.SeatArea, error) {
	return svc.seatAreaRepo.ListBySessionID(ctx, sessionID)
}

func (svc *catalogService) CheckAvailability(ctx context.Context, sessionID, seatAreaID string, quantity int32) (bool, int32, decimal.Decimal, error) {
	area, err := svc.seatAreaRepo.GetByID(ctx, seatAreaID)
	if err != nil {
		return false, 0, decimal.Zero, err
	}

	if area.SessionID != sessionID {
		return false, 0, decimal.Zero, errors.New("seat area does not belong to session")
	}

	available := area.AvailableSeats >= quantity
	return available, area.AvailableSeats, area.Price, nil
}

func (svc *catalogService) ReserveSeats(ctx context.Context, sessionID, seatAreaID string, quantity int32, orderID string) error {
	success, err := svc.seatAreaRepo.CheckAndReserve(ctx, seatAreaID, quantity)
	if err != nil {
		return err
	}
	if !success {
		return model.ErrInsufficientSeats
	}

	svc.logger.Info().
		Str("session_id", sessionID).
		Str("seat_area_id", seatAreaID).
		Int32("quantity", quantity).
		Str("order_id", orderID).
		Msg("Seats reserved")

	return nil
}

func (svc *catalogService) ReleaseSeats(ctx context.Context, sessionID, seatAreaID string, quantity int32, orderID string) error {
	if err := svc.seatAreaRepo.ReleaseSeats(ctx, seatAreaID, quantity); err != nil {
		return err
	}

	svc.logger.Info().
		Str("session_id", sessionID).
		Str("seat_area_id", seatAreaID).
		Int32("quantity", quantity).
		Str("order_id", orderID).
		Msg("Seats released")

	return nil
}
