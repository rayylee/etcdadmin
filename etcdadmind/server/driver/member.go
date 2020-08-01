package driver

import (
	"context"
	"errors"
	"fmt"
	"github.com/rayylee/etcdadmin/etcdadmind/server/driver/command"
	"time"
)

func waitMemberHealth(endpoint string) error {
	// etcd/etcdserver/server.go:
	// healthInterval is the minimum time the cluster should be healthy
	// before accepting add member requests.
	//
	// here is the maximum timeout
	healthInterval := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), healthInterval)
	defer cancel()

	var err error
	c := make(chan int)
	over := false
	go func() {
		for over == false {
			time.Sleep(time.Millisecond * 100)
			err = command.EndpointsHealth(endpoint)
			if err == nil {
				break
			}
		}
		c <- 0
	}()

	select {
	case <-ctx.Done():
		over = true
	case <-c:
		return nil
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
