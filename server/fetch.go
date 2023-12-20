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
	if realtime {
		current, err = forecast.Realtime(query)
		if err != nil {
			return
		}
	}
	forecasts, err := forecast.Forecast(query, n)
	if err != nil {
		return
	}
	if len(forecasts) < n {
		err = fmt.Errorf("bad forecast number: %d", len(forecasts))
		return
	}
	if date := t.Format("01-02"); !strings.HasSuffix(forecasts[0].Date, date) {
		err = fmt.Errorf("the first forecast date(%s) does not match the input(%s)", forecasts[0].Date, date)
		return
	}
	yesterday, err := history.History(query, t.AddDate(0, 0, -1))
	if err != nil {
		return
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

func getAll(q string, n int, aqiType aqi.Type, t time.Time, realtime bool) (
	current weather.Current, days []weather.Day, avg weather.Day, currentAQI aqi.Current, err error) {
	if current, days, err = getWeather(q, n, t, realtime); err != nil {
		return
	}
	if q == *query {
		if avg, err = average(t.Format("01-02"), 2); err != nil {
			return
		}
	}
	currentAQI, err = getAQI(aqiType, q)
	return
}
