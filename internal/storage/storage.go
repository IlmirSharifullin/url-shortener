package storage

import "errors"

type Storage interface {
	GetURL(alias string) (string, error)
	SaveURL(urlToSave, alias string) error
}

var (
	ErrUrlExists = errors.New("This url already exists")
	ErrNoRow     = errors.New("This url not exists")
)
