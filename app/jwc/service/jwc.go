package service

import (
	"context"
	"fmt"
	"hqujwc/app/jwc/utils"
	pb "hqujwc/idl/pb/jwc"
	"hqujwc/types"
	"sync"
)

type JwcSrv struct {
	pb.UnimplementedJwcServiceServer
}

var JwcSrvIns *JwcSrv
var JwcSrvOnce sync.Once

func GetJwcSrv() *JwcSrv {
	JwcSrvOnce.Do(func() {
		JwcSrvIns = &JwcSrv{}
	})
	return JwcSrvIns
}

func (*JwcSrv) GetGsSession(ctx context.Context, req *pb.LoginRequest) (resp *pb.GsSessionResponse, err error) {
	resp = new(pb.GsSessionResponse)
	resp.Code = 200
	loginRequest := types.LoginRequestBody{
		Stunum:   req.StuNum,
		Password: req.StuPass,
	}
	gsSession, err := utils.GetGsSession(&loginRequest)
	fmt.Println(err)
	if err != nil {
		return
	}
	resp.GesSession = gsSession
	resp.Code = 200
	return
}
