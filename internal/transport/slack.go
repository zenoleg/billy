package transport

import (
	"context"
	"os"

	"emperror.dev/errors"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type (
	Bot interface {
		Run(ctx context.Context) error
	}

	SlackConfig struct {
		appToken  string
		authToken string
		channelID string
	}

	SlackServer struct {
		client *socketmode.Client
		logger zerolog.Logger
	}
)

func NewSlackConfigFromEnv() (SlackConfig, error) {
	appToken, exists := os.LookupEnv("SLACK_APP_TOKEN")
	if !exists {
		return SlackConfig{}, errors.New("SLACK_APP_TOKEN environment variable not set")
	}

	authToken, exists := os.LookupEnv("SLACK_AUTH_TOKEN")
	if !exists {
		return SlackConfig{}, errors.New("SLACK_AUTH_TOKEN environment variable not set")
	}

	channelID, exists := os.LookupEnv("SLACK_CHANNEL_ID")
	if !exists {
		return SlackConfig{}, errors.New("SLACK_CHANNEL_ID environment variable not set")
	}

	return SlackConfig{
		appToken:  appToken,
		authToken: authToken,
		channelID: channelID,
	}, nil
}

func NewSlackClient(config SlackConfig, logger zerolog.Logger) *socketmode.Client {
	client := slack.New(
		config.authToken,
		slack.OptionAppLevelToken(config.appToken),
	)

	return socketmode.New(client)
}

func NewSlackBot(client *socketmode.Client, logger zerolog.Logger) Bot {
	return &SlackServer{
		client: client,
		logger: logger,
	}
}

func (s SlackServer) Run(ctx context.Context) error {
	s.logger.Info().Msg("🚀 Starting Slack Server")

	return s.client.RunContext(ctx)
}
