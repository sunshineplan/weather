package visualcrossing

import (
	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/unit"
	"github.com/sunshineplan/weather/unit/wind"
)

type Response struct {
	Latitude          float64  `json:"latitude,omitempty"`
	Longitude         float64  `json:"longitude,omitempty"`
	Address           string   `json:"address,omitempty"`
	ResolvedAddress   string   `json:"resolvedAddress,omitempty"`
	Timezone          string   `json:"timezone,omitempty"`
	TzOffset          float64  `json:"tzoffset,omitempty"`
	CurrentConditions *Current `json:"currentConditions,omitempty"`
	Days              []Day    `json:"days,omitempty"`
}

type Current struct {
	Datetime       string            `json:"datetime,omitempty"`
	DatetimeEpoch  unit.UnixTime     `json:"datetimeEpoch,omitempty"`
	Temp           unit.Celsius      `json:"temp"`
	FeelsLike      unit.Celsius      `json:"feelslike"`
	Humidity       unit.Percent      `json:"humidity"`
	Dew            unit.Celsius      `json:"dew"`
	Precip         float64           `json:"precip,omitempty"`
	PrecipProb     unit.Percent      `json:"precipprob,omitempty"`
	Snow           float64           `json:"snow,omitempty"`
	SnowDepth      float64           `json:"snowdepth,omitempty"`
	PrecipType     []string          `json:"preciptype,omitempty"`
	WindGust       wind.KPH          `json:"windgust,omitempty"`
	WindSpeed      wind.KPH          `json:"windspeed,omitempty"`
	WindDir        wind.Degree       `json:"winddir,omitempty"`
	Pressure       float64           `json:"pressure,omitempty"`
	Visibility     float64           `json:"visibility,omitempty"`
	CloudCover     unit.Percent      `json:"cloudcover"`
	SolarRadiation float64           `json:"solarradiation,omitempty"`
	SolarEnergy    float64           `json:"solarenergy,omitempty"`
	UVIndex        unit.UVIndex      `json:"uvindex,omitempty"`
	Conditions     weather.Condition `json:"conditions,omitempty"`
	Icon           string            `json:"icon,omitempty"`
	Sunrise        string            `json:"sunrise,omitempty"`
	SunriseEpoch   int64             `json:"sunriseEpoch,omitempty"`
	Sunset         string            `json:"sunset,omitempty"`
	SunsetEpoch    int64             `json:"sunsetEpoch,omitempty"`
	MoonPhase      float64           `json:"moonphase,omitempty"`
}

type Day struct {
	Datetime       string            `json:"datetime,omitempty"`
	DatetimeEpoch  unit.UnixTime     `json:"datetimeEpoch,omitempty"`
	TempMax        unit.Celsius      `json:"tempmax"`
	TempMin        unit.Celsius      `json:"tempmin"`
	Temp           unit.Celsius      `json:"temp"`
	FeelsLikeMax   unit.Celsius      `json:"feelslikemax"`
	FeelsLikeMin   unit.Celsius      `json:"feelslikemin"`
	FeelsLike      unit.Celsius      `json:"feelslike"`
	Humidity       unit.Percent      `json:"humidity"`
	Dew            unit.Celsius      `json:"dew"`
	Precip         float64           `json:"precip,omitempty"`
	PrecipProb     unit.Percent      `json:"precipprob,omitempty"`
	PrecipCover    unit.Percent      `json:"precipcover,omitempty"`
	Snow           float64           `json:"snow,omitempty"`
	SnowDepth      float64           `json:"snowdepth,omitempty"`
	PrecipType     []string          `json:"preciptype,omitempty"`
	WindGust       wind.KPH          `json:"windgust,omitempty"`
	WindSpeed      wind.KPH          `json:"windspeed,omitempty"`
	WindDir        wind.Degree       `json:"winddir,omitempty"`
	Pressure       float64           `json:"pressure,omitempty"`
	CloudCover     unit.Percent      `json:"cloudcover"`
	Visibility     float64           `json:"visibility,omitempty"`
	SolarRadiation float64           `json:"solarradiation,omitempty"`
	SolarEnergy    float64           `json:"solarenergy,omitempty"`
	UVIndex        unit.UVIndex      `json:"uvindex,omitempty"`
	SevereRisk     float64           `json:"severerisk,omitempty"`
	Sunrise        string            `json:"sunrise,omitempty"`
	SunriseEpoch   int64             `json:"sunriseEpoch,omitempty"`
	Sunset         string            `json:"sunset,omitempty"`
	SunsetEpoch    int64             `json:"sunsetEpoch,omitempty"`
	MoonPhase      float64           `json:"moonphase,omitempty"`
	Conditions     weather.Condition `json:"conditions,omitempty"`
	Description    string            `json:"description,omitempty"`
	Icon           string            `json:"icon,omitempty"`
	Hours          []Hour            `json:"hours,omitempty"`
}

type Hour struct {
	Datetime       string            `json:"datetime,omitempty"`
	DatetimeEpoch  unit.UnixTime     `json:"datetimeEpoch,omitempty"`
	Temp           unit.Celsius      `json:"temp"`
	FeelsLike      unit.Celsius      `json:"feelslike"`
	Humidity       unit.Percent      `json:"humidity"`
	Dew            unit.Celsius      `json:"dew"`
	Precip         float64           `json:"precip,omitempty"`
	PrecipProb     unit.Percent      `json:"precipprob,omitempty"`
	Snow           float64           `json:"snow,omitempty"`
	SnowDepth      float64           `json:"snowdepth,omitempty"`
	PrecipType     []string          `json:"preciptype,omitempty"`
	WindGust       wind.KPH          `json:"windgust,omitempty"`
	WindSpeed      wind.KPH          `json:"windspeed,omitempty"`
	WindDir        wind.Degree       `json:"winddir,omitempty"`
	Pressure       float64           `json:"pressure,omitempty"`
	Visibility     float64           `json:"visibility,omitempty"`
	CloudCover     unit.Percent      `json:"cloudcover"`
	SolarRadiation float64           `json:"solarradiation,omitempty"`
	SolarEnergy    float64           `json:"solarenergy,omitempty"`
	UVIndex        unit.UVIndex      `json:"uvindex,omitempty"`
	SevereRisk     float64           `json:"severerisk,omitempty"`
	Conditions     weather.Condition `json:"conditions,omitempty"`
	Icon           string            `json:"icon,omitempty"`
}
