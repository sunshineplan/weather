package zoomearth

import (
	"testing"

	"github.com/sunshineplan/weather/unit/coordinates"
)

func TestMap(t *testing.T) {
	coords := coordinates.New(0, 0)
	for i, testcase := range []struct {
		path     string
		opt      *MapOptions
		expected string
	}{
		{"", NewMapOptions(), "https://zoom.earth/maps/satellite/#view=0,0,4z"},
		{"satellite", NewMapOptions().SetZoom(5), "https://zoom.earth/maps/satellite/#view=0,0,5z"},
		{"radar", NewMapOptions().SetZoom(6).SetOverlays([]string{"wind"}), "https://zoom.earth/maps/radar/#view=0,0,6z/overlays=wind"},
		{"wind", NewMapOptions().SetZoom(7).SetOverlays([]string{"radar", "wind"}), "https://zoom.earth/maps/wind/#view=0,0,7z/overlays=radar,wind"},
	} {
		if res := URL(testcase.path, coords, testcase.opt.Zoom(), testcase.opt.Overlays()); res != testcase.expected {
			t.Errorf("%d expected %q; got %q", i, testcase.expected, res)
		}
	}
}
