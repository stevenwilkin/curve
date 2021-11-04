package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"
)

var (
	yields []result
	m      sync.Mutex
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

type result struct {
	Expiration int64   `json:"expiration"`
	Yield      float64 `json:"yield"`
}

func getJSON(path string, params url.Values, response interface{}) error {
	u := url.URL{
		Scheme:   "https",
		Host:     "www.deribit.com",
		Path:     path,
		RawQuery: params.Encode()}

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	json.Unmarshal(body, &response)
	return nil
}

func getYields() ([]result, error) {
	var response instrumentsResponse
	err := getJSON(
		"/api/v2/public/get_instruments",
		url.Values{"currency": {"BTC"}, "kind": {"future"}},
		&response)
	if err != nil {
		return nil, err
	}

	instruments := response.Result
	sort.Slice(instruments, func(i, j int) bool {
		return instruments[i].ExpirationTimestamp < instruments[j].ExpirationTimestamp
	})

	var results []result

	for _, i := range instruments {
		if i.InstrumentName == "BTC-PERPETUAL" {
			continue
		}
		msToExpiration := i.ExpirationTimestamp - time.Now().UnixMilli()

		var response tickerResponse
		err = getJSON(
			"/api/v2/public/ticker",
			url.Values{"instrument_name": {i.InstrumentName}},
			&response)
		if err != nil {
			return nil, err
		}

		yield := (response.Result.MarkPrice - response.Result.IndexPrice) / response.Result.IndexPrice
		annualisedYield := yield / (float64(msToExpiration) / (1000 * 60 * 60 * 24 * 365))

		results = append(results, result{
			Expiration: i.ExpirationTimestamp,
			Yield:      annualisedYield})
	}

	return results, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	m.Lock()
	jsonYields, err := json.Marshal(yields)
	m.Unlock()

	if err != nil {
		log.Println(err.Error())
		return
	}

	w.Write(jsonYields)
}

func main() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)

		for {
			if results, err := getYields(); err != nil {
				log.Println(err.Error())
			} else {
				m.Lock()
				yields = results
				m.Unlock()
			}

			<-ticker.C
		}
	}()

	log.Println("Starting on 0.0.0.0:8080")
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
