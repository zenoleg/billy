package transport

import (
	"context"
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
	topMemes   usecase.TopMemes
	topAuthors usecase.TopAuthors
	logger     zerolog.Logger
}

func NewSlackEventListener(
	client *socketmode.Client,
	initRating usecase.InitRating,
	like usecase.Like,
	dislike usecase.Dislike,
	topMemes usecase.TopMemes,
	topAuthors usecase.TopAuthors,
	logger zerolog.Logger,
) SlackEventListener {
	return SlackEventListener{
		client:     client,
		initRating: initRating,
		like:       like,
		dislike:    dislike,
		topMemes:   topMemes,
		topAuthors: topAuthors,
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
						err := l.topMemes.Handle(usecase.NewTopMemesQuery(data.UserID, data.ChannelID, time.Now().UTC(), usecase.TopDay))
						if err != nil {
							l.logger.Err(err).Msg("Can not fetch top memes rating")
						}

						l.client.Ack(*evt.Request)
					case "/memes_week":
						err := l.topMemes.Handle(usecase.NewTopMemesQuery(data.UserID, data.ChannelID, time.Now().UTC(), usecase.TopWeek))
						if err != nil {
							l.logger.Err(err).Msg("Can not fetch top memes rating")
						}

						l.client.Ack(*evt.Request)
					case "/memes_month":
						err := l.topMemes.Handle(usecase.NewTopMemesQuery(data.UserID, data.ChannelID, time.Now().UTC(), usecase.TopMonth))
						if err != nil {
							l.logger.Err(err).Msg("Can not fetch top memes rating")
						}

						l.client.Ack(*evt.Request)
					case "/memes_ever":
						err := l.topMemes.Handle(usecase.NewTopMemesQuery(data.UserID, data.ChannelID, time.Now().UTC(), usecase.TopEver))
						if err != nil {
							l.logger.Err(err).Msg("Can not fetch top memes rating")
						}

						l.client.Ack(*evt.Request)
					case "/authors_week":
						err := l.topAuthors.Handle(usecase.NewTopAuthorsQuery(data.UserID, data.ChannelID, time.Now().UTC(), usecase.TopWeek))
						if err != nil {
							l.logger.Err(err).Msg("Can not fetch top authors rating")
						}

						l.client.Ack(*evt.Request)
					case "/authors_month":
						err := l.topAuthors.Handle(usecase.NewTopAuthorsQuery(data.UserID, data.ChannelID, time.Now().UTC(), usecase.TopMonth))
						if err != nil {
							l.logger.Err(err).Msg("Can not fetch top authors rating")
						}

						l.client.Ack(*evt.Request)
					case "/authors_ever":
						err := l.topAuthors.Handle(usecase.NewTopAuthorsQuery(data.UserID, data.ChannelID, time.Now().UTC(), usecase.TopEver))
						if err != nil {
							l.logger.Err(err).Msg("Can not fetch top authors rating")
						}

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
