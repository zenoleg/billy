package usecase

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/zenoleg/binomeme/internal/rating"
)

type InitRating struct {
	storage     rating.MemeStorage
	memeScanner rating.MemeScanner
	client      *socketmode.Client
}

func NewInitRating(storage rating.MemeStorage, memeScanner rating.MemeScanner, client *socketmode.Client) InitRating {
	return InitRating{
		storage:     storage,
		memeScanner: memeScanner,
		client:      client,
	}
}

func (r InitRating) Handle(channelID string) error {
	memes, err := r.memeScanner.Scan(channelID)
	if err != nil {
		return err
	}

	err = r.storage.Save(memes...)
	if err != nil {
		return err
	}

	_, _, err = r.client.PostMessage(channelID, slack.MsgOptionText(fmt.Sprintf("Channel scanned for a memes. Found: %d", len(memes)), false))
	if err != nil {
		return err
	}

	return err
}
