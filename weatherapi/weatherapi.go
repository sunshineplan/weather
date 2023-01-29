package weatherapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sunshineplan/weather"
)

const baseURL = "https://api.weatherapi.com/v1"

var _ weather.API = &WeatherAPI{}

type WeatherAPI struct {
	key string
}

func New(key string) *WeatherAPI {
	return &WeatherAPI{key}
}

func (api *WeatherAPI) Request(endpoint, query string) (res Response, err error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s?key=%s&%s", baseURL, endpoint, api.key, query))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("status code: %d", resp.StatusCode)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&res)
	return
}

func (api *WeatherAPI) Realtime(query string) (current weather.Current, err error) {
	resp, err := api.Request("current.json", fmt.Sprintf("q=%s", query))
	if err != nil {
		return
	}
	current = resp.Current.Convert()
	return
}

func (api *WeatherAPI) Forecast(query string, n int) (current weather.Current, days []weather.Day, err error) {
	resp, err := api.Request("forecast.json", fmt.Sprintf("q=%s&days=%d", query, n))
	if err != nil {
		return
	}
	current = resp.Current.Convert()
	if days = resp.Forecast.Convert(); len(days) < n {
		err = fmt.Errorf("bad forecast number: %d", len(days))
	}
	return
}

func (api *WeatherAPI) History(query string, date time.Time) (day weather.Day, err error) {
	resp, err := api.Request("history.json", fmt.Sprintf("q=%s&dt=%s", query, date.Format("2006-01-02")))
	if err != nil {
		return
	}
	if days := resp.Forecast.Convert(); len(days) == 0 {
		err = errors.New("no history result")
	} else {
		day = days[0]
	}
	return
}
