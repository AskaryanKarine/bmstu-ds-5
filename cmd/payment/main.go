package main

import (
	postgres2 "github.com/AskaryanKarine/bmstu-ds-4/internal/payment/repositories/postgres"
	"github.com/AskaryanKarine/bmstu-ds-4/internal/payment/server"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func main() {
	cfg, err := config.ReadConfig("./configs/payment.env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatal(err)
	}

	pr := postgres2.NewPaymentStorage(db)
	s := server.NewServer(pr)

	s.Run(cfg.Port)
}
