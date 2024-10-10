package testdb

import "url-shortener/internal/storage"

type Storage struct {
	db map[string]string // key - alias, value - url
}

func New() *Storage {
	return &Storage{db: make(map[string]string)}
}

func (s *Storage) GetURL(alias string) (string, error) {
	if v, ok := s.db[alias]; ok {
		return v, nil
	}
	return "", storage.ErrNoRow
}

func (s *Storage) SaveURL(urlToSave, alias string) error {
	if _, ok := s.db[alias]; ok {
		return storage.ErrUrlExists
	}
	s.db[alias] = urlToSave
	return nil
}
