package driver

import (
	"github.com/rayylee/etcdadmin/etcdadmind/config"
	"github.com/rayylee/etcdadmin/etcdadmind/server/driver/command"
	"github.com/rayylee/etcdadmin/etcdadmind/server/driver/etcdcfg"
)

func resetEtcdConfig() error {
	etcdCfgFile := config.Init().Get("ETCD_CONF_FILE")

	m, err := etcdcfg.EtcdConfigMapInit()

	if err != nil {
		return err
	}

	return etcdcfg.EtcdConfigWrite(etcdCfgFile, m)
}

func resetEtcd(isStart bool) error {
	var err error

	// Stop etcd but ignore error
	command.CmdEtcdctlStop()

	if err := resetEtcdConfig(); err != nil {
		goto exit
	}

	if err = etcdcfg.EtcdWalDelete(); err != nil {
		goto exit
	}

exit:
	if isStart == true {
		er := command.EtcdctlStart()
		if err == nil {
			err = er
		}
	}
	return err
}
