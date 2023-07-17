package rpc

import (
	"context"
	"errors"
	wxPb "hqujwc/idl/pb/wx"
)

func CheckSignature(ctx context.Context, req *wxPb.WxAccessRequest) (resp *wxPb.WxAccessResponse, err error) {
	r, err := WxClient.WeChatCallBack(ctx, req)
	if err != nil {
		return
	}

	if r.Code != 200 {
		err = errors.New("校验失败")
		return
	}

	return r, nil
}
