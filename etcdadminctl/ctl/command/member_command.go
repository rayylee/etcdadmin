package command

import (
	"errors"
	"fmt"
	"github.com/rayylee/etcdadmin/etcdadmind/client"
	"strings"

	"github.com/spf13/cobra"
)

var (
	memberPeerIP  string
	newMemberName string
)

// NewMemberCommand returns the cobra command for "member".
func NewMemberCommand() *cobra.Command {
	mc := &cobra.Command{
		Use:   "member <subcommand>",
		Short: "Membership related commands",
	}

	mc.AddCommand(NewMemberAddCommand())
	mc.AddCommand(NewMemberListCommand())
	mc.AddCommand(NewMemberRemoveCommand())

	return mc
}

// NewMemberAddCommand returns the cobra command for "member add".
func NewMemberAddCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "add <name>",
		Short: "Adds a member into the cluster",

		Run: memberAddCommandFunc,
	}

	cc.Flags().StringVar(&memberPeerIP, "peer-ip", "", "the ip address of new member.")

	return cc
}

// NewMemberListCommand returns the cobra command for "member list".
func NewMemberListCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "list",
		Short: "Lists all members in the cluster",
		Long: `When --write-out is set to simple, this command prints out comma-separated member lists for each endpoint.
The items in the lists are ID, Status, Name, Peer Addrs, Client Addrs, Is Learner.
`,

		Run: memberListCommandFunc,
	}

	return cc
}

// NewMemberRemoveCommand returns the cobra command for "member remove".
func NewMemberRemoveCommand() *cobra.Command {
	cc := &cobra.Command{
		Use:   "remove <name>",
		Short: "Removes a member from the cluster",

		Run: memberRemoveCommandFunc,
	}

	return cc
}

// memberAddCommandFunc executes the "member add" command.
func memberAddCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		ExitWithError(ExitBadArgs, errors.New("member name not provided"))
	}
	if len(args) > 1 {
		ev := "too many arguments"
		ev += fmt.Sprintf(", did you mean --peer-ip=%s", args[1])
		ExitWithError(ExitBadArgs, errors.New(ev))
	}
	newMemberName = args[0]

	if len(memberPeerIP) == 0 {
		ExitWithError(ExitBadArgs, errors.New("member peer ip not provided"))
	}

	fmt.Printf("peer ip is: %v\n", memberPeerIP)

	gf := getGlobalFlags(cmd)
	fmt.Printf("endpoint: %v\n", gf.Endpoint)

	s := strings.Split(gf.Endpoint, ":")
	if len(s) == 2 {
		c := client.New(s[0], s[1])
		defer client.Release(c)
		c.GrpcClientAddmember(newMemberName, memberPeerIP)
	} else {
		fmt.Printf("error endpoint: %v\n", gf.Endpoint)
	}
}

// memberListCommandFunc executes the "member list" command.
func memberListCommandFunc(cmd *cobra.Command, args []string) {
	fmt.Printf("Display members\n")

	gf := getGlobalFlags(cmd)
	fmt.Printf("endpoint: %v\n", gf.Endpoint)

	s := strings.Split(gf.Endpoint, ":")
	if len(s) == 2 {
		c := client.New(s[0], s[1])
		defer client.Release(c)
		m, err := c.GrpcClientListmember()
		if err == nil {
			fmt.Printf("\nMembers:\n")
			for i := range m {
				fmt.Printf("name:%s ip:%s\n", i, m[i])
			}
		}
	} else {
		fmt.Printf("error endpoint: %v\n", gf.Endpoint)
	}
}

// memberRemoveCommandFunc executes the "member remove" command.
func memberRemoveCommandFunc(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		ExitWithError(ExitBadArgs, fmt.Errorf("member name is not provided"))
	}

	fmt.Printf("Remove member name: %s\n", args[0])
	name := args[0]

	gf := getGlobalFlags(cmd)
	fmt.Printf("endpoint: %v\n", gf.Endpoint)

	s := strings.Split(gf.Endpoint, ":")
	if len(s) == 2 {
		c := client.New(s[0], s[1])
		defer client.Release(c)
		err := c.GrpcClientRemovemember(name)

		if err != nil {
			fmt.Printf("%v\n", err)
		}
	} else {
		fmt.Printf("error endpoint: %v\n", gf.Endpoint)

	}
}
