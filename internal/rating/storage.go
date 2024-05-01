package rating

import (
	"database/sql"
	_ "embed"

	"emperror.dev/errors"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog"
)

//go:embed meme_schema.sql
var memeSchema string

//go:embed member_schema.sql
var memberSchema string

type (
	MemeStorage interface {
		Get(id string) (Meme, error)
		Save(memes ...Meme) error
	}

	MemberStorage interface {
		Get(id MemberID) (Member, error)
		Save(members ...Member) error
	}

	SQLiteMemeStorage struct {
		connection Connection
		logger     zerolog.Logger
	}

	SQLiteMemberStorage struct {
		connection Connection
		logger     zerolog.Logger
	}
)

func NewSqliteMemeStorage(connection Connection, logger zerolog.Logger) (MemeStorage, error) {
	storage := SQLiteMemeStorage{
		connection: connection,

		logger: logger,
	}

	if err := storage.createTable(); err != nil {
		logger.Err(err).Msg("Can not create table")

		return SQLiteMemeStorage{}, err
	}

	return storage, nil
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
	_, err := s.connection.Exec(memeSchema)

	return err
}

func NewSqliteMemberStorage(connection Connection, logger zerolog.Logger) (MemberStorage, error) {
	storage := SQLiteMemberStorage{
		connection: connection,

		logger: logger,
	}

	if err := storage.createTable(); err != nil {
		logger.Err(err).Msg("Can not create table")

		return SQLiteMemberStorage{}, err
	}

	return storage, nil
}

func (s SQLiteMemberStorage) Get(id MemberID) (Member, error) {
	member := Member{}

	row := s.connection.QueryRow("SELECT id, full_name, display_name FROM members WHERE id = ?", id)
	err := row.Scan(&member.id, &member.fullName, &member.displayName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Member{}, errors.Errorf("Member with ID %s not found", id)
		}

		s.logger.Err(err).Msg("Can not fetch a Member")

		return Member{}, err
	}

	return member, nil
}

func (s SQLiteMemberStorage) Save(members ...Member) error {
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

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO members(id, full_name, display_name) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, member := range members {
		_, err = stmt.Exec(member.id, member.fullName, member.displayName)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s SQLiteMemberStorage) createTable() error {
	_, err := s.connection.Exec(memberSchema)

	return err
}
