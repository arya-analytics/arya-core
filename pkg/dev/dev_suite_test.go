package dev_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDev(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dev Suite")
}
