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
		switch current.Precip {
		case 0:
			if i.PrecipProb > 0 {
				hour = &i
				start = true
				return
			}
		default:
			if i.PrecipProb == 0 {
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
		if math.Abs(today.TempMax-i.TempMax) >= difference {
			if today.TempMax < i.TempMax {
				up = true
			}
		} else if math.Abs(today.TempMin-i.TempMin) >= difference {
			if today.TempMin < i.TempMin {
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
