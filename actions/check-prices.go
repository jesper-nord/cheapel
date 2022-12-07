package actions

import (
	"errors"
	"fmt"
	"github.com/jesper-nord/cheapel/integration"
	"github.com/samber/lo"
	"github.com/urfave/cli"
	"github.com/xconstruct/go-pushbullet"
	"log"
	"math"
	"time"
)

func CheckPricesAction(c *cli.Context) error {
	hours := c.Int("hours")
	response, err := integration.GetPrices(c.String("tibber-token"))
	if err != nil {
		return err
	}

	// get current home subscription
	home, found := lo.Find(response.Data.Viewer.Homes, func(item integration.Home) bool {
		return item.Subscription != nil
	})
	if !found {
		return errors.New("unable to find home")
	}

	prices := preparePrices(home.Subscription.PriceInfo)
	cheapest := calculateCheapestPeriod(prices, hours)

	msg := fmt.Sprintf("%s (avg price: %.2f kr/kWh)", cheapest.StartTime.Format("2006-01-02 15:04"), cheapest.TotalPrice/float64(hours))
	log.Print(msg)

	if c.Bool("notify") {
		pb := pushbullet.New(c.String("pb-token"))
		device, err := pb.Device(c.String("pb-device"))
		if err != nil {
			return err
		}
		return pb.PushNote(device.Iden, fmt.Sprintf("cheapel %d hour update", hours), msg)
	}

	return nil
}

// merge today's and tomorrow's prices, filter out past prices
func preparePrices(priceInfo integration.PriceInfo) []integration.Price {
	prices := priceInfo.Today
	prices = append(prices, priceInfo.Tomorrow...)
	prices = lo.Filter(prices, func(item integration.Price, _ int) bool {
		startsAt, _ := time.Parse(time.RFC3339, item.StartsAt)
		return startsAt.After(time.Now())
	})
	return prices
}

// calculate cheapest cohesive period of given hours
func calculateCheapestPeriod(prices []integration.Price, hours int) Period {
	cheapest := Period{
		TotalPrice: math.MaxFloat64,
	}
	for i := range prices {
		sum := 0.0
		for j := 0; j < hours; j++ {
			sum += prices[i+j].Total
		}
		if sum < cheapest.TotalPrice {
			cheapest.TotalPrice = sum
			startsAt, _ := time.Parse(time.RFC3339, prices[i].StartsAt)
			cheapest.StartTime = startsAt
		}
		if i+hours >= len(prices) {
			break
		}
	}
	return cheapest
}

type Period struct {
	TotalPrice float64
	StartTime  time.Time
}
