package usecase

import (
	"github.com/rs/zerolog"
	"github.com/zenoleg/binomeme/internal/rating"
)

type Dislike struct {
	storage rating.MemeStorage
	logger  zerolog.Logger
}

func NewDislike(storage rating.MemeStorage, logger zerolog.Logger) Dislike {
	return Dislike{
		storage: storage,
		logger:  logger,
	}
}

func (l Dislike) Handle(memeID string, reaction rating.Reaction) error {
	meme, err := l.storage.Get(memeID)
	if err != nil {
		return err
	}

	meme = meme.Underrate(reaction.Score())

	err = l.storage.Save(meme)
	if err != nil {
		return err
	}

	l.logger.Info().Str("meme_id", memeID).Int("score", reaction.Score()).Msg("❌ Reaction removed successfully")

	return nil
}
