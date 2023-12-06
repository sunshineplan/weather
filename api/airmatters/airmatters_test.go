package airmatters

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/sunshineplan/weather/aqi"
)

func TestAirMatters(t *testing.T) {
	key := os.Getenv("AIR_MATTERS")
	if key == "" {
		log.Print("Skip air-matters")
		return
	}
	var api aqi.API = New(key)
	if _, err := api.Coordinates("shanghai"); err != nil {
		t.Error(err)
	}
	if _, err := api.Realtime(aqi.China, "shanghai"); err != nil {
		t.Error(err)
	}
	if _, err := api.Forecast(aqi.China, "shanghai", 0); err != nil {
		t.Error(err)
	}
	if _, err := api.History(aqi.China, "shanghai", time.Now().AddDate(0, 0, -10)); err != nil {
		t.Error(err)
	}
}
