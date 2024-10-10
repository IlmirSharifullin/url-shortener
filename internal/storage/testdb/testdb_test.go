package testdb

import (
	"github.com/stretchr/testify/assert"
	"testing"
	_storage "url-shortener/internal/storage"
)

func TestStorage_SaveURL(t *testing.T) {
	storage := New()

	storage.SaveURL("some", "some-test")

	assert.Equal(t, map[string]string{"some-test": "some"}, storage.db)
}

func TestStorage_GetURL(t *testing.T) {
	storage := New()

	storage.SaveURL("some", "some-test")

	url, err := storage.GetURL("some-test")

	assert.NoError(t, err)
	assert.Equal(t, "some", url)
}

func TestStorage_GetURL2(t *testing.T) {
	storage := New()

	url, err := storage.GetURL("some-test")

	assert.ErrorIs(t, err, _storage.ErrNoRow)
	assert.Equal(t, "", url)
}
