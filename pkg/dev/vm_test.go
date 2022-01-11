package dev_test

import (
	"github.com/arya-analytics/aryacore/pkg/dev"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("VM", func() {
	Describe("Multipass VM", func() {
		Describe("Provisioning new VM", func() {
			Context("When the VM doesn'tools exist", func() {
				It("Should assign the correct Name ", func() {
					Expect(vmInfo.Name).To(Equal(vmCfg.Name))
				})
				It("`Should be in a running state", func() {
					Expect(vmInfo.State).To(Equal("Running"))
				})
			})
			Context("When the VM already exists", func() {
				It("Should throw an error", func() {
					existingVm := dev.NewVM(dev.VMConfig{Name: vmCfg.Name})
					err := existingVm.Provision()
					Expect(err).ToNot(BeNil())
				})
			})
		})
		Describe("Deleting VM", func() {
			It("(SLOW) Should delete and purge a VM", func() {
				tempVM := dev.NewVM(dev.VMConfig{Name: "testtempvm1"})
				if err := tempVM.Provision(); err != nil {
					Fail("Failed to provision VM")
				}
				if !vm.Exists() {
					Fail("Failed to provision VM")
				}
				if err := tempVM.Delete(); err != nil {
					Fail("Failed to delete Vm")
				}
				Expect(tempVM.Exists()).To(BeFalse())
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
			Context("When a VM exists", func() {
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
			})
			Context("When a VM doesn'tools exist", func() {
				It("Should return an error", func() {
					nonExistentVm := dev.NewVM(dev.VMConfig{Name: "doesnotexist"})
					_, err := nonExistentVm.Info()
					Expect(err).ToNot(BeNil())
				})
			})
		})
		Describe("Executing commands", func() {
			It("Should return the correct command output", func() {
				cmdString := "pwd"
				o, err := vm.Exec(cmdString)
				if err != nil {
					Fail(err.Error())
				}
				Expect(string(o[:])).To(Equal("/home/ubuntu\n"))
			})
			It("Should return the error output of a failed command", func() {
				cmdString := "lsawdaw"
				o, err := vm.Exec(cmdString)
				Expect(err).ToNot(BeNil())
				Expect(string(o[:])).To(Equal("bash: line 1: lsawdaw: command not found\n"))
			})
		})
	})
})
