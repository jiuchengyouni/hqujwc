package routes

import (
	"github.com/gin-gonic/gin"
	"hqujwc/app/gateway/internal/http"
	"hqujwc/app/gateway/middlewares"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middlewares.Cors())
	ginRouter.POST("/jwc", http.GetSession)
	return ginRouter
}
