package usecase

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/zenoleg/billy/internal/rating"
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

func (p TopPreset) Title() string {
	switch p {
	case TopDay:
		return "♂️️ Топ мемов за сегодня ♀️"
	case TopWeek:
		return "️️️♂️ Топ мемов за неделю ♀️"
	case TopMonth:
		return "️️♂️ Топ мемов за месяц ♀️"
	case TopEver:
		return "️️♂️ Топ мемов за все время ♀️"
	}

	return "️️♂️ Топ мемов ♀️"
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
	message.WriteString(fmt.Sprintf("%s\n\n", query.period.Title()))

	i := 1
	for _, view := range memeViews {
		placement := ""

		switch i {
		case 1:
			placement = "🥇 "
		case 2:
			placement = "🥈 "
		case 3:
			placement = "🥉 "
		default:
			placement = strconv.Itoa(i)
		}

		memeInfo := fmt.Sprintf("%s <%s|От %s> (%d)\n", placement, view.Link, view.MemberFullName, view.Score)
		message.WriteString(memeInfo)

		i++
	}

	_, _, err = h.client.PostMessage(query.channelID, slack.MsgOptionText(message.String(), false))

	return err
}
