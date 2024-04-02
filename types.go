package weather

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/weather/unit"
	"github.com/sunshineplan/weather/unit/wind"
)

type Current struct {
	Datetime      string           `json:"datetime,omitempty"`
	DatetimeEpoch unit.UnixTime    `json:"datetimeEpoch,omitempty"`
	Temp          unit.Temperature `json:"temp"`
	FeelsLike     unit.Temperature `json:"feelslike"`
	Humidity      unit.Percent     `json:"humidity,omitempty"`
	Dew           unit.Temperature `json:"dew"`
	Precip        float64          `json:"precip,omitempty"`
	PrecipType    []string         `json:"preciptype,omitempty"`
	WindGust      wind.Speed       `json:"windgust,omitempty"`
	WindSpeed     wind.Speed       `json:"windspeed,omitempty"`
	WindDegree    wind.Degree      `json:"winddegree,omitempty"`
	Pressure      float64          `json:"pressure,omitempty"`
	Visibility    float64          `json:"visibility,omitempty"`
	CloudCover    unit.Percent     `json:"cloudcover"`
	UVIndex       unit.UVIndex     `json:"uvindex,omitempty"`
	Condition     Condition        `json:"condition,omitempty"`
	Icon          string           `json:"icon,omitempty"`
}

func (current Current) TimeInfoHTML() html.HTML {
	return html.HTML(fmt.Sprintf("%s %s", current.Datetime, current.Condition.Img(current.Icon)))
}

func (current Current) HTML() html.HTML {
	return html.Div().AppendChild(
		html.Span().Style("display:list-item;margin-left:15px;list-style-type:disclosure-open").Content(current.TimeInfoHTML()),
		html.Table().AppendChild(
			html.Tbody().AppendContent(
				html.Tr(
					html.Td("Temp:"), html.Td(current.Temp),
					html.Td("FeelsLike:"), html.Td(current.FeelsLike),
					html.Td("Humidity:"), html.Td(current.Humidity),
				),
				html.Tr(
					html.Td("Pressure:"), html.Td(fmt.Sprintf("%ghPa", current.Pressure)),
					html.Td("Precip:"), html.Td(fmt.Sprintf("%gmm", current.Precip)),
					html.Td("Wind:"), html.Td(current.WindSpeed.HTML()+current.WindDegree.HTML()).Style("display:flex;align-items:center"),
				),
				html.Tr(
					html.Td("CloudCover:"), html.Td(current.CloudCover),
					html.Td("Visibility:"), html.Td(fmt.Sprintf("%gkm", current.Visibility)),
					html.Td("UVIndex:"), html.Td(current.UVIndex),
				))),
	).HTML()
}

type Day struct {
	Date         string           `json:"date,omitempty"`
	DateEpoch    unit.UnixTime    `json:"dateEpoch,omitempty"`
	TempMax      unit.Temperature `json:"tempmax"`
	TempMin      unit.Temperature `json:"tempmin"`
	Temp         unit.Temperature `json:"temp"`
	FeelsLikeMax unit.Temperature `json:"feelslikemax"`
	FeelsLikeMin unit.Temperature `json:"feelslikemin"`
	FeelsLike    unit.Temperature `json:"feelslike"`
	Humidity     unit.Percent     `json:"humidity,omitempty"`
	Dew          unit.Temperature `json:"dew"`
	Precip       float64          `json:"precip,omitempty"`
	PrecipProb   unit.Percent     `json:"precipprob,omitempty"`
	PrecipCover  unit.Percent     `json:"precipcover,omitempty"`
	Snow         float64          `json:"snow,omitempty"`
	SnowDepth    float64          `json:"snowdepth,omitempty"`
	PrecipType   []string         `json:"preciptype,omitempty"`
	WindGust     wind.Speed       `json:"windgust,omitempty"`
	WindSpeed    wind.Speed       `json:"windspeed,omitempty"`
	WindDir      wind.Degree      `json:"winddir,omitempty"`
	Pressure     float64          `json:"pressure,omitempty"`
	CloudCover   unit.Percent     `json:"cloudcover"`
	Visibility   float64          `json:"visibility,omitempty"`
	UVIndex      unit.UVIndex     `json:"uvindex,omitempty"`
	SevereRisk   float64          `json:"severerisk,omitempty"`
	Condition    Condition        `json:"condition,omitempty"`
	Description  string           `json:"description,omitempty"`
	Icon         string           `json:"icon,omitempty"`
	Hours        []Hour           `json:"hours,omitempty"`
}

