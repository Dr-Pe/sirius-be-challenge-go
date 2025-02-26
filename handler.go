package main

import "database/sql"

type Handler struct {
	dbConn *sql.DB
}
