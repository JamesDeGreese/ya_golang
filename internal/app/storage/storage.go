package storage

type Repository interface {
	Get(ID string) string
	Add(ID string, URL string)
}

var instance *Storage = nil

type Storage struct {
	List map[string]string
}

func (s Storage) Get(ID string) string {
	return s.List[ID]
}

func (s Storage) Add(ID string, URL string) {
	s.List[ID] = URL
	return
}

func GetInstance() *Storage {
	if instance == nil {
		instance = &Storage{
			List: make(map[string]string),
		}
	}

	return instance
}
