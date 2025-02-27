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
		return selectMatchesWhere(dbConn, "SELECT * FROM matches WHERE start_time > '"+time.Now().Format("2006-01-02 15:04:05")+"'")
	} else if status == "ongoing" {
		return selectMatchesWhere(dbConn, "SELECT * FROM matches WHERE '"+time.Now().Format("2006-01-02 15:04:05")+"' BETWEEN start_time AND end_time")
	} else if status == "finished" {
		return selectMatchesWhere(dbConn, "SELECT * FROM matches WHERE end_time < '"+time.Now().Format("2006-01-02 15:04:05")+"'")
	} else {
		return nil, MatchError{StatusCode: http.StatusBadRequest, Err: "Invalid status"}
	}
}

func SelectMatchById(dbConn *sql.DB, id string) (Match, error) {
	matches, err := selectMatchesWhere(dbConn, "SELECT * FROM matches WHERE id = "+id)
	if err != nil {
		return Match{}, err
	} else if len(matches) == 0 {
		return Match{}, MatchError{http.StatusNotFound, fmt.Sprintf("Match with id %s not found", id)}
	}

	return matches[0], nil
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

		rows.Scan(&match.Id, &match.Player1id, &match.Player2id, &match.StartTime, &match.EndTime, &match.WinnerId, &match.TableNumber)
		matches = append(matches, match)
	}

	return matches, nil
}

func UpdateMatchById(dbConn *sql.DB, id string, match Match) (sql.Result, error) {
	return dbConn.Exec("UPDATE matches SET player1_id = ?, player2_id = ?, start_time = ?, end_time = ?, winner_id = ?, table_number = ? WHERE id = ?", match.Player1id, match.Player2id, match.StartTime, match.EndTime, match.WinnerId, match.TableNumber, id)
}

func DeleteMatchById(dbConn *sql.DB, id string) (sql.Result, error) {
	return dbConn.Exec("DELETE FROM matches WHERE id = ?", id)
}

type MatchError struct {
	StatusCode int
	Err        string
}

func (e MatchError) Error() string {
	return e.Err
}

type Match struct {
	Id          int       `json:"id" uri:"id"`
	Player1id   int       `json:"player1id" binding:"required"`
	Player2id   int       `json:"player2id" binding:"required"`
	StartTime   time.Time `json:"startTime" binding:"required"`
	EndTime     time.Time `json:"endTime"`
	WinnerId    int       `json:"winnerId"`
	TableNumber int       `json:"tableNumber"`
}

func (m *Match) Create(dbConn *sql.DB) (sql.Result, error) {
	if m.Player1id == m.Player2id {
		return nil, MatchError{StatusCode: http.StatusBadRequest, Err: "Player1 and Player2 must be different"}
	} else if _, err := SelectPlayerById(dbConn, fmt.Sprintf("%d", m.Player1id)); err != nil {
		return nil, MatchError{StatusCode: http.StatusBadRequest, Err: "Player1 does not exist"}
	} else if _, err := SelectPlayerById(dbConn, fmt.Sprintf("%d", m.Player2id)); err != nil {
		return nil, MatchError{StatusCode: http.StatusBadRequest, Err: "Player2 does not exist"}
	} else if m.EndTime == (time.Time{}) {
		m.EndTime = m.StartTime.Add(time.Hour)
	}
	tx, err := dbConn.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query("SELECT * FROM matches WHERE table_number = ? AND ((start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?))", m.TableNumber, m.StartTime, m.EndTime, m.StartTime, m.EndTime)
	if err != nil {
		tx.Rollback()
		return nil, err
	} else if rows.Next() {
		tx.Rollback()
		return nil, MatchError{StatusCode: http.StatusBadRequest, Err: "Table already booked"}
	}
	rows, err = tx.Query("SELECT * FROM matches WHERE (player1_id = ? OR player2_id = ?) AND ((start_time BETWEEN ? AND ?) OR (end_time BETWEEN ? AND ?))", m.Player1id, m.Player1id, m.StartTime, m.EndTime, m.StartTime, m.EndTime)
	if err != nil {
		tx.Rollback()
		return nil, err
	} else if rows.Next() {
		tx.Rollback()
		return nil, MatchError{StatusCode: http.StatusBadRequest, Err: "Players already booked"}
	}
	res, err := tx.Exec("INSERT INTO matches (player1_id, player2_id, start_time, end_time, table_number) VALUES (?, ?, ?, ?, ?)", m.Player1id, m.Player2id, m.StartTime, m.EndTime, m.TableNumber)
	tx.Commit()

	return res, err
}
