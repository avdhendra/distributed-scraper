package storage

import (
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPostgresStorage(url string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	logger, _ := zap.NewProduction()
	s := &PostgresStorage{db: db, logger: logger}
	go s.startBackupRoutine()
	return s, nil
}

func (s *PostgresStorage) Save(data interface{}) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	platform := data.(map[string]interface{})["platform"].(string)
	_, err = s.db.Exec(
		"INSERT INTO scraped_data (platform, data, created_at, updated_at) VALUES ($1, $2, $3, $3)",
		platform, dataBytes, time.Now(),
	)
	return err
}

func (s *PostgresStorage) startBackupRoutine() {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		s.logger.Info("performing database backup")
	}
}

func (s *PostgresStorage) Close() {
	s.db.Close()
}