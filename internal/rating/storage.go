package rating

import (
	"database/sql"
	"fmt"
	"sync"

	_ "embed"

	"emperror.dev/errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
)

//go:embed schema.sql
var sqliteSchema string

type (
	MemeStorage interface {
		Get(id string) (Meme, error)
		Save(memes ...Meme) error
	}

	InMemoryMemeStorage struct {
		memes map[string]Meme
		mx    sync.RWMutex
	}

	SQLiteMemeStorage struct {
		db     *sql.DB
		logger zerolog.Logger
	}
)

func NewSqliteMemeStorage(dbPath string, logger zerolog.Logger) (SQLiteMemeStorage, func(), error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc&_journal_mode=WAL", dbPath))
	if err != nil {
		logger.Err(err).Msg("Can not open SQLite storage")

		return SQLiteMemeStorage{}, func() {}, err
	}

	storage := SQLiteMemeStorage{
		db:     db,
		logger: logger,
	}

	if err = storage.createTable(); err != nil {
		logger.Err(err).Msg("Can not create table")

		return SQLiteMemeStorage{}, func() {}, err
	}

	logger.Info().Msg("📦 SQLite storage created successfully")

	return storage, func() {
		storage.close()
	}, nil
}

func (s SQLiteMemeStorage) Get(id string) (Meme, error) {
	meme := Meme{}

	row := s.db.QueryRow("SELECT id, channel_id, member_id, score, timestamp FROM memes WHERE id = ?", id)
	err := row.Scan(&meme.id, &meme.channelID, &meme.memberID, &meme.score, &meme.timestamp)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Meme{}, errors.Errorf("Meme with ID %s not found", id)
		}

		s.logger.Err(err).Msg("Can not fetch a Meme")

		return Meme{}, err
	}

	return meme, nil
}

func (s SQLiteMemeStorage) Save(memes ...Meme) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()

			return
		}

		err = tx.Commit()
	}()

	stmt, err := tx.Prepare("INSERT INTO memes(id, channel_id, member_id, score, timestamp) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, meme := range memes {
		_, err = stmt.Exec(meme.id, meme.channelID, meme.memberID, meme.score, meme.timestamp)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SQLiteMemeStorage) createTable() error {
	_, err := s.db.Exec(sqliteSchema)
	return err
}

func (s SQLiteMemeStorage) close() {
	err := s.db.Close()
	if err != nil {
		s.logger.Err(err).Msg("Can not close SQLite storage")
	}

	s.logger.Info().Msg("SQLite storage closed")
}
func NewInMemoryMemeStorage() MemeStorage {
	return &InMemoryMemeStorage{memes: map[string]Meme{}, mx: sync.RWMutex{}}
}

func (i *InMemoryMemeStorage) Get(id string) (Meme, error) {
	i.mx.RLock()
	defer i.mx.RUnlock()

	meme, ok := i.memes[id]
	if !ok {
		return Meme{}, errors.New("Meme not found")
	}

	return meme, nil
}

func (i *InMemoryMemeStorage) Save(memes ...Meme) error {
	i.mx.Lock()
	defer i.mx.Unlock()

	for _, meme := range memes {
		i.memes[meme.id] = meme
	}

	return nil
}
