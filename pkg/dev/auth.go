package dev

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type AuthConfig map[string]map[string]string

type Config struct {
	StackOrchestrator string `json:"stackOrchestrator"`
	CredsStore string `json:"credsStore"`
	Auths AuthConfig `json:"auths"`
	Experimental string `json:"experimental"`
}

func CheckFileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

func ConstructConfig() {
	path := filepath.Join(
		os.ExpandEnv("$HOME"),
		".arya",
	)
	fmt.Println(path)
	if err := os.Mkdir(path, os.FileMode(int(0777)),
		); err != nil {
		panic(err)
	}
}

func Login(ghToken string) {
	auth := AuthConfig{"ghcr.io": {"auth": ghToken}}

	cfg := Config{
		StackOrchestrator: "swarm",
		Auths: auth,
		CredsStore: "osxkeychain",
		Experimental: "disabled",
	}
	dat, err := json.Marshal(cfg)
	if err != nil {
		panic(err)
	}
	path := filepath.Join(
		os.ExpandEnv("$HOME"),
		".arya",
		"config.json",
	)
	err = ioutil.WriteFile(path, dat, 0777)
}
