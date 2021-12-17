package dev

import (
	"fmt"
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
	Name: "cluster",
	Usage: "Provision and manage development clusters.",
}

func emoji(s string) string {
	r, _ := strconv.ParseInt(strings.TrimPrefix(s, "\\U"), 16, 32)
	return string(r)
}

var toolingCmd = &cli.Command{
	Name: "tooling",
	Usage: "Install and manage development tools.",
	Subcommands: []*cli.Command{
		{
			Name: "install",
			Usage: "Install development tools.",
			Action: func(c *cli.Context) error {
				fmt.Printf("%s Sit back and relax while we install the required" +
					" development tools \n", emoji("\\U0001F6E0"))
				t := NewTooling()
				for k := range RequiredTools {
					installed := t.CheckInstalled(k)
					if installed {
						fmt.Printf("%s %s is already installed. " +
							"Skipping re-install. \n", emoji("\\U0001F438"), k)
					} else {
						fmt.Printf("%s Installing %s. \n", emoji("\\U0001f525"), k)
						if err := t.Install(k); err != nil {
							panic(err)
						}
						fmt.Printf("%s  Successfully installed %s. \n",
							emoji("\\U0002705"), k)
					}
				}
				fmt.Printf("%s  All Done! \n", emoji("\\U0002705"))
				return nil
			},
		},
	},
}

var reloaderCmd = &cli.Command{
	Name: "reloader",
	Usage: "Operate the development hot-reloader.",
	Action: func(c *cli.Context) error {
		return nil
	},
}

var configCmd = &cli.Command{
	Name: "config",
	Usage: "Manage development configurations.",
}

var loginCmd = &cli.Command{
	Name: "login",
	Usage: "Login to Github and register credentials in config",
}
