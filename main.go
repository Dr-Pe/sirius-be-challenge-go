package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"example.com/m/v2/handlers"
	"example.com/m/v2/models"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

func main() {
	var err error
	var handler handlers.Handler
	var router *gin.Engine

	// Dotenv
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}

	// Database
	dbConn := setupDatabaseConnection(os.Getenv("DB_NAME"))
	defer dbConn.Close()

	// AWS S3
	s3Client := setupS3Client(os.Getenv("AWS_REGION"))

	// Handler
	handler = handlers.Handler{DbConn: dbConn, S3Client: s3Client, BucketName: os.Getenv("AWS_BUCKET_NAME"), Region: os.Getenv("AWS_REGION")}
	err = handler.CreateBucket(context.TODO())
	if err != nil {
		fmt.Println("Error creating bucket")
		panic(err)
	}

	// Router
	router = setupRouter(handler)
	router.Run() // listen and serve on 0.0.0.0:8080
}

func setupDatabaseConnection(dbName string) *sql.DB {
	var dbConn *sql.DB
	var err error

	dbConn, err = sql.Open("sqlite", dbName)
	if err != nil {
		fmt.Println("Could not open database connection")
		panic(err)
	}

	_, err = models.CreatePlayersTable(dbConn)
	if err != nil {
		fmt.Println("Could not create players table")
		panic(err)
	}
	_, err = models.CreateMatchesTable(dbConn)
	if err != nil {
		fmt.Println("Could not create matches table")
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

	router.POST("/presign/:filename", h.PostPresign)

	return router
}

func setupS3Client(region string) *s3.Client {
	ctx := context.TODO()

	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		fmt.Println("Failed to load AWS config")
		panic(err)
	}

	return s3.NewFromConfig(cfg)
}
