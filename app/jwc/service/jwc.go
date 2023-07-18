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
func (*JwcSrv) GetEmaphome_WEU(ctx context.Context, req *pb.Emaphome_WEURequest) (resp *pb.Emaphome_WEUResponse, err error) {
	resp = new(pb.Emaphome_WEUResponse)
	resp.Code = 200
	gsSession := req.GesSession
	emaphome_WEU, err := utils.GetEmaphome_WEU(gsSession)
	if err != nil {
		resp.Code = 404
		return
	}
	resp.Emaphome_WEU = emaphome_WEU
	return
}
