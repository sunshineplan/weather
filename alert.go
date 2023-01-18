package weather

import "math"

func (res *Response) WillRainSnow() (hour *ForecastHour, start bool) {
	if res.Forecast == nil {
		panic("forecast is nil")
	}
	if res.Current == nil {
		panic("current is nil")
	}

	var hours []ForecastHour
	for _, i := range res.Forecast.Forecastday {
		hours = append(hours, i.Hour...)
	}

	for _, i := range hours {
		switch res.Current.PrecipMm {
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

func (res *Response) WillUpDown(difference float64) (day *ForecastForecastday, up bool) {
	if res.Forecast == nil {
		panic("forecast is nil")
	}
	if res.Current == nil {
		panic("current is nil")
	}
	if len(res.Forecast.Forecastday) == 0 {
		panic("length of forecastday is zero")
	}

	today := res.Forecast.Forecastday[0]
	for _, i := range res.Forecast.Forecastday[1:] {
		if math.Abs(today.Day.Maxtemp-i.Day.Maxtemp) >= difference {
			if today.Day.Maxtemp < i.Day.Maxtemp {
				up = true
			}
		} else if math.Abs(today.Day.Mintemp-i.Day.Mintemp) >= difference {
			if today.Day.Mintemp < i.Day.Mintemp {
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
