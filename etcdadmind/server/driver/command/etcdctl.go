package command

import (
	"encoding/json"
	"errors"
	_ "fmt"
	simplejson "github.com/bitly/go-simplejson"
	_ "reflect"
	"strings"
)

type EtcdMember struct {
	Name   string
	Ipaddr string
	Id     string
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
				id, _ := member["ID"].(json.Number)

				urls, _ := member["peerURLs"].([]interface{})
				// Using the first url: urls[0]
				ip := strings.Split(urls[0].(string), "//")[1]
				ip = strings.Split(ip, ":")[0]

				mslice = append(mslice,
					&EtcdMember{
						Name:   name,
						Id:     string(id),
						Ipaddr: ip,
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
