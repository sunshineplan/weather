package visualcrossing

import "github.com/sunshineplan/weather"

func (current *Current) Convert() weather.Current {
	return weather.Current{
		Datetime:      current.Datetime,
		DatetimeEpoch: current.DatetimeEpoch,
		Temp:          current.Temp,
		FeelsLike:     current.FeelsLike,
		WindSpeed:     current.WindSpeed,
		WindDegree:    current.WindDir,
		Pressure:      current.Pressure,
		Precip:        current.Precip,
		PrecipType:    current.PrecipType,
		Humidity:      current.Humidity,
		CloudCover:    current.CloudCover,
		Visibility:    current.Visibility,
		UVIndex:       current.UVIndex,
		WindGust:      current.WindGust,
		Condition:     current.Conditions,
		Icon:          ConvertIcon(current.Icon),
	}
}

func ConvertDays(src []Day) (days []weather.Day) {
	for _, i := range src {
		days = append(days, weather.Day{
			Date:         i.Datetime,
			DateEpoch:    i.DatetimeEpoch,
			TempMax:      i.TempMax,
			TempMin:      i.TempMin,
			Temp:         i.Temp,
			FeelsLikeMax: i.FeelsLikeMax,
			FeelsLikeMin: i.FeelsLikeMin,
			FeelsLike:    i.FeelsLike,
			Humidity:     i.Humidity,
			Dew:          i.Dew,
			Precip:       i.Precip,
			PrecipProb:   i.PrecipProb,
			PrecipCover:  i.PrecipCover,
			Snow:         i.Snow,
			SnowDepth:    i.SnowDepth,
			PrecipType:   i.PrecipType,
			WindGust:     i.WindGust,
			WindSpeed:    i.WindSpeed,
			WindDir:      i.WindDir,
			Pressure:     i.Pressure,
			CloudCover:   i.CloudCover,
			Visibility:   i.Visibility,
			UVIndex:      i.UVIndex,
			SevereRisk:   i.SevereRisk,
			Condition:    i.Conditions,
			Icon:         ConvertIcon(i.Icon),
			Hours:        ConvertHours(i.Hours),
		})
	}
	return
}

func ConvertHours(src []Hour) (hours []weather.Hour) {
	for _, i := range src {
		hours = append(hours, weather.Hour{
			Time:           i.Datetime,
			TimeEpoch:      i.DatetimeEpoch,
			Temp:           i.Temp,
			FeelsLike:      i.FeelsLike,
			Humidity:       i.Humidity,
			Dew:            i.Dew,
			Precip:         i.Precip,
			PrecipProb:     i.PrecipProb,
			Snow:           i.Snow,
			SnowDepth:      i.SnowDepth,
			PrecipType:     i.PrecipType,
			WindGust:       i.WindGust,
			WindSpeed:      i.WindSpeed,
			WindDir:        i.WindDir,
			Pressure:       i.Pressure,
			Visibility:     i.Visibility,
			CloudCover:     i.CloudCover,
			SolarRadiation: i.SolarRadiation,
			SolarEnergy:    i.SolarEnergy,
			UVIndex:        i.UVIndex,
			SevereRisk:     i.SevereRisk,
			Condition:      i.Conditions,
			Icon:           ConvertIcon(i.Icon),
		})
	}
	return
}

func ConvertIcon(icon string) string {
	//return "https://www.visualcrossing.com/img/" + icon + ".svg"
	//return "https://cdn.jsdelivr.net/gh/visualcrossing/WeatherIcons@main/PNG/2nd%20Set%20-%20Color/" + icon + ".png"
	//return "https://raw.githubusercontent.com/visualcrossing/WeatherIcons/main/PNG/2nd%20Set%20-%20Color/" + icon + ".png"
	return "https://weather.sunshineplan.cc/img/" + icon + ".png"
}
