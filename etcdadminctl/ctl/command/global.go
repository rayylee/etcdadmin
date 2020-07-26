package command

import (
	"github.com/spf13/cobra"
)

// GlobalFlags are flags that defined globally
// and are inherited to all sub-commands.
type GlobalFlags struct {
	Endpoint     string
	OutputFormat string
	IsHex        bool

	User     string
	Password string

	Debug bool
}

func endpointsFromCmd(cmd *cobra.Command) (string, error) {
	ep, err := cmd.Flags().GetString("endpoint")
	return ep, err
}

func getGlobalFlags(cmd *cobra.Command) *GlobalFlags {
	gf := &GlobalFlags{}

	ep, err := endpointsFromCmd(cmd)
	if err != nil {
		return gf
	}
	gf.Endpoint = ep

	return gf
}
