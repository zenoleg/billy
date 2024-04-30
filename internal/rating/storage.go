package rating

import (
	"database/sql"
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
		connection Connection
		logger     zerolog.Logger
	}
)

func NewSqliteMemeStorage(connection Connection, logger zerolog.Logger) (MemeStorage, func(), error) {
	storage := SQLiteMemeStorage{
		connection: connection,

		logger: logger,
	}

	if err := storage.createTable(); err != nil {
		logger.Err(err).Msg("Can not create table")

		return SQLiteMemeStorage{}, func() {}, err
	}

	return storage, func() {
		storage.close()
	}, nil
}

func (s SQLiteMemeStorage) Get(id string) (Meme, error) {
	meme := Meme{}

	row := s.connection.QueryRow("SELECT id, channel_id, member_id, score, timestamp, link FROM memes WHERE id = ?", id)
	err := row.Scan(&meme.id, &meme.channelID, &meme.memberID, &meme.score, &meme.timestamp, &meme.link)

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
	tx, err := s.connection.Begin()
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

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO memes(id, channel_id, member_id, score, timestamp, link) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, meme := range memes {
		_, err = stmt.Exec(meme.id, meme.channelID, meme.memberID, meme.score, meme.timestamp, meme.link)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SQLiteMemeStorage) createTable() error {
	_, err := s.connection.Exec(sqliteSchema)

	return err
}

func (s SQLiteMemeStorage) close() {
	s.connection.Close()
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
