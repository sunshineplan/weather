package weather

import "github.com/sunshineplan/weather/unit/coordinates"

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

type MapOption interface {
	Value() any
	Compatibility(MapAPI) bool
}

type MapAPI interface {
	Realtime(t MapType, coords coordinates.Coordinates, zoom float64, quality int, opt MapOption) ([]byte, error)
}
