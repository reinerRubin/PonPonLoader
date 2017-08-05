package main

import (
	"os"

	ponpon "github.com/PonPonLoader"
	"github.com/codegangsta/cli"
)

var (
	cliFlags = []cli.Flag{
		cli.BoolFlag{
			Name:  "watch",
			Usage: "watch for new images",
		},
	}
)

func main() {
	app := cli.NewApp()
	app.Flags = cliFlags

	app.Action = func(c *cli.Context) error {
		app, err := ponpon.NewApp(c)
		if err != nil {
			return err
		}

		return app.Run()
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
