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
	client.GrpcClientAddmember("name-test", "127.0.0.1")
}

// Test Entry
func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	teardown()
	os.Exit(ret)
}
