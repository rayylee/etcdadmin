package command

import (
	"encoding/json"
	"errors"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	_ "reflect"
	"strconv"
	"strings"
)

type EtcdMember struct {
	Name   string
	Ipaddr string
	Id     string
	Status string
}

func EtcdctlStart() error {
	result := CmdEtcdctlStart()
	if result.err != nil {
		return result.err
	}
	return nil
}

func MemberList() ([]*EtcdMember, error) {
	mslice := []*EtcdMember{}

	var err error
	result := CmdEtcdctlMemberList()

	if len(result.stdout) > 0 {
		j, err := simplejson.NewJson([]byte(result.stdout))
		if err == nil {
			members, err := j.Get("members").Array()
			if err != nil {
				return mslice, nil
			}

			for _, i := range members {
				member, _ := i.(map[string]interface{})

				name, _ := member["name"].(string)
				idnum, _ := member["ID"].(json.Number)
				id, _ := strconv.ParseUint(string(idnum), 10, 64)

				urls, _ := member["peerURLs"].([]interface{})
				// Using the first url: urls[0]
				ip := strings.Split(urls[0].(string), "//")[1]
				ip = strings.Split(ip, ":")[0]

				// reference etcd/etcdctl/ctlv3/command/printer.go
				status := "started"
				if len(name) == 0 {
					status = "unstarted"
				}

				mslice = append(mslice,
					&EtcdMember{
						Name:   name,
						Id:     fmt.Sprintf("%x", id),
						Ipaddr: ip,
						Status: status,
					})
			}
		}
	}
	if result.err != nil {
		err = result.err
	}
	if len(result.stderr) > 0 {
		err = errors.New(result.stderr)
	}

	return mslice, err
}

func MemberAdd(name string, ip string) error {
	result := CmdEtcdctlMemberAdd(name, ip)

	var err error
	if result.err != nil {
		err = result.err
	}
	if len(result.stderr) > 0 {
		err = errors.New(result.stderr)
	}
	return err
}

func MemberRemove(id string) error {
	result := CmdEtcdctlMemberRemove(id)

	var err error
	if result.err != nil {
		err = result.err
	}
	if len(result.stderr) > 0 {
		err = errors.New(result.stderr)
	}
	return err
}
