package server

import (
	"fmt"
	"github.com/rayylee/etcdadmin/etcdadmind/log"
	pb "github.com/rayylee/etcdadmin/etcdadmind/pb/etcdadminpb"
	"github.com/rayylee/etcdadmin/etcdadmind/server/driver"
	"google.golang.org/grpc"
	"net"
)

func Init(port string) error {
	etcdS := &ImplEtcdAdminServer{
		logger: log.GetLogger(),
		port:   port,
		drv:    driver.New(),
	}

	listen, err := net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%s", port))

	if err != nil {
		etcdS.logger.Error(fmt.Sprintf("failed to listen: %v", err))
		return err
	}

	s := grpc.NewServer()
	pb.RegisterGrpcEtcdAdminServer(s, etcdS)

	etcdS.logger.Info(fmt.Sprintf("gRpc sever listening:%s", port))
	go s.Serve(listen)

	return nil
}
