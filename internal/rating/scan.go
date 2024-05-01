package rating

import (
	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
)

type (
	MemeScanner interface {
		Scan(channelID string) ([]Meme, []Member, error)
	}

	LoggedScanner struct {
		logger      zerolog.Logger
		memeScanner MemeScanner
	}

	SlackMemeScanner struct {
		client      *socketmode.Client
		linkFetcher LinkFetcher
	}
)

func NewSlackMemeScanner(client *socketmode.Client, linkFetcher LinkFetcher, logger zerolog.Logger) MemeScanner {
	return LoggedScanner{
		logger: logger,
		memeScanner: SlackMemeScanner{
			client:      client,
			linkFetcher: linkFetcher,
		},
	}
}

func (s SlackMemeScanner) Scan(channelID string) ([]Meme, []Member, error) {
	historyParams := slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     999,
	}

	memes := make([]Meme, 0, 100)
	members := make([]Member, 0, 10)

	for {
		conversationResponse, err := s.client.GetConversationHistory(&historyParams)
		if err != nil {
			return []Meme{}, []Member{}, err
		}

		for _, message := range conversationResponse.Messages {
			if len(message.Files) == 0 {
				// ignore message if there is no meme inside
				continue
			}

			reactions := make([]Reaction, 0, 5)
			for _, reactionInfo := range message.Reactions {
				reactions = append(reactions, NewReaction(reactionInfo.Name, reactionInfo.Count))
			}

			memes = append(
				memes,
				NewMeme(
					message.Timestamp,
					channelID,
					NewMemberID(message.User),
					NewReactions(reactions),
					message.Timestamp,
					s.linkFetcher.Fetch(message.Timestamp, channelID),
				),
			)
		}

		if !conversationResponse.HasMore {
			break
		}

		historyParams.Cursor = conversationResponse.ResponseMetadata.Cursor
	}

	return memes, members, nil
}

func (l LoggedScanner) Scan(channelID string) ([]Meme, []Member, error) {
	memes, members, err := l.memeScanner.Scan(channelID)
	if err != nil {
		l.logger.Err(err).Str("channel_id", channelID).Msg("❌ Can not scan channel conversation")

		return []Meme{}, []Member{}, err
	}

	l.logger.Info().Int("meme_count", len(memes)).Str("channel_id", channelID).Msg("✅ Channel conversation scanned successfully")

	return memes, members, nil
}
