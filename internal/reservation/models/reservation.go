package models

import (
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"time"
)

type ReservationResponse struct {
	ReservationUid string
	StartDate      string
	EndDate        string
	Status         models.PaymentStatusType
	PaymentUID     string
}

func (r *ReservationResponse) ToResponse(info models.HotelInfo) (models.ExtendedReservationResponse, error) {
	layout := "2006-01-02T15:04:05Z"
	tStart, err := time.Parse(layout, r.StartDate)
	if err != nil {
		return models.ExtendedReservationResponse{}, err
	}
	tEnd, err := time.Parse(layout, r.EndDate)
	if err != nil {
		return models.ExtendedReservationResponse{}, err
	}
	return models.ExtendedReservationResponse{
		ReservationResponse: models.ReservationResponse{
			ReservationUid: r.ReservationUid,
			StartDate:      tStart.Format("2006-01-02"),
			EndDate:        tEnd.Format("2006-01-02"),
			Status:         r.Status,
			Hotel:          info,
		},
		PaymentUID: r.PaymentUID,
	}, nil
}

type ReservationTable struct {
	ReservationUid string
	Username       string
	PaymentUid     string
	HotelID        int
	Status         models.PaymentStatusType
	StartDate      string
	EndDate        string
}

func ToReservationTable(model models.ExtendedCreateReservationResponse) ReservationTable {
	return ReservationTable{
		PaymentUid: model.PaymentUid,
		StartDate:  model.StartDate,
		EndDate:    model.EndDate,
	}
}
