package rpc

import (
	"context"
	"errors"
	"fmt"
	jwcPb "hqujwc/idl/pb/jwc"
)

// 统一身份认证登录
func GetGsSession(ctx context.Context, req *jwcPb.LoginRequest) (resp *jwcPb.GsSessionResponse, err error) {
	fmt.Println(req)
	r, err := JwcClient.GetGsSession(ctx, req)
	fmt.Println(err)
	if err != nil {
		return
	}

	if r.Code != 200 {
		err = errors.New("登陆失败")
		return
	}
	return r, nil
}

func GetGradeList(ctx context.Context, req *jwcPb.Emaphome_WEURequest) (resp *jwcPb.Emaphome_WEUResponse, err error) {
	r, err := JwcClient.GetEmaphome_WEU(ctx, req)
	if err != nil || resp.Code != 200 {
		return
	}
	return r, nil
}
