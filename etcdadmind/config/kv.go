package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// config file with key=value
type KvConfig struct {
	FilePath string
	KvMap    map[string]string
	rwMutex  *sync.RWMutex
}

func Load(path string) *KvConfig {
	kv := new(KvConfig)

	kv.FilePath = path
	kv.rwMutex = new(sync.RWMutex)
	configMap := make(map[string]string)

	f, err := os.OpenFile(path, os.O_RDONLY, 0666)
	defer f.Close()
	if err != nil {
		return nil
	}

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil
		}
		s := strings.TrimSpace(string(b))
		index := strings.Index(s, "=")
		if index < 0 || s[0] == '#' {
			continue
		}
		key := strings.TrimSpace(s[:index])
		if len(key) == 0 {
			continue
		}
		value := strings.TrimSpace(s[index+1:])
		if len(value) == 0 {
			continue
		}
		configMap[key] = value
	}
	kv.KvMap = configMap
	return kv
}

func (kv *KvConfig) Get(key string) string {
	kv.rwMutex.RLock()
	value := kv.KvMap[key]
	kv.rwMutex.RUnlock()
	return value
}

func (kv *KvConfig) Set(key string, value string) error {
	kv.rwMutex.Lock()
	kv.KvMap[key] = value
	f, err := os.OpenFile(kv.FilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer f.Close()
	if err != nil {
		goto exit
	}
	for k := range kv.KvMap {
		f.WriteString(fmt.Sprintf("%s=%s\n", k, kv.KvMap[k]))
	}

exit:
	kv.rwMutex.Unlock()
	return err
}
