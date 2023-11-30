package weather

import (
	"time"

	"github.com/sunshineplan/weather/unit/coordinates"
)

type API interface {
	coordinates.GeoLocator
	Realtime(string) (Current, error)
	Forecast(string, int) ([]Day, error)
	History(string, time.Time) (Day, error)
}
