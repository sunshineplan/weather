package visualcrossing

import "github.com/sunshineplan/weather"

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
	Datetime       string              `json:"datetime,omitempty"`
	DatetimeEpoch  int64               `json:"datetimeEpoch,omitempty"`
	Temp           weather.Temperature `json:"temp"`
	FeelsLike      weather.Temperature `json:"feelslike"`
	Humidity       weather.Percent     `json:"humidity"`
	Dew            weather.Temperature `json:"dew"`
	Precip         float64             `json:"precip,omitempty"`
	PrecipProb     weather.Percent     `json:"precipprob,omitempty"`
	Snow           float64             `json:"snow,omitempty"`
	SnowDepth      float64             `json:"snowdepth,omitempty"`
	PrecipType     []string            `json:"preciptype,omitempty"`
	WindGust       float64             `json:"windgust,omitempty"`
	WindSpeed      float64             `json:"windspeed,omitempty"`
	WindDir        float64             `json:"winddir,omitempty"`
	Pressure       float64             `json:"pressure,omitempty"`
	Visibility     float64             `json:"visibility,omitempty"`
	CloudCover     weather.Percent     `json:"cloudcover"`
	SolarRadiation float64             `json:"solarradiation,omitempty"`
	SolarEnergy    float64             `json:"solarenergy,omitempty"`
	UVIndex        float64             `json:"uvindex,omitempty"`
	Conditions     weather.Condition   `json:"conditions,omitempty"`
	Icon           string              `json:"icon,omitempty"`
	Sunrise        string              `json:"sunrise,omitempty"`
	SunriseEpoch   int64               `json:"sunriseEpoch,omitempty"`
	Sunset         string              `json:"sunset,omitempty"`
	SunsetEpoch    int64               `json:"sunsetEpoch,omitempty"`
	MoonPhase      float64             `json:"moonphase,omitempty"`
}

type Day struct {
	Datetime       string              `json:"datetime,omitempty"`
	DatetimeEpoch  int64               `json:"datetimeEpoch,omitempty"`
	TempMax        weather.Temperature `json:"tempmax"`
	TempMin        weather.Temperature `json:"tempmin"`
	Temp           weather.Temperature `json:"temp"`
	FeelsLikeMax   weather.Temperature `json:"feelslikemax"`
	FeelsLikeMin   weather.Temperature `json:"feelslikemin"`
	FeelsLike      weather.Temperature `json:"feelslike"`
	Humidity       weather.Percent     `json:"humidity"`
	Dew            weather.Temperature `json:"dew"`
	Precip         float64             `json:"precip,omitempty"`
	PrecipProb     weather.Percent     `json:"precipprob,omitempty"`
	PrecipCover    weather.Percent     `json:"precipcover,omitempty"`
	Snow           float64             `json:"snow,omitempty"`
	SnowDepth      float64             `json:"snowdepth,omitempty"`
	PrecipType     []string            `json:"preciptype,omitempty"`
	WindGust       float64             `json:"windgust,omitempty"`
	WindSpeed      float64             `json:"windspeed,omitempty"`
	WindDir        float64             `json:"winddir,omitempty"`
	Pressure       float64             `json:"pressure,omitempty"`
	CloudCover     weather.Percent     `json:"cloudcover"`
	Visibility     float64             `json:"visibility,omitempty"`
	SolarRadiation float64             `json:"solarradiation,omitempty"`
	SolarEnergy    float64             `json:"solarenergy,omitempty"`
	UVIndex        float64             `json:"uvindex,omitempty"`
	SevereRisk     float64             `json:"severerisk,omitempty"`
	Sunrise        string              `json:"sunrise,omitempty"`
	SunriseEpoch   int64               `json:"sunriseEpoch,omitempty"`
	Sunset         string              `json:"sunset,omitempty"`
	SunsetEpoch    int64               `json:"sunsetEpoch,omitempty"`
	MoonPhase      float64             `json:"moonphase,omitempty"`
	Conditions     weather.Condition   `json:"conditions,omitempty"`
	Description    string              `json:"description,omitempty"`
	Icon           string              `json:"icon,omitempty"`
	Hours          []Hour              `json:"hours,omitempty"`
}

type Hour struct {
	Datetime       string              `json:"datetime,omitempty"`
	DatetimeEpoch  int64               `json:"datetimeEpoch,omitempty"`
	Temp           weather.Temperature `json:"temp"`
	FeelsLike      weather.Temperature `json:"feelslike"`
	Humidity       weather.Percent     `json:"humidity"`
	Dew            weather.Temperature `json:"dew"`
	Precip         float64             `json:"precip,omitempty"`
	PrecipProb     weather.Percent     `json:"precipprob,omitempty"`
	Snow           float64             `json:"snow,omitempty"`
	SnowDepth      float64             `json:"snowdepth,omitempty"`
	PrecipType     []string            `json:"preciptype,omitempty"`
	WindGust       float64             `json:"windgust,omitempty"`
	WindSpeed      float64             `json:"windspeed,omitempty"`
	WindDir        float64             `json:"winddir,omitempty"`
	Pressure       float64             `json:"pressure,omitempty"`
	Visibility     float64             `json:"visibility,omitempty"`
	CloudCover     weather.Percent     `json:"cloudcover"`
	SolarRadiation float64             `json:"solarradiation,omitempty"`
	SolarEnergy    float64             `json:"solarenergy,omitempty"`
	UVIndex        float64             `json:"uvindex,omitempty"`
	SevereRisk     float64             `json:"severerisk,omitempty"`
	Conditions     weather.Condition   `json:"conditions,omitempty"`
	Icon           string              `json:"icon,omitempty"`
}
