package storage

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/JamesDeGreese/ya_golang/internal/app"
	"github.com/jackc/pgx/v4"
)

type Repository interface {
	GetURL(ID string) (string, error)
	AddURL(ID string, URL string, userID string) error
	GetUserURLs(userID string) []string
	DB() *pgx.Conn
}

type Storage struct {
	ShortenURLs map[string]string
	UserLinks   map[string][]string
	DBConn      *pgx.Conn
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
	userURLs := s.GetUserURLs(userID)
	if len(userURLs) == 0 {
		s.UserLinks[userID] = make([]string, 0)
	}
	s.UserLinks[userID] = append(s.UserLinks[userID], ID)

	return nil
}

func (s Storage) DB() *pgx.Conn {
	return s.DBConn
}

func InitStorage(c app.Config) *Storage {
	conn, err := pgx.Connect(context.Background(), c.DatabaseDSN)
	defer conn.Close(context.Background())

	s := &Storage{
		ShortenURLs: make(map[string]string),
		UserLinks:   make(map[string][]string),
		DBConn:      conn,
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
