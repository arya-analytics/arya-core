package chanstream_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestChanstream(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chanstream Suite")
}
