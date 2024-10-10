package gin_router

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/lib/logger"
	"url-shortener/internal/storage/testdb"
)

func TestPingRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := InitEngine(logger.SetupLogger("local"), testdb.New())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := `{"message":"pong"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestAddAlias_1(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := InitEngine(logger.SetupLogger("local"), testdb.New())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/add", bytes.NewBufferString(`{"url": "some", "alias": "sm.ck"}`))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := `{"message":"some added to storage with sm.ck alias"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestAddAlias_2(t *testing.T) {
	gin.SetMode(gin.TestMode)

	storage := testdb.New()
	storage.SaveURL("some", "sm.ck")
	r := InitEngine(logger.SetupLogger("local"), storage)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/add", bytes.NewBufferString(`{"url": "some", "alias": "sm.ck"}`))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"message":"This url already exists"}`
	assert.JSONEq(t, expected, w.Body.String())
}

func TestGetRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	storage := testdb.New()
	_ = storage.SaveURL("some", "sm.ck")

	r := InitEngine(logger.SetupLogger("local"), storage)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/sm.ck", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := `{"url":"some"}`
	assert.JSONEq(t, expected, w.Body.String())
}
