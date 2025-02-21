package main

import (
	"database/sql"
)

func createPlayersTable(dbConn *sql.DB) (sql.Result, error) {
	return dbConn.Exec("CREATE TABLE IF NOT EXISTS players (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, ranking INTEGER, preferred_cue TEXT, profile_picture_url TEXT)")
}

type Player struct {
	id                int
	Name              string `json:"name" binding:"required"`
	Ranking           int    `json:"ranking"`
	PreferredCue      string `json:"preferredCue"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
}

func (p Player) create(dbConn *sql.DB) (sql.Result, error) {
	return dbConn.Exec("")
}
