package dev

import (
	"github.com/urfave/cli/v2"
	"strconv"
	"strings"
)

var Cmd = &cli.Command{
	Name: "dev",
	Usage: "Provides access to development services such as cluster provisioning, " +
		"tooling installs, and configuration management.",
	Subcommands: []*cli.Command{
		clusterCmd,
		toolingCmd,
		reloaderCmd,
		configCmd,
		loginCmd,
	},
}

var clusterCmd = &cli.Command{
	Name:  "cluster",
	Usage: "Provision and manage development clusters.",
}

func emoji(s string) string {
	r, _ := strconv.ParseInt(strings.TrimPrefix(s, "\\U"), 16, 32)
	return string(r)
}

var toolingCmd = &cli.Command{
	Name:  "tooling",
	Usage: "Install and manage development tools.",
	Subcommands: []*cli.Command{
		{
			Name:  "install",
			Usage: "Install development tools.",
			Action: func(c *cli.Context) error {
				return InstallRequired()
			},
		},
		{
			Name:  "uninstall",
			Usage: "Uninstall development tools.",
			Action: func(c *cli.Context) error {
				return UninstallRequired()
			},
		},
	},
}

var reloaderCmd = &cli.Command{
	Name:  "reloader",
	Usage: "Operate the development hot-reloader.",
	Action: func(c *cli.Context) error {
		return nil
	},
}

var configCmd = &cli.Command{
	Name:  "config",
	Usage: "Manage development configurations.",
}

var loginCmd = &cli.Command{
	Name:  "login",
	Usage: "Login to Github and register credentials in config",
}
