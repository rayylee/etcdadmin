package etcdadmind

import (
	"fmt"
	"github.com/rayylee/etcdadmin/etcdadmind/config"
	"github.com/rayylee/etcdadmin/etcdadmind/daemon"
	"github.com/rayylee/etcdadmin/etcdadmind/log"
	"github.com/rayylee/etcdadmin/etcdadmind/server"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

const (
	// name of the service
	serviceName = "etcdadmind"

	// service description
	serviceDescription = "Etcd admin service"
)

var (
	logger              *zap.Logger
	serviceDependencies = []string{}
)

// Manage run the daemon
func (service *Service) Manage() (string, error) {
	usage := "Usage: " + serviceName +
		" install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			config.CreateConfigFile()
			return service.Install()
		case "remove":
			config.RemoveConfigFile()
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	// Do something, call goroutines, etc
	run()

	// Set up channel on which to send signal notifications.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// loop work cycle with accept interrupt by system signal
	for {
		select {
		case killSignal := <-interrupt:
			fmt.Println("Got signal:", killSignal)
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	// never happen, but need to complete code
	return "Exit", nil
}

func run() {
	cfg := config.Init()

	log.Init(log.Config{
		Level: cfg.Get("LOG_LEVEL"),
		File:  cfg.Get("LOG_FILE"),
	})
	logger = log.GetLogger()

	logger.Info("etcdadmind start.")
	server.Init(cfg.Get("GRPC_PORT"))
}

func Main() {
	srv, err := daemon.New(serviceName, serviceDescription,
		daemon.SystemDaemon, serviceDependencies...)
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}

	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		fmt.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
