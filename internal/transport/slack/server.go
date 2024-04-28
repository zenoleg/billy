package slack

import (
	"context"
	"os"

	"emperror.dev/errors"
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/zenoleg/binomeme/internal/transport"
)

type (
	Config struct {
		appToken  string
		authToken string
		channelID string
	}
	Server struct {
		client    *socketmode.Client
		logger    zerolog.Logger
		channelID string
	}
	debugLogger struct {
		logger zerolog.Logger
	}
)

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

func NewBot(config Config, logger zerolog.Logger) transport.Bot {
	debugLog := newDebugLogger(logger.With().Str("bot", "slack_socket").Logger())

	client := slack.New(
		config.authToken,
		slack.OptionLog(debugLog),
		slack.OptionDebug(true),
		slack.OptionAppLevelToken(config.appToken),
	)

	socketClient := socketmode.New(
		client,
		socketmode.OptionDebug(true),
		socketmode.OptionLog(debugLog),
	)

	return &Server{
		client:    socketClient,
		logger:    logger,
		channelID: config.channelID,
	}
}

func newDebugLogger(logger zerolog.Logger) debugLogger {
	return debugLogger{logger: logger}
}

func (s Server) Run(ctx context.Context) error {
	s.logger.Info().Msg("🚀 Starting Slack Server")

	return s.client.RunContext(ctx)
}

func (l debugLogger) Output(i int, msg string) error {
	l.logger.Debug().Msg(msg)

	return nil
}
