package storage

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"

	"github.com/JamesDeGreese/ya_golang/internal/app"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

type Repository interface {
	GetURLByID(ID string) (string, error)
	GetURLByOriginalURL(OriginalURL string) (string, error)
	AddURL(link ShortLink) error
	AddURLBatch(links []ShortLink) error
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

type RecordDuplicateError struct {
	param string
	value string
}

func (e *RecordDuplicateError) Error() string {
	return fmt.Sprintf("Record with same param %s with value %s already exists", e.param, e.value)
}

func (s MemoryStorage) GetURLByID(ID string) (string, error) {
	item := s.ShortenURLs[ID]

	if item == "" {
		return "", fmt.Errorf("item not found")
	}

	return item, nil
}

func (s MemoryStorage) GetURLByOriginalURL(OriginalURL string) (string, error) {
	rev := make(map[string]string, len(s.ShortenURLs))
	for ID, URL := range s.ShortenURLs {
		rev[URL] = ID
	}

	item := rev[OriginalURL]

	if item == "" {
		return "", fmt.Errorf("item not found")
	}

	return item, nil
}

func (s MemoryStorage) GetUserURLs(userID string) []ShortLink {
	var res []ShortLink
	userURLs := s.UserLinks[userID]
	if len(userURLs) == 0 {
		return res
	}

	for _, shortID := range userURLs {
		URL, _ := s.GetURLByID(shortID)
		res = append(res, ShortLink{
			shortID,
			URL,
			userID,
		})
	}

	return res
}

func (s MemoryStorage) AddURL(link ShortLink) error {
	existing, _ := s.GetURLByOriginalURL(link.OriginalURL)
	if existing != "" {
		return &RecordDuplicateError{param: "OriginalID", value: link.OriginalURL}
	}
	s.ShortenURLs[link.ID] = link.OriginalURL
	userURLs := s.GetUserURLs(link.UserID)
	if len(userURLs) == 0 {
		s.UserLinks[link.UserID] = make([]string, 0)
	}
	s.UserLinks[link.UserID] = append(s.UserLinks[link.UserID], link.ID)

	return nil
}

func (s MemoryStorage) AddURLBatch(links []ShortLink) error {
	for _, link := range links {
		err := s.AddURL(link)
		if err != nil {
			return err
		}
	}
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

func (s DBStorage) GetURLByID(ID string) (string, error) {
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

func (s DBStorage) GetURLByOriginalURL(OriginalURL string) (string, error) {
	var res string
	err := s.DBConn.QueryRow(context.Background(), "SELECT id FROM shorten_urls WHERE original_url = $1", OriginalURL).Scan(&res)
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

func (s DBStorage) AddURL(link ShortLink) error {
	_, err := s.DBConn.Exec(context.Background(), "INSERT INTO shorten_urls VALUES ($1, $2, $3)", link.ID, link.OriginalURL, link.UserID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return &RecordDuplicateError{param: "original_url", value: link.OriginalURL}
			}
		}
		return err
	}

	return nil
}

func (s DBStorage) AddURLBatch(links []ShortLink) error {
	rows := make([][]interface{}, 0)
	for _, link := range links {
		rows = append(rows, []interface{}{link.ID, link.OriginalURL, link.UserID})
	}
	_, err := s.DBConn.CopyFrom(
		context.Background(),
		pgx.Identifier{"shorten_urls"},
		[]string{"id", "original_url", "user_id"},
		pgx.CopyFromRows(rows),
	)

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
		_, err = dbSt.DBConn.Exec(context.Background(), "CREATE UNIQUE INDEX IF NOT EXISTS original_url_idx ON shorten_urls (original_url);")
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
