package handler

import (
	"context"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"

	catalogv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/catalog/v1"
	commonv1 "github.com/wylu1037/go-micro-boilerplate/gen/go/common/v1"
	"github.com/wylu1037/go-micro-boilerplate/pkg/tools"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/errors"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/model"
	"github.com/wylu1037/go-micro-boilerplate/services/catalog/internal/service"
)

type CatalogHandler struct {
	svc service.CatalogService
}

func NewCatalogHandler(
	svc service.CatalogService,
) catalogv1.CatalogServiceHandler {
	return &CatalogHandler{svc: svc}
}

func (h *CatalogHandler) CreateShow(ctx context.Context, req *catalogv1.CreateShowRequest, rsp *catalogv1.CreateShowResponse) error {
	show := &model.Show{
		Title:       req.Title,
		Description: req.Description,
		Artist:      req.Artist,
		Category:    req.Category.String(), // Store enum as string
		PosterURL:   req.PosterUrl,
		Status:      catalogv1.ShowStatus_SHOW_STATUS_DRAFT.String(),
	}

	if err := h.svc.CreateShow(ctx, show); err != nil {
		return errors.ToMicroError(err)
	}

	rsp.Show = h.convertShow(show)
	return nil
}

func (h *CatalogHandler) GetShow(ctx context.Context, req *catalogv1.GetShowRequest, rsp *catalogv1.GetShowResponse) error {
	show, err := h.svc.GetShow(ctx, req.ShowId)
	if err != nil {
		return errors.ToMicroError(err)
	}

	rsp.Show = h.convertShow(show)
	return nil
}

func (h *CatalogHandler) ListShows(ctx context.Context, req *catalogv1.ListShowsRequest, rsp *catalogv1.ListShowsResponse) error {
	var category, status, city *string
	if req.GetCategory() != catalogv1.ShowCategory_SHOW_CATEGORY_UNSPECIFIED {
		c := req.GetCategory().String()
		category = &c
	}
	if req.GetStatus() != catalogv1.ShowStatus_SHOW_STATUS_UNSPECIFIED {
		s := req.GetStatus().String()
		status = &s
	}
	if req.City != nil {
		city = req.City
	}

	page := lo.Ternary(req.Page < 1, 1, req.Page)
	pageSize := lo.Ternary(req.PageSize < 1, 10, req.PageSize)
	offset := int((page - 1) * pageSize)
	limit := int(pageSize)

	shows, total, err := h.svc.ListShows(ctx, category, status, city, offset, limit)
	if err != nil {
		return errors.ToMicroError(err)
	}

	rsp.Shows = make([]*catalogv1.Show, len(shows))
	for i, s := range shows {
		rsp.Shows[i] = h.convertShow(s)
	}

	rsp.Pagination = &commonv1.PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int32((total + int64(limit) - 1) / int64(limit)),
		TotalCount: total,
	}

	return nil
}

func (h *CatalogHandler) UpdateShow(ctx context.Context, req *catalogv1.UpdateShowRequest, rsp *catalogv1.UpdateShowResponse) error {
	show := &model.Show{
		ID: req.ShowId,
	}

	if req.Title != nil {
		show.Title = *req.Title
	}
	if req.Description != nil {
		show.Description = *req.Description
	}
	if req.Artist != nil {
		show.Artist = *req.Artist
	}
	if req.Category != nil {
		show.Category = req.Category.String()
	}
	if req.PosterUrl != nil {
		show.PosterURL = *req.PosterUrl
	}
	if req.Status != nil {
		show.Status = req.Status.String()
	}

	if err := h.svc.UpdateShow(ctx, show); err != nil {
		return errors.ToMicroError(err)
	}

	rsp.Show = h.convertShow(show)
	return nil
}

func (h *CatalogHandler) DeleteShow(ctx context.Context, req *catalogv1.DeleteShowRequest, rsp *catalogv1.DeleteShowResponse) error {
	if err := h.svc.DeleteShow(ctx, req.ShowId); err != nil {
		return errors.ToMicroError(err)
	}
	rsp.Message = "Show deleted successfully"
	return nil
}

