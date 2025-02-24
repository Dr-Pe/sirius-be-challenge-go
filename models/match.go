package models

import (
	"database/sql"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateMatchesTable(dbConn *sql.DB) (sql.Result, error) {
	return dbConn.Exec("CREATE TABLE IF NOT EXISTS matches (id INTEGER PRIMARY KEY AUTOINCREMENT, playder1_id INTEGER, player2_id INTEGER, winner_id INTEGER, table_number INTEGER)")
}

type Match struct {
	id           int
	player1_id   int
	player2_id   int
	start_time   timestamppb.Timestamp
	end_time     timestamppb.Timestamp
	winner_id    int
	table_number int
}