func (day Day) DateInfo(condition bool) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Date: %s %s", day.Date, day.DateEpoch.Weekday())
	if condition {
		fmt.Fprintf(&b, " (%s)", day.Condition)
	}
	return b.String()
}

func (day Day) DateInfoHTML() html.HTML {
	return html.HTML(fmt.Sprintf("%s %s %s", day.Date, day.DateEpoch.Weekday(), day.Condition.Img(day.Icon)))
}

func (day Day) Temperature() string {
	var b strings.Builder
	fmt.Fprintf(&b, "TempMax: %s, TempMin: %s, Temp: %s\n", day.TempMax, day.TempMin, day.Temp)
	fmt.Fprintf(&b, "FeelsLikeMax: %s, FeelsLikeMin: %s, FeelsLike: %s", day.FeelsLikeMax, day.FeelsLikeMin, day.FeelsLike)
	return b.String()
}

func (day Day) TemperatureHTML() html.HTML {
	return html.Background().AppendChild(
		html.Tr(
			html.Td("TempMax:"), html.Td(day.TempMax),
			html.Td("TempMin:"), html.Td(day.TempMin),
			html.Td("Temp:"), html.Td(day.Temp),
		),
		html.Tr(
			html.Td("FeelsLikeMax:"), html.Td(day.FeelsLikeMax),
			html.Td("FeelsLikeMin:"), html.Td(day.FeelsLikeMin),
			html.Td("FeelsLike:"), html.Td(day.FeelsLike),
		),
	).HTML()
}

func (day Day) Precipitation() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Precip: %gmm, PrecipProb: %s, PrecipCover: %s\n", day.Precip, day.PrecipProb, day.PrecipCover)
	_, precipHours := day.PrecipHours()
	if len(precipHours) == 0 {
		fmt.Fprint(&b, "PrecipHours: none")
	} else {
		fmt.Fprint(&b, "PrecipHours: ", strings.Join(precipHours, ", "))
	}
	if len(day.PrecipType) > 0 {
		fmt.Fprint(&b, "\nPrecipType: ", strings.Join(day.PrecipType, ", "))
	}
	return b.String()
}

func (day Day) PrecipitationHTML(highlight ...int) html.HTML {
	div := html.Div().Style("display:grid")
	div.AppendChild(
		html.Span().
			Contentf("Precip: %gmm, PrecipProb: %s, PrecipCover: %s", day.Precip, day.PrecipProb, day.PrecipCover))
	if hours, precipHours := day.PrecipHours(); len(precipHours) == 0 {
		div.AppendChild(html.Span().Content("PrecipHours: none"))
	} else {
		span := html.Span().Content("PrecipHours: ")
		for i, hour := range hours {
			if slices.Contains(highlight, hour) {
				span.AppendChild(html.Span().Style("color:red").Content(precipHours[i]))
			} else {
				span.AppendContent(precipHours[i])
			}
			if i < len(hours)-1 {
				span.AppendContent(", ")
			}
		}
		div.AppendChild(span)
	}
	if len(day.PrecipType) > 0 {
		div.AppendChild(html.Span().Contentf("PrecipType: %s", strings.Join(day.PrecipType, ", ")))
	}
	return div.HTML()
}

func (day Day) PrecipHours() (hours []int, output []string) {
	for _, i := range day.Hours {
		if i.Precip > 0 {
			hour := i.TimeEpoch.Time().Hour()
			hours = append(hours, hour)
			output = append(output, fmt.Sprintf("%02d(%gmm,%s)", hour, i.Precip, i.PrecipProb))
		}
	}
	return
}

func (day Day) String() string {
	var b strings.Builder
	fmt.Fprintln(&b, day.DateInfo(true))
	fmt.Fprintln(&b, day.Temperature())
	fmt.Fprintf(&b, "Humidity: %s, Dew Point: %s, Pressure: %ghPa\n", day.Humidity, day.Dew, day.Pressure)
	fmt.Fprintf(&b, "Precip: %gmm, PrecipProb: %s, PrecipCover: %s\n", day.Precip, day.PrecipProb, day.PrecipCover)
	fmt.Fprintf(&b, "WindGust: %s, WindSpeed: %s, WindDir: %s\n", day.WindGust, day.WindSpeed, day.WindDir)
	fmt.Fprintf(&b, "CloudCover: %s, Visibility: %gkm, UVIndex: %s", day.CloudCover, day.Visibility, day.UVIndex)
	return b.String()
}

