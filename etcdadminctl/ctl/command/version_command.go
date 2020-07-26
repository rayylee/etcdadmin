package command

import (
	"fmt"
	"github.com/rayylee/etcdadmin/etcdadminctl/ctl/version"
	"github.com/spf13/cobra"
)

// NewVersionCommand prints out the version of etcd.
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version",
		Run:   versionCommandFunc,
	}
}

func versionCommandFunc(cmd *cobra.Command, args []string) {
	fmt.Println("Version:", version.Version)
	fmt.Println("API version:", version.APIVersion)
}
