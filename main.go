package main

import (
	"flag"
	"fmt"
	"github.com/jesper-nord/cheapel/service"
	"github.com/spf13/viper"
	"log"
	"os"
)

var (
	hoursFlag  = flag.Int("hours", 1, "length of period to check for (hours)")
	notifyFlag = flag.Bool("notify", false, "send a notification with result via pushbullet")
	configFlag = flag.String("config", "./cheapel-config.yaml", "config file location")
	helpFlag   = flag.Bool("help", false, "show available commands")
)

func main() {
	flag.Parse()

	if *helpFlag {
		flag.PrintDefaults()
		os.Exit(1)
	}

	readConfig()
	tibberToken := viper.GetString("tibber-token")
	pbToken := viper.GetString("pb-token")
	pbDevice := viper.GetString("pb-device")

	err := service.CheckPrices(*hoursFlag, tibberToken, pbToken, pbDevice, *notifyFlag)
	if err != nil {
		log.Fatal(err)
	}
}

func readConfig() {
	viper.SetConfigFile(*configFlag)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}
