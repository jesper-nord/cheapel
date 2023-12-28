package main

import (
	"flag"
	"github.com/jesper-nord/cheapel/v2/service"
	"log"
	"os"
)

var (
	tokenFlag  = flag.String("token", "", "Tibber API token")
	hoursFlag  = flag.Int("hours", 1, "Length of period to check for (hours)")
	notifyFlag = flag.Bool("notify", false, "Send a Tibber notification with result")
	helpFlag   = flag.Bool("help", false, "Show available commands")
)

func main() {
	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := service.CheckPrices(*hoursFlag, *tokenFlag, *notifyFlag)
	if err != nil {
		log.Fatal(err)
	}
}
