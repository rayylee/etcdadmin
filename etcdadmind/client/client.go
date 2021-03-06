package client

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/rayylee/etcdadmin/etcdadmind/pb/etcdadminpb"
	"google.golang.org/grpc"
)

type EtcdMember struct {
	Name     string
	Ip       string
	Id       string
	IsLeader string
	IsHealth string
}

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

func (c *GrpcClient) GrpcClientListmember() ([]*EtcdMember, error) {
	mslice := []*EtcdMember{}

	r, err := c.caller.GrpcListMember(context.Background(),
		&pb.ListMemberRequest{})

	if err != nil {
		fmt.Printf("%v %v\n", r.Errcode, r.Errmsg)
		return mslice, err
	}

	for _, member := range r.Members {
		mslice = append(mslice,
			&EtcdMember{
				Id:       member.Id,
				Name:     member.Name,
				Ip:       member.Ip,
				IsLeader: member.Isleader,
				IsHealth: member.Ishealth,
			})
	}
	return mslice, nil
}

func (c *GrpcClient) GrpcClientRemovemember(id string) error {
	r, err := c.caller.GrpcRemoveMember(context.Background(),
		&pb.RemoveMemberRequest{
			Id: id,
		})

	if err != nil {
		fmt.Printf("%v %v\n", r.Errcode, r.Errmsg)
		return err
	}

	if len(r.Errmsg) > 0 {
		return errors.New(r.Errmsg)
	}

	return nil
}
