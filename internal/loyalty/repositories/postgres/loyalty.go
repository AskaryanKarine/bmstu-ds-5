package postgres

import (
	"context"
	"fmt"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"gorm.io/gorm"
)

const (
	loyaltyTable = "loyalty"
)

type storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) *storage {
	return &storage{db: db}
}

func (s *storage) GetByUser(ctx context.Context, username string) (models.LoyaltyInfoResponse, error) {
	result := models.LoyaltyInfoResponse{}
	err := s.db.WithContext(ctx).Table(loyaltyTable).Where("username = ?", username).Take(&result).Error
	if err != nil {
		return models.LoyaltyInfoResponse{}, fmt.Errorf("failed to get loyalty by username %w", err)
	}
	return result, nil
}

func (s *storage) UpdateByUser(ctx context.Context, username string, usersLoyalty models.LoyaltyInfoResponse) error {
	err := s.db.WithContext(ctx).Table(loyaltyTable).Where("username = ?", username).Updates(&usersLoyalty).Error
	if err != nil {
		return fmt.Errorf("failed to update loyalty by username %w", err)
	}
	return nil
}
