package models

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

func CreateMatchesTable(dbConn *sql.DB) (sql.Result, error) {
	return dbConn.Exec("CREATE TABLE IF NOT EXISTS matches (id INTEGER PRIMARY KEY AUTOINCREMENT, player1_id INTEGER, player2_id INTEGER, start_time DATETIME, end_time DATETIME, winner_id INTEGER, table_number INTEGER)")
}

func SelectAllMatches(dbConn *sql.DB) ([]Match, error) {
	return selectMatchesWhere(dbConn, "SELECT * FROM matches")
}

func SelectMatchesByStatus(dbConn *sql.DB, status string) ([]Match, error) {
	if status == "upcoming" {
		return selectMatchesWhere(dbConn, "SELECT * FROM matches WHERE start_time > "+time.Now().Format("2006-01-02 15:04:05"))
	} else if status == "ongoing" {
		return selectMatchesWhere(dbConn, "SELECT * FROM matches WHERE "+time.Now().Format("2006-01-02 15:04:05")+" BETWEEN start_time AND end_time")
	} else if status == "finished" {
		return selectMatchesWhere(dbConn, "SELECT * FROM matches WHERE end_time < "+time.Now().Format("2006-01-02 15:04:05"))
	} else {
		return nil, MatchError{StatusCode: http.StatusBadRequest, Err: "Invalid status"}
	}
}

func selectMatchesWhere(dbConn *sql.DB, query string) ([]Match, error) {
	matches := []Match{}
	rows, err := dbConn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var match Match

		rows.Scan(&match.id, &match.Player1id, &match.Player2id, &match.StartTime, &match.endTime, &match.winnerId, &match.tableNumber)
		matches = append(matches, match)
	}

	return matches, nil
}

type MatchError struct {
	StatusCode int
	Err        string
}

func (e MatchError) Error() string {
	return e.Err
}

type Match struct {
	id          int
	Player1id   int       `json:"player1id" binding:"required"`
	Player2id   int       `json:"player2id" binding:"required"`
	StartTime   time.Time `json:"startTime" binding:"required"`
	endTime     time.Time
	winnerId    int
	tableNumber int
}

func (m *Match) Create(dbConn *sql.DB) (sql.Result, error) {
	var err error

	if m.Player1id == m.Player2id {
		return nil, MatchError{StatusCode: http.StatusBadRequest, Err: "Player1 and Player2 must be different"}
	}
	if _, err = SelectPlayerById(dbConn, fmt.Sprintf("%d", m.Player1id)); err != nil {
		return nil, MatchError{StatusCode: http.StatusBadRequest, Err: "Player1 does not exist"}
	}
	if _, err = SelectPlayerById(dbConn, fmt.Sprintf("%d", m.Player2id)); err != nil {
		return nil, MatchError{StatusCode: http.StatusBadRequest, Err: "Player2 does not exist"}
	}
	if m.endTime == (time.Time{}) {
		m.endTime = m.StartTime.Add(time.Hour)
	}
	// TODO: Implement double-booking check
	return dbConn.Exec("INSERT INTO matches (player1_id, player2_id, start_time, table_number) VALUES (?, ?, ?, ?)", m.Player1id, m.Player2id, m.StartTime, m.tableNumber)
}
