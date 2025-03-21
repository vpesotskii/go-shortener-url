package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
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
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}
	logger.Log.Info("insert row", zap.String("short", record.ShortURL), zap.String("original", record.OriginalURL))
	_, err = db.DB.Exec(`INSERT INTO shorten_urls (short_url, original_url) VALUES ($1, $2);`,
		record.ShortURL,
		record.OriginalURL)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *DsStorageAdapter) GetByID(url string) (models.URL, bool) {
	var (
		UUID    int
		ID      string
		FullURL string
	)
	logger.Log.Info("select row", zap.String("url", url))
	row := db.DB.QueryRow(
		`SELECT uuid, short_url, original_url FROM shorten_urls WHERE short_url = $1`, url)

	err := row.Scan(&UUID, &ID, &FullURL)
	if err != nil {
		logger.Log.Debug(err.Error())
		return models.URL{}, false
	}
	logger.Log.Info("selected row", zap.String("ID", ID), zap.String("FullURL", FullURL))
	result := models.URL{
		UUID:        UUID,
		OriginalURL: FullURL,
		ShortURL:    ID,
	}
	return result, true
}

func (db *DsStorageAdapter) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return db.DB.PingContext(ctx)
}

func (db *DsStorageAdapter) InsertBatch(records []models.BatchRequest) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return err
	}

	for _, record := range records {
		logger.Log.Info("insert row from Batch", zap.String("short", record.ShortURL), zap.String("original", record.OriginalURL))
		_, err = db.DB.Exec(`INSERT INTO shorten_urls (short_url, original_url) VALUES ($1, $2);`,
			record.ShortURL,
			record.OriginalURL)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
