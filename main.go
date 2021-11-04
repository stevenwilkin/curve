package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type instrument struct {
	InstrumentName      string `json:"instrument_name"`
	ExpirationTimestamp int64  `json:"expiration_timestamp"`
}

type instrumentsResponse struct {
	Result []instrument `json:"result"`
}

type tickerResponse struct {
	Result struct {
		MarkPrice  float64 `json:"mark_price"`
		IndexPrice float64 `json:"index_price"`
	} `json:"result"`
}

func getJSON(path string, params url.Values, response interface{}) {
	u := url.URL{
		Scheme:   "https",
		Host:     "www.deribit.com",
		Path:     path,
		RawQuery: params.Encode()}

	resp, err := http.Get(u.String())
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(body, &response)
}

func main() {
	var response instrumentsResponse
	getJSON(
		"/api/v2/public/get_instruments",
		url.Values{"currency": {"BTC"}, "kind": {"future"}},
		&response)

	instruments := response.Result
	sort.Slice(instruments, func(i, j int) bool {
		return instruments[i].ExpirationTimestamp < instruments[j].ExpirationTimestamp
	})

	for _, i := range instruments {
		if i.InstrumentName == "BTC-PERPETUAL" {
			continue
		}
		fmt.Println(i.InstrumentName)
		msToExpiration := i.ExpirationTimestamp - time.Now().UnixMilli()

		var response tickerResponse
		getJSON(
			"/api/v2/public/ticker",
			url.Values{"instrument_name": {i.InstrumentName}},
			&response)

		yield := (response.Result.MarkPrice - response.Result.IndexPrice) / response.Result.IndexPrice
		annualisedYield := yield / (float64(msToExpiration) / (1000 * 60 * 60 * 24 * 365))
		fmt.Printf("%.2f%%\n", annualisedYield*100)
	}
}
