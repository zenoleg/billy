package usecase

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/zenoleg/billy/internal/rating"
)

type InitRating struct {
	memeStorage   rating.MemeStorage
	memberStorage rating.MemberStorage
	memeScanner   rating.MemeScanner
	client        *socketmode.Client
}

func NewInitRating(
	memeStorage rating.MemeStorage,
	memberStorage rating.MemberStorage,
	memeScanner rating.MemeScanner,
	client *socketmode.Client,
) InitRating {
	return InitRating{
		memeStorage:   memeStorage,
		memberStorage: memberStorage,
		memeScanner:   memeScanner,
		client:        client,
	}
}

func (r InitRating) Handle(channelID string) error {
	memes, err := r.memeScanner.Scan(channelID)
	if err != nil {
		return err
	}

	err = r.memeStorage.Save(memes...)
	if err != nil {
		return err
	}

	_, _, err = r.client.PostMessage(channelID, slack.MsgOptionText(fmt.Sprintf("Дико извеняюс пробегал кобанчиком увидел ваш канал чё тут? АХАХХААХ бляяя 20+ лет люди сидят картинки смешные оценивают, я в ваших годах уже качалку окончил)) ладно до встречи Задроты бляя)))\n\nНашел %d кринжовых мема и составил списки на увольнение: \n\n/memes_day\n/memes_week\n/memes_month\n/memes_ever\n/authros_week\n/authors_month\n/authors_ever", len(memes)), false))
	if err != nil {
		return err
	}

	return err
}
