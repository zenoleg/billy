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

type (
	TopAuthorsQuery struct {
		requestMemberID rating.MemberID
		channelID       string
		now             time.Time
		period          TopPreset
	}

	TopAuthors struct {
		fetcher rating.TopFetcher
		client  *socketmode.Client
		logger  zerolog.Logger
	}
)

func NewTopAuthorsQuery(requestMemberID string, channelID string, now time.Time, period TopPreset) TopAuthorsQuery {
	return TopAuthorsQuery{
		requestMemberID: rating.NewMemberID(requestMemberID),
		channelID:       channelID,
		now:             now,
		period:          period,
	}
}

func NewTopAuthors(fetcher rating.TopFetcher, client *socketmode.Client, logger zerolog.Logger) TopAuthors {
	return TopAuthors{
		fetcher: fetcher,
		client:  client,
		logger:  logger,
	}
}

func (h TopAuthors) Handle(query TopAuthorsQuery) error {
	from, to := query.period.MakeFromAndTo(time.Now().UTC())

	authorViews, err := h.fetcher.FetchTopAuthors(rating.NewTopMemeCriterion(from, to, defaultLimit))
	if err != nil {
		return err
	}

	message := strings.Builder{}
	message.WriteString(fmt.Sprintf("%s\n\n", query.period.Title()))

	i := 1
	for _, view := range authorViews {
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

		memeInfo := fmt.Sprintf("%s %s (%d)\n", placement, view.MemberFullName, view.Score)
		message.WriteString(memeInfo)

		i++
	}

	_, err = h.client.PostEphemeral(query.channelID, string(query.requestMemberID), slack.MsgOptionText(message.String(), false))

	return err
}
