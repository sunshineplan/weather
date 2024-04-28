package weatherapi

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/sunshineplan/weather"
)

func TestWeatherAPI(t *testing.T) {
	key := os.Getenv("WEATHERAPI")
	if key == "" {
		log.Print("Skip weatherapi")
		return
	}
	var api weather.API = New(key)
	coords, err := api.Coordinates("shanghai")
	if err != nil {
		t.Error(err)
	}
	if _, err := api.Realtime("shanghai"); err != nil {
		t.Error(err)
	}
	if _, err := api.Forecast("shanghai", 1); err != nil {
		t.Error(err)
	}
	date := time.Now().AddDate(0, 0, -3)
	if _, err := api.History("shanghai", date); err != nil {
		t.Error(err)
	}
	if _, err := api.RealtimeByCoordinates(coords); err != nil {
		t.Error(err)
	}
	if _, err := api.ForecastByCoordinates(coords, 1); err != nil {
		t.Error(err)
	}
	if _, err := api.HistoryByCoordinates(coords, date); err != nil {
		t.Error(err)
	}
}
