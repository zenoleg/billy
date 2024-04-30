package rating

import (
	"time"

	"github.com/rs/zerolog"
)

type (
	TopMemeView struct {
		Link           string
		MemberFullName string
		Score          int
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
		connection Connection
		logger     zerolog.Logger
	}
)

func NewTopMemeCriterion(from time.Time, to time.Time, limit int) TopMemeCriterion {
	return TopMemeCriterion{
		from:  from,
		to:    to,
		limit: limit,
	}
}

func NewSQLiteTopMemeFetcher(connection Connection, logger zerolog.Logger) TopMemeFetcher {
	return SQLiteTopMemeFetcher{connection: connection, logger: logger}
}

func (s SQLiteTopMemeFetcher) Fetch(criterion TopMemeCriterion) ([]TopMemeView, error) {
	rows, err := s.connection.Query(
		"SELECT meme.link, COALESCE(member.full_name, 'Неизвестно кто'), meme.score FROM memes meme LEFT JOIN members member ON meme.member_id = member.id WHERE timestamp BETWEEN ? AND ? ORDER BY score DESC LIMIT ?",
		criterion.from.Unix(),
		criterion.to.Unix(),
		criterion.limit,
	)

	if err != nil {
		return []TopMemeView{}, err
	}

	defer rows.Close()

	result := make([]TopMemeView, 0, criterion.limit)

	for rows.Next() {
		memeView := TopMemeView{}

		scanErr := rows.Scan(&memeView.Link, &memeView.MemberFullName, &memeView.Score)
		if scanErr != nil {
			s.logger.Err(scanErr).Msg("Can not scan row")

			return nil, scanErr
		}

		result = append(result, memeView)
	}

	return result, nil
}
