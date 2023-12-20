package aqi

import (
	"encoding"
	"errors"
	"strings"

	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/weather/unit"
)

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

var (
	_ encoding.TextMarshaler   = Type(0)
	_ encoding.TextUnmarshaler = new(Type)
)

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

func (t Type) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Type) UnmarshalText(b []byte) error {
	s := strings.ToLower(strings.TrimSpace(string(b)))
	for k, v := range map[Type][]string{
		Australia:  {"australia", "au"},
		Canada:     {"canada", "ca"},
		China:      {"china", "cn"},
		Europe:     {"europe", "eu"},
		India:      {"india", "in"},
		Netherland: {"netherlands", "nl"},
		UK:         {"uk", "gb"},
		US:         {"us", "usa"},
	} {
		for _, i := range v {
			if s == i {
				*t = k
				return nil
			}
		}
	}
	return errors.New("unknown AQI type")
}

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

func (k Kind) HTML() html.HTML {
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
