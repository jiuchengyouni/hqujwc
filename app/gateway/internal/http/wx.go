package http

import (
	"github.com/gin-gonic/gin"
	"hqujwc/app/gateway/rpc"
	pb "hqujwc/idl/pb/wx"
	"net/http"
)

func WeChatCallBack(c *gin.Context) {
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")
	wxAccessRequest := pb.WxAccessRequest{
		Signature: signature,
		Timestamp: timestamp,
		Nonce:     nonce,
		Echoster:  echostr,
	}
	r, err := rpc.CheckSignature(c, &wxAccessRequest)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 403,
			"msg":  err,
			"data": "",
		})
		return
	}
	c.String(http.StatusOK, r.Echoster)
}
