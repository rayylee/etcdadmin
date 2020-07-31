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
