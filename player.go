package main

import (
	"database/sql"
	"fmt"
)

func createPlayersTable(dbConn *sql.DB) (sql.Result, error) {
	return dbConn.Exec("CREATE TABLE IF NOT EXISTS players (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, ranking INTEGER, preferred_cue TEXT, profile_picture_url TEXT)")
}

func selectAllPlayers(dbConn *sql.DB) ([]Player, error) {
	return selectPlayersWhere(dbConn, "SELECT * FROM players")
}

func selectPlayersByName(dbConn *sql.DB, name string) ([]Player, error) {
	return selectPlayersWhere(dbConn, "SELECT * FROM players WHERE name LIKE '%"+name+"%'")
}

func selectPlayersWhere(dbConn *sql.DB, query string) ([]Player, error) {
	players := []Player{}
	rows, err := dbConn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var player Player

		rows.Scan(&player.id, &player.Name, &player.Ranking, &player.PreferredCue, &player.ProfilePictureUrl)
		players = append(players, player)
	}

	return players, nil
}

type Player struct {
	id                int    `uri:"id"`
	Name              string `json:"name" binding:"required"`
	Ranking           int    `json:"ranking"`
	PreferredCue      string `json:"preferredCue"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
}

func (p Player) create(dbConn *sql.DB) (sql.Result, error) {
	rows, err := dbConn.Query("SELECT * FROM players WHERE ranking = ?", p.Ranking)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		return nil, fmt.Errorf("Player with ranking %d already exists", p.Ranking)
	}

	return dbConn.Exec(
		"INSERT INTO players (name, ranking, preferred_cue, profile_picture_url) VALUES (?, ?, ?, ?)",
		p.Name, p.Ranking, p.PreferredCue, p.ProfilePictureUrl,
	)
}
