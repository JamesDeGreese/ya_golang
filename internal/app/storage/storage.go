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
	GetUserURLs(userID string) []ShortLink
	CleanUp(c app.Config)
}

type MemoryStorage struct {
	ShortenURLs map[string]string
	UserLinks   map[string][]string
}

type DBStorage struct {
	DBConn *pgx.Conn
}

type ShortLink struct {
	ID          string
	OriginalURL string
	UserID      string
}

func (s MemoryStorage) GetURL(ID string) (string, error) {
	item := s.ShortenURLs[ID]

	if item == "" {
		return "", fmt.Errorf("item not found")
	}

	return s.ShortenURLs[ID], nil
}

func (s MemoryStorage) GetUserURLs(userID string) []ShortLink {
	var res []ShortLink
	userURLs := s.UserLinks[userID]
	if len(userURLs) == 0 {
		return res
	}

	for _, shortID := range userURLs {
		URL, _ := s.GetURL(shortID)
		res = append(res, ShortLink{
			shortID,
			URL,
			userID,
		})
	}

	return res
}

func (s MemoryStorage) AddURL(ID string, URL string, userID string) error {
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

func (s MemoryStorage) CleanUp(c app.Config) {
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

type ShortenURLEntity struct {
	ShortURL    string
	OriginalURL string
	UserID      string
}

func (s DBStorage) GetURL(ID string) (string, error) {
	var res string
	err := s.DBConn.QueryRow(context.Background(), "SELECT original_url FROM shorten_urls WHERE id = $1", ID).Scan(&res)
	if err == pgx.ErrNoRows {
		return res, nil
	}
	if err != nil {
		return "", err
	}

	return res, nil
}

func (s DBStorage) GetUserURLs(userID string) []ShortLink {
	res := make([]ShortLink, 0)
	rows, err := s.DBConn.Query(context.Background(), "SELECT id, original_url FROM shorten_urls WHERE user_id = $1", userID)
	if err != nil {
		return res
	}
	defer rows.Close()

	for rows.Next() {
		var r ShortLink
		err := rows.Scan(&r.ID, &r.OriginalURL)
		if err != nil {
			return nil
		}
		res = append(res, r)
	}

	return res
}

func (s DBStorage) AddURL(ID string, URL string, userID string) error {
	_, err := s.DBConn.Exec(context.Background(), "INSERT INTO shorten_urls VALUES ($1, $2, $3)", ID, URL, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s DBStorage) CleanUp(c app.Config) {
	err := s.DBConn.Close(context.Background())
	if err != nil {
		return
	}
}

func InitStorage(c app.Config) Repository {
	if c.DatabaseDSN != "" {
		conn, err := pgx.Connect(context.Background(), c.DatabaseDSN)
		if err != nil {
			return initMemoryStorage(c)
		}
		dbSt := &DBStorage{
			DBConn: conn,
		}

		_, err = dbSt.DBConn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS shorten_urls (id varchar(36), original_url varchar(255), user_id varchar(36));")
		if err != nil {
			return initMemoryStorage(c)
		}

		return dbSt
	}

	return initMemoryStorage(c)
}

func initMemoryStorage(c app.Config) Repository {
	memSt := &MemoryStorage{
		ShortenURLs: make(map[string]string),
		UserLinks:   make(map[string][]string),
	}

	file, err := os.OpenFile(c.FileStoragePath, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		return memSt
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return memSt
	}

	for _, line := range records {
		memSt.ShortenURLs[line[0]] = line[1]
	}

	return memSt
}
