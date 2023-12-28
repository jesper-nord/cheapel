package service

import (
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/jesper-nord/cheapel/v2/integration"
	"github.com/samber/lo"
	"log"
	"math"
	"time"
)

func CheckPrices(hours int, token string, notify bool) error {
	checkPricesRetryable := func() error {
		response, err := integration.GetPrices(token)
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

		msg := fmt.Sprintf("%s-%s (avg price %.2f kr/kWh)", cheapest.StartTime.Format("15:04"), cheapest.EndTime.Format("15:04"), cheapest.TotalPrice/float64(hours))
		log.Print(msg)

		if notify {
			_, err = integration.SendNotification(token, "cheapel", msg)
			if err != nil {
				return err
			}
		}

		return nil
	}

	return retry.Do(checkPricesRetryable, retry.Attempts(5), retry.Delay(time.Second*10), retry.DelayType(retry.FixedDelay))
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
			cheapest.EndTime = startsAt.Add(time.Duration(hours) * time.Hour)
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
	EndTime    time.Time
}
