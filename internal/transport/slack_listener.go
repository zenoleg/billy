package transport

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/zenoleg/binomeme/internal/rating"
	"github.com/zenoleg/binomeme/internal/rating/usecase"
)

type SlackEventListener struct {
	client     *socketmode.Client
	initRating usecase.InitRating
	rate       usecase.Rate
	logger     zerolog.Logger
}

func NewSlackEventListener(
	client *socketmode.Client,
	initRating usecase.InitRating,
	rate usecase.Rate,
	logger zerolog.Logger,
) SlackEventListener {
	return SlackEventListener{
		client:     client,
		initRating: initRating,
		rate:       rate,
		logger:     logger,
	}
}

func (l SlackEventListener) Start(ctx context.Context) {
	go func() {
		l.logger.Info().Msg("👂 Slack event listener started")

		for {
			select {
			case evt := <-l.client.Events:
				l.logger.Debug().Msgf("Got event from Slack listener %#v", evt)

				switch evt.Type {
				case socketmode.EventTypeSlashCommand:
					data := evt.Data.(slack.SlashCommand)
					switch data.Command {
					case "/init":
						err := l.initRating.Handle(data.ChannelID)
						if err != nil {
							l.logger.Err(err).Msg("Can not initialize rating")
						}
					}

				case socketmode.EventTypeEventsAPI:
					eventsAPIEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
					if !ok {
						l.logger.Error().Msgf("unexpected event type: %s", evt.Type)
					}

					if eventsAPIEvent.Type == slackevents.CallbackEvent {
						innerEvent := eventsAPIEvent.InnerEvent

						switch ev := innerEvent.Data.(type) {
						case *slackevents.ReactionAddedEvent:
							err := l.rate.Handle(ev.Item.Timestamp, rating.NewReaction(ev.Reaction, 1))
							if err != nil {
								l.logger.Err(err).Str("meme_id", ev.Item.Timestamp).Str("reaction", ev.Reaction).Msg("Can not rate a meme")
							}
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
