package usecase

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/zenoleg/binomeme/internal/rating"
)

const (
	defaultLimit = 10

	TopDay TopPreset = iota
	TopWeek
	TopMonth
	TopEver
)

type (
	TopPreset uint8

	TopMemesQuery struct {
		now       time.Time
		period    TopPreset
		channelID string
	}

	TopMemes struct {
		fetcher rating.TopMemeFetcher
		client  *socketmode.Client
		logger  zerolog.Logger
	}
)

func (p TopPreset) MakeFromAndTo(now time.Time) (time.Time, time.Time) {
	switch p {
	case TopDay:
		return now.Add(-time.Hour * 24), now
	case TopWeek:
		return now.Add(-time.Hour * 24 * 7), now
	case TopMonth:
		return now.Add(-time.Hour * 24 * 30), now
	case TopEver:
		return time.Unix(0, 0), now
	}

	return time.Time{}, time.Time{}
}

func NewTopMemesQuery(now time.Time, period TopPreset, channelID string) TopMemesQuery {
	return TopMemesQuery{
		now:       now,
		period:    period,
		channelID: channelID,
	}
}

func NewTop(fetcher rating.TopMemeFetcher, client *socketmode.Client, logger zerolog.Logger) TopMemes {
	return TopMemes{
		fetcher: fetcher,
		client:  client,
		logger:  logger,
	}
}

func (h TopMemes) Handle(query TopMemesQuery) error {
	from, to := query.period.MakeFromAndTo(time.Now().UTC())

	memeViews, err := h.fetcher.Fetch(rating.NewTopMemeCriterion(from, to, defaultLimit))
	if err != nil {
		return err
	}

	message := strings.Builder{}
	message.WriteString("🏆 Топ мемов 🏆\n\n")

	placement := 1
	for _, view := range memeViews {
		memeInfo := fmt.Sprintf("%d. <%s|От %s> (%d)\n", placement, view.Link, view.MemberFullName, view.Score)
		message.WriteString(memeInfo)

		placement++
	}

	_, _, err = h.client.PostMessage(query.channelID, slack.MsgOptionText(message.String(), false))

	return err
}
