package main

import (
	"fmt"
	"net"

	"github.com/isd-sgcu/johnjud-backend/src/config"
	"github.com/isd-sgcu/johnjud-backend/src/database"
	"github.com/rs/zerolog/log"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "backend load config").
			Msg("Failed to start service")
	}

	db, err := database.InitPostgresDatabase(&conf.Database, conf.App.Debug)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "backend database").
			Msg("Failed to start service")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", conf.App.Port))
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", "backend server").
			Msg("Failed to start service")
	}
}