func (h *CatalogHandler) convertShow(s *model.Show) *catalogv1.Show {
	return &catalogv1.Show{
		ShowId:      s.ID,
		Title:       s.Title,
		Description: s.Description,
		Artist:      s.Artist,
		Category:    catalogv1.ShowCategory(catalogv1.ShowCategory_value[s.Category]),
		PosterUrl:   s.PosterURL,
		Status:      catalogv1.ShowStatus(catalogv1.ShowStatus_value[s.Status]),
		CreatedAt:   tools.ToProtoTimestamp(s.CreatedAt),
		UpdatedAt:   tools.ToProtoTimestamp(s.UpdatedAt),
	}
}

// Venue

func (h *CatalogHandler) CreateVenue(ctx context.Context, req *catalogv1.CreateVenueRequest, rsp *catalogv1.CreateVenueResponse) error {
	venue := &model.Venue{
		Name:     req.Name,
		City:     req.City,
		Address:  req.Address,
		Capacity: req.Capacity,
	}

	if err := h.svc.CreateVenue(ctx, venue); err != nil {
		return errors.ToMicroError(err)
	}

	rsp.Venue = h.convertVenue(venue)
	return nil
}

func (h *CatalogHandler) GetVenue(ctx context.Context, req *catalogv1.GetVenueRequest, rsp *catalogv1.GetVenueResponse) error {
	venue, err := h.svc.GetVenue(ctx, req.VenueId)
	if err != nil {
		return errors.ToMicroError(err)
	}

	rsp.Venue = h.convertVenue(venue)
	return nil
}

func (h *CatalogHandler) ListVenues(ctx context.Context, req *catalogv1.ListVenuesRequest, rsp *catalogv1.ListVenuesResponse) error {
	page := lo.Ternary(req.Page < 1, 1, req.Page)
	pageSize := lo.Ternary(req.PageSize < 1, 10, req.PageSize)
	offset := int((page - 1) * pageSize)
	limit := int(pageSize)

	venues, total, err := h.svc.ListVenues(ctx, req.City, offset, limit)
	if err != nil {
		return errors.ToMicroError(err)
	}

	rsp.Venues = make([]*catalogv1.Venue, len(venues))
	for i, v := range venues {
		rsp.Venues[i] = h.convertVenue(v)
	}

	rsp.Pagination = &commonv1.PaginationResponse{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int32((total + int64(limit) - 1) / int64(limit)),
		TotalCount: total,
	}

	return nil
}

func (h *CatalogHandler) convertVenue(v *model.Venue) *catalogv1.Venue {
	return &catalogv1.Venue{
		VenueId:   v.ID,
		Name:      v.Name,
		City:      v.City,
		Address:   v.Address,
		Capacity:  v.Capacity,
		CreatedAt: tools.ToProtoTimestamp(v.CreatedAt),
	}
}

func (h *CatalogHandler) CreateSession(ctx context.Context, req *catalogv1.CreateSessionRequest, rsp *catalogv1.CreateSessionResponse) error {
	session := &model.Session{
		ShowID:        req.ShowId,
		VenueID:       req.VenueId,
		StartTime:     tools.ToTime(req.StartTime),
		EndTime:       tools.ToTimePtr(req.EndTime),
		SaleStartTime: tools.ToTimePtr(req.SaleStartTime),
		SaleEndTime:   tools.ToTimePtr(req.SaleEndTime),
		Status:        catalogv1.SessionStatus_SESSION_STATUS_SCHEDULED.String(),
	}

	if err := h.svc.CreateSession(ctx, session); err != nil {
		return errors.ToMicroError(err)
	}

	rsp.Session = h.convertSession(session)
	return nil
}

func (h *CatalogHandler) GetSession(ctx context.Context, req *catalogv1.GetSessionRequest, rsp *catalogv1.GetSessionResponse) error {
	session, err := h.svc.GetSession(ctx, req.SessionId)
	if err != nil {
		return errors.ToMicroError(err)
	}

	rsp.Session = h.convertSession(session)

	seatAreas, err := h.svc.ListSeatAreas(ctx, req.SessionId)
	if err == nil {
		rsp.SeatAreas = make([]*catalogv1.SeatArea, len(seatAreas))
		for i, sa := range seatAreas {
			rsp.SeatAreas[i] = h.convertSeatArea(sa)
		}
	}

	return nil
}

func (h *CatalogHandler) ListSessions(ctx context.Context, req *catalogv1.ListSessionsRequest, rsp *catalogv1.ListSessionsResponse) error {
	sessions, err := h.svc.ListSessions(ctx, req.ShowId)
	if err != nil {
		return errors.ToMicroError(err)
	}

	rsp.Sessions = make([]*catalogv1.Session, len(sessions))
	for i, s := range sessions {
		rsp.Sessions[i] = h.convertSession(s)
	}

	return nil
}

