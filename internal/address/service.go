package address

import (
	"context"

	apperrors "github.com/gumla-hds/gumla-backend/pkg/errors"
	db "github.com/gumla-hds/gumla-backend/sqlc/generated"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type AddressResponse struct {
	ID          int64   `json:"id"`
	Label       string  `json:"label"`
	FullAddress string  `json:"full_address"`
	Landmark    string  `json:"landmark"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	IsDefault   bool    `json:"is_default"`
}

type CreateAddressRequest struct {
	Label       string  `json:"label" binding:"required"`
	FullAddress string  `json:"full_address" binding:"required"`
	Landmark    string  `json:"landmark"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	IsDefault   bool    `json:"is_default"`
}

type UpdateAddressRequest struct {
	Label       string  `json:"label" binding:"required"`
	FullAddress string  `json:"full_address" binding:"required"`
	Landmark    string  `json:"landmark"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	IsDefault   bool    `json:"is_default"`
}

func toAddressResponse(a db.Address) AddressResponse {
	lat, _ := a.Latitude.Float64Value()
	lng, _ := a.Longitude.Float64Value()
	return AddressResponse{
		ID:          a.ID,
		Label:       a.Label,
		FullAddress: a.FullAddress,
		Landmark:    a.Landmark,
		Latitude:    lat.Float64,
		Longitude:   lng.Float64,
		IsDefault:   a.IsDefault,
	}
}

func (s *Service) List(ctx context.Context, userID int64) ([]AddressResponse, *apperrors.AppError) {
	addrs, err := s.repo.ListAddresses(ctx, userID)
	if err != nil {
		return nil, apperrors.InternalWrap("failed to list addresses", "पते लोड करने में विफल", err)
	}
	resp := make([]AddressResponse, len(addrs))
	for i, a := range addrs {
		resp[i] = toAddressResponse(a)
	}
	return resp, nil
}

func (s *Service) GetByID(ctx context.Context, id, userID int64) (*AddressResponse, *apperrors.AppError) {
	a, err := s.repo.GetAddressByID(ctx, id, userID)
	if err != nil {
		return nil, apperrors.NotFoundWrap("address not found", "पता नहीं मिला", err)
	}
	resp := toAddressResponse(a)
	return &resp, nil
}

func (s *Service) Create(ctx context.Context, userID int64, req CreateAddressRequest) (*AddressResponse, *apperrors.AppError) {
	if req.IsDefault {
		if err := s.repo.UnsetDefaultAddress(ctx, userID); err != nil {
			return nil, apperrors.InternalWrap("failed to update default address", "डिफ़ॉल्ट पता अपडेट करने में विफल", err)
		}
	}

	a, err := s.repo.CreateAddress(ctx, db.CreateAddressParams{
		UserID:      userID,
		Label:       req.Label,
		FullAddress: req.FullAddress,
		Landmark:    req.Landmark,
		Latitude:    float64ToNumeric(req.Latitude),
		Longitude:   float64ToNumeric(req.Longitude),
		IsDefault:   req.IsDefault,
	})
	if err != nil {
		return nil, apperrors.InternalWrap("failed to create address", "पता बनाने में विफल", err)
	}
	resp := toAddressResponse(a)
	return &resp, nil
}

func (s *Service) Update(ctx context.Context, id, userID int64, req UpdateAddressRequest) (*AddressResponse, *apperrors.AppError) {
	if req.IsDefault {
		if err := s.repo.UnsetDefaultAddress(ctx, userID); err != nil {
			return nil, apperrors.InternalWrap("failed to update default address", "डिफ़ॉल्ट पता अपडेट करने में विफल", err)
		}
	}

	a, err := s.repo.UpdateAddress(ctx, db.UpdateAddressParams{
		ID:          id,
		UserID:      userID,
		Label:       req.Label,
		FullAddress: req.FullAddress,
		Landmark:    req.Landmark,
		Latitude:    float64ToNumeric(req.Latitude),
		Longitude:   float64ToNumeric(req.Longitude),
		IsDefault:   req.IsDefault,
	})
	if err != nil {
		return nil, apperrors.NotFoundWrap("address not found", "पता नहीं मिला", err)
	}
	resp := toAddressResponse(a)
	return &resp, nil
}

func (s *Service) Delete(ctx context.Context, id, userID int64) *apperrors.AppError {
	if err := s.repo.DeleteAddress(ctx, id, userID); err != nil {
		return apperrors.NotFoundWrap("address not found", "पता नहीं मिला", err)
	}
	return nil
}
