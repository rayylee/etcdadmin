package driver

import (
	"context"
	"errors"
	"fmt"
	"github.com/rayylee/etcdadmin/etcdadmind/server/driver/command"
	"time"
)

func waitMemberHealth(endpoint string, timeout uint) error {
	// time.Second*timeout second
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Second*time.Duration(timeout))
	defer cancel()

	var err error
	over := false
	c := make(chan string)
	go func() {
		for over == false {
			err = command.EndpointsHealth(endpoint)
			time.Sleep(time.Millisecond * 100)
		}
		c <- err.Error()
	}()

	select {
	case <-ctx.Done():
		over = true
	case ret := <-c:
		err = errors.New(ret)
	}

	return err
}

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
