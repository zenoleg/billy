package usecase

import (
	"github.com/rs/zerolog"
	"github.com/zenoleg/binomeme/internal/rating"
)

type (
	DislikeCommand struct {
		MemeID    string
		ChannelID string
		MemberID  rating.MemberID
		Reaction  rating.Reaction
		Timestamp string
	}

	Dislike struct {
		storage rating.MemeStorage
		logger  zerolog.Logger
	}
)

func NewDislikeCommand(
	memeID string,
	channelID string,
	memberID rating.MemberID,
	reaction string,
	timestamp string,
) DislikeCommand {
	return DislikeCommand{
		MemeID:    memeID,
		ChannelID: channelID,
		MemberID:  memberID,
		Reaction:  rating.NewReaction(reaction, 1),
		Timestamp: timestamp,
	}
}

func NewDislike(storage rating.MemeStorage, logger zerolog.Logger) Dislike {
	return Dislike{
		storage: storage,
		logger:  logger,
	}
}

func (l Dislike) Handle(command DislikeCommand) error {
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

	meme = meme.Underrate(command.Reaction.Score())

	err = l.storage.Save(meme)
	if err != nil {
		return err
	}

	l.logger.Info().
		Str("meme_id", command.MemeID).
		Str("reaction", command.Reaction.String()).
		Int("score", command.Reaction.Score()).
		Msg("❌ Reaction removed successfully")

	return nil
}
