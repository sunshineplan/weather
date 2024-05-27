package zoomearth

import (
	"testing"
	"time"

	"github.com/sunshineplan/weather/unit/coordinates"
)

func TestMap(t *testing.T) {
	coords := coordinates.New(0, 0)
	for i, testcase := range []struct {
		path     string
		t        time.Time
		opt      *MapOptions
		expected string
	}{
		{
			"",
			time.Time{},
			NewMapOptions(),
			"https://zoom.earth/maps/satellite/#view=0,0,4z",
		},
		{
			"satellite",
			time.Date(2006, 1, 2, 3, 4, 0, 0, time.UTC),
			NewMapOptions().SetZoom(5),
			"https://zoom.earth/maps/satellite/#view=0,0,5z/date=2006-01-02,03:04",
		},
		{
			"radar",
			time.Time{},
			NewMapOptions().SetZoom(6).SetOverlays([]string{"wind"}),
			"https://zoom.earth/maps/radar/#view=0,0,6z/overlays=wind",
		},
		{
			"wind",
			time.Time{},
			NewMapOptions().SetZoom(7).SetOverlays([]string{"radar", "wind"}),
			"https://zoom.earth/maps/wind/#view=0,0,7z/overlays=radar,wind",
		},
	} {
		if res := URL(testcase.path, testcase.t, coords, testcase.opt.Zoom(), testcase.opt.Overlays()); res != testcase.expected {
			t.Errorf("%d expected %q; got %q", i, testcase.expected, res)
		}
	}
}
