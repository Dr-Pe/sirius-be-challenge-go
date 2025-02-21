package main

import (
	"database/sql"
	"fmt"
	"net/http"

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

	_, err = createPlayersTable(dbConn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	_, err = createMatchesTable(dbConn)
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.POST("/players", postPlayer)
	router.Run() // listen and serve on 0.0.0.0:8080
}

func postPlayer(ctx *gin.Context) {
	var player Player
	if err := ctx.ShouldBindJSON(&player); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	player.create(dbConn)
}
