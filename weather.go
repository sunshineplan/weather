package weather

import (
	"time"
)

type API interface {
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
