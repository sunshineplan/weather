package zoomearth

import (
	"image"
	"time"

	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/option"
	"github.com/sunshineplan/weather/storm"
	"github.com/sunshineplan/weather/unit/coordinates"
)

const root = "https://zoom.earth"

var (
	_ storm.API      = ZoomEarthAPI{}
	_ weather.MapAPI = ZoomEarthAPI{}
)

func (ZoomEarthAPI) GetStorms(t time.Time) ([]storm.Storm, error) {
	return GetStorms(t)
}

func (ZoomEarthAPI) URL(t weather.MapType, coords coordinates.Coordinates, opt any) string {
	zoom := defaultMapOptions.zoom
	if opt, ok := opt.(option.Zoom); ok {
		zoom = opt.Zoom()
	}
	overlays := defaultMapOptions.overlays
	if opt, ok := opt.(option.Overlays); ok {
		overlays = opt.Overlays()
	}
	return URL(mapPath[t], coords, zoom, overlays)
}

func (ZoomEarthAPI) Realtime(t weather.MapType, coords coordinates.Coordinates, opt any) (time.Time, image.Image, error) {
	if opt == nil {
		return Realtime(mapPath[t], coords, nil)
	}
	o := defaultMapOptions
	if opt, ok := opt.(option.Size); ok {
		o.width, o.height = opt.Size()
	}
	if opt, ok := opt.(option.Zoom); ok {
		o.zoom = opt.Zoom()
	}
	if opt, ok := opt.(option.Overlays); ok {
		o.overlays = opt.Overlays()
	}
	if opt, ok := opt.(option.TimeZone); ok {
		o.timezone = opt.TimeZone()
	}
	return Realtime(mapPath[t], coords, &o)
}
