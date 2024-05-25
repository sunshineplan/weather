package zoomearth

import (
	"testing"
	"time"

	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/unit/coordinates"
)

func TestZoomEarth(t *testing.T) {
	api := ZoomEarthAPI{}
	if _, err := api.GetStorms(time.Now()); err != nil {
		t.Error(err)
	}
	if _, _, err := api.Realtime(weather.Satellite, coordinates.New(0, 0), nil); err != nil {
		t.Error(err)
	}
}
