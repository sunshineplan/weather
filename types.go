package weather

type Current struct {
	LastUpdatedEpoch int64
	LastUpdated      string
	Temp             float64
	FeelsLike        float64
	WindKph          float64
	WindDegree       int
	WindDir          string
	PressureMb       float64
	PrecipMm         float64
	Humidity         int
	Cloud            int
	VisKm            float64
	Uv               float64
	GustKph          float64
	Condition        string
	Icon             string
}

type Day struct {
	Date              string
	DateEpoch         int64 `json:",omitempty"`
	MaxTemp           float64
	MinTemp           float64
	AvgTemp           float64
	MaxWindKph        float64
	TotalPrecipMm     float64
	AvgVisKm          float64
	AvgHumidity       float64
	DailyWillItRain   int     `json:",omitempty"`
	DailyChanceOfRain int     `json:",omitempty"`
	DailyWillItSnow   int     `json:",omitempty"`
	DailyChanceOfSnow int     `json:",omitempty"`
	Uv                float64 `json:",omitempty"`
	Condition         string
	Icon              string `json:",omitempty"`
	Hours             []Hour `json:",omitempty"`
}

type Hour struct {
	TimeEpoch    int64
	Time         string
	Temp         float64
	WindKph      float64
	WindDegree   int
	WindDir      string
	PressureMb   float64
	PrecipMm     float64
	Humidity     int
	Cloud        int
	FeelsLike    float64
	WindChill    float64
	HeatIndex    float64
	DewPoint     float64
	WillItRain   int
	ChanceOfRain int
	WillItSnow   int
	ChanceOfSnow int
	VisKm        float64
	GustKph      float64
	Uv           float64
	Condition    string
	Icon         string
}
