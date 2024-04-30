package rating

import (
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type (
	LinkFetcher interface {
		Fetch(memeID string, channelID string) string
	}

	SlackLinkFetcher struct {
		client *socketmode.Client
		logger zerolog.Logger
	}
)

func NewSlackLinkFetcher(client *socketmode.Client, logger zerolog.Logger) LinkFetcher {
	return SlackLinkFetcher{
		client: client,
		logger: logger,
	}
}

func (s SlackLinkFetcher) Fetch(memeID string, channelID string) string {
	permalink, err := s.client.GetPermalink(&slack.PermalinkParameters{
		Channel: channelID,
		Ts:      memeID,
	})

	if err != nil {
		s.logger.Err(err).
			Str("meme_id", memeID).
			Str("channel_id", channelID).
			Msg("⏳ Failed to fetch permalink to message")

		return ""
	}

	return permalink
}
