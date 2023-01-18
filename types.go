package weather

type Response struct {
	Location *Location `json:"location,omitempty"`
	Current  *Current  `json:"current,omitempty"`
	Forecast *Forecast `json:"forecast,omitempty"`
}

type Location struct {
	Name           string  `json:"name,omitempty"`
	Region         string  `json:"region,omitempty"`
	Country        string  `json:"country,omitempty"`
	Lat            float64 `json:"lat,omitempty"`
	Lon            float64 `json:"lon,omitempty"`
	TzId           string  `json:"tz_id,omitempty"`
	LocaltimeEpoch int64   `json:"localtime_epoch,omitempty"`
	Localtime      string  `json:"localtime,omitempty"`
}

type Current struct {
	LastUpdatedEpoch int64      `json:"last_updated_epoch,omitempty"`
	LastUpdated      string     `json:"last_updated,omitempty"`
	Temp             float64    `json:"temp_c,omitempty"`
	IsDay            int        `json:"is_day,omitempty"`
	Condition        *Condition `json:"condition,omitempty"`
	WindKph          float64    `json:"wind_kph,omitempty"`
	WindDegree       int        `json:"wind_degree,omitempty"`
	WindDir          string     `json:"wind_dir,omitempty"`
	PressureMb       float64    `json:"pressure_mb,omitempty"`
	PrecipMm         float64    `json:"precip_mm,omitempty"`
	Humidity         int        `json:"humidity,omitempty"`
	Cloud            int        `json:"cloud,omitempty"`
	Feelslike        float64    `json:"feelslike_c,omitempty"`
	VisKm            float64    `json:"vis_km,omitempty"`
	Uv               float64    `json:"uv,omitempty"`
	GustKph          float64    `json:"gust_kph,omitempty"`
}

type Condition struct {
	Text string `json:"text,omitempty"`
	Icon string `json:"icon,omitempty"`
	Code int    `json:"code,omitempty"`
}

type Forecast struct {
	Forecastday []ForecastForecastday `json:"forecastday,omitempty"`
}

type ForecastForecastday struct {
	Date      string         `json:"date,omitempty"`
	DateEpoch int64          `json:"date_epoch,omitempty"`
	Day       *ForecastDay   `json:"day,omitempty"`
	Astro     *ForecastAstro `json:"astro,omitempty"`
	Hour      []ForecastHour `json:"hour,omitempty"`
}

type ForecastDay struct {
	Maxtemp           float64    `json:"maxtemp_c,omitempty"`
	Mintemp           float64    `json:"mintemp_c,omitempty"`
	Avgtemp           float64    `json:"avgtemp_c,omitempty"`
	MaxwindKph        float64    `json:"maxwind_kph,omitempty"`
	TotalprecipMm     float64    `json:"totalprecip_mm,omitempty"`
	AvgvisKm          float64    `json:"avgvis_km,omitempty"`
	Avghumidity       float64    `json:"avghumidity,omitempty"`
	DailyWillItRain   int        `json:"daily_will_it_rain,omitempty"`
	DailyChanceOfRain int        `json:"daily_chance_of_rain,omitempty"`
	DailyWillItSnow   int        `json:"daily_will_it_snow,omitempty"`
	DailyChanceOfSnow int        `json:"daily_chance_of_snow,omitempty"`
	Condition         *Condition `json:"condition,omitempty"`
	Uv                float64    `json:"uv,omitempty"`
}

type ForecastAstro struct {
	Sunrise          string `json:"sunrise,omitempty"`
	Sunset           string `json:"sunset,omitempty"`
	Moonrise         string `json:"moonrise,omitempty"`
	Moonset          string `json:"moonset,omitempty"`
	MoonPhase        string `json:"moon_phase,omitempty"`
	MoonIllumination string `json:"moon_illumination,omitempty"`
}

type ForecastHour struct {
	TimeEpoch    int64      `json:"time_epoch,omitempty"`
	Time         string     `json:"time,omitempty"`
	Temp         float64    `json:"temp_c,omitempty"`
	IsDay        int        `json:"is_day,omitempty"`
	Condition    *Condition `json:"condition,omitempty"`
	WindKph      float64    `json:"wind_kph,omitempty"`
	WindDegree   int        `json:"wind_degree,omitempty"`
	WindDir      string     `json:"wind_dir,omitempty"`
	PressureMb   float64    `json:"pressure_mb,omitempty"`
	PrecipMm     float64    `json:"precip_mm,omitempty"`
	Humidity     int        `json:"humidity,omitempty"`
	Cloud        int        `json:"cloud,omitempty"`
	Feelslike    float64    `json:"feelslike_c,omitempty"`
	Windchill    float64    `json:"windchill_c,omitempty"`
	Heatindex    float64    `json:"heatindex_c,omitempty"`
	Dewpoint     float64    `json:"dewpoint_c,omitempty"`
	WillItRain   int        `json:"will_it_rain,omitempty"`
	ChanceOfRain int        `json:"chance_of_rain,omitempty"`
	WillItSnow   int        `json:"will_it_snow,omitempty"`
	ChanceOfSnow int        `json:"chance_of_snow,omitempty"`
	VisKm        float64    `json:"vis_km,omitempty"`
	GustKph      float64    `json:"gust_kph,omitempty"`
	Uv           float64    `json:"uv,omitempty"`
}
