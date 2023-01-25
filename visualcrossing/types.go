package visualcrossing

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
	Datetime       string   `json:"datetime,omitempty"`
	DatetimeEpoch  int64    `json:"datetimeEpoch,omitempty"`
	Temp           float64  `json:"temp"`
	FeelsLike      float64  `json:"feelslike"`
	Humidity       float64  `json:"humidity"`
	Dew            float64  `json:"dew"`
	Precip         float64  `json:"precip,omitempty"`
	PrecipProb     float64  `json:"precipprob,omitempty"`
	Snow           float64  `json:"snow,omitempty"`
	SnowDepth      float64  `json:"snowdepth,omitempty"`
	PrecipType     []string `json:"preciptype,omitempty"`
	WindGust       float64  `json:"windgust,omitempty"`
	WindSpeed      float64  `json:"windspeed,omitempty"`
	WindDir        float64  `json:"winddir,omitempty"`
	Pressure       float64  `json:"pressure,omitempty"`
	Visibility     float64  `json:"visibility,omitempty"`
	CloudCover     float64  `json:"cloudcover"`
	SolarRadiation float64  `json:"solarradiation,omitempty"`
	SolarEnergy    float64  `json:"solarenergy,omitempty"`
	UVIndex        float64  `json:"uvindex,omitempty"`
	Conditions     string   `json:"conditions,omitempty"`
	Icon           string   `json:"icon,omitempty"`
	Sunrise        string   `json:"sunrise,omitempty"`
	SunriseEpoch   int64    `json:"sunriseEpoch,omitempty"`
	Sunset         string   `json:"sunset,omitempty"`
	SunsetEpoch    int64    `json:"sunsetEpoch,omitempty"`
	MoonPhase      float64  `json:"moonphase,omitempty"`
}

type Day struct {
	Datetime       string   `json:"datetime,omitempty"`
	DatetimeEpoch  int64    `json:"datetimeEpoch,omitempty"`
	TempMax        float64  `json:"tempmax"`
	TempMin        float64  `json:"tempmin"`
	Temp           float64  `json:"temp"`
	FeelsLikeMax   float64  `json:"feelslikemax"`
	FeelsLikeMin   float64  `json:"feelslikemin"`
	FeelsLike      float64  `json:"feelslike"`
	Humidity       float64  `json:"humidity"`
	Dew            float64  `json:"dew"`
	Precip         float64  `json:"precip,omitempty"`
	PrecipProb     float64  `json:"precipprob,omitempty"`
	PrecipCover    float64  `json:"precipcover,omitempty"`
	Snow           float64  `json:"snow,omitempty"`
	SnowDepth      float64  `json:"snowdepth,omitempty"`
	PrecipType     []string `json:"preciptype,omitempty"`
	WindGust       float64  `json:"windgust,omitempty"`
	WindSpeed      float64  `json:"windspeed,omitempty"`
	WindDir        float64  `json:"winddir,omitempty"`
	Pressure       float64  `json:"pressure,omitempty"`
	CloudCover     float64  `json:"cloudcover"`
	Visibility     float64  `json:"visibility,omitempty"`
	SolarRadiation float64  `json:"solarradiation,omitempty"`
	SolarEnergy    float64  `json:"solarenergy,omitempty"`
	UVIndex        float64  `json:"uvindex,omitempty"`
	SevereRisk     float64  `json:"severerisk,omitempty"`
	Sunrise        string   `json:"sunrise,omitempty"`
	SunriseEpoch   int64    `json:"sunriseEpoch,omitempty"`
	Sunset         string   `json:"sunset,omitempty"`
	SunsetEpoch    int64    `json:"sunsetEpoch,omitempty"`
	MoonPhase      float64  `json:"moonphase,omitempty"`
	Conditions     string   `json:"conditions,omitempty"`
	Description    string   `json:"description,omitempty"`
	Icon           string   `json:"icon,omitempty"`
	Hours          []Hour   `json:"hours,omitempty"`
}

type Hour struct {
	Datetime       string   `json:"datetime,omitempty"`
	DatetimeEpoch  int64    `json:"datetimeEpoch,omitempty"`
	Temp           float64  `json:"temp"`
	FeelsLike      float64  `json:"feelslike"`
	Humidity       float64  `json:"humidity"`
	Dew            float64  `json:"dew"`
	Precip         float64  `json:"precip,omitempty"`
	PrecipProb     float64  `json:"precipprob,omitempty"`
	Snow           float64  `json:"snow,omitempty"`
	SnowDepth      float64  `json:"snowdepth,omitempty"`
	PrecipType     []string `json:"preciptype,omitempty"`
	WindGust       float64  `json:"windgust,omitempty"`
	WindSpeed      float64  `json:"windspeed,omitempty"`
	WindDir        float64  `json:"winddir,omitempty"`
	Pressure       float64  `json:"pressure,omitempty"`
	Visibility     float64  `json:"visibility,omitempty"`
	CloudCover     float64  `json:"cloudcover"`
	SolarRadiation float64  `json:"solarradiation,omitempty"`
	SolarEnergy    float64  `json:"solarenergy,omitempty"`
	UVIndex        float64  `json:"uvindex,omitempty"`
	SevereRisk     float64  `json:"severerisk,omitempty"`
	Conditions     string   `json:"conditions,omitempty"`
	Icon           string   `json:"icon,omitempty"`
}
