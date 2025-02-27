package main

import (
	"database/sql"
	"fmt"
	"os"

	"example.com/m/v2/handlers"
	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	var dbConn *sql.DB
	var handler handlers.Handler
	var router *gin.Engine

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	dbConn = setupDatabaseConnection(os.Getenv("DB_NAME"))
	defer dbConn.Close()

	handler = handlers.Handler{DbConn: dbConn}

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

func setupRouter(h handlers.Handler) *gin.Engine {
	router := gin.Default()

	router.POST("/players", h.PostPlayer)
	router.GET("/players", h.GetPlayers)
	router.GET("/players/:id", h.GetPlayer)
	router.PUT("/players/:id", h.PutPlayer)
	router.DELETE("/players/:id", h.DeletePlayer)

	router.POST("/matches", h.PostMatch)
	router.GET("/matches", h.GetMatches)
	router.GET("/matches/:id", h.GetMatch)
	router.PUT("/matches/:id", h.PutMatch)
	router.DELETE("/matches/:id", h.DeleteMatch)

	return router
}
