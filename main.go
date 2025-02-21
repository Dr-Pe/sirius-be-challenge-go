package main

import (
	"database/sql"
	"fmt"

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
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listen and serve on 0.0.0.0:8080
}
