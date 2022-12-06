package main

import (
	"cheapel/actions"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "cheapel"
	app.Usage = "find the cheapest cohesive period of given hours with regards to electricity costs"
	app.Commands = []cli.Command{
		{
			Name:     "find",
			HelpName: "find",
			Action:   actions.CheckPricesAction,
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:     "hours",
					Usage:    "period of hours to check for",
					Required: true,
					Value:    4,
				},
				&cli.StringFlag{
					Name:     "tibber-token",
					Usage:    "tibber access token",
					Required: true,
				},
				&cli.BoolFlag{
					Name:     "notify",
					Required: false,
					Usage:    "send a notification with result via pushbullet",
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
