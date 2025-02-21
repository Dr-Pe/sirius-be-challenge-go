package main

import (
	"database/sql"
)

func createPlayersTable(dbConn *sql.DB) (sql.Result, error) {
	return dbConn.Exec("CREATE TABLE IF NOT EXISTS players (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, ranking INTEGER, preferred_cue TEXT, profile_picture_url TEXT)")
}

// func selectPlayerById(dbConn *sql.DB, id int) (Player, error) {
// 	if row, err := dbConn.Exec("SELECT * FROM players WHERE id = ?", id); err != nil {
// 		return nil, err
// 	}
// }

type Player struct {
	id                int    `uri:"id"`
	Name              string `json:"name" binding:"required"`
	Ranking           int    `json:"ranking"`
	PreferredCue      string `json:"preferredCue"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
}

func (p Player) create(dbConn *sql.DB) (sql.Result, error) {
	return dbConn.Exec(
		"INSERT INTO players (name, ranking, preferred_cue, profile_picture_id) VALUES (?)",
		p.Name, p.Ranking, p.PreferredCue, p.ProfilePictureUrl,
	)
}
