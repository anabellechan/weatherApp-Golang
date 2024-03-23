package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const apiKey = "6d18b7538ec43712f53da5049d65351b"

func fetchWeather(city string, ch chan<- string, wg *sync.WaitGroup) interface{} {
	var data struct {
		Main struct {
			Temp float64 `json:"temp"`
		} `json:"main"`
	}

	defer wg.Done()

	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching weather for %s: %s\n", city, err)
		return data
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Printf("Error decoding weather data for %s: %s\n", city, err)
		return data
	}
	ch <- fmt.Sprintf("The temperature in %s is %v Kelvin", city, data.Main.Temp)

	return data
}

func main() {
	startNow := time.Now()
	cities := []string{"Toronto", "London", "Paris", "Tokyo"}
	ch := make(chan string)
	var wg sync.WaitGroup

	for _, city := range cities {
		wg.Add(1)
		go fetchWeather(city, ch, &wg)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		fmt.Println(result)
	}

	fmt.Println("This operation task took", time.Since(startNow))
}
