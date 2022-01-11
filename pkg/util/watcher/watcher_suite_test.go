package watcher_test

import (
	log "github.com/sirupsen/logrus"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestWatcher(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Watcher Suite")
}

const tmpDirBasePath = "/tmp/"
var tmpDir string

var _ = BeforeSuite(func() {
	var err error
	tmpDir, err = os.MkdirTemp(tmpDirBasePath, "*")
	if err != nil {
		log.Fatalln(err)
	}
})

var _ = AfterSuite(func() {
	if err := os.Remove(tmpDir); err != nil {
		log.Fatalln(err)
	}
})