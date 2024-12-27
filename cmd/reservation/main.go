package main

import (
	repository "github.com/AskaryanKarine/bmstu-ds-4/internal/reservation/repositories/postgres"
	"github.com/AskaryanKarine/bmstu-ds-4/internal/reservation/server"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func main() {
	cfg, err := config.ReadConfig("./configs/reservation.env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatal(err)
	}

	hr := repository.NewHotelStorage(db)
	rr := repository.NewReservationStorage(db)
	s := server.New(hr, rr)

	s.Run(cfg.Port)
}
