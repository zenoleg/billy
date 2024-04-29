package usecase

import "github.com/zenoleg/binomeme/internal/rating"

type InitRating struct {
	storage     rating.MemeStorage
	memeScanner rating.MemeScanner
}

func NewInitRating(storage rating.MemeStorage, memeScanner rating.MemeScanner) InitRating {
	return InitRating{storage: storage, memeScanner: memeScanner}
}

func (r InitRating) Handle(channelID string) error {
	memes, err := r.memeScanner.Scan(channelID)
	if err != nil {
		return err
	}

	return r.storage.Save(memes...)
}
