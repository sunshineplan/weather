package aqi

var _ Level = level{}

type level struct {
	level string
	color string
}

func NewLevel(lv, color string) Level { return level{lv, color} }
func (i level) String() string        { return i.level }
func (i level) Color() string         { return i.color }

var _ AQI = aqi{}

type aqi struct {
	typ   Type
	value int
	level Level
}

func NewAQI(typ Type, value int, level Level) AQI { return aqi{typ, value, level} }
func (i aqi) Type() Type                          { return i.typ }
func (i aqi) Value() int                          { return i.value }
func (i aqi) Level() Level                        { return i.level }

var _ Pollutant = pollutant{}

type pollutant struct {
	kind  Kind
	unit  string
	value float64
	level Level
}

func NewPollutant(kind Kind, unit string, value float64, level Level) Pollutant {
	return pollutant{kind, unit, value, level}
}
func (i pollutant) Kind() Kind     { return i.kind }
func (i pollutant) Unit() string   { return i.unit }
func (i pollutant) Value() float64 { return i.value }
func (i pollutant) Level() Level   { return i.level }
