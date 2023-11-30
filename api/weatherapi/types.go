package weatherapi

import (
	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/unit"
)

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
	LastUpdatedEpoch unit.UnixTime   `json:"last_updated_epoch,omitempty"`
	LastUpdated      string          `json:"last_updated,omitempty"`
	Temp             unit.Celsius    `json:"temp_c"`
	IsDay            int             `json:"is_day"`
	Condition        *Condition      `json:"condition,omitempty"`
	WindKph          unit.WindKPH    `json:"wind_kph,omitempty"`
	WindDegree       float64         `json:"wind_degree,omitempty"`
	WindDir          string          `json:"wind_dir,omitempty"`
	PressureMb       float64         `json:"pressure_mb,omitempty"`
	PrecipMm         float64         `json:"precip_mm,omitempty"`
	Humidity         weather.Percent `json:"humidity,omitempty"`
	Cloud            weather.Percent `json:"cloud"`
	FeelsLike        unit.Celsius    `json:"feelslike_c"`
	VisKm            float64         `json:"vis_km,omitempty"`
	UV               unit.UVIndex    `json:"uv,omitempty"`
	GustKph          unit.WindKPH    `json:"gust_kph,omitempty"`
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
	DateEpoch unit.UnixTime  `json:"date_epoch,omitempty"`
	Day       *ForecastDay   `json:"day,omitempty"`
	Astro     *ForecastAstro `json:"astro,omitempty"`
	Hour      []ForecastHour `json:"hour,omitempty"`
}

type ForecastDay struct {
	MaxTemp           unit.Celsius    `json:"maxtemp_c"`
	MinTemp           unit.Celsius    `json:"mintemp_c"`
	AvgTemp           unit.Celsius    `json:"avgtemp_c"`
	MaxWindKph        unit.WindKPH    `json:"maxwind_kph,omitempty"`
	TotalPrecipMm     float64         `json:"totalprecip_mm"`
	AvgVisKm          float64         `json:"avgvis_km,omitempty"`
	AvgHumidity       weather.Percent `json:"avghumidity,omitempty"`
	DailyWillItRain   int             `json:"daily_will_it_rain,omitempty"`
	DailyChanceOfRain weather.Percent `json:"daily_chance_of_rain,omitempty"`
	DailyWillItSnow   int             `json:"daily_will_it_snow,omitempty"`
	DailyChanceOfSnow weather.Percent `json:"daily_chance_of_snow,omitempty"`
	Condition         *Condition      `json:"condition,omitempty"`
	UV                unit.UVIndex    `json:"uv,omitempty"`
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
	Time         string          `json:"time,omitempty"`
	TimeEpoch    unit.UnixTime   `json:"time_epoch,omitempty"`
	Temp         unit.Celsius    `json:"temp_c"`
	IsDay        int             `json:"is_day"`
	Condition    *Condition      `json:"condition,omitempty"`
	WindKph      unit.WindKPH    `json:"wind_kph,omitempty"`
	WindDegree   float64         `json:"wind_degree,omitempty"`
	WindDir      string          `json:"wind_dir,omitempty"`
	PressureMb   float64         `json:"pressure_mb,omitempty"`
	PrecipMm     float64         `json:"precip_mm"`
	Humidity     weather.Percent `json:"humidity,omitempty"`
	Cloud        weather.Percent `json:"cloud"`
	FeelsLike    unit.Celsius    `json:"feelslike_c"`
	WindChill    unit.Celsius    `json:"windchill_c"`
	HeatIndex    unit.Celsius    `json:"heatindex_c"`
	DewPoint     unit.Celsius    `json:"dewpoint_c"`
	WillItRain   int             `json:"will_it_rain"`
	ChanceOfRain weather.Percent `json:"chance_of_rain"`
	WillItSnow   int             `json:"will_it_snow"`
	ChanceOfSnow weather.Percent `json:"chance_of_snow"`
	VisKm        float64         `json:"vis_km,omitempty"`
	GustKph      unit.WindKPH    `json:"gust_kph,omitempty"`
	UV           unit.UVIndex    `json:"uv,omitempty"`
}
