package usecase

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"github.com/zenoleg/billy/internal/rating"
)

const (
	defaultLimit = 10
)

type (
	TopMemesQuery struct {
		requestMemberID rating.MemberID
		channelID       string
		now             time.Time
		period          rating.PeriodPreset
	}

	TopMemes struct {
		fetcher rating.TopFetcher
		client  *socketmode.Client
		logger  zerolog.Logger
	}
)

func NewTopMemesQuery(requestMemberID string, channelID string, now time.Time, period rating.PeriodPreset) TopMemesQuery {
	return TopMemesQuery{
		requestMemberID: rating.NewMemberID(requestMemberID),
		channelID:       channelID,
		now:             now,
		period:          period,
	}
}

func NewTopMemes(fetcher rating.TopFetcher, client *socketmode.Client, logger zerolog.Logger) TopMemes {
	return TopMemes{
		fetcher: fetcher,
		client:  client,
		logger:  logger,
	}
}

func (h TopMemes) Handle(query TopMemesQuery) error {
	from, to := query.period.MakeFromAndTo(time.Now().UTC())

	memeViews, err := h.fetcher.FetchTopMemes(rating.NewTopMemeCriterion(from, to, defaultLimit))
	if err != nil {
		return err
	}

	_, err = h.client.PostEphemeral(
		query.channelID,
		string(query.requestMemberID),
		slack.MsgOptionText(rating.NewTopMemesTemplate(memeViews, query.period).String(), false),
	)

	return err
}
