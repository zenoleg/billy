package usecase

import (
	"github.com/rs/zerolog"
	"github.com/zenoleg/binomeme/internal/rating"
)

type (
	LikeCommand struct {
		MemeID    string
		ChannelID string
		MemberID  rating.MemberID
		Reaction  rating.Reaction
		Timestamp string
	}

	Like struct {
		storage rating.MemeStorage
		logger  zerolog.Logger
	}
)

func NewLikeCommand(
	memeID string,
	channelID string,
	memberID rating.MemberID,
	reaction string,
	timestamp string,
) LikeCommand {
	return LikeCommand{
		MemeID:    memeID,
		ChannelID: channelID,
		MemberID:  memberID,
		Reaction:  rating.NewReaction(reaction, 1),
		Timestamp: timestamp,
	}
}

func NewLike(storage rating.MemeStorage, logger zerolog.Logger) Like {
	return Like{
		storage: storage,
		logger:  logger,
	}
}

func (l Like) Handle(command LikeCommand) error {
	meme, err := l.storage.Get(command.MemeID)
	if err != nil {
		meme = rating.NewMeme(
			command.MemeID,
			command.ChannelID,
			command.MemberID,
			rating.NewReactions([]rating.Reaction{command.Reaction}),
			command.Timestamp,
		)
	}

	meme = meme.Rate(command.Reaction.Score())

	err = l.storage.Save(meme)
	if err != nil {
		return err
	}

	l.logger.Info().
		Str("meme_id", command.MemeID).
		Str("reaction", command.Reaction.String()).
		Int("score", command.Reaction.Score()).
		Msg("👍 Reaction added successfully")

	return nil
}
