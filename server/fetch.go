package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/api/airmatters"
	"github.com/sunshineplan/weather/aqi"
)

func getWeather(query string, n int, t time.Time) ([]weather.Day, error) {
	forecasts, err := forecast.Forecast(query, n)
	if err != nil {
		return nil, err
	}
	if len(forecasts) < n {
		return nil, fmt.Errorf("bad forecast number: %d", len(forecasts))
	}
	if date := t.Format("01-02"); !strings.HasSuffix(forecasts[0].Date, date) {
		return nil, fmt.Errorf("the first forecast date(%s) does not match the input(%s)", forecasts[0].Date, date)
	}
	yesterday, err := history.History(query, t.AddDate(0, 0, -1))
	if err != nil {
		return nil, err
	}
	return append([]weather.Day{yesterday}, forecasts...), nil
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

func getAll(q string, n int, aqiType aqi.Type, t time.Time) (days []weather.Day, avg weather.Day, current aqi.Current, err error) {
	if days, err = getWeather(q, n, t); err != nil {
		return
	}
	if q == *query {
		if avg, err = average(t.Format("01-02"), 2); err != nil {
			return
		}
	}
	current, err = getAQI(aqiType, q)
	return
}
