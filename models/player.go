package models

import (
	"database/sql"
	"fmt"
	"net/http"
)

func CreatePlayersTable(dbConn *sql.DB) (sql.Result, error) {
	return dbConn.Exec("CREATE TABLE IF NOT EXISTS players (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, ranking INTEGER, preferred_cue TEXT, profile_picture_url TEXT)")
}

func SelectAllPlayers(dbConn *sql.DB) ([]Player, error) {
	return selectPlayersWhere(dbConn, "SELECT * FROM players")
}

func SelectPlayersByName(dbConn *sql.DB, name string) ([]Player, error) {
	return selectPlayersWhere(dbConn, "SELECT * FROM players WHERE name LIKE '%"+name+"%'")
}

func SelectPlayerById(dbConn *sql.DB, id string) (Player, error) {
	players, err := selectPlayersWhere(dbConn, "SELECT * FROM players WHERE id = "+id)
	if err != nil {
		return Player{}, err
	} else if len(players) == 0 {
		return Player{}, PlayerError{http.StatusNotFound, fmt.Sprintf("Player with id %s not found", id)}
	}

	return players[0], nil
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

		rows.Scan(&player.Id, &player.Name, &player.Ranking, &player.PreferredCue, &player.ProfilePictureUrl)
		players = append(players, player)
	}

	return players, nil
}

func UpdatePlayerById(dbConn *sql.DB, id string, player Player) (sql.Result, error) {
	return dbConn.Exec("UPDATE players SET name = ?, ranking = ?, preferred_cue = ?, profile_picture_url = ? WHERE id = ?", player.Name, player.Ranking, player.PreferredCue, player.ProfilePictureUrl, id)
}

func DeletePlayerById(dbConn *sql.DB, id string) (sql.Result, error) {
	return dbConn.Exec("DELETE FROM players WHERE id = ?", id)
}

type PlayerError struct {
	StatusCode int
	Err        string
}

func (e PlayerError) Error() string {
	return e.Err
}

type Player struct {
	Id                int    `json:"id" uri:"id"`
	Name              string `json:"name" binding:"required"`
	Ranking           int    `json:"ranking"` // 0 means no ranking, 1 means the best player
	PreferredCue      string `json:"preferredCue"`
	ProfilePictureUrl string `json:"profilePictureUrl"`
}

func (p Player) Create(dbConn *sql.DB) (sql.Result, error) {
	if p.Ranking != 0 {
		rows, err := dbConn.Query("SELECT * FROM players WHERE ranking = ?", p.Ranking)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		if rows.Next() {
			return nil, PlayerError{http.StatusConflict, fmt.Sprintf("Player with ranking %d already exists", p.Ranking)}
		}
	}

	return dbConn.Exec(
		"INSERT INTO players (name, ranking, preferred_cue, profile_picture_url) VALUES (?, ?, ?, ?)",
		p.Name, p.Ranking, p.PreferredCue, p.ProfilePictureUrl,
	)
}
