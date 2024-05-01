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

	TopAuthorView struct {
		MemberFullName string
		Score          int
	}

	TopCriterion struct {
		from  time.Time
		to    time.Time
		limit int
	}

	TopFetcher interface {
		FetchTopMemes(criterion TopCriterion) ([]TopMemeView, error)
		FetchTopAuthors(criterion TopCriterion) ([]TopAuthorView, error)
	}

	SQLiteTopMemeFetcher struct {
		connection Connection
		logger     zerolog.Logger
	}
)

func NewTopMemeCriterion(from time.Time, to time.Time, limit int) TopCriterion {
	return TopCriterion{
		from:  from,
		to:    to,
		limit: limit,
	}
}

func NewSQLiteTopFetcher(connection Connection, logger zerolog.Logger) TopFetcher {
	return SQLiteTopMemeFetcher{connection: connection, logger: logger}
}

func (s SQLiteTopMemeFetcher) FetchTopMemes(criterion TopCriterion) ([]TopMemeView, error) {
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

func (s SQLiteTopMemeFetcher) FetchTopAuthors(criterion TopCriterion) ([]TopAuthorView, error) {
	query := `
		SELECT member.full_name, sum(meme.score) AS score FROM memes meme
		LEFT JOIN main.members member on member.id = meme.member_id
		WHERE meme.timestamp BETWEEN ? AND ?
		GROUP BY member.full_name
		ORDER BY score DESC
		LIMIT ?
	`

	rows, err := s.connection.Query(
		query,
		criterion.from.Unix(),
		criterion.to.Unix(),
		criterion.limit,
	)

	if err != nil {
		return []TopAuthorView{}, err
	}

	defer rows.Close()

	result := make([]TopAuthorView, 0, criterion.limit)

	for rows.Next() {
		authorView := TopAuthorView{}

		scanErr := rows.Scan(&authorView.MemberFullName, &authorView.Score)
		if scanErr != nil {
			s.logger.Err(scanErr).Msg("Can not scan row")

			return nil, scanErr
		}

		result = append(result, authorView)
	}

	return result, nil
}
