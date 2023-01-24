package weatherapi

import "github.com/sunshineplan/weather"

func (current *Current) Convert() weather.Current {
	return weather.Current{
		LastUpdatedEpoch: current.LastUpdatedEpoch,
		LastUpdated:      current.LastUpdated,
		Temp:             current.Temp,
		FeelsLike:        current.FeelsLike,
		WindKph:          current.WindKph,
		WindDegree:       current.WindDegree,
		WindDir:          current.WindDir,
		PressureMb:       current.PressureMb,
		PrecipMm:         current.PrecipMm,
		Humidity:         current.Humidity,
		Cloud:            current.Cloud,
		VisKm:            current.VisKm,
		Uv:               current.Uv,
		GustKph:          current.GustKph,
		Condition:        current.Condition.Text,
		Icon:             current.Condition.Icon,
	}
}

func (forecast *Forecast) Convert() (days []weather.Day) {
	for _, i := range forecast.Forecastday {
		days = append(days, weather.Day{
			Date:              i.Date,
			DateEpoch:         i.DateEpoch,
			MaxTemp:           i.Day.MaxTemp,
			MinTemp:           i.Day.MinTemp,
			AvgTemp:           i.Day.AvgTemp,
			MaxWindKph:        i.Day.MaxWindKph,
			TotalPrecipMm:     i.Day.TotalPrecipMm,
			AvgVisKm:          i.Day.AvgVisKm,
			AvgHumidity:       i.Day.AvgHumidity,
			DailyWillItRain:   i.Day.DailyWillItRain,
			DailyChanceOfRain: i.Day.DailyChanceOfRain,
			DailyWillItSnow:   i.Day.DailyWillItSnow,
			DailyChanceOfSnow: i.Day.DailyChanceOfSnow,
			Uv:                i.Day.Uv,
			Condition:         i.Day.Condition.Text,
			Icon:              i.Day.Condition.Icon,
			Hours:             ConvertHours(i.Hour),
		})
	}
	return
}

func ConvertHours(src []ForecastHour) (hours []weather.Hour) {
	for _, i := range src {
		hours = append(hours, weather.Hour{
			TimeEpoch:    i.TimeEpoch,
			Time:         i.Time,
			Temp:         i.Temp,
			WindKph:      i.WindKph,
			WindDegree:   i.WindDegree,
			WindDir:      i.WindDir,
			PressureMb:   i.PressureMb,
			PrecipMm:     i.PrecipMm,
			Humidity:     i.Humidity,
			Cloud:        i.Cloud,
			FeelsLike:    i.FeelsLike,
			WindChill:    i.WindChill,
			HeatIndex:    i.HeatIndex,
			DewPoint:     i.DewPoint,
			WillItRain:   i.WillItRain,
			ChanceOfRain: i.ChanceOfRain,
			WillItSnow:   i.WillItSnow,
			ChanceOfSnow: i.ChanceOfSnow,
			VisKm:        i.VisKm,
			GustKph:      i.GustKph,
			Uv:           i.Uv,
			Condition:    i.Condition.Text,
			Icon:         i.Condition.Icon,
		})
	}
	return
}
