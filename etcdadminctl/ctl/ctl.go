package ctl

import (
	"github.com/rayylee/etcdadmin/etcdadminctl/ctl/command"
	"github.com/spf13/cobra"
	"time"
)

const (
	cliName        = "etcdadminctl"
	cliDescription = "A simple command line of etcdadminctl."

	defaultDialTimeout      = 2 * time.Second
	defaultCommandTimeOut   = 5 * time.Second
	defaultKeepAliveTime    = 2 * time.Second
	defaultKeepAliveTimeOut = 6 * time.Second
)

var (
	globalFlags = command.GlobalFlags{}
)

var (
	rootCmd = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"etcdadminctl"},
	}
)

func init() {

	rootCmd.PersistentFlags().StringVar(&globalFlags.Endpoint,
		"endpoint", "127.0.0.1:2390", "gRPC endpoints")

	rootCmd.PersistentFlags().StringVar(&globalFlags.OutputFormat,
		"write-out", "json", "Output format")

	rootCmd.AddCommand(
		command.NewVersionCommand(),
		command.NewMemberCommand(),
	)
}

func Start() {
	if err := rootCmd.Execute(); err != nil {
		command.ExitWithError(command.ExitError, err)
	}
}

func init() {
	cobra.EnablePrefixMatching = true
}
