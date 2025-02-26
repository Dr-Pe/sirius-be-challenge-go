package main

import (
	"errors"
	"net/http"

	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
)

func (h Handler) postMatch(ctx *gin.Context) {
	var err error
	var match models.Match

	err = ctx.ShouldBindJSON(&match)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid match data"})
		return
	}
	_, err = match.Create(h.dbConn)
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

func (h Handler) getMatches(ctx *gin.Context) {
	var err error
	var matches []models.Match
	var query = struct {
		Status string `form:"status"`
	}{}

	err = ctx.ShouldBindQuery(&query)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if query.Status != "" {
		matches, err = models.SelectMatchesByStatus(h.dbConn, query.Status)
	} else {
		matches, err = models.SelectAllMatches(h.dbConn)
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, matches)
}
