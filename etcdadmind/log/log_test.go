package log

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"testing"
)

var (
	gLogFile  = "test.log"
	gLogLevel = "debug"
)

func setup() {
	logCfg := Config{
		Level: gLogLevel,
		File:  gLogFile,
	}
	Init(logCfg)
}

func teardown() {
	_, err := os.Stat(gLogFile)
	if err == nil {
		os.Remove(gLogFile)
	}
}

func TestCase(t *testing.T) {
	logger := GetLogger()

	logger.Debug("hello debug")
	logger.Info("hello info")
	logger.Warn("hello warn")
	logger.Error("hello error")

	f, err := os.OpenFile(gLogFile, os.O_RDONLY, 0666)
	defer f.Close()
	if err != nil {
		t.Fatal(fmt.Sprintf("Open file: %s", gLogFile))
	}

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		fmt.Printf("%v\n", string(b))
	}
}

// Test Entry
func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	teardown()
	os.Exit(ret)
}
