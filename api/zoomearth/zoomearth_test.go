package zoomearth

import (
	"image/png"
	"os"
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
	_, img, err := api.Realtime(weather.Satellite, coordinates.New(0, 0), nil)
	if err != nil {
		t.Fatal(err)
	}
	f, err := os.Create("test.png")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		t.Fatal(err)
	}
}
