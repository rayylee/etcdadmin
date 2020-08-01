package command

import (
	"fmt"
	"testing"
)

func TestGetMemberList(t *testing.T) {
	mslice, err := MemberList()
	if err == nil {
		for i := range mslice {
			fmt.Printf("Member slice: %v\n", mslice[i])
		}
	} else {
		fmt.Printf("Err: %v\n", err)
	}
}

func TestGetEndpointsStatus(t *testing.T) {
	t.Skip("skipping: need input of endpoints")
	endpoints := ""

	eslice, err := EndpointsStatus(endpoints)
	if err == nil {
		for i := range eslice {
			fmt.Printf("Endpoint Status: %v\n", eslice[i])
		}
	} else {
		fmt.Printf("Err: %v\n", err)
	}
}

func TestGetEndpointsHealth(t *testing.T) {
	endpoints := "127.0.0.1:2379"

	err := EndpointsHealth(endpoints)

	if err == nil {
		fmt.Printf("%s health\n", endpoints)
	} else {
		fmt.Printf("%s unhealth: %s\n", endpoints, err)
	}
}
