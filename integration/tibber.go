package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const endpointUrl = "https://api.tibber.com/v1-beta/gql"
const pricesQuery = `{"query": "{viewer {homes {currentSubscription {priceInfo {today {total startsAt} tomorrow {total startsAt}}}}}}" }`
const notificationMutation = `{"query": "mutation{ sendPushNotification(input: { title: \"%s\", message: \"%s\", screenToOpen: HOME }){ successful }}"}`

var (
	client = &http.Client{
		Timeout: 5 * time.Second,
	}
)

type Response struct {
	Data ResponseData `json:"data"`
}

type ResponseData struct {
	Viewer Viewer `json:"viewer"`
}

type Viewer struct {
	Homes []Home `json:"homes"`
}

type Home struct {
	Subscription *Subscription `json:"currentSubscription"`
}

type Subscription struct {
	PriceInfo PriceInfo `json:"priceInfo"`
}

type PriceInfo struct {
	Today    []Price `json:"today"`
	Tomorrow []Price `json:"tomorrow"`
}

type Price struct {
	Total    float64 `json:"total"`
	StartsAt string  `json:"startsAt"`
}

func GetPrices(accessToken string) (Response, error) {
	response, err := doPost(accessToken, pricesQuery)
	if err != nil {
		return Response{}, err
	}
	defer response.Body.Close()

	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		return Response{}, err
	}

	var responseData Response
	err = json.Unmarshal(jsonData, &responseData)
	if err != nil {
		return Response{}, err
	}

	return responseData, nil
}

func SendNotification(accessToken, title, msg string) (Response, error) {
	response, err := doPost(accessToken, fmt.Sprintf(notificationMutation, title, msg))
	if err != nil {
		return Response{}, err
	}
	defer response.Body.Close()

	jsonData, err := io.ReadAll(response.Body)
	if err != nil {
		return Response{}, err
	}

	var responseData Response
	err = json.Unmarshal(jsonData, &responseData)
	if err != nil {
		return Response{}, err
	}

	return responseData, nil
}

func doPost(accessToken, query string) (*http.Response, error) {
	request, err := http.NewRequest("POST", endpointUrl, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return nil, err
	}
	request.Header = map[string][]string{
		"User-Agent":    {"cheapel by jesper@nord.pm"},
		"Authorization": {fmt.Sprintf("Bearer %s", accessToken)},
		"Content-Type":  {"application/json"},
	}
	response, err := client.Do(request)

	if err != nil {
		return response, err
	}
	if response.StatusCode != 200 {
		return response, errors.New("response status code not 200")
	}
	return response, nil
}
