package caseconv_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCaseconv(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Caseconv Suite")
}