func (h *CatalogHandler) convertSession(s *model.Session) *catalogv1.Session {
	var venue *catalogv1.Venue
	if s.Venue != nil {
		venue = h.convertVenue(s.Venue)
	}

	return &catalogv1.Session{
		SessionId:     s.ID,
		ShowId:        s.ShowID,
		VenueId:       s.VenueID,
		Venue:         venue,
		StartTime:     tools.ToProtoTimestamp(s.StartTime),
		EndTime:       tools.ToProtoTimestampPtr(s.EndTime),
		SaleStartTime: tools.ToProtoTimestampPtr(s.SaleStartTime),
		SaleEndTime:   tools.ToProtoTimestampPtr(s.SaleEndTime),
		Status:        catalogv1.SessionStatus(catalogv1.SessionStatus_value[s.Status]),
		CreatedAt:     tools.ToProtoTimestamp(s.CreatedAt),
	}
}

func (h *CatalogHandler) CreateSeatArea(ctx context.Context, req *catalogv1.CreateSeatAreaRequest, rsp *catalogv1.CreateSeatAreaResponse) error {
	if req.Price == "" {
		req.Price = "0"
	}
	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		return errors.ToMicroError(errors.ErrInvalidPrice)
	}

	seatArea := &model.SeatArea{
		SessionID:  req.SessionId,
		Name:       req.Name,
		Price:      price,
		TotalSeats: req.TotalSeats,
	}

	if err := h.svc.CreateSeatArea(ctx, seatArea); err != nil {
		return errors.ToMicroError(err)
	}

	rsp.SeatArea = h.convertSeatArea(seatArea)
	return nil
}

func (h *CatalogHandler) ListSeatAreas(ctx context.Context, req *catalogv1.ListSeatAreasRequest, rsp *catalogv1.ListSeatAreasResponse) error {
	seatAreas, err := h.svc.ListSeatAreas(ctx, req.SessionId)
	if err != nil {
		return errors.ToMicroError(err)
	}

	rsp.SeatAreas = make([]*catalogv1.SeatArea, len(seatAreas))
	for i, sa := range seatAreas {
		rsp.SeatAreas[i] = h.convertSeatArea(sa)
	}

	return nil
}

func (h *CatalogHandler) convertSeatArea(sa *model.SeatArea) *catalogv1.SeatArea {
	return &catalogv1.SeatArea{
		SeatAreaId:     sa.ID,
		SessionId:      sa.SessionID,
		Name:           sa.Name,
		Price:          sa.Price.String(),
		TotalSeats:     sa.TotalSeats,
		AvailableSeats: sa.AvailableSeats,
		CreatedAt:      tools.ToProtoTimestamp(sa.CreatedAt),
	}
}

func (h *CatalogHandler) CheckAvailability(ctx context.Context, req *catalogv1.CheckAvailabilityRequest, rsp *catalogv1.CheckAvailabilityResponse) error {
	available, count, price, err := h.svc.CheckAvailability(ctx, req.SessionId, req.SeatAreaId, req.Quantity)
	if err != nil {
		return errors.ToMicroError(err)
	}

	rsp.Available = available
	rsp.AvailableSeats = count
	rsp.Price = price.String()
	return nil
}

func (h *CatalogHandler) ReserveSeats(ctx context.Context, req *catalogv1.ReserveSeatsRequest, rsp *catalogv1.ReserveSeatsResponse) error {
	err := h.svc.ReserveSeats(ctx, req.SessionId, req.SeatAreaId, req.Quantity, req.OrderId)
	if err != nil {
		return errors.ToMicroError(err)
	}

	rsp.Success = true
	rsp.Message = "Seats reserved"
	return nil
}

func (h *CatalogHandler) ReleaseSeats(ctx context.Context, req *catalogv1.ReleaseSeatsRequest, rsp *catalogv1.ReleaseSeatsResponse) error {
	err := h.svc.ReleaseSeats(ctx, req.SessionId, req.SeatAreaId, req.Quantity, req.OrderId)
	if err != nil {
		return errors.ToMicroError(err)
	}

	rsp.Success = true
	rsp.Message = "Seats released"
	return nil
}
