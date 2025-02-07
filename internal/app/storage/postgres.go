package storage

import (
	"context"
	"database/sql"
	"github.com/vpesotskii/go-shortener-url/internal/app/logger"
	"github.com/vpesotskii/go-shortener-url/internal/app/models"
	"go.uber.org/zap"
	"time"
)

type DsStorageAdapter struct {
	DB *sql.DB
}

func NewDatabase(dbConfig string) (*DsStorageAdapter, error) {
	db, err := sql.Open("pgx", dbConfig)
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(context.Background(),
		`CREATE TABLE IF NOT EXISTS shorten_urls (
		"uuid" SERIAL PRIMARY KEY,
		"short_url" VARCHAR(50),
		"original_url" TEXT
	)`)
	if err != nil {
		return nil, err
	}

	return &DsStorageAdapter{db}, nil
}

func (db *DsStorageAdapter) Create(record *models.URL) error {
	logger.Log.Info("insert row", zap.String("short", record.ShortURL), zap.String("original", record.OriginalURL))
	_, err := db.DB.ExecContext(context.Background(),
		"INSERT INTO shorten_urls (short_url, original_url) VALUES ($1, $2);",
		record.ShortURL,
		record.OriginalURL)
	if err != nil {
		return err
	}
	return nil
}

func (db *DsStorageAdapter) GetByID(url string) (models.URL, bool) {
	var (
		UUID    int
		ID      string
		FullURL string
	)
	logger.Log.Info("select row", zap.String("url", url))
	row := db.DB.QueryRowContext(context.Background(),
		"SELECT uuid, short_url, original_url FROM shorten_urls WHERE short_url = $1", url)

	err := row.Scan(&UUID, &ID, &FullURL)
	if err != nil {
		return models.URL{}, false
	}
	logger.Log.Info("selected row", zap.String("ID", ID), zap.String("FullURL", FullURL))
	result := models.URL{
		UUID:        UUID,
		OriginalURL: ID,
		ShortURL:    FullURL,
	}
	return result, true
}

func (db *DsStorageAdapter) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.DB.PingContext(ctx)
}
