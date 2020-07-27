package driver

import (
	"fmt"
	"github.com/rayylee/etcdadmin/etcdadmind/config"
	"github.com/rayylee/etcdadmin/etcdadmind/log"
	pb "github.com/rayylee/etcdadmin/etcdadmind/pb/etcdadminpb"
	"github.com/rayylee/etcdadmin/etcdadmind/server/driver/client"
	"github.com/rayylee/etcdadmin/etcdadmind/server/driver/command"
	"github.com/rayylee/etcdadmin/etcdadmind/server/driver/etcdcfg"
	"github.com/rayylee/etcdadmin/etcdadmind/utils"
	"go.uber.org/zap"
	"strconv"
)

type DriverInterface interface {
	// Attempts to add a member into the cluster.
	AddMember(m *pb.AddMemberRequest_Member) error

	// Manager config files include wal of etcd, and alse manger etcd service.
	ManagerEtcd(cmd pb.EtcdCmd, clearwal bool, cfgs []*pb.ManagerEtcdRequest_Config) error

	// Lists all the members in the cluster.
	ListMember() ([]*pb.ListMemberReply_Member, error)

	// Removes a member from the cluster.
	RemoveMember(name string) error
}

type DriverImpl struct {
	logger   *zap.Logger
	portGrpc string
}

func resetEtcdConfig() error {
	etcdCfgFile := config.Init().Get("ETCD_CONF_FILE")

	m, err := etcdcfg.EtcdConfigMapInit()

	if err != nil {
		return err
	}

	return etcdcfg.EtcdConfigWrite(etcdCfgFile, m)
}

func resetEtcd(isStart bool) error {
	var err error

	// Stop etcd but ignore error
	command.CmdEtcdctlStop()

	if err := resetEtcdConfig(); err != nil {
		goto exit
	}

	if err = etcdcfg.EtcdWalDelete(); err != nil {
		goto exit
	}

exit:
	if isStart == true {
		er := command.EtcdctlStart()
		if err == nil {
			err = er
		}
	}
	return err
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
		return err
	}

	// reset etcd.cfg, remove wal and stop etcd by remote
	cfgs, _ := etcdcfg.EtcdConfigMapInit()
	c.GrpcClientManagerEtcd(cfgs, true, client.EtcdCmdStop)

	cfgs["ETCD_NAME"] = m.Name
	cfgs["ETCD_ADVERTISE_CLIENT_URLS"] = fmt.Sprintf("http://%s:%s", m.Ip, clientPort)
	cfgs["ETCD_INITIAL_ADVERTISE_PEER_URLS"] = fmt.Sprintf("http://%s:%s", m.Ip, peerPort)
	if utils.ContainsString(ips, m.Ip) >= 0 {
		cfgs["ETCD_INITIAL_CLUSTER_STATE"] = "new"
		cfgs["ETCD_INITIAL_CLUSTER"] = fmt.Sprintf("%s=http://%s:%s", m.Name, m.Ip, peerPort)
		c.GrpcClientManagerEtcd(cfgs, false, client.EtcdCmdStart)
	} else {
		cfgs["ETCD_INITIAL_CLUSTER_STATE"] = "existing"

		members, _ := command.MemberList()
		initCluster := ""
		for _, i := range members {
			initCluster = fmt.Sprintf("%s%s=http://%s:%s,", initCluster,
				i.Name, i.Ipaddr, peerPort)
		}
		initCluster = fmt.Sprintf("%s%s=http://%s:%s", initCluster, m.Name, m.Ip, peerPort)
		cfgs["ETCD_INITIAL_CLUSTER"] = initCluster

		// Add member to cluster
		command.MemberAdd(m.Name, m.Ip)

		// Remote writing etcd.cfg and start etcd
		c.GrpcClientManagerEtcd(cfgs, false, client.EtcdCmdStart)
	}
	return nil
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
			Name: m.Name,
			Ip:   m.Ipaddr,
		}
		memberList = append(memberList, memberInfo)
	}

	return memberList, nil
}

func (drv *DriverImpl) RemoveMember(name string) error {
	var err error

	members, _ := command.MemberList()
	for _, m := range members {
		if m.Name == name {

			// stop etcd if the cluster only contains one member
			if len(members) == 1 {
				c := client.New(m.Ipaddr, drv.portGrpc)
				defer client.Release(c)

				// reset etcd.cfg, remove wal and stop etcd by remote
				cfgs, _ := etcdcfg.EtcdConfigMapInit()
				c.GrpcClientManagerEtcd(cfgs, true, client.EtcdCmdStop)
			} else {
				id, _ := strconv.ParseUint(m.Id, 10, 64)
				err = command.MemberRemove(fmt.Sprintf("%x", id))
			}

			break
		}
	}

	return err
}
