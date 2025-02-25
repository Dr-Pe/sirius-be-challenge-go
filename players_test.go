package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/m/v2/models"
	"github.com/stretchr/testify/assert"
)

func TestPostPlayer(t *testing.T) {
	dbConn = setupDatabaseConnection("test.db")
	router := setupRouter()

	w := httptest.NewRecorder()
	// Create an example user for testing
	examplePlayer := models.Player{
		Name:    "test_name",
		Ranking: 0,
	}
	playerJson, _ := json.Marshal(examplePlayer)
	req, _ := http.NewRequest("POST", "/players", strings.NewReader(string(playerJson)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
