package rating

import (
	"slices"

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
		client             *socketmode.Client
		channelInfoFetcher ChannelInfoFetcher
		logger             zerolog.Logger
	}
)

func NewSlackMemeScanner(client *socketmode.Client, channelInfoFetcher ChannelInfoFetcher, logger zerolog.Logger) MemeScanner {
	return LoggedScanner{
		logger: logger,
		memeScanner: SlackMemeScanner{
			client:             client,
			channelInfoFetcher: channelInfoFetcher,
			logger:             logger,
		},
	}
}

func (s SlackMemeScanner) Scan(channelID string) ([]Meme, []Member, error) {
	historyParams := slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     500,
	}

	memes := make([]Meme, 0, 500)

	existedMembers := make([]MemberID, 0, 10)
	members := make([]Member, 0, 10)

	batchCount := 0

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
					s.channelInfoFetcher.FetchMemeLink(message.Timestamp, channelID),
				),
			)

			if !slices.Contains(existedMembers, NewMemberID(message.User)) {
				members = append(members, s.channelInfoFetcher.FetchMember(NewMemberID(message.User)))
				existedMembers = append(existedMembers, NewMemberID(message.User))
			}
		}

		if !conversationResponse.HasMore {
			break
		}

		historyParams.Cursor = conversationResponse.ResponseMetadata.Cursor

		batchCount++

		s.logger.Info().Int("batch", batchCount).Msg("🤣 Conversation history batch fetched")
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
