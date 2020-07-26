package utils

import (
	"fmt"
	"testing"
)

func TestGetHostIP4(t *testing.T) {
	ips, _ := GetHostIP4()
	fmt.Printf("IPs: %v\n", ips)
}
