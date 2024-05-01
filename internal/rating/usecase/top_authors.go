package usecase

import (
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
		period          rating.PeriodPreset
	}

	TopAuthors struct {
		fetcher rating.TopFetcher
		client  *socketmode.Client
		logger  zerolog.Logger
	}
)

func NewTopAuthorsQuery(requestMemberID string, channelID string, now time.Time, period rating.PeriodPreset) TopAuthorsQuery {
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

	_, err = h.client.PostEphemeral(
		query.channelID,
		string(query.requestMemberID),
		slack.MsgOptionText(rating.NewTopAuthorsTemplate(authorViews, query.period).String(), false),
	)

	return err
}
