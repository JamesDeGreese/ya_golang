package storage

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/JamesDeGreese/ya_golang/internal/app"
)

type Repository interface {
	Get(ID string) (string, error)
	Add(ID string, URL string) error
}

type Storage struct {
	List map[string]string
}

func (s Storage) Get(ID string) (string, error) {
	item := s.List[ID]

	if item == "" {
		return "", fmt.Errorf("item not found")
	}

	return s.List[ID], nil
}

func (s Storage) Add(ID string, URL string) error {
	existing, _ := s.Get(ID)
	if existing != "" {
		return fmt.Errorf("item with ID %s already exsists", ID)
	}
	s.List[ID] = URL

	return nil
}

func ConstructStorage(c app.Config) *Storage {
	s := &Storage{
		List: make(map[string]string),
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
		s.List[line[0]] = line[1]
	}

	return s
}

func DestructStorage(c app.Config, s *Storage) {
	file, err := os.OpenFile(c.FileStoragePath, os.O_WRONLY, 0664)
	if err != nil {
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	var pairs [][]string
	for key, value := range s.List {
		pairs = append(pairs, []string{key, value})
	}

	err = writer.WriteAll(pairs)
	if err != nil {
		return
	}

	writer.Flush()
}
