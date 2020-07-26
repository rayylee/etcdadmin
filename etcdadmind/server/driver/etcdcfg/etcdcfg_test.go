package etcdcfg

import (
	"fmt"
	"testing"
)

func TestCreateEtcdConfigMap(t *testing.T) {
	m, _ := EtcdConfigMapInit()
	m["ETCD_INITIAL_CLUSTER"] = "modify-value"
	m["ADD_KEY"] = "add-value"
	delete(m, "ETCD_DATA_DIR")
	fmt.Printf("rest map: %v\n", m)
}
