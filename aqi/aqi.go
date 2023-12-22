package aqi

import (
	"time"

	"github.com/sunshineplan/weather/unit/coordinates"
)

type API interface {
	coordinates.GeoLocator
	Standard(Type) ([]AQI, error)
	Realtime(Type, string) (Current, error)
	Forecast(Type, string, int) ([]Day, error)
	History(Type, string, time.Time) (Day, error)
}
