package command

import (
	"fmt"
	"testing"
)

func TestCmdExecNoBlock(t *testing.T) {
	env := []string{}
	name := "date"
	result := cmdExec(env, name, "-R")
	if len(result.stdout) > 0 {
		fmt.Printf("stdout: %v\n", result.stdout)
	}
	if len(result.stderr) > 0 {
		fmt.Printf("stderr: %v\n", result.stderr)
	}
	if result.err != nil {
		fmt.Printf("err: %v\n", result.err)
	}
}

func TestCmdExecBlock(t *testing.T) {
	t.Skip("skipping: take so long to test")

	env := []string{}
	name := "ping"
	result := cmdExec(env, name, "127.0.0.1")
	if len(result.stdout) > 0 {
		fmt.Printf("stdout: %v\n", result.stdout)
	}
	if len(result.stderr) > 0 {
		fmt.Printf("stderr: %v\n", result.stderr)
	}
	if result.err != nil {
		fmt.Printf("err: %v\n", result.err)
	}
}
func TestCmdEtcdVersion(t *testing.T) {
	result := CmdEtcdctl("version")
	if len(result.stdout) > 0 {
		fmt.Printf("stdout: %v\n", result.stdout)
	}
	if len(result.stderr) > 0 {
		fmt.Printf("stderr: %v\n", result.stderr)
	}
	if result.err != nil {
		fmt.Printf("err: %v\n", result.err)
	}
}

func TestCmdEtcdMemberAdd(t *testing.T) {
	t.Skip("skipping: do not add member")

	result := CmdEtcdctlMemberAdd("test-name", "127.0.0.2")

	if len(result.stdout) > 0 {
		fmt.Printf("stdout: %v\n", result.stdout)
	}
	if len(result.stderr) > 0 {
		fmt.Printf("stderr: %v\n", result.stderr)
	}
	if result.err != nil {
		fmt.Printf("err: %v\n", result.err)
	}
}

func TestCmdEtcdctlMemberList(t *testing.T) {
	result := CmdEtcdctlMemberList()

	if len(result.stdout) > 0 {
		fmt.Printf("stdout: %v\n", result.stdout)
	}
	if len(result.stderr) > 0 {
		fmt.Printf("stderr: %v\n", result.stderr)
	}
	if result.err != nil {
		fmt.Printf("err: %v\n", result.err)
	}

}
