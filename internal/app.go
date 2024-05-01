package internal

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/zenoleg/billy/internal/rating"
	"github.com/zenoleg/billy/internal/rating/usecase"
	"github.com/zenoleg/billy/internal/transport"
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

	sqliteMemeStorage, err := rating.NewSqliteMemeStorage(connection, logger)
	if err != nil {
		return App{}, func() {
			connection.Close()
		}, err
	}

	sqliteMemberStorage, err := rating.NewSqliteMemberStorage(connection, logger)
	if err != nil {
		return App{}, func() {
			connection.Close()
		}, err
	}

	linkFetcher := rating.NewSlackLinkFetcher(client, logger)
	topFetcher := rating.NewSQLiteTopFetcher(connection, logger)

	initRating := usecase.NewInitRating(
		sqliteMemeStorage,
		sqliteMemberStorage,
		rating.NewSlackMemeScanner(client, linkFetcher, logger),
		client,
	)
	like := usecase.NewLike(sqliteMemeStorage, linkFetcher, logger)
	topMemes := usecase.NewTopMemes(topFetcher, client, logger)
	topAuthors := usecase.NewTopAuthors(topFetcher, client, logger)
	dislike := usecase.NewDislike(sqliteMemeStorage, linkFetcher, logger)

	listener := transport.NewSlackEventListener(
		client,
		initRating,
		like,
		dislike,
		topMemes,
		topAuthors,
		logger,
	)

	return App{bot: bot, listener: listener}, func() {
		connection.Close()
	}, nil
}

func (app App) Start(ctx context.Context) error {
	app.listener.Start(ctx)

	return app.bot.Run(ctx)
}
