package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
)

// @Summary Post match
// @Description Create a new match
// @Tags matches
// @Accept json
// @Produce json
// @Param match body models.Match true "Match object"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /matches [post]
func (h Handler) PostMatch(ctx *gin.Context) {
	var err error
	var match models.Match

	err = ctx.ShouldBindJSON(&match)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match data"})
		return
	}
	_, err = match.Create(h.DbConn)
	if err != nil {
		var matchErr models.MatchError
		if errors.As(err, &matchErr) {
			ctx.JSON(matchErr.StatusCode, gin.H{"error": matchErr.Err})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Match created successfully"})
}

// @Summary Get matches
// @Description Get all matches
// @Tags matches
// @Accept json
// @Produce json
// @Param status query string false "Match status"
// @Success 200 {object} []models.Match
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /matches [get]
func (h Handler) GetMatches(ctx *gin.Context) {
	var err error
	var matches []models.Match
	var query = struct {
		Status string `form:"status"`
	}{}

	err = ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else if query.Status != "" {
		matches, err = models.SelectMatchesByStatus(h.DbConn, query.Status)
	} else {
		matches, err = models.SelectAllMatches(h.DbConn)
	}

	if err != nil {
		var matchErr models.MatchError
		if errors.As(err, &matchErr) {
			ctx.JSON(matchErr.StatusCode, gin.H{"error": matchErr.Err})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, matches)
}

// @Summary Get match
// @Description Get match by id
// @Tags matches
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} models.Match
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /matches/{id} [get]
func (h Handler) GetMatch(ctx *gin.Context) {
	var err error
	var match models.Match
	var id = ctx.Param("id")

	match, err = models.SelectMatchById(h.DbConn, id)
	if err != nil {
		var matchErr models.MatchError
		if errors.As(err, &matchErr) {
			ctx.JSON(matchErr.StatusCode, gin.H{"error": matchErr.Err})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	ctx.JSON(http.StatusOK, match)
}

// @Summary Put match
// @Description Update match by id
// @Tags matches
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param match body models.Match true "Match object"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /matches/{id} [put]
func (h Handler) PutMatch(ctx *gin.Context) {
	var err error
	var id = ctx.Param("id")
	var match models.Match

	match, err = models.SelectMatchById(h.DbConn, id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = ctx.ShouldBindJSON(&match)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = updateMatch(h.DbConn, id, match)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Match updated successfully"})
}

// @Summary Delete match
// @Description Delete match by id
// @Tags matches
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /matches/{id} [delete]
func (h Handler) DeleteMatch(ctx *gin.Context) {
	var err error
	var id = ctx.Param("id")
	var res sql.Result

	res, err = models.DeleteMatchById(h.DbConn, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if rowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Match not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Match deleted successfully"})
}

func updateMatch(dbConn *sql.DB, id string, match models.Match) error {
	var err error
	var loserId int

	_, err = models.UpdateMatchById(dbConn, id, match)
	if err != nil {
		return err
	} else if match.WinnerId != 0 {
		if match.WinnerId == match.Player1id {
			loserId = match.Player2id
		} else {
			loserId = match.Player1id
		}
		err = models.UpdatePoints(dbConn, match.WinnerId, loserId)
		if err != nil {
			return err
		}
	}

	return nil
}
