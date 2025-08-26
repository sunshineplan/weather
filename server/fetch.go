package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/api/airmatters"
	"github.com/sunshineplan/weather/aqi"
	"github.com/sunshineplan/weather/unit/coordinates"
)

var coordsMap sync.Map

func getCoords(query string, api coordinates.GeoLocator) (res coordinates.Coordinates, err error) {
	if v, ok := coordsMap.Load(query); ok {
		res = v.(coordinates.Coordinates)
		return
	}
	if api != nil {
		res, err = api.Coordinates(query)
	} else {
		res, err = forecast.Coordinates(query)
		if err != nil {
			res, err = realtime.Coordinates(query)
		}
	}
	if err != nil {
		return
	}
	coordsMap.Store(query, res)
	return
}

func getAQIStandard() (standard int, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()
	std, err := aqiAPI.Standard(aqiType)
	if err != nil {
		return
	}
	index := std[1]
	for i := 2; i < len(std); i++ {
		if index.Level().String() != std[i].Level().String() {
			index = std[i]
			break
		}
		index = std[i]
	}
	standard = index.Value()
	return
}

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

func getWeatherByCoordinates(coords coordinates.Coordinates, n int, t time.Time, realtime bool,
) (current weather.Current, days []weather.Day, err error) {
	c := make(chan error, 3)
	go func() {
		var err error
		if realtime {
			current, err = forecast.RealtimeByCoordinates(coords)
		}
		c <- err
	}()
	var forecasts []weather.Day
	go func() {
		var err error
		forecasts, err = forecast.ForecastByCoordinates(coords, n)
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
		yesterday, err = history.HistoryByCoordinates(coords, t.AddDate(0, 0, -1))
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
	coords, err := getCoords(q, forecast)
	if err != nil {
		return nil, err
	}
	return getAQIByCoordinates(aqiType, coords)
}

func getAQIByCoordinates(aqiType aqi.Type, coords coordinates.Coordinates) (aqi.Current, error) {
	_, res, err := aqiAPI.(*airmatters.AirMatters).RealtimeNearby(aqiType, coords)
	return res, err
}

func getAll(q string, n int, aqiType aqi.Type, t time.Time, realtime bool,
) (current weather.Current, days []weather.Day, lastYear, avg weather.Day, currentAQI aqi.Current, err error) {
	c := make(chan error, 3)
	go func() {
		var err error
		current, days, err = getWeather(q, n, t, realtime)
		c <- err
	}()
	go func() {
		var err error
		if q == *query {
			lastYear, avg, err = historyRecord(t, 2)
		}
		c <- err
	}()
	go func() {
		var err error
		if realtime {
			currentAQI, err = getAQI(aqiType, q)
		}
		c <- err
	}()
	for range 3 {
		if err = <-c; err != nil {
			return
		}
	}
	return
}

func getAllByCoordinates(coords coordinates.Coordinates, n int, aqiType aqi.Type, t time.Time, realtime bool,
) (current weather.Current, days []weather.Day, lastYear, avg weather.Day, currentAQI aqi.Current, err error) {
	c := make(chan error, 2)
	go func() {
		var err error
		current, days, err = getWeatherByCoordinates(coords, n, t, realtime)
		c <- err
	}()
	go func() {
		var err error
		if realtime {
			currentAQI, err = getAQIByCoordinates(aqiType, coords)
		}
		c <- err
	}()
	for i := 0; i < 2; i++ {
		if err = <-c; err != nil {
			return
		}
	}
	return
}
