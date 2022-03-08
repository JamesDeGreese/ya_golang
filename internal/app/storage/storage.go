package storage

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/JamesDeGreese/ya_golang/internal/app"
)

type Repository interface {
	GetURL(ID string) (string, error)
	AddURL(ID string, URL string, userID string) error
	GetUserURLs(userID string) []string
}

type Storage struct {
	ShortenURLs map[string]string
	UserLinks   map[string][]string
}

func (s Storage) GetURL(ID string) (string, error) {
	item := s.ShortenURLs[ID]

	if item == "" {
		return "", fmt.Errorf("item not found")
	}

	return s.ShortenURLs[ID], nil
}

func (s Storage) GetUserURLs(userID string) []string {
	return s.UserLinks[userID]
}

func (s Storage) AddURL(ID string, URL string, userID string) error {
	existing, _ := s.GetURL(ID)
	if existing != "" {
		return fmt.Errorf("item with ID %s already exsists", ID)
	}
	s.ShortenURLs[ID] = URL
	_ = append(s.UserLinks[userID], ID)

	return nil
}

func InitStorage(c app.Config) *Storage {
	s := &Storage{
		ShortenURLs: make(map[string]string),
	}

	file, err := os.OpenFile(c.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		return s
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return s
	}

	for _, line := range records {
		s.ShortenURLs[line[0]] = line[1]
	}

	return s
}

func CleanupStorage(c app.Config, s *Storage) {
	file, err := os.OpenFile(c.FileStoragePath, os.O_WRONLY, 0664)
	if err != nil {
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	var pairs [][]string
	for key, value := range s.ShortenURLs {
		pairs = append(pairs, []string{key, value})
	}

	err = writer.WriteAll(pairs)
	if err != nil {
		return
	}

	writer.Flush()
}
