package main

import (
	"fmt"
	"log"
	"os"

	"github.com/benwebber/craftctl/config"
	"github.com/benwebber/craftctl/rcon"
	"github.com/codegangsta/cli"
)

//go:generate ./scripts/generate.py -i commands.txt -o init.go --format cli.go

func main() {
	app := cli.NewApp()
	app.Name = "craftctl"
	app.Usage = "Command-line Minecraft RCON client."
	app.Version = fmt.Sprintf("%v (%v)", Version, GitCommit)
	// Conflicts with console help command.
	app.HideHelp = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "host, H",
			Value:  "localhost",
			Usage:  "RCON host",
			EnvVar: "CRAFTCTL_HOST",
		},
		cli.IntFlag{
			Name:   "port, P",
			Value:  25575,
			Usage:  "RCON port",
			EnvVar: "CRAFTCTL_PORT",
		},
		cli.StringFlag{
			Name:   "password, p",
			Value:  "password",
			Usage:  "RCON password",
			EnvVar: "CRAFTCTL_PASSWORD",
		},
		// Hidden when app.HideHelp == true.
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "show help and exit",
		},
	}

	app.Action = func(ctx *cli.Context) {
		cfg, err := config.NewConfigFromContext(ctx)
		if err != nil {
			log.Fatal(err.Error())
		}

		client, err := rcon.NewClient(cfg)
		if err != nil {
			log.Fatal(err.Error())
		}

		if len(ctx.Args()) < 1 {
			os.Stderr.WriteString("craftctl: enter a command\n")
			os.Exit(1)
		}

		client.Auth()
		command, err := rcon.NewCommand(ctx.Args()...)
		if err != nil {
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(1)
		}

		resp := client.Execute(command)
		if resp.Err != nil {
			os.Stderr.WriteString(resp.Err.Error() + "\n")
			os.Exit(1)
		}
		fmt.Println(resp)
	}

	app.Run(os.Args)
}
