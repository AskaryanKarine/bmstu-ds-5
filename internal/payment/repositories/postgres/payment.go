package postgres

import (
	"context"
	"fmt"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type paymentStorage struct {
	db *gorm.DB
}

func NewPaymentStorage(db *gorm.DB) *paymentStorage {
	return &paymentStorage{
		db: db,
	}
}

const (
	paymentTable = "payment"
)

func (p *paymentStorage) GetPaymentInfoByUUID(ctx context.Context, uuid string) (models.PaymentInfo, error) {
	var paymentInfo models.PaymentInfo

	err := p.db.WithContext(ctx).Table(paymentTable).Where("payment_uid = ?", uuid).Take(&paymentInfo).Error
	if err != nil {
		return models.PaymentInfo{}, fmt.Errorf("failed to get reservation by uuid %s: %w", uuid, err)
	}

	return paymentInfo, nil
}

func (p *paymentStorage) Delete(ctx context.Context, uuid string) error {
	err := p.db.WithContext(ctx).Table(paymentTable).Where("payment_uid = ?", uuid).Update("status", models.CANCELED).Error
	if err != nil {
		return fmt.Errorf("failed deleting payment %s: %w", uuid, err)
	}
	return nil
}

func (p *paymentStorage) Create(ctx context.Context, payment models.PaymentInfo) (string, error) {
	newPayment := struct {
		PaymentUid string
		models.PaymentInfo
	}{}
	newPayment.PaymentUid = uuid.New().String()
	newPayment.PaymentInfo = payment
	err := p.db.WithContext(ctx).Table(paymentTable).Create(&newPayment).Error
	if err != nil {
		return "", fmt.Errorf("failed creating payment: %w", err)
	}
	return newPayment.PaymentUid, nil
}
