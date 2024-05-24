package zoomearth

import (
	"testing"

	"github.com/sunshineplan/weather/unit/coordinates"
)

func TestMap(t *testing.T) {
	coords := coordinates.New(0, 0)
	for i, testcase := range []struct {
		path     string
		zoom     float64
		overlays Overlays
		expected string
	}{
		{"satellite", 4, nil, "https://zoom.earth/maps/satellite/#view=0,0,4z"},
		{"radar", 5, []string{"wind"}, "https://zoom.earth/maps/radar/#view=0,0,5z/overlays=wind"},
		{"wind", 6, []string{"radar", "wind"}, "https://zoom.earth/maps/wind/#view=0,0,6z/overlays=radar,wind"},
	} {
		if res := url(testcase.path, coords, testcase.zoom, testcase.overlays); res != testcase.expected {
			t.Errorf("%d expected %q; got %q", i, testcase.expected, res)
		}
	}
	if _, err := Realtime("", nil, coords, 4, 95); err == nil {
		t.Error("expected error; got nil")
	}
}
