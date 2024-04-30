package rating

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"
)

type Connection struct {
	*sql.DB
	Close func()
}

func NewConnection(dbPath string, logger zerolog.Logger) (Connection, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc&_journal_mode=WAL", dbPath))
	if err != nil {
		logger.Err(err).Msg("Can not open SQLite connection")

		return Connection{}, err
	}

	logger.Info().Msg("📦 SQLite connection created successfully")

	return Connection{
		DB: db,
		Close: func() {
			err := db.Close()
			if err != nil {
				logger.Err(err).Msg("Can not close SQLite connection")
			}

			logger.Info().Msg("SQLite connection closed")
		},
	}, nil
}
