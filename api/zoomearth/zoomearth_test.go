package zoomearth

import (
	"testing"
	"time"

	"github.com/sunshineplan/weather/storm"
)

func TestZoomEarth(t *testing.T) {
	var api storm.API = ZoomEarthAPI{}
	if _, err := api.GetStorms(time.Now()); err != nil {
		t.Fatal(err)
	}
}
