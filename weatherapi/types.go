package weatherapi

import "github.com/sunshineplan/weather"

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
	Localtime      string  `json:"localtime,omitempty"`
	LocaltimeEpoch int64   `json:"localtime_epoch,omitempty"`
}

type Current struct {
	LastUpdatedEpoch int64               `json:"last_updated_epoch,omitempty"`
	LastUpdated      string              `json:"last_updated,omitempty"`
	Temp             weather.Temperature `json:"temp_c"`
	IsDay            int                 `json:"is_day"`
	Condition        *Condition          `json:"condition,omitempty"`
	WindKph          float64             `json:"wind_kph,omitempty"`
	WindDegree       float64             `json:"wind_degree,omitempty"`
	WindDir          string              `json:"wind_dir,omitempty"`
	PressureMb       float64             `json:"pressure_mb,omitempty"`
	PrecipMm         float64             `json:"precip_mm,omitempty"`
	Humidity         weather.Percent     `json:"humidity,omitempty"`
	Cloud            weather.Percent     `json:"cloud"`
	FeelsLike        weather.Temperature `json:"feelslike_c"`
	VisKm            float64             `json:"vis_km,omitempty"`
	UV               float64             `json:"uv,omitempty"`
	GustKph          float64             `json:"gust_kph,omitempty"`
}

type Condition struct {
	Text weather.Condition `json:"text,omitempty"`
	Icon string            `json:"icon,omitempty"`
	Code int               `json:"code,omitempty"`
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
	MaxTemp           weather.Temperature `json:"maxtemp_c"`
	MinTemp           weather.Temperature `json:"mintemp_c"`
	AvgTemp           weather.Temperature `json:"avgtemp_c"`
	MaxWindKph        float64             `json:"maxwind_kph,omitempty"`
	TotalPrecipMm     float64             `json:"totalprecip_mm"`
	AvgVisKm          float64             `json:"avgvis_km,omitempty"`
	AvgHumidity       weather.Percent     `json:"avghumidity,omitempty"`
	DailyWillItRain   int                 `json:"daily_will_it_rain,omitempty"`
	DailyChanceOfRain weather.Percent     `json:"daily_chance_of_rain,omitempty"`
	DailyWillItSnow   int                 `json:"daily_will_it_snow,omitempty"`
	DailyChanceOfSnow weather.Percent     `json:"daily_chance_of_snow,omitempty"`
	Condition         *Condition          `json:"condition,omitempty"`
	UV                float64             `json:"uv,omitempty"`
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
	Time         string              `json:"time,omitempty"`
	TimeEpoch    int64               `json:"time_epoch,omitempty"`
	Temp         weather.Temperature `json:"temp_c"`
	IsDay        int                 `json:"is_day"`
	Condition    *Condition          `json:"condition,omitempty"`
	WindKph      float64             `json:"wind_kph,omitempty"`
	WindDegree   float64             `json:"wind_degree,omitempty"`
	WindDir      string              `json:"wind_dir,omitempty"`
	PressureMb   float64             `json:"pressure_mb,omitempty"`
	PrecipMm     float64             `json:"precip_mm"`
	Humidity     weather.Percent     `json:"humidity,omitempty"`
	Cloud        weather.Percent     `json:"cloud"`
	FeelsLike    weather.Temperature `json:"feelslike_c"`
	WindChill    weather.Temperature `json:"windchill_c"`
	HeatIndex    weather.Temperature `json:"heatindex_c"`
	DewPoint     weather.Temperature `json:"dewpoint_c"`
	WillItRain   int                 `json:"will_it_rain"`
	ChanceOfRain weather.Percent     `json:"chance_of_rain"`
	WillItSnow   int                 `json:"will_it_snow"`
	ChanceOfSnow weather.Percent     `json:"chance_of_snow"`
	VisKm        float64             `json:"vis_km,omitempty"`
	GustKph      float64             `json:"gust_kph,omitempty"`
	UV           float64             `json:"uv,omitempty"`
}
