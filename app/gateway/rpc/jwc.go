package rpc

import (
	"context"
	"errors"
	"fmt"
	jwcPb "hqujwc/idl/pb/jwc"
)

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
