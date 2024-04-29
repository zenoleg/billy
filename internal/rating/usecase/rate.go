package usecase

import (
	"github.com/rs/zerolog"
	"github.com/zenoleg/binomeme/internal/rating"
)

type Rate struct {
	storage rating.MemeStorage
	logger  zerolog.Logger
}

func NewRate(storage rating.MemeStorage, logger zerolog.Logger) Rate {
	return Rate{
		storage: storage,
		logger:  logger,
	}
}

func (r Rate) Handle(memeID string, reaction rating.Reaction) error {
	meme, err := r.storage.Get(memeID)
	if err != nil {
		return err
	}

	meme = meme.Rate(reaction.Score())

	err = r.storage.Save(meme)
	if err != nil {
		return err
	}

	r.logger.Info().Str("meme_id", memeID).Int("score", reaction.Score()).Msg("👍 Meme rated successfully")

	return nil
}
