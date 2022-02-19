package storage

import (
	"fmt"
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

func ConstructStorage() *Storage {
	return &Storage{
		List: make(map[string]string),
	}
}
