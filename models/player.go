package models

import (
	"database/sql"
	"fmt"
	"net/http"
)

func CreatePlayersTable(dbConn *sql.DB) (sql.Result, error) {
	return dbConn.Exec("CREATE TABLE IF NOT EXISTS players (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, ranking INTEGER, preferred_cue TEXT, profile_picture_url TEXT, points INTEGER)")
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

		rows.Scan(&player.Id, &player.Name, &player.Ranking, &player.PreferredCue, &player.ProfilePictureUrl, &player.Points)
		players = append(players, player)
	}

	return players, nil
}

func UpdatePlayerById(dbConn *sql.DB, id string, player Player) (sql.Result, error) {
	return dbConn.Exec("UPDATE players SET name = ?, ranking = ?, preferred_cue = ?, profile_picture_url = ?, points = ? WHERE id = ?", player.Name, player.Ranking, player.PreferredCue, player.ProfilePictureUrl, player.Points, id)
}

func DeletePlayerById(dbConn *sql.DB, id string) (sql.Result, error) {
	return dbConn.Exec("DELETE FROM players WHERE id = ?", id)
}

func UpdatePoints(dbConn *sql.DB, winnerId int, loserId int) error {
	var winner Player
	var loser Player
	var err error
	winner, err = SelectPlayerById(dbConn, fmt.Sprintf("%d", winnerId))
	if err != nil {
		return err
	}
	loser, err = SelectPlayerById(dbConn, fmt.Sprintf("%d", loserId))
	if err != nil {
		return err
	}
	fmt.Println("update points", winnerId, loserId)

	return winner.win(dbConn, loser)
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
	Points            int    `json:"points"` // 1 point for each win, 2 points for winning a better player
}

func (p Player) Create(dbConn *sql.DB) (sql.Result, error) {
	return dbConn.Exec(
		"INSERT INTO players (name, ranking, preferred_cue, profile_picture_url) VALUES (?, ?, ?, ?)",
		p.Name, p.Ranking, p.PreferredCue, p.ProfilePictureUrl,
	)
}

func (p Player) win(dbConn *sql.DB, loser Player) error {
	if loser.Points > p.Points {
		p.Points += 2
	} else {
		p.Points++
	}
	_, err := UpdatePlayerById(dbConn, fmt.Sprintf("%d", p.Id), p)

	return err
}
