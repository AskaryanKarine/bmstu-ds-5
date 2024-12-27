package main

import (
	repository "github.com/AskaryanKarine/bmstu-ds-4/internal/loyalty/repositories/postgres"
	"github.com/AskaryanKarine/bmstu-ds-4/internal/loyalty/server"
	"github.com/AskaryanKarine/bmstu-ds-4/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func main() {
	cfg, err := config.ReadConfig("./configs/loyalty.env")
	if err != nil {
		log.Fatal(err)
	}

	db, err := gorm.Open(postgres.Open(cfg.PostgresDSN), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatal(err)
	}

	lr := repository.NewStorage(db)
	s := server.NewServer(lr)

	s.Run(cfg.Port)
}
