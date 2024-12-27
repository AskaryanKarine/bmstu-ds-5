package main

import (
	"github.com/AskaryanKarine/bmstu-ds-4/internal/gateway/config"
	"github.com/AskaryanKarine/bmstu-ds-4/internal/gateway/server"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(cfg)
	s.Run(cfg.Port)
}
