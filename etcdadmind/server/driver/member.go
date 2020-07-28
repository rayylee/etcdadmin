package driver

import (
	"errors"
	"fmt"
	"github.com/rayylee/etcdadmin/etcdadmind/server/driver/command"
)

func getMemberByIpaddr(ipaddr string) (*command.EtcdMember, error) {
	members, err := command.MemberList()
	if err != nil {
		return &command.EtcdMember{}, err
	}

	for _, i := range members {
		if i.Ipaddr == ipaddr {
			return i, nil
		}
	}

	errmsg := fmt.Sprintf("member not found by ipaddr: %s", ipaddr)
	return &command.EtcdMember{}, errors.New(errmsg)
}

func getMemberByName(name string) (*command.EtcdMember, error) {
	members, err := command.MemberList()
	if err != nil {
		return &command.EtcdMember{}, err
	}

	for _, i := range members {
		if i.Name == name {
			return i, nil
		}
	}

	errmsg := fmt.Sprintf("member not found by name: %s", name)
	return &command.EtcdMember{}, errors.New(errmsg)
}
