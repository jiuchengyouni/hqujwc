package main

import "hqujwc/app/jwc/repository/dao"

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"hqujwc/app/jwc/service"
	"hqujwc/config"
	jwcPb "hqujwc/idl/pb/jwc"
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
	grpcAddress := config.Conf.Services["jwc"].Addr[0]
	defer etcdRegister.Stop()
	taskNode := discovery.Server{
		Name: config.Conf.Domain["jwc"].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()
	// 绑定service
	jwcPb.RegisterJwcServiceServer(server, service.GetJwcSrv())
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