func (day Day) HTML() html.HTML {
	return html.Div().AppendChild(
		html.Span().Style("display:list-item;margin-left:15px;list-style-type:disclosure-open").Content(day.DateInfoHTML()),
		html.Table().AppendChild(
			html.Tbody().AppendContent(
				day.TemperatureHTML(),
				html.Tr(
					html.Td("Humidity:"), html.Td(day.Humidity),
					html.Td("Dew Point:"), html.Td(day.Dew),
					html.Td("Pressure:"), html.Td(fmt.Sprintf("%ghPa", day.Pressure)),
				),
				html.Tr(
					html.Td("Precip:"), html.Td(fmt.Sprintf("%gmm", day.Precip)),
					html.Td("PrecipProb:"), html.Td(day.PrecipProb),
					html.Td("PrecipCover:"), html.Td(day.PrecipCover),
				),
				html.Tr(
					html.Td("WindGust:"), html.Td(day.WindGust),
					html.Td("WindSpeed:"), html.Td(day.WindSpeed),
					html.Td("WindDir:"), html.Td(day.WindDir.String()),
				),
				html.Tr(
					html.Td("CloudCover:"), html.Td(day.CloudCover),
					html.Td("Visibility:"), html.Td(fmt.Sprintf("%gkm", day.Visibility)),
					html.Td("UVIndex:"), html.Td(day.UVIndex),
				))),
	).HTML()
}

type Hour struct {
	Time           string           `json:"time,omitempty"`
	TimeEpoch      unit.UnixTime    `json:"timeEpoch,omitempty"`
	Temp           unit.Temperature `json:"temp"`
	FeelsLike      unit.Temperature `json:"feelslike"`
	Humidity       unit.Percent     `json:"humidity"`
	Dew            unit.Temperature `json:"dew"`
	Precip         float64          `json:"precip,omitempty"`
	PrecipProb     unit.Percent     `json:"precipprob,omitempty"`
	Snow           float64          `json:"snow,omitempty"`
	SnowDepth      float64          `json:"snowdepth,omitempty"`
	PrecipType     []string         `json:"preciptype,omitempty"`
	WindGust       wind.Speed       `json:"windgust,omitempty"`
	WindSpeed      wind.Speed       `json:"windspeed,omitempty"`
	WindDir        wind.Degree      `json:"winddir,omitempty"`
	Pressure       float64          `json:"pressure,omitempty"`
	Visibility     float64          `json:"visibility,omitempty"`
	CloudCover     unit.Percent     `json:"cloudcover"`
	SolarRadiation float64          `json:"solarradiation,omitempty"`
	SolarEnergy    float64          `json:"solarenergy,omitempty"`
	UVIndex        unit.UVIndex     `json:"uvindex,omitempty"`
	SevereRisk     float64          `json:"severerisk,omitempty"`
	Condition      Condition        `json:"condition,omitempty"`
	Icon           string           `json:"icon,omitempty"`
}

func (hour Hour) String() string {
	var format struct {
		Time       string       `json:"time"`
		Temp       string       `json:"temp"`
		FeelsLike  string       `json:"feelslike"`
		Humidity   unit.Percent `json:"humidity"`
		Dew        string       `json:"dew"`
		Precip     float64      `json:"precip"`
		PrecipProb unit.Percent `json:"precipprob"`
		WindGust   string       `json:"windgust"`
		WindSpeed  string       `json:"windspeed"`
		WindDir    float64      `json:"winddir"`
		Pressure   float64      `json:"pressure"`
		Visibility float64      `json:"visibility"`
		CloudCover unit.Percent `json:"cloudcover"`
		UVIndex    unit.UVIndex `json:"uvindex"`
		SevereRisk float64      `json:"severerisk"`
		Condition  Condition    `json:"condition"`
	}
	b, _ := json.Marshal(hour)
	json.Unmarshal(b, &format)
	b, _ = json.Marshal(format)
	return string(b)
}

type Condition string

func (s Condition) Short() string {
	condition := strings.Split(string(s), ",")[0]
	switch condition {
	case "Partially cloudy":
		return "Partly"
	default:
		return condition
	}
}

func (s Condition) Img(icon string) html.HTML {
	return html.Img().Style("height:20px").Src(icon).Title(string(s)).HTML()
}
