package zoomearth

import (
	"errors"
	"time"

	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/storm"
	"github.com/sunshineplan/weather/unit/coordinates"
)

var (
	_ storm.API      = ZoomEarthAPI{}
	_ weather.MapAPI = ZoomEarthAPI{}
)

func (ZoomEarthAPI) GetStorms(t time.Time) ([]storm.Storm, error) {
	return GetStorms(t)
}

func (ZoomEarthAPI) Realtime(
	t weather.MapType, coords coordinates.Coordinates, zoom float64, quality int, opt weather.MapOption,
) ([]byte, error) {
	path, ok := mapPath[t]
	if !ok {
		return nil, errors.New("unsupported map type")
	}
	var overlays Overlays
	if opt != nil && opt.Compatibility(ZoomEarthAPI{}) {
		overlays = opt.Value().(Overlays)
	}
	return Realtime(path, overlays, coords, zoom, quality)
}
