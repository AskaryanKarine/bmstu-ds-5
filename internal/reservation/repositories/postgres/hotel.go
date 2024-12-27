package postgres

import (
	"context"
	"fmt"
	innermodels "github.com/AskaryanKarine/bmstu-ds-4/internal/reservation/models"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"

	"gorm.io/gorm"
)

type hotelStorage struct {
	db *gorm.DB
}

const (
	hotelTable = "hotels h"
)

func NewHotelStorage(db *gorm.DB) *hotelStorage {
	return &hotelStorage{db: db}
}

func (h *hotelStorage) GetAllHotels(ctx context.Context, pagination innermodels.PaginationParams) ([]models.HotelResponse, int, error) {
	var (
		results    []models.HotelResponse
		totalCount int
	)
	err := h.db.WithContext(ctx).Table(hotelTable).Order("id").
		Limit(pagination.Size).Offset((pagination.Page - 1) * pagination.Size).
		Find(&results).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get hotels: %w", err)
	}

	err = h.db.Table(hotelTable).
		Select("COUNT(*)").
		Take(&totalCount).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get count of hotels: %w", err)
	}

	return results, totalCount, nil
}

func (h *hotelStorage) GetHotelsInfoByUser(ctx context.Context, username string) ([]models.HotelResponse, error) {
	var results []models.HotelResponse
	err := h.db.WithContext(ctx).Table(reservationTable).
		Joins("JOIN "+hotelTable+" ON h.id = r.hotel_id").
		Select("h.*").Where("username = ?", username).Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get hotels by user %s: %w", username, err)
	}
	return results, nil
}

func (h *hotelStorage) GetHotelInfoByUUID(ctx context.Context, uuid string) (models.HotelResponse, error) {
	var result models.HotelResponse
	err := h.db.WithContext(ctx).Table(hotelTable).
		Where("hotel_uid = ?", uuid).Take(&result).Error
	if err != nil {
		return models.HotelResponse{}, fmt.Errorf("failed to get hotel by uuid %s: %w", uuid, err)
	}
	return result, nil
}
