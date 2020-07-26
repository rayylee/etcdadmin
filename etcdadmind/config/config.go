package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const (
	configFile = "/etc/etcdadmin/etcdadmind.conf"
)

var (
	configContext = []string{
		"LOG_FILE=/var/log/etcdadmind/etcdadmind.log",
		"LOG_LEVEL=DEBUG",
		"GRPC_PORT=2390",
		"ETCD_WAL_DIR=/var/lib/etcd/default.etcd/member/wal",
		"ETCD_SNAP_DIR=/var/lib/etcd/default.etcd/member/snap",
		"ETCD_PEER_PORT=2380",
		"ETCD_CLIENT_PORT=2379",
		"ETCD_CONF_FILE=/etc/etcd/etcd.conf",
	}
)

var (
	once sync.Once
	cfg  *KvConfig
)

func CreateConfigFile() error {
	path := configFile
	context := configContext

	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	dirpath := filepath.Dir(path)
	_, err = os.Stat(path)
	if err != nil {
		err = os.MkdirAll(dirpath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		return err
	}
	for _, c := range context {
		f.WriteString(fmt.Sprintf("%s\n", c))
	}

	return err
}

func Init() *KvConfig {
	once.Do(func() {
		CreateConfigFile()
		cfg = Load(configFile)
	})
	return cfg
}

func RemoveConfigFile() error {
	path := configFile

	_, err := os.Stat(path)
	if err == nil {
		err = os.Remove(path)
	}

	return err
}
