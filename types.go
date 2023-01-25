package weather

type Current struct {
	Datetime      string   `json:"datetime,omitempty"`
	DatetimeEpoch int64    `json:"datetimeEpoch,omitempty"`
	Temp          float64  `json:"temp"`
	FeelsLike     float64  `json:"feelslike"`
	Humidity      float64  `json:"humidity,omitempty"`
	Dew           float64  `json:"dew"`
	Precip        float64  `json:"precip,omitempty"`
	PrecipType    []string `json:"preciptype,omitempty"`
	WindGust      float64  `json:"windgust,omitempty"`
	WindSpeed     float64  `json:"windspeed,omitempty"`
	WindDegree    float64  `json:"winddegree,omitempty"`
	WindDir       string   `json:"winddir,omitempty"`
	Pressure      float64  `json:"pressure,omitempty"`
	Visibility    float64  `json:"visibility,omitempty"`
	CloudCover    float64  `json:"cloudcover"`
	UVIndex       float64  `json:"uvindex,omitempty"`
	Condition     string   `json:"condition,omitempty"`
	Icon          string   `json:"icon,omitempty"`
}

type Day struct {
	Date         string   `json:"date,omitempty"`
	DateEpoch    int64    `json:"dateEpoch,omitempty"`
	TempMax      float64  `json:"tempmax"`
	TempMin      float64  `json:"tempmin"`
	Temp         float64  `json:"temp"`
	FeelsLikeMax float64  `json:"feelslikemax"`
	FeelsLikeMin float64  `json:"feelslikemin"`
	FeelsLike    float64  `json:"feelslike"`
	Humidity     float64  `json:"humidity,omitempty"`
	Dew          float64  `json:"dew"`
	Precip       float64  `json:"precip,omitempty"`
	PrecipProb   float64  `json:"precipprob,omitempty"`
	PrecipCover  float64  `json:"precipcover,omitempty"`
	Snow         float64  `json:"snow,omitempty"`
	SnowDepth    float64  `json:"snowdepth,omitempty"`
	PrecipType   []string `json:"preciptype,omitempty"`
	WindGust     float64  `json:"windgust,omitempty"`
	WindSpeed    float64  `json:"windspeed,omitempty"`
	WindDir      float64  `json:"winddir,omitempty"`
	Pressure     float64  `json:"pressure,omitempty"`
	CloudCover   float64  `json:"cloudcover"`
	Visibility   float64  `json:"visibility,omitempty"`
	UVIndex      float64  `json:"uvindex,omitempty"`
	SevereRisk   float64  `json:"severerisk,omitempty"`
	Condition    string   `json:"condition,omitempty"`
	Icon         string   `json:"icon,omitempty"`
	Hours        []Hour   `json:"hours,omitempty"`
}

type Hour struct {
	Time           string   `json:"time,omitempty"`
	TimeEpoch      int64    `json:"timeEpoch,omitempty"`
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
	Condition      string   `json:"condition,omitempty"`
	Icon           string   `json:"icon,omitempty"`
}
