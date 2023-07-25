package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	zerolog.New(os.Stderr).With().
		Str("service", "test").
		Str("hahaha", "hohoho").
		Logger()

	subLog := log.With().
		Str("method", "GET").
		Str("url", "asdf").
		Logger()

	subLog.Info().Msg("test")
}
