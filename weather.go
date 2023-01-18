package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var ApiKey string

const baseURL = "https://api.weatherapi.com/v1"

func weather(api, query string) (res Response, err error) {
	resp, err := http.Get(baseURL + fmt.Sprintf("/%s?key=%s&%s", api, ApiKey, query))
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

func RealtimeWeather(q string) (Response, error) {
	return weather("current.json", fmt.Sprintf("q=%s", q))
}

func ForecastWeather(q string, days int) (Response, error) {
	return weather("forecast.json", fmt.Sprintf("q=%s&days=%d", q, days))
}

func HistoryWeather(q string, day time.Time) (Response, error) {
	return weather("history.json", fmt.Sprintf("q=%s&dt=%s", q, day.Format("2006-01-02")))
}
