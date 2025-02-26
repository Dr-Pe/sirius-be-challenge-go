package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/m/v2/handlers"
	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var dbConn *sql.DB
var router *gin.Engine

func setupTestingSuit() (*sql.DB, *gin.Engine) {
	dbConn := setupDatabaseConnection("test" + time.Now().Format("20060102_150405") + ".db")
	handler := handlers.Handler{DbConn: dbConn}
	router := setupRouter(handler)

	return dbConn, router
}

func TestPlayers(t *testing.T) {
	dbConn, router = setupTestingSuit()
	defer dbConn.Close()

	t.Run("PostPlayer", testPostPlayer)
	t.Run("PostPlayerWithSameRanking", testPostPlayerWithSameRanking)
	t.Run("GetPlayers", testGetPlayers)
	t.Run("GetPlayer", testGetPlayer)
	t.Run("PutPlayer", testPutPlayer)
	t.Run("DeletePlayer", testDeletePlayer)
}

func testPostPlayer(t *testing.T) {
	// Create an example user for testing
	examplePlayer := models.Player{
		Name:    "TestPostPlayer",
		Ranking: 1,
	}
	playerJson, _ := json.Marshal(examplePlayer)
	req, _ := http.NewRequest("POST", "/players", strings.NewReader(string(playerJson)))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func testPostPlayerWithSameRanking(t *testing.T) {
	// Intend to create a player with the same ranking (it should fail)
	examplePlayer := models.Player{
		Name:    "TestPostPlayerWithSameRanking",
		Ranking: 1,
	}
	playerJson, _ := json.Marshal(examplePlayer)
	req, _ := http.NewRequest("POST", "/players", strings.NewReader(string(playerJson)))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 409, w.Code)
}

func testGetPlayers(t *testing.T) {
	// Get all users
	req, _ := http.NewRequest("GET", "/players", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func testGetPlayer(t *testing.T) {
	// Get the created user
	req, _ := http.NewRequest("GET", "/players/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func testPutPlayer(t *testing.T) {
	// Update the created user
	examplePlayer := models.Player{
		Name:    "TestPutPlayer",
		Ranking: 2,
	}
	playerJson, _ := json.Marshal(examplePlayer)
	req, _ := http.NewRequest("PUT", "/players/1", strings.NewReader(string(playerJson)))
	w := httptest.NewRecorder()
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

func testDeletePlayer(t *testing.T) {
	// Delete the created user
	req, _ := http.NewRequest("DELETE", "/players/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Check if the user was deleted
	req, _ = http.NewRequest("GET", "/players/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}
