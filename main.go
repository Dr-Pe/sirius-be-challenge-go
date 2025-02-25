package main

import (
	"database/sql"
	"fmt"

	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

var dbConn *sql.DB

func main() {
	dbConn = setupDatabaseConnection("database.db")
	defer dbConn.Close()

	router := setupRouter()
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

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.POST("/players", postPlayer)
	router.GET("/players", getPlayers)
	router.GET("/players/:id", getPlayer)
	router.PUT("/players/:id", putPlayer)
	router.DELETE("/players/:id", deletePlayer)

	return router
}
