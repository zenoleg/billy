package transport

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/zenoleg/binomeme/internal/rating/usecase"
)

type SlackEventListener struct {
	channelID  string
	client     *socketmode.Client
	initRating usecase.InitRating
	logger     zerolog.Logger
}

func NewSlackEventListener(
	channelID string,
	client *socketmode.Client,
	initRating usecase.InitRating,
	logger zerolog.Logger,
) SlackEventListener {
	return SlackEventListener{
		channelID:  channelID,
		client:     client,
		initRating: initRating,
		logger:     logger,
	}
}

func (l SlackEventListener) Start(ctx context.Context) {
	go func() {
		l.logger.Info().Msg("Slack event listener started")

		for {
			select {
			case evt := <-l.client.Events:
				l.logger.Debug().Msgf("Got event from Slack listener %#v", evt)

				switch evt.Type {
				case socketmode.EventTypeSlashCommand:
					data := evt.Data.(slack.SlashCommand)
					switch data.Command {
					case "/init":
						err := l.initRating.Handle(l.channelID)
						if err != nil {
							l.logger.Err(err).Msg("Can not initialize rating")
						}
					}

					l.client.Ack(*evt.Request)
				}
			case <-ctx.Done():
				l.logger.Info().Msg("Slack event listener shutting down")

				return
			}
		}
	}()
}
