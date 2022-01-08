package dev

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type AuthConfig map[string]map[string]string

type Config struct {
	StackOrchestrator string     `json:"stackOrchestrator"`
	CredsStore        string     `json:"credsStore"`
	Auths             AuthConfig `json:"auths"`
	Experimental      string     `json:"experimental"`
}

// Login creates/updates arya's config.
// json file and supplies the appropriate credentials.
func Login(ghToken string) error {
	auth := AuthConfig{"ghcr.io": {"auth": ghToken}}
	cfg := Config{
		StackOrchestrator: "swarm",
		Auths:             auth,
		CredsStore:        "osxkeychain",
		Experimental:      "disabled",
	}
	dat, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	path := filepath.Join(
		os.ExpandEnv("$HOME"),
		".arya",
		"config.json",
	)
	err = ioutil.WriteFile(path, dat, 0777)
	return err
}
