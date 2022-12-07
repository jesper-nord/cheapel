package actions

import (
	"github.com/jesper-nord/cheapel/integration"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCalculateCheapestPeriod(t *testing.T) {
	prices := []integration.Price{
		{
			Total:    2.7,
			StartsAt: "2022-12-06T00:00:00.000+01:00",
		},
		{
			Total:    2.9,
			StartsAt: "2022-12-06T01:00:00.000+01:00",
		},
		{
			Total:    3.1,
			StartsAt: "2022-12-06T02:00:00.000+01:00",
		},
		{
			Total:    3.2,
			StartsAt: "2022-12-06T03:00:00.000+01:00",
		},
		{
			Total:    1.4,
			StartsAt: "2022-12-06T04:00:00.000+01:00",
		},
		{
			Total:    1.6,
			StartsAt: "2022-12-06T05:00:00.000+01:00",
		},
		{
			Total:    9.2,
			StartsAt: "2022-12-06T06:00:00.000+01:00",
		},
		{
			Total:    4.7,
			StartsAt: "2022-12-06T07:00:00.000+01:00",
		},
	}

	expectedStartTime, _ := time.Parse(time.RFC3339, "2022-12-06T02:00:00.000+01:00")
	expected := Period{
		TotalPrice: 9.3,
		StartTime:  expectedStartTime,
	}

	result := calculateCheapestPeriod(prices, 4)

	assert.Equal(t, expected, result)
}
