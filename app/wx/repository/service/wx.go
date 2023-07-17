package service

import (
	"context"
	"hqujwc/app/wx/utils"
	"hqujwc/config"
	pb "hqujwc/idl/pb/wx"
	"sort"
	"sync"
)

type WxSrv struct {
	pb.UnimplementedWxServiceServer
}

var WxSrvIns *WxSrv

var WxSrvOnce sync.Once

func GetWxSrv() *WxSrv {
	WxSrvOnce.Do(func() {
		WxSrvIns = &WxSrv{}
	})
	return WxSrvIns
}

func (*WxSrv) CheckSignature(ctx context.Context, req *pb.WxAccessRequest) (resp *pb.WxAccessResponse) {
	resp.Code = 200
	tmpArr := []string{config.Toekn, req.Timestamp, req.Nonce}
	sort.Strings(tmpArr)
	tmpStr := utils.Sha1String(utils.Implode(tmpArr))
	if tmpStr == req.Signature {
		resp.Echoster = req.Echoster
		return
	} else {
		resp.Code = 404
		return
	}
}
