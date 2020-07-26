package daemon

import (
	"errors"
	"os"
	"runtime"
	"strings"
)

// Status constants.
const (
	statNotInstalled = "Service not installed"
)

// Daemon interface has a standard set of methods/commands
type Daemon interface {
	// GetTemplate - gets service config template
	GetTemplate() string

	// SetTemplate - sets service config template
	SetTemplate(string) error

	// Install the service into the system
	Install(args ...string) (string, error)

	// Remove the service and all corresponding files from the system
	Remove() (string, error)

	// Start the service
	Start() (string, error)

	// Stop the service
	Stop() (string, error)

	// Status - check the service status
	Status() (string, error)

	// Run - run executable service
	Run(e Executable) (string, error)
}

// Executable interface defines controlling methods of executable service
type Executable interface {
	// Start - non-blocking start service
	Start()
	// Stop - non-blocking stop service
	Stop()
	// Run - blocking run service
	Run()
}

// Kind is type of the daemon
type Kind string

const (
	// SystemDaemon is a system daemon that runs as the root user.
	SystemDaemon Kind = "SystemDaemon"
)

// New - Create a new daemon
//
// name: name of the service
//
// description: any explanation, what is the service, its purpose
//
// kind: what kind of daemon to create
func New(name, description string, kind Kind,
	dependencies ...string) (Daemon, error) {

	switch runtime.GOOS {
	case "linux":
		if kind != SystemDaemon {
			return nil, errors.New("Invalid daemon kind specified")
		}
	}

	return newDaemon(strings.Join(strings.Fields(name), "_"), description,
		kind, dependencies)
}

// Get the daemon properly
func newDaemon(name, description string, kind Kind,
	dependencies []string) (Daemon, error) {

	// newer subsystem must be checked first
	if _, err := os.Stat("/run/systemd/system"); err == nil {
		return &systemDRecord{name, description, kind, dependencies}, nil
	}
	return &systemDRecord{name, description, kind, dependencies}, nil
}
