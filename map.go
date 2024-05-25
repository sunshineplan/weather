package weather

import (
	"image"
	"time"

	"github.com/sunshineplan/weather/unit/coordinates"
)

type MapType int

const (
	Satellite MapType = iota + 1
	Radar
	Precipitation
	Wind
	Temperature
	Humidity
	DewPoint
	Pressure
)

type MapAPI interface {
	URL(t MapType, coords coordinates.Coordinates, opt any) string
	Realtime(t MapType, coords coordinates.Coordinates, opt any) (time.Time, image.Image, error)
}
