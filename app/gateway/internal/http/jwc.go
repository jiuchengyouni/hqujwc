package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hqujwc/app/gateway/rpc"
	pb "hqujwc/idl/pb/jwc"
	"hqujwc/types/wx"
	"net/http"
	"time"
)

func WxReply(c *gin.Context) {
	var userReq pb.LoginRequest
	if err := c.Bind(&userReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "绑定参数错误",
		})
		return
	}
	fmt.Println(userReq.StuPass)
	r, err := rpc.GetGsSession(c, &userReq)
	bangdingjiebang := "由于新版教务处更新,请前往<a href=\"https://apps.hqu.edu.cn/wechat-hqu/wechatauth.html?proxyTo=authoauth&sendUrl=/connect/oauth2/authorize?appid=wxfe035b066fb1158b&redirect_uri=http%3A%2F%2Fwx.sends.cc%2Ftemporary%2Flogin&encode_flag=Y&response_type=code&scope=snsapi_base#wechat_redirect\">绑定与更新至新版教务处密码（点这里~）</a> 进行更新操作,请确保所绑定账号密码能登录<a href=\"https://id.hqu.edu.cn\">统一身份认证测试</a>。\n如果本微信绑定错学号，请输入{强制解绑}，以重新绑定。不包括{}\n如有问题，请添加咨询QQ群：712174205"
	bangdingjiebang += "\n现已支持新版教务处绩点排名查询功能~"
	c.Header("Content-Type", "application/xml; charset=utf-8")
	if err != nil {
		if err.Error() == "登录失败" {
			xml := wx.WxXmlResponse{
				//ToUserName:   req.FromUserName,
				//FromUserName: req.ToUserName,
				CreateTime: time.Now().Unix(),
				MsgType:    "text",
				Content:    bangdingjiebang,
			}
			c.Header("Content-Type", "application/xml; charset=utf-8")
			c.XML(http.StatusOK, xml)
			return
		}
		if err.Error() == "由于数据问题，暂不提供研究生查成绩服务" {
			xml := wx.WxXmlResponse{
				//ToUserName:   req.FromUserName,
				//FromUserName: req.ToUserName,
				CreateTime: time.Now().Unix(),
				MsgType:    "text",
				Content:    "由于数据问题，暂不提供研究生查成绩服务",
			}
			c.Header("Content-Type", "application/xml; charset=utf-8")
			c.XML(http.StatusOK, xml)
			return
		}
		return
	}
	var gradeReq = pb.Emaphome_WEURequest{GesSession: r.GetGesSession()}
	gradeResponse, err := rpc.GetGradeList(c, &gradeReq)
	if err != nil {
		xml := wx.WxXmlResponse{
			//ToUserName:   req.FromUserName,
			//FromUserName: req.ToUserName,
			CreateTime: time.Now().Unix(),
			MsgType:    "text",
			Content:    "由于数据问题，暂不提供研究生查成绩服务",
		}
		c.XML(http.StatusOK, xml)
		return
	}
	xml := wx.WxXmlResponse{
		//ToUserName:   req.FromUserName,
		//FromUserName: req.ToUserName,
		CreateTime: time.Now().Unix(),
		MsgType:    "text",
		Content:    gradeResponse.Emaphome_WEU,
	}
	c.XML(http.StatusOK, xml)
	return
}
