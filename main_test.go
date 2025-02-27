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
	t.Run("GetPlayers", testGetPlayers)
	t.Run("GetPlayer", testGetPlayer)
	t.Run("PutPlayer", testPutPlayer)
	t.Run("DeletePlayer", testDeletePlayer)
}

func TestMatches(t *testing.T) {
	dbConn, router = setupTestingSuit()
	defer dbConn.Close()

	t.Run("PostMatch", testPostMatch)
	t.Run("GetMatches", testGetMatches)
	t.Run("GetMatch", testGetMatch)
	t.Run("PutMatch", testPutMatch)
	t.Run("DeleteMatch", testDeleteMatch)
}

func testPostPlayer(t *testing.T) {
	// Create an example user for testing
	examplePlayer := models.Player{
		Name: "TestPostPlayer",
	}
	playerJson, _ := json.Marshal(examplePlayer)
	req, _ := http.NewRequest("POST", "/players", strings.NewReader(string(playerJson)))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func testGetPlayers(t *testing.T) {
	// Get all users
	req, _ := http.NewRequest("GET", "/players", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Eval if the created user is in the list
	var players []models.Player
	json.Unmarshal(w.Body.Bytes(), &players)

	assert.Greater(t, len(players), 0)
	assert.Equal(t, "TestPostPlayer", players[0].Name)
}

func testGetPlayer(t *testing.T) {
	// Get the created user
	req, _ := http.NewRequest("GET", "/players/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Eval if the user is the created user
	var player models.Player
	json.Unmarshal(w.Body.Bytes(), &player)

	assert.Equal(t, "TestPostPlayer", player.Name)
}

func testPutPlayer(t *testing.T) {
	// Update the created user
	examplePlayer := models.Player{
		Name: "TestPutPlayer",
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

func testPostMatch(t *testing.T) {
	var w = httptest.NewRecorder()
	// Create two players for testing
	examplePlayer1 := models.Player{
		Name: "TestPostMatch1",
	}
	playerJson, _ := json.Marshal(examplePlayer1)
	req, _ := http.NewRequest("POST", "/players", strings.NewReader(string(playerJson)))
	router.ServeHTTP(w, req)
	examplePlayer2 := models.Player{
		Name: "TestPostMatch2",
	}
	playerJson, _ = json.Marshal(examplePlayer2)
	req, _ = http.NewRequest("POST", "/players", strings.NewReader(string(playerJson)))
	router.ServeHTTP(w, req)

	// Create an example match for testing
	exampleMatch := models.Match{
		Player1id: 2,
		Player2id: 3,
		StartTime: time.Now(),
	}
	matchJson, _ := json.Marshal(exampleMatch)
	req, _ = http.NewRequest("POST", "/matches", strings.NewReader(string(matchJson)))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func testGetMatches(t *testing.T) {
	// Get all matches
	req, _ := http.NewRequest("GET", "/matches", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Eval if the created match is in the list
	var matches []models.Match
	json.Unmarshal(w.Body.Bytes(), &matches)

	assert.Greater(t, len(matches), 0)
	assert.Equal(t, 2, matches[0].Player1id)
}

func testGetMatch(t *testing.T) {
	// Get the created match
	req, _ := http.NewRequest("GET", "/matches/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Eval if the match is the created match
	var match models.Match
	json.Unmarshal(w.Body.Bytes(), &match)

	assert.Equal(t, 2, match.Player1id)
}

func testPutMatch(t *testing.T) {
	// Update the created match
	exampleMatch := models.Match{
		Player1id: 2,
		Player2id: 3,
		StartTime: time.Now(),
		EndTime:   time.Now(),
		WinnerId:  2,
	}
	matchJson, _ := json.Marshal(exampleMatch)
	req, _ := http.NewRequest("PUT", "/matches/1", strings.NewReader(string(matchJson)))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Get the updated match
	req, _ = http.NewRequest("GET", "/matches/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check if the match was updated
	var match models.Match
	json.Unmarshal(w.Body.Bytes(), &match)

	assert.Equal(t, exampleMatch.EndTime.Format("2006-01-02 15:04:05"), match.EndTime.Format("2006-01-02 15:04:05"))

	// Assert that the winner got the points
	req, _ = http.NewRequest("GET", "/players/2", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	var player models.Player
	json.Unmarshal(w.Body.Bytes(), &player)

	assert.Equal(t, 1, player.Points)
}
func testDeleteMatch(t *testing.T) {
	// Delete the created match
	req, _ := http.NewRequest("DELETE", "/matches/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	// Check if the match was deleted
	req, _ = http.NewRequest("GET", "/matches/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}
