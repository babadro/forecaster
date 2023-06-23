package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	bot "github.com/babadro/forecaster/internal/core/forecaster_bot"
	"github.com/babadro/forecaster/internal/infra/postgres"
	"github.com/caarlos0/env"
	"github.com/jackc/pgx/v4/pgxpool"
)

var envVars = struct {
	HTTPAddr string `env:"HTTP_ADDR" envDefault:":8080"`
	DBConn   string `env:"DB_CONN"`
}{}

func main() {
	// listen to os signals
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	if err := env.Parse(&envVars); err != nil {
		log.Fatalf("Unable to parse env vars: %v\n", err)
	}

	dbPool, err := pgxpool.Connect(context.Background(), envVars.DBConn)
	if err != nil {
		log.Fatalf("Unable to connection to database :%v\n", err)
	}
	defer dbPool.Close()

	forecastDB := postgres.NewForecastDB(dbPool)

	_ = bot.NewService(forecastDB)

	// wait for os signal
	<-c
}
