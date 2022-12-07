package main

import (
	"github.com/jesper-nord/cheapel/actions"
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
					Usage:    "send a notification with result via pushbullet",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "pb-token",
					Usage:    "pushbullet access token, required to notify",
					Required: false,
				},
				&cli.StringFlag{
					Name:     "pb-device",
					Usage:    "pushbullet device name to notify, required to notify",
					Required: false,
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
