package visualcrossing

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sunshineplan/weather"
)

const baseURL = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline"

var _ weather.API = &VisualCrossing{}

type VisualCrossing struct {
	key string
}

func New(key string) *VisualCrossing {
	return &VisualCrossing{key}
}

func (api *VisualCrossing) Request(endpoint, include, query string) (res Response, err error) {
	url := fmt.Sprintf("%s/%s/%s?unitGroup=metric&key=%s", baseURL, query, endpoint, api.key)
	if include != "" {
		url += "&include=" + include
	}
	resp, err := http.Get(url)
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

func (api *VisualCrossing) Realtime(query string) (current weather.Current, err error) {
	resp, err := api.Request("today", "current", query)
	if err != nil {
		return
	}
	current = resp.CurrentConditions.Convert()
	return
}

func (api *VisualCrossing) Forecast(query string, n int) (current weather.Current, days []weather.Day, err error) {
	var endpoint string
	switch n {
	case 1:
		endpoint = "next24hours"
	case 7:
		endpoint = "next7days"
	case 15:
	case 30:
		endpoint = "next30days"
	default:
		now := time.Now()
		endpoint = fmt.Sprintf("%s/%s", now.Format("2006-01-02"), now.AddDate(0, 0, n).Format("2006-01-02"))
	}
	resp, err := api.Request(endpoint, "hours%2Ccurrent", query)
	if err != nil {
		return
	}
	current = resp.CurrentConditions.Convert()
	if days = ConvertDays(resp.Days); len(days) < n {
		err = fmt.Errorf("bad forecast number: %d", len(days))
	}
	return
}

func (api *VisualCrossing) History(query string, date time.Time) (day weather.Day, err error) {
	resp, err := api.Request(date.Format("2006-01-02"), "days", query)
	if err != nil {
		return
	}
	if days := ConvertDays(resp.Days); len(days) == 0 {
		err = errors.New("no history result")
	} else {
		day = days[0]
	}
	return
}
