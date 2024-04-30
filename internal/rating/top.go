package rating

import (
	"database/sql"
	"time"

	"github.com/rs/zerolog"
)

type (
	TopMemeView struct {
		ID        string
		ChannelID string
		MemberID  MemberID
		Score     int
	}

	TopMemeCriterion struct {
		from  time.Time
		to    time.Time
		limit int
	}

	TopMemeFetcher interface {
		Fetch(criterion TopMemeCriterion) ([]TopMemeView, error)
	}

	SQLiteTopMemeFetcher struct {
		db     *sql.DB
		logger zerolog.Logger
	}
)

func NewTopMemeCriterion(from time.Time, to time.Time, limit int) TopMemeCriterion {
	return TopMemeCriterion{
		from:  from,
		to:    to,
		limit: limit,
	}
}
