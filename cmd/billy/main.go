package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/zenoleg/billy/internal"
	"github.com/zenoleg/billy/third_party/logger"
)

var version = "unknown"

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: %v\n", err)

		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	zerolog.DurationFieldUnit = time.Microsecond
	log := logger.NewLogger(logger.NewConfig(), version)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("Can not load .env")
	}

	app, closeFunc, err := internal.MakeApp(log)
	defer closeFunc()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create app")
	}

	ctx, cancel := context.WithCancel(context.Background())

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		defer cancel()
		close(stopped)
	}()

	if err = app.Start(ctx); err != nil {
		log.Err(err).Msg("Bot stopped")
	}

	<-stopped

	log.Info().Msg("Bye! 👋")

	return nil
}
