package client

import (
	"context"
	"fmt"
	pb "github.com/rayylee/etcdadmin/etcdadmind/pb/etcdadminpb"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	caller pb.GrpcEtcdAdminClient
	conn   *grpc.ClientConn
	ip     string
	port   string
}

func New(ip string, port string) *GrpcClient {
	c := &GrpcClient{
		ip:   ip,
		port: port,
	}

	addr := fmt.Sprintf("%s:%s", c.ip, c.port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())

	if err != nil {
		fmt.Printf("Dial error: %v\n", addr)
		return c
	}
	c.conn = conn
	c.caller = pb.NewGrpcEtcdAdminClient(conn)

	return c
}

func Release(c *GrpcClient) error {
	return c.conn.Close()
}

func (c *GrpcClient) GrpcClientAddmember(name string, ip string) {
	m := &pb.AddMemberRequest_Member{Name: name, Ip: ip}

	r, err := c.caller.GrpcAddMember(context.Background(),
		&pb.AddMemberRequest{Member: m})

	if err == nil {
		fmt.Printf("%v \n", r.Errcode)
	} else {
		fmt.Printf("%v %v\n", r.Errcode, r.Errmsg)
	}
}

func (c *GrpcClient) GrpcClientListmember() (map[string]string, error) {
	m := make(map[string]string)

	r, err := c.caller.GrpcListMember(context.Background(),
		&pb.ListMemberRequest{})

	if err != nil {
		fmt.Printf("%v %v\n", r.Errcode, r.Errmsg)
		return m, err
	}

	for _, member := range r.Members {
		m[member.Name] = member.Ip
	}
	return m, nil
}

func (c *GrpcClient) GrpcClientRemovemember(name string) error {
	r, err := c.caller.GrpcRemoveMember(context.Background(),
		&pb.RemoveMemberRequest{
			Name: name,
		})

	if err != nil {
		fmt.Printf("%v %v\n", r.Errcode, r.Errmsg)
		return err
	}

	return nil
}
