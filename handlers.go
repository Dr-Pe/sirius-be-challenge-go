package main

import (
	"net/http"

	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
)

func postPlayer(ctx *gin.Context) {
	var err error
	var player models.Player

	err = ctx.ShouldBindJSON(&player)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = player.Create(dbConn)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{"message": "Player created successfully"})
	}
}

func getPlayers(ctx *gin.Context) {
	var err error
	var players []models.Player
	var query = struct {
		Name string `form:"name"`
	}{}

	err = ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if query.Name != "" {
		players, err = models.SelectPlayersByName(dbConn, query.Name)
	} else {
		players, err = models.SelectAllPlayers(dbConn)
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, players)
}

func getPlayer(ctx *gin.Context) {
	var err error
	var player models.Player
	var id = ctx.Param("id")

	player, err = models.SelectPlayerById(dbConn, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, player)
}

func putPlayer(ctx *gin.Context) {
	var err error
	var id = ctx.Param("id")
	var player models.Player

	player, err = models.SelectPlayerById(dbConn, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = ctx.ShouldBindJSON(&player)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = models.UpdatePlayerById(dbConn, id, player)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Player updated successfully"})
}

func deletePlayer(ctx *gin.Context) {
	var err error
	var id = ctx.Param("id")

	_, err = models.DeletePlayerById(dbConn, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Player deleted successfully"})
}
