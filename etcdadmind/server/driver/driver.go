package driver

import (
	"errors"
	"fmt"
	"github.com/rayylee/etcdadmin/etcdadmind/config"
	"github.com/rayylee/etcdadmin/etcdadmind/log"
	pb "github.com/rayylee/etcdadmin/etcdadmind/pb/etcdadminpb"
	"github.com/rayylee/etcdadmin/etcdadmind/server/driver/client"
	"github.com/rayylee/etcdadmin/etcdadmind/server/driver/command"
	"github.com/rayylee/etcdadmin/etcdadmind/server/driver/etcdcfg"
	"github.com/rayylee/etcdadmin/etcdadmind/utils"
	"go.uber.org/zap"
)

type EtcdMember struct {
	Name     string
	Ipaddr   string
	Id       string
	IsLeader string
	IsHealth string
}

type DriverInterface interface {
	// Attempts to add a member into the cluster.
	AddMember(m *pb.AddMemberRequest_Member) error

	// Manager config files include wal of etcd, and alse manger etcd service.
	ManagerEtcd(cmd pb.EtcdCmd, clearwal bool, cfgs []*pb.ManagerEtcdRequest_Config) error

	// Lists all the members in the cluster.
	ListMember() ([]*pb.ListMemberReply_Member, error)

	// Removes a member from the cluster.
	RemoveMember(id string) error
}

type DriverImpl struct {
	logger   *zap.Logger
	portGrpc string
}

func New() *DriverImpl {
	cfg := config.Init()

	drv := &DriverImpl{
		portGrpc: cfg.Get("GRPC_PORT"),
		logger:   log.GetLogger(),
	}
	return drv
}

func (drv *DriverImpl) AddMember(m *pb.AddMemberRequest_Member) error {
	cfgStore := config.Init()
	clientPort := cfgStore.Get("ETCD_CLIENT_PORT")
	peerPort := cfgStore.Get("ETCD_PEER_PORT")

	c := client.New(m.Ip, drv.portGrpc)
	defer client.Release(c)

	drv.logger.Info(fmt.Sprintf("add member: %v %v", m.Name, m.Ip))
	ips, err := utils.GetHostIP4()
	if err != nil {
		drv.logger.Error(fmt.Sprintf("get host ipv4: %v", err))
		return err
	}

	// reset etcd.cfg, remove wal and stop etcd by remote
	cfgs, _ := etcdcfg.EtcdConfigMapInit()
	_, err = c.GrpcClientManagerEtcd(cfgs, true, client.EtcdCmdStop)
	if err != nil {
		drv.logger.Error(fmt.Sprintf("remote manage etcd: %v", err))
		goto exit
	}

	cfgs["ETCD_NAME"] = m.Name
	cfgs["ETCD_ADVERTISE_CLIENT_URLS"] = fmt.Sprintf("http://%s:%s", m.Ip, clientPort)
	cfgs["ETCD_INITIAL_ADVERTISE_PEER_URLS"] = fmt.Sprintf("http://%s:%s", m.Ip, peerPort)

	if utils.ContainsString(ips, m.Ip) >= 0 {
		cfgs["ETCD_INITIAL_CLUSTER_STATE"] = "new"
		cfgs["ETCD_INITIAL_CLUSTER"] = fmt.Sprintf("%s=http://%s:%s", m.Name, m.Ip, peerPort)

		c.GrpcClientManagerEtcd(cfgs, false, client.EtcdCmdStart)
	} else {
		cfgs["ETCD_INITIAL_CLUSTER_STATE"] = "existing"

		members, err := command.MemberList()
		if err != nil {
			drv.logger.Error(fmt.Sprintf("get member list: %v", err))
			goto exit
		}
		if len(members) < 1 {
			err = errors.New("cluster no any member")
			drv.logger.Error(fmt.Sprintf("get member list: %v", err))
			goto exit
		}

		initCluster := ""
		for _, i := range members {
			initCluster = fmt.Sprintf("%s%s=http://%s:%s,", initCluster,
				i.Name, i.Ipaddr, peerPort)
		}
		initCluster = fmt.Sprintf("%s%s=http://%s:%s", initCluster, m.Name, m.Ip, peerPort)
		cfgs["ETCD_INITIAL_CLUSTER"] = initCluster

		// Add member to cluster
		err = command.MemberAdd(m.Name, m.Ip)
		if err != nil {
			drv.logger.Error(fmt.Sprintf("add member: %v", err))
			goto exit
		}

		// writing etcd.cfg and start etcd by remote
		_, err = c.GrpcClientManagerEtcd(cfgs, false, client.EtcdCmdStart)
		if err != nil {
			drv.logger.Error(fmt.Sprintf("remote manage etcd: %v", err))
			// rollback
			member, er := getMemberByName(m.Name)
			if er == nil && len(member.Id) > 0 {
				command.MemberRemove(member.Id)
			}
		}
	}
exit:
	return err
}

func (drv *DriverImpl) ManagerEtcd(cmd pb.EtcdCmd, clearwal bool,
	cfgs []*pb.ManagerEtcdRequest_Config) error {
	cfgStore := config.Init()

	if cmd != pb.EtcdCmd_NONE {
		command.CmdEtcdctlStop()
	}

	if len(cfgs) > 0 {
		etcdCfgFile := cfgStore.Get("ETCD_CONF_FILE")
		m := map[string]string{}
		for _, c := range cfgs {
			m[c.Key] = c.Value
		}

		etcdcfg.EtcdConfigWrite(etcdCfgFile, m)
	}

	if clearwal == true {
		etcdcfg.EtcdWalDelete()
	}

	if cmd == pb.EtcdCmd_START || cmd == pb.EtcdCmd_RESTART {
		command.CmdEtcdctlStart()
	}
	return nil
}

func (drv *DriverImpl) ListMember() ([]*pb.ListMemberReply_Member, error) {
	memberList := []*pb.ListMemberReply_Member{}

	members, _ := command.MemberList()
	for _, m := range members {
		memberInfo := &pb.ListMemberReply_Member{
			Id:   m.Id,
			Name: m.Name,
			Ip:   m.Ipaddr,
		}
		memberList = append(memberList, memberInfo)
	}

	return memberList, nil
}

func (drv *DriverImpl) RemoveMember(id string) error {
	var err error

	members, _ := command.MemberList()
	for _, m := range members {
		if m.Id == id {

			// stop etcd if the cluster only contains one member
			if len(members) == 1 {
				c := client.New(m.Ipaddr, drv.portGrpc)
				defer client.Release(c)

				// reset etcd.cfg, remove wal and stop etcd by remote
				cfgs, _ := etcdcfg.EtcdConfigMapInit()
				c.GrpcClientManagerEtcd(cfgs, true, client.EtcdCmdStop)
			} else {
				err = command.MemberRemove(m.Id)
			}

			break
		}
	}

	return err
}
