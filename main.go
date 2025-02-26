package main

import (
	"database/sql"
	"fmt"

	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

func main() {
	var dbConn *sql.DB
	var handler Handler
	var router *gin.Engine

	dbConn = setupDatabaseConnection("database.db")
	defer dbConn.Close()

	handler = Handler{dbConn: dbConn}

	router = setupRouter(handler)
	router.Run() // listen and serve on 0.0.0.0:8080
}

func setupDatabaseConnection(dbName string) *sql.DB {
	var dbConn *sql.DB
	var err error

	dbConn, err = sql.Open("sqlite", dbName)
	if err != nil {
		panic(err)
	}

	_, err = models.CreatePlayersTable(dbConn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	_, err = models.CreateMatchesTable(dbConn)
	if err != nil {
		panic(err)
	}

	return dbConn
}

func setupRouter(h Handler) *gin.Engine {
	router := gin.Default()

	router.POST("/players", h.postPlayer)
	router.GET("/players", h.getPlayers)
	router.GET("/players/:id", h.getPlayer)
	router.PUT("/players/:id", h.putPlayer)
	router.DELETE("/players/:id", h.deletePlayer)

	router.POST("/matches", h.postMatch)
	router.GET("/matches", h.getMatches)

	return router
}
