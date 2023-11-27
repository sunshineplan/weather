package weather

import "time"

type API interface {
	Coordinates(string) (latitude, longitude float64, err error)
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
