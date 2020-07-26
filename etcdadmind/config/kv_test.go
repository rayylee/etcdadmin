package config

import (
	"os"
	"testing"
)

var (
	gConfigFile = "test.cfg"

	gIpaddr   = "192.168.1.100"
	gPort     = "123"
	gHostname = "host-test"
	gContext  = []string{
		"# this is a comment   \n",
		"  ip=192.168.1.100  \n",
		"\n",
		"# listen=on\n",
		"  port=123 \n",
		"\n"}
)

func setup() {
	f, err := os.OpenFile(gConfigFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer f.Close()
	if err == nil {
		for _, s := range gContext {
			f.WriteString(string(s))
		}
	}
}

func teardown() {
	_, err := os.Stat(gConfigFile)
	if err == nil {
		os.Remove(gConfigFile)
	}
}

func TestCase(t *testing.T) {
	kv := Load(gConfigFile)

	ip := kv.Get("ip")
	if ip != gIpaddr {
		t.Fatal("Get ip")
	}

	port := kv.Get("port")
	if port != gPort {
		t.Fatal("Get port")
	}

	err := kv.Set("hostname", gHostname)
	if err != nil {
		t.Fatal("Set hostname")
	}

	hostname := kv.Get("hostname")
	if hostname != gHostname {
		t.Fatal("Get hostname")
	}

}

// Test Entry
func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	teardown()
	os.Exit(ret)
}
