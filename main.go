package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/benwebber/craftctl/config"
	"github.com/benwebber/craftctl/rcon"
	"github.com/codegangsta/cli"
)

//go:generate ./scripts/generate.py -i commands.txt -o init.go --format cli.go

func handleError(err error) int {
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		return 1
	}
	return 0
}

func checkArgs(ctx *cli.Context) error {
	if len(ctx.Args()) < 1 {
		return errors.New("craftctl: enter a command")
	}
	return nil
}

func realMain() (rc int) {
	app := cli.NewApp()
	app.Name = "craftctl"
	app.Usage = "Command-line Minecraft RCON client."
	app.Version = fmt.Sprintf("%v (%v)", Version, GitCommit)
	// Conflicts with console help command.
	app.HideHelp = true
	app.Before = checkArgs

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

		resp, err := client.Auth()
		if rc = handleError(err); rc != 0 {
			return
		}

		command, err := rcon.NewCommand(ctx.Args()...)
		if rc = handleError(err); rc != 0 {
			return
		}

		resp, err = client.Execute(command)
		if rc = handleError(err); rc != 0 {
			return
		}

		fmt.Println(resp)
	}

	err := app.Run(os.Args)
	if rc = handleError(err); rc != 0 {
		return
	}

	return
}

func main() {
	os.Exit(realMain())
}
