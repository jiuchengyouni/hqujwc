package routes

import (
	"github.com/gin-gonic/gin"
	"hqujwc/app/gateway/internal/http"
	"hqujwc/app/gateway/middlewares"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middlewares.Cors())
	ginRouter.POST("/jwc", http.WxReply)
	ginRouter.GET("/match", http.WeChatCallBack)
	return ginRouter
}
