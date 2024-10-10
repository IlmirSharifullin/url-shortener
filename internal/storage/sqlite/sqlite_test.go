package sqlite

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	_storage "url-shortener/internal/storage"
)

func TestNew(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Expectation for creating the table
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS url").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("CREATE INDEX IF NOT EXISTS idx_alias").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	storage, err := New(":memory:")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Assert
	assert.NotNil(t, storage)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveURL(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	// Expectation for inserting a new URL
	mock.ExpectPrepare("INSERT INTO url").ExpectExec().
		WithArgs("https://example.com", "example").
		WillReturnResult(sqlmock.NewResult(1, 1))

	// Act
	err = storage.SaveURL("https://example.com", "example")

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveURL_AlreadyExists(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	mock.ExpectPrepare("INSERT INTO url").ExpectExec().
		WithArgs("https://example.com", "example")

	_ = storage.SaveURL("https://example.com", "example")

	//s, err := storage.GetURL("example")
	//slog.Info(s, err)
	// Expectation for inserting a URL that already exists
	mock.ExpectPrepare("INSERT INTO url").ExpectExec().
		WithArgs("https://example.com", "example").
		WillReturnError(&sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintUnique})

	err = storage.SaveURL("https://example.com", "example")

	assert.ErrorIs(t, err, _storage.ErrUrlExists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetURL(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	// Expectation for selecting a URL
	mock.ExpectPrepare("SELECT url FROM url WHERE alias=?").ExpectQuery().
		WithArgs("example").
		WillReturnRows(sqlmock.NewRows([]string{"url"}).AddRow("https://example.com"))

	// Act
	url, err := storage.GetURL("example")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "https://example.com", url)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetURL_NotFound(t *testing.T) {
	// Arrange
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	storage := &Storage{db: db}

	// Expectation for selecting a URL that doesn't exist
	mock.ExpectPrepare("SELECT url FROM url WHERE alias=?").ExpectQuery().
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	// Act
	url, err := storage.GetURL("nonexistent")

	// Assert
	assert.ErrorIs(t, err, _storage.ErrNoRow)
	assert.Empty(t, url)
	assert.NoError(t, mock.ExpectationsWereMet())
}
