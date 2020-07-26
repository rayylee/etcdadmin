package main

import (
	"fmt"
	"github.com/rayylee/etcdadmin/etcdadminctl/ctl"
	"os"
)

const (
	apiEnv = "ETCDADMINCTL_API"
)

func main() {
	apiv := os.Getenv(apiEnv)

	// unset apiEnv to avoid side-effect for future env and flag parsing.
	os.Unsetenv(apiEnv)

	if len(apiv) == 0 || apiv == "1" {
		ctl.Start()
		return
	}

	fmt.Fprintln(os.Stderr, "unsupported API version", apiv)
	os.Exit(1)
}
