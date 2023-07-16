package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hqujwc/app/gateway/rpc"
	pb "hqujwc/idl/pb/jwc"
	"net/http"
)

func GetSession(ctx *gin.Context) {
	var userReq pb.LoginRequest
	if err := ctx.Bind(&userReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "绑定参数错误",
		})
		return
	}
	fmt.Println(userReq.StuPass)
	r, err := rpc.GetGsSession(ctx, &userReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "调用失败",
			"data": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": r,
	})
}
