package models

import "github.com/AskaryanKarine/bmstu-ds-4/pkg/models"

type ExpandedLoyalty struct {
	models.LoyaltyInfoResponse
	Discount int
}
