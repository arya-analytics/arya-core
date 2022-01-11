package main

import (
	"github.com/arya-analytics/aryacore/pkg/dev"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "arya-core",
		Usage: "Hello",
		Commands: []*cli.Command{
			dev.Cmd,
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
