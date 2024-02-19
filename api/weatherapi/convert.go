package weatherapi

import "github.com/sunshineplan/weather"

func (current *Current) Convert() weather.Current {
	return weather.Current{
		Datetime:      current.LastUpdated,
		DatetimeEpoch: current.LastUpdatedEpoch,
		Temp:          current.Temp,
		FeelsLike:     current.FeelsLike,
		WindSpeed:     current.WindKph,
		WindDegree:    current.WindDegree,
		Pressure:      current.PressureMb,
		Precip:        current.PrecipMm,
		Humidity:      current.Humidity,
		CloudCover:    current.Cloud,
		Visibility:    current.VisKm,
		UVIndex:       current.UV,
		WindGust:      current.GustKph,
		Condition:     current.Condition.Text,
		Icon:          current.Condition.Icon,
	}
}

func (forecast *Forecast) Convert() (days []weather.Day) {
	for _, i := range forecast.Forecastday {
		days = append(days, weather.Day{
			Date:       i.Date,
			DateEpoch:  i.DateEpoch,
			TempMax:    i.Day.MaxTemp,
			TempMin:    i.Day.MinTemp,
			Temp:       i.Day.AvgTemp,
			WindSpeed:  i.Day.MaxWindKph,
			Precip:     i.Day.TotalPrecipMm,
			PrecipProb: i.Day.DailyChanceOfRain.Max(i.Day.DailyChanceOfSnow),
			Visibility: i.Day.AvgVisKm,
			Humidity:   i.Day.AvgHumidity,
			UVIndex:    i.Day.UV,
			Condition:  i.Day.Condition.Text,
			Icon:       i.Day.Condition.Icon,
			Hours:      ConvertHours(i.Hour),
		})
	}
	return
}

func ConvertHours(src []ForecastHour) (hours []weather.Hour) {
	for _, i := range src {
		hours = append(hours, weather.Hour{
			Time:       i.Time,
			TimeEpoch:  i.TimeEpoch,
			Temp:       i.Temp,
			FeelsLike:  i.FeelsLike,
			Humidity:   i.Humidity,
			Dew:        i.DewPoint,
			Precip:     i.PrecipMm,
			PrecipProb: i.ChanceOfRain.Max(i.ChanceOfSnow),
			WindGust:   i.GustKph,
			WindSpeed:  i.WindKph,
			WindDir:    i.WindDegree,
			Pressure:   i.PressureMb,
			Visibility: i.VisKm,
			CloudCover: i.Cloud,
			UVIndex:    i.UV,
			Condition:  i.Condition.Text,
			Icon:       i.Condition.Icon,
		})
	}
	return
}
