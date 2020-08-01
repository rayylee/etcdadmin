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

type EtcdEndpoints struct {
	Endpoint string
	Id       string
	Leader   string
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

func EndpointsStatus(endpoints string) ([]*EtcdEndpoints, error) {
	eslice := []*EtcdEndpoints{}

	var err error
	result := CmdEtcdctlEndpointsStatus(endpoints)

	if len(result.stdout) <= 0 {
		if result.err != nil {
			err = result.err
		}
		if len(result.stderr) > 0 {
			err = errors.New(result.stderr)
		}
		return eslice, nil
	}

	j, err := simplejson.NewJson([]byte(result.stdout))
	if err != nil {
		return eslice, nil
	}

	maps, _ := j.Array()
	for _, element := range maps {
		ept := &EtcdEndpoints{}

		member, _ := element.(map[string]interface{})

		for key, value := range member {
			if key == "Endpoint" {
				// endpoint
				ept.Endpoint = value.(string)

			} else if key == "Status" {
				status := value.(map[string]interface{})

				for key2, value2 := range status {
					if key2 == "header" {
						header, _ := value2.(map[string]interface{})

						memberidstr :=
							header["member_id"].(json.Number).String()

						memberid, _ := strconv.ParseUint(memberidstr, 10, 64)

						// member_id
						ept.Id = fmt.Sprintf("%x", memberid)

					} else if key2 == "leader" {
						leaderstr := value2.(json.Number).String()
						leader, _ := strconv.ParseUint(leaderstr, 10, 64)

						// leader
						ept.Leader = fmt.Sprintf("%x", leader)

					}
				}
			}
		}
		// append endpoint
		eslice = append(eslice, ept)

	}
	return eslice, nil
}

func EndpointsHealth(endpoints string) error {
	var err error

	result := CmdEtcdctlEndpointsHealth(endpoints)

	if result.err != nil {
		err = result.err
	}
	if len(result.stderr) > 0 {
		err = errors.New(result.stderr)
	}
	return err
}
