package weatherapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/unit/coordinates"
)

const baseURL = "https://api.weatherapi.com/v1"

var _ weather.API = &WeatherAPI{}

type WeatherAPI struct {
	key string
}

func New(key string) *WeatherAPI {
	return &WeatherAPI{key}
}

func (api *WeatherAPI) Request(endpoint string, query url.Values) (res Response, err error) {
	query.Set("key", api.key)
	resp, err := http.Get(fmt.Sprintf("%s/%s?%s", baseURL, endpoint, query.Encode()))
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

func (api *WeatherAPI) Coordinates(query string) (coordinates.Coordinates, error) {
	resp, err := api.Request("current.json", url.Values{"q": {query}})
	if err != nil {
		return nil, err
	}
	if location := resp.Location; location != nil {
		return coordinates.New(location.Lat, location.Lon), nil
	}
	return nil, errors.New("location is nil")
}

func (api *WeatherAPI) Realtime(query string) (current weather.Current, err error) {
	resp, err := api.Request("current.json", url.Values{"q": {query}})
	if err != nil {
		return
	}
	current = resp.Current.Convert()
	return
}

func (api *WeatherAPI) Forecast(query string, n int) (days []weather.Day, err error) {
	resp, err := api.Request("forecast.json", url.Values{"q": {query}, "days": {strconv.Itoa(n)}})
	if err != nil {
		return
	}
	if days = resp.Forecast.Convert(); len(days) < n {
		err = fmt.Errorf("bad forecast number: %d", len(days))
	}
	return
}

func (api *WeatherAPI) History(query string, date time.Time) (day weather.Day, err error) {
	resp, err := api.Request("history.json", url.Values{"q": {query}, "dt": {date.Format("2006-01-02")}})
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
