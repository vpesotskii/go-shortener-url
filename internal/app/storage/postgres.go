package storage

import (
	"context"
	"database/sql"
	"github.com/vpesotskii/go-shortener-url/internal/app/models"
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

	row := db.DB.QueryRowContext(context.Background(),
		"SELECT uuid, short_url, original_url FROM shorten_urls WHERE short_url = $1", url)

	err := row.Scan(&UUID, &ID, &FullURL)
	if err != nil {
		return models.URL{}, false
	}

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
