package internal

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/zenoleg/binomeme/internal/rating"
	"github.com/zenoleg/binomeme/internal/rating/usecase"
	"github.com/zenoleg/binomeme/internal/transport"
)

type App struct {
	bot      transport.Bot
	listener transport.SlackEventListener
}

func MakeApp(logger zerolog.Logger) (App, error) {
	config, err := transport.NewSlackConfigFromEnv()
	if err != nil {
		return App{}, err
	}

	client := transport.NewSlackClient(config, logger)
	bot := transport.NewSlackBot(client, logger)

	initRating := usecase.NewInitRating(rating.NewInMemoryMemeStorage(), rating.NewSlackMemeScanner(client, logger))

	listener := transport.NewSlackEventListener(client, initRating, logger)

	return App{bot: bot, listener: listener}, nil
}

func (app App) Start(ctx context.Context) error {
	app.listener.Start(ctx)

	return app.bot.Run(ctx)
}
