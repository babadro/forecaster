package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/babadro/forecaster/internal/infra/postgres"
	"github.com/caarlos0/env"
	"github.com/jackc/pgx/v4/pgxpool"
)

var envVars = struct {
	HTTPAddr string `env:"HTTP_ADDR" envDefault:":8080"`
}{}

func main() {
	// listen to os signals
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	if err := env.Parse(&envVars); err != nil {
		log.Fatalf("Unable to parse env vars: %v\n", err)
	}

	// read env vars
	connString := "postgresql://postgres:postgres@localhost:5432/forecaster"

	dbPool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		log.Fatalf("Unable to connection to database :%v\n", err)
	}
	defer dbPool.Close()

	forecastDB := postgres.NewForecastDB(dbPool)

	_ = forecastDB

	// wait for os signal
	<-c
}
