package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Storage struct {
	config *DbConfig
	db     *sql.DB
}

func New(config *DbConfig) *Storage {
	return &Storage{
		config: config,
	}
}

func (s *Storage) Open() error {
	db, err := sql.Open("postgres", s.config.DbURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Storage) Query(query string, args ...any) (*sql.Rows, error) {
	return s.db.Query(query, args...)
}

func (s *Storage) Exec(query string, args ...any) (sql.Result, error) {
	return s.db.Exec(query, args...)
}

func (s *Storage) Close() {
	s.db.Close()
}
