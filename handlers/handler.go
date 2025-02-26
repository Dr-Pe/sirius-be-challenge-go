package handlers

import "database/sql"

type Handler struct {
	DbConn *sql.DB
}
