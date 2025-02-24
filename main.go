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
	var err error

	dbConn, err = sql.Open("sqlite", "database.db")
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	_, err = models.CreatePlayersTable(dbConn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	_, err = models.CreateMatchesTable(dbConn)
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.POST("/players", postPlayer)
	router.GET("/players", getPlayers)
	router.GET("/players/:id", getPlayer)
	router.PUT("/players/:id", putPlayer)
	router.DELETE("/players/:id", deletePlayer)

	router.Run() // listen and serve on 0.0.0.0:8080
}
