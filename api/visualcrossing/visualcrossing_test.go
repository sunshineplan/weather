package visualcrossing

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/sunshineplan/weather"
)

func TestVisualCrossing(t *testing.T) {
	key := os.Getenv("VISUALCROSSING")
	if key == "" {
		log.Print("Skip visualcrossing")
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
	if _, err := api.History("shanghai", time.Now()); err != nil {
		t.Error(err)
	}
	if _, err := api.RealtimeByCoordinates(coords); err != nil {
		t.Error(err)
	}
	if _, err := api.ForecastByCoordinates(coords, 1); err != nil {
		t.Error(err)
	}
	if _, err := api.HistoryByCoordinates(coords, time.Now()); err != nil {
		t.Error(err)
	}
}
