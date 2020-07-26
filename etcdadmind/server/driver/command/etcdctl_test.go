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
