package etcdcfg

import (
	"fmt"
	"github.com/rayylee/etcdadmin/etcdadmind/config"
	"os"
	"path/filepath"
	"sync"
)

var (
	mutexEtcdConfig sync.Mutex
)

func EtcdConfigWrite(path string, m map[string]string) error {
	var err error

	dirname := filepath.Dir(path)

	mutexEtcdConfig.Lock()

	_, err = os.Stat(path)
	if err != nil {
		os.MkdirAll(dirname, os.ModePerm)
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		goto exit
	}
	for k := range m {
		f.WriteString(fmt.Sprintf("%s=\"%s\"\n", k, m[k]))
	}

exit:
	mutexEtcdConfig.Unlock()

	return err
}

func EtcdConfigMapInit() (map[string]string, error) {

	cfgServer := config.Init()

	lisP := fmt.Sprintf("http://0.0.0.0:%s", cfgServer.Get("ETCD_PEER_PORT"))
	lisC := fmt.Sprintf("http://0.0.0.0:%s", cfgServer.Get("ETCD_CLIENT_PORT"))
	name, err := os.Hostname()
	if err != nil {
		return map[string]string{}, err
	}

	ip := "127.0.0.1"
	advP := fmt.Sprintf("http://%s:%s", ip, cfgServer.Get("ETCD_PEER_PORT"))
	advC := fmt.Sprintf("http://%s:%s", ip, cfgServer.Get("ETCD_CLIENT_PORT"))
	initCluster := fmt.Sprintf("%s=http://%s:%s", name, ip,
		cfgServer.Get("ETCD_PEER_PORT"))

	m := map[string]string{
		"ETCD_DATA_DIR":                    "/var/lib/etcd/default.etcd",
		"ETCD_LISTEN_PEER_URLS":            lisP,
		"ETCD_LISTEN_CLIENT_URLS":          lisC,
		"ETCD_NAME":                        name,
		"ETCD_INITIAL_ADVERTISE_PEER_URLS": advP,
		"ETCD_ADVERTISE_CLIENT_URLS":       advC,
		"ETCD_INITIAL_CLUSTER":             initCluster,
		"ETCD_INITIAL_CLUSTER_STATE":       "new",
		"ETCD_INITIAL_CLUSTER_TOKEN":       "etcd-cluster",
	}

	return m, nil
}

func EtcdConfigReset(file string) error {
	m, err := EtcdConfigMapInit()

	if err != nil {
		return err
	}

	return EtcdConfigWrite(file, m)
}

func EtcdWalDelete() error {
	cfgServer := config.Init()

	walDir := cfgServer.Get("ETCD_WAL_DIR")
	snapDir := cfgServer.Get("ETCD_SNAP_DIR")

	_, err := os.Stat(walDir)
	if err == nil {
		if os.RemoveAll(walDir); err != nil {
			goto exit
		}
	}

	_, err = os.Stat(snapDir)
	if err == nil {
		if os.RemoveAll(snapDir); err != nil {
			goto exit
		}
	}

exit:
	return err
}
