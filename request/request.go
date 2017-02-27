package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	url := os.Getenv("PING_URL")
	if url == "" {
		url = "http://localhost:3000/api/events/count"
	}

	ticker := time.NewTicker(time.Second * 2)
	go func() {
		for t := range ticker.C {
			makeRequest(url)
			fmt.Println(t)
		}
	}()

	time.Sleep(time.Hour * 999999) // run for long time
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

func makeRequest(url string) {
	start := time.Now()

	res, err := http.Get(url)

	defer func() {
		// res.Close = true
		res.Body.Close()
	}()

	if err != nil {
		fmt.Println("err", err)
		return
	}

	// fmt.Println(res)

	dur := time.Since(start)

	data := map[string]interface{}{}

	json.NewDecoder(res.Body).Decode(&data)

	// fmt.Println(data)

	count, ok := data["count"].(float64)
	if !ok {
		fmt.Printf("unexpected count type: %T\n", data["count"])
	}
	saveResults(dur, count)

}

func saveResults(dur time.Duration, count float64) {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080/api/watchman"
	}

	results := map[string]interface{}{
		"event_count": count,
		"resp_time":   dur.Nanoseconds() / 1000000, //ms
		"created":     time.Now().UnixNano() / 1000000,
	}

	data, err := json.Marshal(results)

	res, err := http.Post(apiURL, "application/json", bytes.NewBuffer(data))

	defer func() {
		// res.Close = true
		res.Body.Close()
	}()

	if err != nil {
		fmt.Println("err", err)
		return
	}

	// fmt.Println(res)
}
