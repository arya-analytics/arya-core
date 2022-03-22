package tsquery_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTsquery(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tsquery Suite")
}
