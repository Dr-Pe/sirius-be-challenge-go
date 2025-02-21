package main

import (
	"database/sql"
	"fmt"
)

func createPlayersTable(dbConn *sql.DB) (sql.Result, error) {
	return dbConn.Exec("CREATE TABLE IF NOT EXISTS players (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, ranking INTEGER, preferred_cue TEXT, profile_picture_url TEXT)")
}

type Player struct {
	id                  int
	name                string
	ranking             int
	preferred_cue       string
	profile_picture_url string
}

func (p Player) create(name string, ranking int, preferred_cue string, profile_picture_url string) {
	// TODO: Actual query. For noew I just print on stdout what the method receives.
	fmt.Println(name, ranking, preferred_cue, profile_picture_url)
}
