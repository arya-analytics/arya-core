package main

import (
	"fmt"
	"github.com/arya-analytics/aryacore"
	"github.com/arya-analytics/aryacore/cmd/migrations"
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
					core, ctx := aryacore.NewCore(c.Context, aryacore.GetConfig())
					db := core.ConnManager.GetOrCreate("default").(*bun.DB)
					migrator := migrate.NewMigrator(db, migrations)
					return migrator.Init(ctx)
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go Migration",
				Action: func(c *cli.Context) error {
					core, ctx := aryacore.NewCore(c.Context, aryacore.GetConfig())
					db := core.ConnManager.GetOrCreate("default").(*bun.DB)
					migrator := migrate.NewMigrator(db, migrations)
					name := strings.Join(c.Args().Slice(), "_")
					mf, err := migrator.CreateGoMigration(ctx, name)
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
					core, ctx := aryacore.NewCore(c.Context, aryacore.GetConfig())
					db := core.ConnManager.GetOrCreate("default").(*bun.DB)
					migrator := migrate.NewMigrator(db, migrations)
					group, err := migrator.Migrate(ctx)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Println("there are no new migrations to run")
						return nil
					}
					fmt.Printf("migrated to %s \n", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					core, ctx := aryacore.NewCore(c.Context, aryacore.GetConfig())
					db := core.ConnManager.GetOrCreate("default").(*bun.DB)
					migrator := migrate.NewMigrator(db, migrations)
					group, err := migrator.Rollback(ctx)
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

func isServerClosed(err error) bool {
	return err.Error() == "http: Server closed"
}
