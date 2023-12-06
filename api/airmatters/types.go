package airmatters

import "github.com/sunshineplan/weather/aqi"

type Place struct {
	ID          string  `json:"place_id,omitempty"`
	Type        string  `json:"type,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Latitude    float64 `json:"lat,omitempty"`
	Longitude   float64 `json:"lon,omitempty"`
}

type Current struct {
	AQIType aqi.Type
	Time    string `json:"update_time,omitempty"`
	Items   []Item `json:"readings,omitempty"`
}

type Item struct {
	AQIType aqi.Type
	Color   string `json:"color,omitempty"`
	Kind    string `json:"kind,omitempty"`
	Level   string `json:"level,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Value   string `json:"value,omitempty"`
	Unit    string `json:"unit,omitempty"`
	Time    string `json:"time,omitempty"`
}
