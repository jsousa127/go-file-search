package main

import (
	"github.com/gin-gonic/gin"
)

var statusOK, statusBadRequest, statusInternalServerError = 200, 400, 500

func WriteResponse(ctx *gin.Context, data interface{}) {
	resp := gin.H{"data": data, "status": statusOK}
	ctx.JSON(statusOK, resp)
}

func WriteError(ctx *gin.Context, err string, status int) {
	resp := gin.H{"error": err, "status": status}
	ctx.JSON(status, resp)
}

func Search(ctx *gin.Context) {
	path := ctx.Query("path")
	keyword := ctx.Query("keyword")

	if path == "" || keyword == "" {
		WriteError(ctx, "Malformed Request", statusBadRequest)
		return
	}

	results, err := searchDir(ctx, path, keyword)
	if err != nil {
		WriteError(ctx, err.Error(), statusInternalServerError)
	} else {
		WriteResponse(ctx, results)
	}
}

func Engine() *gin.Engine {
	engine := gin.Default()
	engine.GET("/search", Search)
	return engine
}
