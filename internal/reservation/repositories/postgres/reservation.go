package postgres

import (
	"context"
	"fmt"
	inner_models "github.com/AskaryanKarine/bmstu-ds-4/internal/reservation/models"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reservationStorage struct {
	db *gorm.DB
}

const (
	reservationTable = "reservation r"
)

func NewReservationStorage(db *gorm.DB) *reservationStorage {
	return &reservationStorage{
		db: db,
	}
}

func (r *reservationStorage) GetReservationByUUID(ctx context.Context, uuid, username string) (models.ExtendedReservationResponse, error) {
	var (
		hotelDB         models.HotelResponse
		reservationResp models.ExtendedReservationResponse
		reservationDB   struct {
			inner_models.ReservationResponse
			Username string
		}
	)
	err := r.db.WithContext(ctx).Table(reservationTable).
		Joins("JOIN "+hotelTable+" ON h.id = r.hotel_id").
		Select("h.*").Where("r.reservation_uid = ?", uuid).Take(&hotelDB).Error
	if err != nil {
		return models.ExtendedReservationResponse{}, fmt.Errorf("failed to get hotels by reservation_uuid %s: %w", uuid, err)
	}
	hotelInfo := models.HotelInfo{
		HotelUid:    hotelDB.HotelUid,
		Name:        hotelDB.Name,
		FullAddress: fmt.Sprintf("%s, %s, %s", hotelDB.Country, hotelDB.City, hotelDB.Address),
		Stars:       hotelDB.Stars,
	}

	err = r.db.WithContext(ctx).Table(reservationTable).Where("r.reservation_uid = ?", uuid).Take(&reservationDB).Error
	if err != nil {
		return models.ExtendedReservationResponse{}, fmt.Errorf("failed to get reservation by uuid %s: %w", uuid, err)
	}
	reservationResp, err = reservationDB.ToResponse(hotelInfo)
	if err != nil {
		return models.ExtendedReservationResponse{}, fmt.Errorf("failed to get reservation response %s: %w", uuid, err)
	}
	if username != reservationDB.Username {
		return models.ExtendedReservationResponse{}, fmt.Errorf("wrong username: %w", models.WrongUsernameError)
	}

	return reservationResp, nil
}

func (r *reservationStorage) GetAllReservationByUsername(ctx context.Context, username string) ([]models.ExtendedReservationResponse, error) {
	var (
		reservationDB       []inner_models.ReservationResponse
		reservationResponse []models.ExtendedReservationResponse
		hotelDB             models.HotelResponse
	)
	err := r.db.WithContext(ctx).Table(reservationTable).Where("r.username = ?", username).Find(&reservationDB).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation by username %s: %w", username, err)
	}
	reservationResponse = make([]models.ExtendedReservationResponse, 0, len(reservationDB))
	for i := range reservationDB {
		err := r.db.WithContext(ctx).Table(reservationTable).
			Joins("JOIN "+hotelTable+" ON h.id = r.hotel_id").
			Select("h.*").Where("r.reservation_uid = ?", reservationDB[i].ReservationUid).Take(&hotelDB).Error
		if err != nil {
			return nil, fmt.Errorf("failed to get hotels by reservation_uuid %s: %w", reservationDB[i].ReservationUid, err)
		}
		hotelInfo := models.HotelInfo{
			HotelUid:    hotelDB.HotelUid,
			Name:        hotelDB.Name,
			FullAddress: fmt.Sprintf("%s, %s, %s", hotelDB.Country, hotelDB.City, hotelDB.Address),
			Stars:       hotelDB.Stars,
		}
		res, err := reservationDB[i].ToResponse(hotelInfo)
		if err != nil {
			return nil, fmt.Errorf("failed to get reservation response %s: %w", reservationDB[i].ReservationUid, err)
		}
		reservationResponse = append(reservationResponse, res)
	}

	return reservationResponse, nil
}

func (r *reservationStorage) Delete(ctx context.Context, uuid string) error {
	err := r.db.WithContext(ctx).Table(reservationTable).Where("reservation_uid = ?", uuid).Update("status", models.CANCELED).Error
	if err != nil {
		return fmt.Errorf("failed deleting reservation %s: %w", uuid, err)
	}
	return nil
}

func (r *reservationStorage) Create(ctx context.Context, reservation models.ExtendedCreateReservationResponse, username string) (string, error) {
	reservationDB := inner_models.ToReservationTable(reservation)

	err := r.db.WithContext(ctx).Table(hotelTable).Where("h.hotel_uid = ?", reservation.HotelUid).
		Select("id").Take(&reservationDB.HotelID).Error
	if err != nil {
		return "", fmt.Errorf("failed creating reservation: %w", err)
	}

	reservationDB.ReservationUid = uuid.New().String()
	reservationDB.Status = models.PAID
	reservationDB.Username = username

	err = r.db.WithContext(ctx).Table("reservation").Create(&reservationDB).Error
	if err != nil {
		return "", fmt.Errorf("failed creating reservation: %w", err)
	}
	return reservationDB.ReservationUid, nil
}
