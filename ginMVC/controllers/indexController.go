package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexHome(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index/index.html", gin.H{
		"msg": "easy gin",
	})
}
