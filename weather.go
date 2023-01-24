package weather

import (
	"errors"
	"math"
	"time"
)

type API interface {
	Realtime(string) (Current, error)
	Forecast(string, int) (Current, []Day, error)
	History(string, time.Time) (Day, error)
}

type Weather struct {
	API
}

func New(api API) *Weather {
	return &Weather{api}
}

func (w Weather) WillRainSnow(query string, n int) (hour *Hour, start bool, err error) {
	current, days, err := w.Forecast(query, n)
	if err != nil {
		return
	}

	var hours []Hour
	for _, i := range days {
		for _, ii := range i.Hours {
			if ii.TimeEpoch > i.DateEpoch {
				hours = append(hours, ii)
			}
		}
	}

	for _, i := range hours {
		switch current.PrecipMm {
		case 0:
			if i.WillItRain+i.WillItSnow > 0 {
				hour = &i
				start = true
				return
			}
		default:
			if i.WillItRain+i.WillItSnow == 0 {
				hour = &i
				return
			}
		}
	}
	return
}

func (w Weather) WillUpDown(difference float64, query string, n int) (day *Day, up bool, err error) {
	_, days, err := w.Forecast(query, n)
	if err != nil {
		return
	}
	if len(days) == 0 {
		err = errors.New("length of forecast days is zero")
		return
	}

	today := days[0]
	for _, i := range days[1:] {
		if math.Abs(today.MaxTemp-i.MaxTemp) >= difference {
			if today.MaxTemp < i.MaxTemp {
				up = true
			}
		} else if math.Abs(today.MinTemp-i.MinTemp) >= difference {
			if today.MinTemp < i.MinTemp {
				up = true
			}
		} else {
			continue
		}
		day = &i
		return
	}
	return
}
