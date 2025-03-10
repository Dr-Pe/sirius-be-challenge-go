package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
)

// @Summary Post player
// @Description Create a new player
// @Tags players
// @Accept json
// @Produce json
// @Param player body models.Player true "Player object"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /players [post]
func (h Handler) PostPlayer(ctx *gin.Context) {
	var err error
	var player models.Player
	var res sql.Result
	var id int64
	var presignedUrl string

	err = ctx.ShouldBindJSON(&player)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player data"})
		return
	}
	res, err = player.Create(h.DbConn)
	if err != nil {
		var playerErr models.PlayerError
		if errors.As(err, &playerErr) {
			ctx.JSON(playerErr.StatusCode, gin.H{"error": playerErr.Err})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	id, err = res.LastInsertId()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	presignedUrl, err = h.createPresignedUrl(context.TODO(), fmt.Sprintf("%d_%s", id, player.Name))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	player.ProfilePictureUrl = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%d_%s", h.BucketName, h.Region, id, player.Name)
	_, err = models.UpdatePlayerById(h.DbConn, fmt.Sprintf("%d", id), player)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Player created successfully, upload your profile picture to the following URL", "url": presignedUrl})
}

// @Summary Get players
// @Description Get all players or players by name
// @Tags players
// @Accept json
// @Produce json
// @Param name query string false "Player name"
// @Success 200 {array} models.Player
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /players [get]
func (h Handler) GetPlayers(ctx *gin.Context) {
	var err error
	var players []models.Player
	var query = struct {
		Name string `form:"name"`
	}{}

	err = ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if query.Name != "" {
		players, err = models.SelectPlayersByName(h.DbConn, query.Name)
	} else {
		players, err = models.SelectAllPlayers(h.DbConn)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, players)
}

// @Summary Get player
// @Description Get player by id
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Success 200 {object} models.Player
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /players/{id} [get]
func (h Handler) GetPlayer(ctx *gin.Context) {
	var err error
	var player models.Player
	var id = ctx.Param("id")

	player, err = models.SelectPlayerById(h.DbConn, id)
	if err != nil {
		var playerErr models.PlayerError
		if errors.As(err, &playerErr) {
			ctx.JSON(playerErr.StatusCode, gin.H{"error": playerErr.Err})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, player)
}

// @Summary Put player
// @Description Update player by id
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Param player body models.Player true "Player object"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /players/{id} [put]
func (h Handler) PutPlayer(ctx *gin.Context) {
	var err error
	var id = ctx.Param("id")
	var player models.Player
	var presignedUrl string

	player, err = models.SelectPlayerById(h.DbConn, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = ctx.ShouldBindJSON(&player)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = models.UpdatePlayerById(h.DbConn, id, player)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	presignedUrl, err = h.createPresignedUrl(context.TODO(), fmt.Sprintf("%s_%s", id, player.Name))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Player updated successfully, you can update your profile picture using the following URL", "url": presignedUrl})
}

// @Summary Delete player
// @Description Delete player by id
// @Tags players
// @Accept json
// @Produce json
// @Param id path string true "Player ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /players/{id} [delete]
func (h Handler) DeletePlayer(ctx *gin.Context) {
	var err error
	var id = ctx.Param("id")
	var player models.Player

	player, err = models.SelectPlayerById(h.DbConn, id)
	if err != nil {
		var playerErr models.PlayerError
		if errors.As(err, &playerErr) {
			ctx.JSON(playerErr.StatusCode, gin.H{"error": playerErr.Err})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	_, err = models.DeletePlayerById(h.DbConn, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = h.deleteObject(context.TODO(), fmt.Sprintf("%s_%s", id, player.Name))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Player deleted successfully"})
}
