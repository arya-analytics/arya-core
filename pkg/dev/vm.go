package dev

import (
	"bytes"
	"fmt"
	"os"
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
	Exec(cmdString string) ([]byte, error)
	Transfer(direction TransferDirection, srcPath string, destPath string) error
}

type VMInfo struct {
	Name      string
	State     string
	IPv4      string
	Release   string
	ImageHash string
	Load      string
	Disk      string
	Memory    string
}

type VMConfig struct {
	Name    string
	Memory  int
	Cores   int
	Storage int
}

// || MULTIPASS VM ||

const multipassBaseCmd string = "multipass"

type MultipassVM struct {
	cfg VMConfig
}

func NewVM(cfg VMConfig) VM {
	return &MultipassVM{cfg}
}

func (vm MultipassVM) command(args ...string) *exec.Cmd {
	c := exec.Command(multipassBaseCmd, args...)
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c
}

func (vm MultipassVM) Provision() error {
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
	return vm.command(args...).Run()
}

func (vm MultipassVM) Exists() bool {
	_, err := vm.Info()
	if err != nil {
		return false
	}
	return true
}

func (vm MultipassVM) Info() (VMInfo, error) {
	o, err := exec.Command("multipass", "info", vm.cfg.Name).Output()
	if err != nil {
		return VMInfo{}, fmt.Errorf("could not find vm with name %s",
			vm.cfg.Name)
	}
	infoStrings := strings.Split(string(o[:]), "\n")
	var parsedInfo = []string{}
	for _, v := range infoStrings[:len(infoStrings)-1] {
		splitV := strings.Split(v, ":")
		if len(splitV) == 2 {
			parsedInfo = append(parsedInfo, strings.Trim(splitV[1], " "))
		}
	}
	i := VMInfo{
		Name:      parsedInfo[0],
		State:     parsedInfo[1],
		IPv4:      parsedInfo[2],
		Release:   parsedInfo[3],
		ImageHash: parsedInfo[4],
		Load:      parsedInfo[5],
		Disk:      parsedInfo[6],
		Memory:    parsedInfo[7],
	}
	return i, nil
}

func (vm MultipassVM) Delete() error {
	if err := vm.command("delete", vm.cfg.Name).Run(); err != nil {
		return err
	}
	if err := vm.command("purge").Run(); err != nil {
		return err
	}
	return nil
}

func (vm MultipassVM) Exec(cmdString string) ([]byte, error) {
	var outb, errb bytes.Buffer
	cmd := vm.command("exec", vm.cfg.Name, "--", "bash")
	cmdWriter, _ := cmd.StdinPipe()
	err := cmd.Start()
	if err != nil {
		return outb.Bytes(), err
	}
	_, err = cmdWriter.Write([]byte(cmdString + "\n"))
	if err != nil {
		return errb.Bytes(), err
	}
	_, err = cmdWriter.Write([]byte("exit" + "\n"))
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
	fmt.Println(srcPath, destPath)
	cmd := vm.command("transfer", srcPath, destPath)
	return cmd.Run()
}
