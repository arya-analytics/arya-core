package chanchunk_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestChanchunk(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Chanchunk Suite")
}
