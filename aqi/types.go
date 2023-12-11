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
	Australia Type = iota + 1
	Canada
	China
	Europe
	India
	Netherland
	UK
	US
)

func (t Type) String() string {
	switch t {
	case Australia:
		return "AQI(Australia)"
	case Canada:
		return "AQHI(Canada)"
	case China:
		return "AQI(China)"
	case Europe:
		return "CAQI(Europe)"
	case India:
		return "AQI(India)"
	case Netherland:
		return "AQI(Netherland)"
	case UK:
		return "DAQI(UK)"
	case US:
		return "AQI(US)"
	}
	return "Unknown AQI Type"
}

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
	CO Kind = iota + 1
	NO2
	O3
	PM2Dot5
	PM10
	SO2
)

func (k Kind) String() string {
	switch k {
	case CO:
		return "CO"
	case NO2:
		return "NO2"
	case O3:
		return "O3"
	case PM2Dot5:
		return "PM2.5"
	case PM10:
		return "PM10"
	case SO2:
		return "SO2"
	}
	return "Unknown Kind"
}

func (k Kind) HTML() string {
	switch k {
	case CO:
		return "CO"
	case NO2:
		return "NO<sub>2</sub>"
	case O3:
		return "O<sub>3</sub>"
	case PM2Dot5:
		return "PM<sub>2.5</sub>"
	case PM10:
		return "PM<sub>10</sub>"
	case SO2:
		return "SO<sub>2</sub>"
	}
	return "Unknown Kind"
}