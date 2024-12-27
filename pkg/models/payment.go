package models

type PaymentStatusType string

const (
	PAID     PaymentStatusType = "PAID"
	CANCELED PaymentStatusType = "CANCELED"
)

type PaymentInfo struct {
	// Status - статус операции оплаты
	Status PaymentStatusType `json:"status"`
	// Price - сумма операции
	Price int `json:"price"`
}

type PaymentCreateRequest struct {
	Price     int    `json:"price" validate:"required,gt=0"`
	Discount  int    `json:"discount" validate:"required,gt=0"`
	StartDate string `json:"startDate" validate:"required,IsISO8601"`
	EndDate   string `json:"endDate" validate:"required,IsISO8601"`
}

type ExtendedPaymentInfo struct {
	PaymentUid string `json:"paymentUid"`
	PaymentInfo
}
