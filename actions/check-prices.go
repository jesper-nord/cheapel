package actions

import (
	"cheapel/integration"
	"errors"
	"github.com/samber/lo"
	"github.com/urfave/cli"
	"log"
	"math"
	"time"
)

func CheckPricesAction(c *cli.Context) error {
	// notifiy := c.IsSet("notify")

	accessToken := c.String("tibber-token")
	hours := c.Int("hours")
	// fetch prices
	response, err := integration.GetPrices(accessToken)
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

	// merge today's and tomorrow's prices, filter out past prices
	prices := home.Subscription.PriceInfo.Today
	prices = append(prices, home.Subscription.PriceInfo.Tomorrow...)
	prices = lo.Filter(prices, func(item integration.Price, _ int) bool {
		startsAt, _ := time.Parse(time.RFC3339, item.StartsAt)
		return startsAt.After(time.Now())
	})

	// calculate cheapest cohesive period of given hours
	cheapest := math.MaxFloat64
	startsAt := ""
	for i := range prices {
		sum := 0.0
		for j := 0; j < hours; j++ {
			sum += prices[i+j].Total
		}
		if sum < cheapest {
			cheapest = sum
			startsAt = prices[i].StartsAt
		}
		if i+hours >= len(prices) {
			break
		}
	}
	log.Printf("cheapest starts at: %s (avg price: %.2f kr/kWh)", startsAt, cheapest/float64(hours))
	return nil
}
