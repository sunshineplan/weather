package maps

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
	URL(mt MapType, t time.Time, coords coordinates.Coordinates, opt any) string
	Map(mt MapType, t time.Time, coords coordinates.Coordinates, opt any) (time.Time, image.Image, error)
	Realtime(mt MapType, coords coordinates.Coordinates, opt any) (time.Time, image.Image, error)
}
