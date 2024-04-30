package transport

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"github.com/zenoleg/billy/internal/rating"
	"github.com/zenoleg/billy/internal/rating/usecase"
)

type SlackEventListener struct {
	client     *socketmode.Client
	initRating usecase.InitRating
	like       usecase.Like
	dislike    usecase.Dislike
	top        usecase.TopMemes
	logger     zerolog.Logger
}

func NewSlackEventListener(
	client *socketmode.Client,
	initRating usecase.InitRating,
	like usecase.Like,
	dislike usecase.Dislike,
	top usecase.TopMemes,
	logger zerolog.Logger,
) SlackEventListener {
	return SlackEventListener{
		client:     client,
		initRating: initRating,
		like:       like,
		dislike:    dislike,
		top:        top,
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

						l.client.Ack(*evt.Request)

					case "/memes_day":
						err := l.top.Handle(usecase.NewTopMemesQuery(time.Now().UTC(), usecase.TopDay, data.ChannelID))
						if err != nil {
							l.logger.Err(err).Msg("Can not fetch top rating")
						}

						l.client.Ack(*evt.Request)
					case "/memes_week":
						err := l.top.Handle(usecase.NewTopMemesQuery(time.Now().UTC(), usecase.TopWeek, data.ChannelID))
						if err != nil {
							l.logger.Err(err).Msg("Can not fetch top rating")
						}

						l.client.Ack(*evt.Request)
					case "/memes_month":
						err := l.top.Handle(usecase.NewTopMemesQuery(time.Now().UTC(), usecase.TopMonth, data.ChannelID))
						if err != nil {
							l.logger.Err(err).Msg("Can not fetch top rating")
						}

						l.client.Ack(*evt.Request)
					case "/memes_ever":
						err := l.top.Handle(usecase.NewTopMemesQuery(time.Now().UTC(), usecase.TopEver, data.ChannelID))
						if err != nil {
							l.logger.Err(err).Msg("Can not fetch top rating")
						}

						l.client.Ack(*evt.Request)
					case "/top_authors_week":
						fmt.Println("Called top_authors_week")
						l.client.Ack(*evt.Request)
					case "/top_authors_month":
						fmt.Println("Called top_authors_month")
						l.client.Ack(*evt.Request)
					case "/top_authors_ever":
						fmt.Println("Called top_authors_ever")
						l.client.Ack(*evt.Request)
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
							err := l.like.Handle(usecase.NewLikeCommand(
								ev.Item.Timestamp,
								ev.Item.Channel,
								rating.NewMemberID(ev.ItemUser),
								ev.Reaction,
								ev.Item.Timestamp,
							))

							if err != nil {
								l.logger.Err(err).Str("meme_id", ev.Item.Timestamp).Str("reaction", ev.Reaction).Msg("Can not like a meme")
							}

							l.client.Ack(*evt.Request)
						case *slackevents.ReactionRemovedEvent:
							err := l.dislike.Handle(usecase.NewDislikeCommand(
								ev.Item.Timestamp,
								ev.Item.Channel,
								rating.NewMemberID(ev.ItemUser),
								ev.Reaction,
								ev.Item.Timestamp,
							))

							if err != nil {
								l.logger.Err(err).Str("meme_id", ev.Item.Timestamp).Str("reaction", ev.Reaction).Msg("Can not dislike a meme")
							}

							l.client.Ack(*evt.Request)
						}
					}
				}

			case <-ctx.Done():
				l.logger.Info().Msg("Slack event listener shutting down")

				return
			}
		}
	}()
}
