package slack

import (
	"context"
	"log"
	"os"

	"emperror.dev/errors"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type Config struct {
	appToken  string
	authToken string
	channelID string
}

func NewConfigFromEnv() (Config, error) {
	appToken, exists := os.LookupEnv("SLACK_APP_TOKEN")
	if !exists {
		return Config{}, errors.New("SLACK_APP_TOKEN environment variable not set")
	}

	authToken, exists := os.LookupEnv("SLACK_AUTH_TOKEN")
	if !exists {
		return Config{}, errors.New("SLACK_AUTH_TOKEN environment variable not set")
	}

	channelID, exists := os.LookupEnv("SLACK_CHANNEL_ID")
	if !exists {
		return Config{}, errors.New("SLACK_CHANNEL_ID environment variable not set")
	}

	return Config{
		appToken:  appToken,
		authToken: authToken,
		channelID: channelID,
	}, nil
}

type Server struct {
	client    *socketmode.Client
	logger    zerolog.Logger
	channelID string
}

func (s Server) Run(ctx context.Context) error {
	s.logger.Info().Msg("🚀 Starting Slack Server")

	return s.client.RunContext(ctx)
}

func NewBot(config Config, logger zerolog.Logger) *Server {
	client := slack.New(
		config.authToken,
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(config.appToken),
	)

	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(log.New(os.Stdout, "socketmode: ", log.Lshortfile|log.LstdFlags)),
	)

	return &Server{
		client:    socketClient,
		logger:    logger,
		channelID: config.channelID,
	}
}
