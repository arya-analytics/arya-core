package dev_test

import (
	"github.com/arya-analytics/aryacore/pkg/dev"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var vm dev.VM
var vmInfo dev.VMInfo

var vmCfg = dev.VMConfig{
	Name:    "adtst1",
	Memory:  2,
	Cores:   3,
	Storage: 4,
}

var _ = BeforeSuite(func() {
	vm = dev.NewVM(vmCfg)
	if vm.Exists() {
		if err := vm.Delete(); err != nil {
			log.Fatalln("Failed to delete test VM")
		}
	}
	if err := vm.Provision(); err != nil {
		log.Fatalln("Failed to provision test VM")
	}
	vi, err := vm.Info()
	vmInfo = vi
	if err != nil {
		log.Fatalln("Failed to pull info from test vm")
	}
})

var _ = AfterSuite(func() {
	vm = dev.NewVM(vmCfg)
	if vm.Exists() {
		if err := vm.Delete(); err != nil {
			log.Fatalln("Failed to delete test VM")
		}
	}
})

var _ = Describe("Vm", func() {
	Context("Multipass VM", func() {
		Describe("Provisioning new VM", func() {
			It("Should assign the correct name ", func() {
				Expect(vmInfo.Name).To(Equal(vmCfg.Name))
			})
			It("Should be in a running state", func() {
				Expect(vmInfo.State).To(Equal("Running"))
			})
		})
		Describe("Checking if a VM exists", func() {
			It("Should return false when the VM doesnt exists", func() {
				nonExistentVm := dev.NewVM(dev.VMConfig{Name: "definitelydoesnotexist"})
				Expect(nonExistentVm.Exists()).To(BeFalse())
			})
			It("Should return true when the VM does exist", func() {
				Expect(vm.Exists()).To(BeTrue())
			})
		})
		Describe("Getting VM info", func() {
			It("Should return an IPV4 address", func() {
				Expect(vmInfo.IPv4).To(MatchRegexp("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$"))
			})
			It("Should return a valid release", func() {
				Expect(vmInfo.Release).To(ContainSubstring("Ubuntu"))
			})
			It("Should return a correctly formatted storage usage", func() {
				Expect(vmInfo.Storage).To(ContainSubstring("out of"))
			})
			It("Should return a correctly formatted memory usage", func() {
				Expect(vmInfo.Memory).To(ContainSubstring("out of"))
			})
			It("Should return a valid image hash", func() {
				Expect(vmInfo.ImageHash).To(ContainSubstring("(Ubuntu"))
			})
			It("Should return an error when the VM can't be found", func() {
				nonExistentVm := dev.NewVM(dev.VMConfig{Name: "definitelydoesnotexist"})
				_, err := nonExistentVm.Info()
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
