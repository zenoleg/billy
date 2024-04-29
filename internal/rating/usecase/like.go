package usecase

import (
	"github.com/rs/zerolog"
	"github.com/zenoleg/binomeme/internal/rating"
)

type Like struct {
	storage rating.MemeStorage
	logger  zerolog.Logger
}

func NewLike(storage rating.MemeStorage, logger zerolog.Logger) Like {
	return Like{
		storage: storage,
		logger:  logger,
	}
}

func (l Like) Handle(memeID string, reaction rating.Reaction) error {
	meme, err := l.storage.Get(memeID)
	if err != nil {
		return err
	}

	meme = meme.Rate(reaction.Score())

	err = l.storage.Save(meme)
	if err != nil {
		return err
	}

	l.logger.Info().Str("meme_id", memeID).Int("score", reaction.Score()).Msg("👍 Reaction added successfully")

	return nil
}
