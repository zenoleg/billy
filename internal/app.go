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

func MakeApp(logger zerolog.Logger) (App, func(), error) {
	config, err := transport.NewSlackConfigFromEnv()
	if err != nil {
		return App{}, func() {}, err
	}

	client := transport.NewSlackClient(config, logger)
	bot := transport.NewSlackBot(client, logger)
	connection, err := rating.NewConnection("./memes.db", logger)
	if err != nil {
		return App{}, func() {}, err
	}

	sqliteMemeStorage, closeFunc, err := rating.NewSqliteMemeStorage(connection, logger)

	if err != nil {
		return App{}, closeFunc, err
	}

	linkFetcher := rating.NewSlackLinkFetcher(client, logger)
	initRating := usecase.NewInitRating(sqliteMemeStorage, rating.NewSlackMemeScanner(client, linkFetcher, logger), client)
	like := usecase.NewLike(sqliteMemeStorage, linkFetcher, logger)
	dislike := usecase.NewDislike(sqliteMemeStorage, linkFetcher, logger)

	listener := transport.NewSlackEventListener(client, initRating, like, dislike, logger)

	return App{bot: bot, listener: listener}, closeFunc, nil
}

func (app App) Start(ctx context.Context) error {
	app.listener.Start(ctx)

	return app.bot.Run(ctx)
}
