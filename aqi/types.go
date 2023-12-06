package aqi

import "github.com/sunshineplan/weather/unit"

type Current interface {
	Day
	Pollutants() []Pollutant
}

type Day interface {
	Date() string
	Unix() unit.UnixTime
	AQI() AQI
}

type AQI interface {
	Type() Type
	Value() int
	Level() Level
}

type Type int

const (
	Australia Type = iota
	Canada
	China
	Europe
	India
	Netherland
	UK
	US
)

type Level interface {
	String() string
	Color() string
}

type Pollutant interface {
	Kind() Kind
	Unit() string
	Value() float64
	Level() Level
}

type Kind int

const (
	CO Kind = iota
	NO2
	O3
	PM2Dot5
	PM10
	SO2
)
