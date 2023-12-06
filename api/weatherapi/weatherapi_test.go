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
	if _, err := api.Coordinates("shanghai"); err != nil {
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
}
