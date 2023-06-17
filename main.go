package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/midoks/nezha/cmd"
	"github.com/midoks/nezha/service/singleton"
)

func main() {
	app := cli.NewApp()
	app.Name = "Nezha Monitoring"
	app.Version = singleton.Version
	app.Usage = "Self-hostable, lightweight, servers and websites monitoring and O&M tool."
	app.Commands = []cli.Command{
		cmd.Web,
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
