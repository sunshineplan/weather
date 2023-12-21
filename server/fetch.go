package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/api/airmatters"
	"github.com/sunshineplan/weather/aqi"
)

func getWeather(query string, n int, t time.Time, realtime bool) (current weather.Current, days []weather.Day, err error) {
	c := make(chan error, 3)
	go func() {
		var err error
		if realtime {
			current, err = forecast.Realtime(query)
		}
		c <- err
	}()
	var forecasts []weather.Day
	go func() {
		var err error
		forecasts, err = forecast.Forecast(query, n)
		if err != nil {
			c <- err
		} else if len(forecasts) < n {
			c <- fmt.Errorf("bad forecast number: %d", len(forecasts))
		} else if date := t.Format("01-02"); !strings.HasSuffix(forecasts[0].Date, date) {
			c <- fmt.Errorf("the first forecast date(%s) does not match the input(%s)", forecasts[0].Date, date)
		} else {
			c <- nil
		}
	}()
	var yesterday weather.Day
	go func() {
		var err error
		yesterday, err = history.History(query, t.AddDate(0, 0, -1))
		c <- err
	}()
	for i := 0; i < 3; i++ {
		if err = <-c; err != nil {
			return
		}
	}
	days = append([]weather.Day{yesterday}, forecasts...)
	return
}

func getAQI(aqiType aqi.Type, q string) (aqi.Current, error) {
	if res, err := aqiAPI.Realtime(aqiType, q); err == nil {
		return res, nil
	}
	coords, err := getCoords(q)
	if err != nil {
		return nil, err
	}
	_, res, err := aqiAPI.(*airmatters.AirMatters).RealtimeNearby(aqiType, coords)
	return res, err
}

func getAll(q string, n int, aqiType aqi.Type, t time.Time, realtime bool,
) (current weather.Current, days []weather.Day, avg weather.Day, currentAQI aqi.Current, err error) {
	c := make(chan error, 3)
	go func() {
		var err error
		current, days, err = getWeather(q, n, t, realtime)
		c <- err
	}()
	go func() {
		var err error
		if q == *query {
			avg, err = average(t.Format("01-02"), 2)
		}
		c <- err
	}()
	go func() {
		var err error
		currentAQI, err = getAQI(aqiType, q)
		c <- err
	}()
	for i := 0; i < 3; i++ {
		if err = <-c; err != nil {
			return
		}
	}
	return
}
