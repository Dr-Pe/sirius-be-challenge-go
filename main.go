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
	router.GET("/players", getPlayers)
	router.Run() // listen and serve on 0.0.0.0:8080
}

func postPlayer(ctx *gin.Context) {
	var err error
	var player Player

	err = ctx.ShouldBindJSON(&player)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = player.create(dbConn)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Player created successfully"})
	}
}

func getPlayers(ctx *gin.Context) {
	var err error
	var players []Player
	var query = struct {
		Name string `form:"name"`
	}{}

	err = ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if query.Name != "" {
		players, err = selectPlayersByName(dbConn, query.Name)
	} else {
		players, err = selectAllPlayers(dbConn)
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, players)
}
