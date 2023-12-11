package airmatters

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/sunshineplan/weather/aqi"
	"github.com/sunshineplan/weather/unit/coordinates"
)

func TestAirMatters(t *testing.T) {
	key := os.Getenv("AIR_MATTERS")
	if key == "" {
		log.Print("Skip air-matters")
		return
	}
	api := New(key)
	if _, err := api.Coordinates("shanghai"); err != nil {
		t.Error(err)
	}
	if _, err := api.Realtime(aqi.China, "shanghai"); err != nil {
		t.Error(err)
	}
	if _, _, err := api.RealtimeNearby(aqi.China, coordinates.New(31.17, 121.47)); err != nil {
		t.Error(err)
	}
	if _, err := api.Forecast(aqi.China, "shanghai", 0); err != nil {
		t.Error(err)
	}
	if _, err := api.History(aqi.China, "shanghai", time.Now().AddDate(0, 0, -10)); err != nil {
		t.Error(err)
	}
}
