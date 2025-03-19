package storage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/vpesotskii/go-shortener-url/internal/app/logger"
	"github.com/vpesotskii/go-shortener-url/internal/app/models"
	"go.uber.org/zap"
	"os"
)

type Storage struct {
	db   map[string]models.URL
	file *os.File
	scan *bufio.Scanner
}

func (s *Storage) InsertBatch(records []models.BatchRequest) error {
	//TODO implement me
	panic("implement me")
}

func NewStorage(db map[string]models.URL) *Storage {
	return &Storage{
		db: db,
	}
}

type Repository interface {
	Create(record *models.URL) error
	GetByID(id string) (models.URL, bool)
	Ping() error
	InsertBatch(records []models.BatchRequest) error
}

func (s *Storage) SetFile(f *os.File) {
	s.file = f
}

func (s *Storage) Create(record *models.URL) error {
	s.db[record.ShortURL] = *record
	record.UUID = len(s.db)

	if s.file != nil {
		name := s.file.Name()
		err := saveToFile(*record, name)
		if err != nil {
			logger.Log.Info("Failed to save URL", zap.String("url", name), zap.Error(err))
			return err
		}
	}
	return nil
}

func (s *Storage) GetByID(id string) (models.URL, bool) {
	url, ok := s.db[id]
	fmt.Println("Storage GetByID: ", url.UUID, url.ShortURL, url.OriginalURL, ok)
	return url, ok
}

func (s *Storage) FillFromFile(file *os.File) error {
	url := &models.URL{}
	s.scan = bufio.NewScanner(file)

	for s.scan.Scan() {
		err := json.Unmarshal(s.scan.Bytes(), url)
		if err != nil {
			return err
		}
		s.db[url.ShortURL] = *url
	}

	return nil
}

func (s *Storage) Ping() error {
	return nil
}

func saveToFile(url models.URL, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(url)
	if err != nil {
		return err
	}

	data = append(data, '\n')
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
