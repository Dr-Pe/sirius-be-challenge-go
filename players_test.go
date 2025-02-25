package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine
var w *httptest.ResponseRecorder

func setupTestingSuit() (*sql.DB, *gin.Engine, *httptest.ResponseRecorder) {
	dbConn := setupDatabaseConnection("test" + time.Now().Format("20060102_150405") + ".db")
	router := setupRouter()
	w := httptest.NewRecorder()

	return dbConn, router, w
}

func TestPostPlayer(t *testing.T) {
	dbConn, router, w = setupTestingSuit()
	defer dbConn.Close()

	// Create an example user for testing
	examplePlayer := models.Player{
		Name:    "TestPostPlayer",
		Ranking: 1,
	}
	playerJson, _ := json.Marshal(examplePlayer)
	req, _ := http.NewRequest("POST", "/players", strings.NewReader(string(playerJson)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestPostPlayerWithSameRanking(t *testing.T) {
	dbConn, router, w = setupTestingSuit()
	defer dbConn.Close()

	// Intend to create a player with the same ranking (it should fail)
	examplePlayer := models.Player{
		Name:    "TestPostPlayerWithSameRanking",
		Ranking: 1,
	}
	playerJson, _ := json.Marshal(examplePlayer)
	req, _ := http.NewRequest("POST", "/players", strings.NewReader(string(playerJson)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 409, w.Code)
}

func TestGetPlayers(t *testing.T) {
	dbConn, router, w = setupTestingSuit()
	defer dbConn.Close()

	req, _ := http.NewRequest("GET", "/players", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestGetPlayer(t *testing.T) {
	dbConn, router, w = setupTestingSuit()
	defer dbConn.Close()

	// Get the created user
	req, _ := http.NewRequest("GET", "/players/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestPutPlayer(t *testing.T) {
	dbConn, router, w = setupTestingSuit()
	defer dbConn.Close()

	// Update the created user
	examplePlayer := models.Player{
		Name:    "TestPutPlayer",
		Ranking: 2,
	}
	playerJson, _ := json.Marshal(examplePlayer)
	req, _ := http.NewRequest("PUT", "/players/1", strings.NewReader(string(playerJson)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Get the updated user
	req, _ = http.NewRequest("GET", "/players/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check if the user was updated
	var player models.Player
	json.Unmarshal(w.Body.Bytes(), &player)
	assert.Equal(t, examplePlayer.Name, player.Name)
	assert.Equal(t, examplePlayer.Ranking, player.Ranking)
}

func TestDeletePlayer(t *testing.T) {
	dbConn, router, w = setupTestingSuit()
	defer dbConn.Close()

	// Delete the created user
	req, _ := http.NewRequest("DELETE", "/players/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Check if the user was deleted
	req, _ = http.NewRequest("GET", "/players/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}
