package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"hqujwc/app/wx/repository/dao"
	"hqujwc/app/wx/service"
	"hqujwc/config"
	wxPb "hqujwc/idl/pb/wx"
	"hqujwc/pkg/discovery"
	"net"
)

func main() {
	config.InitConfig()
	dao.InitDB()
	// etcd 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcAddress := config.Conf.Services["wx"].Addr[0]
	defer etcdRegister.Stop()
	taskNode := discovery.Server{
		Name: config.Conf.Domain["wx"].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定service
	wxPb.RegisterWxServiceServer(server, service.GetWxSrv())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err := etcdRegister.Register(taskNode, 10); err != nil {
		panic(fmt.Sprintf("start server failed, err: %v", err))
	}
	logrus.Info("server started listen on ", grpcAddress)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
