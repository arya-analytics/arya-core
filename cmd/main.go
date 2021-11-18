package main

import (
	"fmt"
	"github.com/arya-analytics/aryacore/cmd/migrations"
	"github.com/arya-analytics/aryacore/config"
	"github.com/arya-analytics/aryacore/server"
	"github.com/arya-analytics/aryacore/telem/live"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"strings"
)

func main() {
	app := &cli.App{
		Name: "aryacore",
		Commands: []*cli.Command{
			newDBCommand(migrations.Migrations),
			serverCommand,
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func newDBCommand(migrations *migrate.Migrations) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					sv := server.New(config.GetConfig())
					db := sv.Context.Pooler.GetOrCreate("aryadb").(*bun.DB)
					migrator := migrate.NewMigrator(db, migrations)
					return migrator.Init(c.Context)
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go Migration",
				Action: func(c *cli.Context) error {
					sv := server.New(config.GetConfig())
					db := sv.Context.Pooler.GetOrCreate("aryadb").(*bun.DB)
					migrator := migrate.NewMigrator(db, migrations)
					name := strings.Join(c.Args().Slice(), "_")
					mf, err := migrator.CreateGoMigration(c.Context, name)
					if err != nil {
						return err
					}
					fmt.Printf("Created migratijons %s (%s) \n", mf.Name, mf.Path)
					return nil
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					sv := server.New(config.GetConfig())
					db := sv.Context.Pooler.GetOrCreate("aryadb").(*bun.DB)
					migrator := migrate.NewMigrator(db, migrations)
					group, err := migrator.Migrate(c.Context)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Println("there are no new migrations to run")
						return nil
					}
					fmt.Printf("migrated to #{group} \n")
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					sv := server.New(config.GetConfig())
					db := sv.Context.Pooler.GetOrCreate("aryadb").(*bun.DB)
					migrator := migrate.NewMigrator(db, migrations)
					group, err := migrator.Rollback(c.Context)
					if err != nil {
						return err
					}
					if group.ID == 0 {
						fmt.Println("there are no groups to roll back")
						return nil
					}
					fmt.Printf("rolled back %s \n", group)
					return nil
				},
			},
		},
	}
}

var serverCommand = &cli.Command{
	Name:  "server",
	Usage: "control server",
	Subcommands: []*cli.Command{{
		Name:  "start",
		Usage: "start server",
		Action: func(c *cli.Context) error {
			sv := server.New(config.GetConfig())
			sv.BindSlice(live.API)
			sv.Start()
			return nil
		},
	},
	},
}
