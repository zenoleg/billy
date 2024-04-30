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

	_, _, err = r.client.PostMessage(channelID, slack.MsgOptionText(fmt.Sprintf("Дико извеняюс пробегал кобанчиком увидел ваш канал чё тут? АХАХХААХ бляяя 20+ лет люди сидят картинки смешные оценивают, я в ваших годах уже старший слесарь был)) ладно до встречи Задроты бляя)))\n\nНашел %d кринжовых мема и составил список на увольнение: \n\n/top_posts_day\n/top_posts_week\n/top_posts_month\n/top_authros_week\n/top_authors_month", len(memes)), false))
	if err != nil {
		return err
	}

	return err
}
