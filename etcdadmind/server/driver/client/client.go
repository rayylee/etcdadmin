package client

import (
	"context"
	"fmt"
	"github.com/rayylee/etcdadmin/etcdadmind/log"
	pb "github.com/rayylee/etcdadmin/etcdadmind/pb/etcdadminpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type EtcdCmdType int32

const (
	EtcdCmdNone    EtcdCmdType = 0
	EtcdCmdStart   EtcdCmdType = 1
	EtcdCmdStop    EtcdCmdType = 2
	EtcdCmdRestart EtcdCmdType = 3
)

type GrpcClient struct {
	remote pb.GrpcEtcdAdminClient
	conn   *grpc.ClientConn
	ip     string
	port   string
	logger *zap.Logger
}

func New(ip string, port string) *GrpcClient {
	c := &GrpcClient{
		ip:     ip,
		port:   port,
		logger: log.GetLogger(),
	}

	addr := fmt.Sprintf("%s:%s", c.ip, c.port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		fmt.Printf("Dial error: %v\n", addr)
		return c
	}
	c.conn = conn
	c.remote = pb.NewGrpcEtcdAdminClient(conn)

	return c
}

func Release(c *GrpcClient) error {
	return c.conn.Close()
}

func (c *GrpcClient) GrpcClientManagerEtcd(cfg map[string]string, clearwal bool,
	cmd EtcdCmdType) (*pb.ManagerEtcdReply, error) {

	if c.logger != nil {
		c.logger.Info("call GrpcClientManagerEtcd")
	}

	var cfgs []*pb.ManagerEtcdRequest_Config
	for key := range cfg {
		cfg := pb.ManagerEtcdRequest_Config{Key: key, Value: cfg[key]}
		cfgs = append(cfgs, &cfg)
	}

	var etcdCmd pb.EtcdCmd
	switch cmd {
	case EtcdCmdStop:
		etcdCmd = pb.EtcdCmd_STOP
	case EtcdCmdStart:
		etcdCmd = pb.EtcdCmd_START
	case EtcdCmdRestart:
		etcdCmd = pb.EtcdCmd_RESTART
	default:
		etcdCmd = pb.EtcdCmd_NONE
	}

	r, err := c.remote.GrpcManagerEtcd(context.Background(),
		&pb.ManagerEtcdRequest{Cmd: etcdCmd, Clearwal: clearwal, Cfgs: cfgs})

	return r, err
}
