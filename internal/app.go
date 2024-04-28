package internal

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/zenoleg/binomeme/internal/transport"
)

type App struct {
	bot transport.Bot
}

func MakeApp(logger zerolog.Logger) (App, error) {
	config, err := transport.NewSlackConfigFromEnv()
	if err != nil {
		return App{}, err
	}

	bot := transport.NewSlackBot(config, logger)

	return App{bot: bot}, nil
}

func (app App) Start(ctx context.Context) error {
	return app.bot.Run(ctx)
}