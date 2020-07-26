package client

import (
	"os"
	"testing"
)

const (
	grpcIP   = "127.0.0.1"
	grpcPort = "2390"
)

var (
	client *GrpcClient
)

func setup() {
	client = New(grpcIP, grpcPort)
}

func teardown() {
	Release(client)
}

func TestGrpcClientAddmember(t *testing.T) {
	m := map[string]string{
		"name1": "hostname1-test",
		"name2": "hostname2-test",
	}
	client.GrpcClientManagerEtcd(m, false, EtcdCmdNone)
}

// Test Entry
func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	teardown()
	os.Exit(ret)
}
