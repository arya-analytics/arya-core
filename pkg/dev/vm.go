package dev

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strconv"
	"strings"
)

// || GENERAL VM ||

type VM interface {
	Provision() error
	Exists() bool
	Delete() error
	Info() (VMInfo, error)
	Exec(cmdStr string) ([]byte, error)
	Transfer(direction TransferDirection, srcPath string, destPath string) error
}

type VMInfo struct {
	Name      string
	State     string
	IPv4      string
	Release   string
	ImageHash string
	Load      string
	Storage   string
	Memory    string
}

type VMConfig struct {
	Name    string
	Memory  int
	Cores   int
	Storage int
}

// || MULTIPASS VM ||

// Command to access multipass executable
const multipassExec = "multipass"

type MultipassVM struct {
	cfg VMConfig
}

func NewVM(cfg VMConfig) VM {
	return &MultipassVM{cfg}
}

func (vm MultipassVM) command(args ...string) *exec.Cmd {
	return exec.Command(multipassExec, args...)
}

func (vm MultipassVM) logFields(vb bool) log.Fields {
	f := log.Fields{
		"name": vm.cfg.Name,
	}
	if vb {
		f["memory"] = vm.cfg.Memory
		f["cores"] = vm.cfg.Cores
		f["storage"] = vm.cfg.Storage
	}
	return f
}

func (vm MultipassVM) Provision() error {
	log.WithFields(vm.logFields(true)).Trace("Provisioning a new multipass VM")
	args := []string{"launch", "--name", vm.cfg.Name}
	if vm.cfg.Memory != 0 {
		args = append(args, "--mem", strconv.Itoa(vm.cfg.Memory)+"g")
	}
	if vm.cfg.Cores != 0 {
		args = append(args, "--cpus", strconv.Itoa(vm.cfg.Cores))
	}
	if vm.cfg.Storage != 0 {
		args = append(args, "--disk", strconv.Itoa(vm.cfg.Storage)+"g")
	}
	if err := vm.command(args...).Run(); err != nil {
		log.WithFields(vm.logFields(true)).Error("Failed to provision new VM")
		return err
	}
	return nil
}

func (vm MultipassVM) Exists() bool {
	_, err := vm.Info()
	if err != nil {
		return false
	}
	return true
}

func (vm MultipassVM) Info() (VMInfo, error) {
	var info VMInfo
	o, err := vm.command("info", vm.cfg.Name).Output()
	if err != nil {
		log.WithFields(vm.logFields(false)).Trace("Couldn't find VM")
		return info, fmt.Errorf("couldn't find VM named %s", vm.cfg.Name)
	}
	rawInfoChain := strings.Split(string(o[:]), "\n")
	var parsedInfo [12]string

	// Using a manually defined i to ensure values are placed at correct index
	i := 0
	for _, ri := range rawInfoChain[:len(rawInfoChain)-1] {
		kv := strings.Split(ri, ":")
		if len(kv) == 2 {
			parsedInfo[i] = strings.Trim(kv[1], " ")
			i += 1
		} else {
			log.WithFields(vm.logFields(false)).Warn("Encountered unknown VM info")
		}
	}
	info.Name = parsedInfo[0]
	info.State = parsedInfo[1]
	info.IPv4 = parsedInfo[2]
	info.Release = parsedInfo[3]
	info.ImageHash = parsedInfo[4]
	info.Load = parsedInfo[5]
	info.Storage = parsedInfo[6]
	info.Memory = parsedInfo[7]
	return info, nil
}

func (vm MultipassVM) Delete() error {
	f := vm.logFields(false)
	if err := vm.command("delete", vm.cfg.Name).Run(); err != nil {
		log.WithFields(f).Error("Failed to delete VM")
		return err
	}
	if err := vm.command("purge").Run(); err != nil {
		log.WithFields(f).Error("Failed to purge deleted VM")
		return err
	}
	log.WithFields(f).Trace("Successfully deleted VM")
	return nil
}

func (vm MultipassVM) Exec(cmdStr string) ([]byte, error) {
	var outb, errb bytes.Buffer
	cmd := vm.command("exec", vm.cfg.Name, "--", "bash")
	w, _ := cmd.StdinPipe()
	cmd.Stdout, cmd.Stderr = &outb, &errb
	err := cmd.Start()
	if err != nil {
		return errb.Bytes(), err
	}
	_, err = w.Write([]byte(cmdStr + "\n"))
	if err != nil {
		return errb.Bytes(), err
	}
	_, err = w.Write([]byte("exit" + "\n"))
	if err != nil {
		return errb.Bytes(), err
	}
	err = cmd.Wait()
	if err != nil {
		return errb.Bytes(), err
	}
	return outb.Bytes(), nil
}

type TransferDirection int

const (
	TransferTo TransferDirection = iota
	TransferFrom
)

func (vm MultipassVM) Transfer(transfer TransferDirection, srcPath, destPath string) error {
	if transfer == TransferTo {
		destPath = fmt.Sprintf("%s:%s", vm.cfg.Name, destPath)
	} else {
		srcPath = fmt.Sprintf("%s:%s", vm.cfg.Name, srcPath)
	}
	var errb bytes.Buffer
	cmd := vm.command("transfer", srcPath, destPath)
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		log.WithFields(vm.logFields(false)).Error(errb.String())
		return err
	}
	return nil
}
