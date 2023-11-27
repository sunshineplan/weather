package weather

import (
	"time"

	"github.com/sunshineplan/weather/unit/coordinates"
)

type API interface {
	Coordinates(string) (coordinates.Coordinates, error)
	Realtime(string) (Current, error)
	Forecast(string, int) ([]Day, error)
	History(string, time.Time) (Day, error)
}

type Weather struct {
	API
}

func New(api API) *Weather {
	return &Weather{api}
}
