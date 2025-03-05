package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
)

func (h Handler) PostPresign(ctx *gin.Context) {
	filename := ctx.Param("filename")
	presignedUrl, err := h.createPresignedUrl(context.TODO(), filename)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"url": presignedUrl})
}
