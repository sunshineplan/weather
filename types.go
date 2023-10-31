package weather

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"strings"
	"time"

	"github.com/sunshineplan/weather/unit"
)

type Current struct {
	Datetime      string           `json:"datetime,omitempty"`
	DatetimeEpoch int64            `json:"datetimeEpoch,omitempty"`
	Temp          unit.Temperature `json:"temp"`
	FeelsLike     unit.Temperature `json:"feelslike"`
	Humidity      Percent          `json:"humidity,omitempty"`
	Dew           unit.Temperature `json:"dew"`
	Precip        float64          `json:"precip,omitempty"`
	PrecipType    []string         `json:"preciptype,omitempty"`
	WindGust      unit.WindSpeed   `json:"windgust,omitempty"`
	WindSpeed     unit.WindSpeed   `json:"windspeed,omitempty"`
	WindDegree    float64          `json:"winddegree,omitempty"`
	WindDir       string           `json:"winddir,omitempty"`
	Pressure      float64          `json:"pressure,omitempty"`
	Visibility    float64          `json:"visibility,omitempty"`
	CloudCover    Percent          `json:"cloudcover"`
	UVIndex       unit.UVIndex     `json:"uvindex,omitempty"`
	Condition     Condition        `json:"condition,omitempty"`
	Icon          string           `json:"icon,omitempty"`
}

type Day struct {
	Date         string           `json:"date,omitempty"`
	DateEpoch    int64            `json:"dateEpoch,omitempty"`
	TempMax      unit.Temperature `json:"tempmax"`
	TempMin      unit.Temperature `json:"tempmin"`
	Temp         unit.Temperature `json:"temp"`
	FeelsLikeMax unit.Temperature `json:"feelslikemax"`
	FeelsLikeMin unit.Temperature `json:"feelslikemin"`
	FeelsLike    unit.Temperature `json:"feelslike"`
	Humidity     Percent          `json:"humidity,omitempty"`
	Dew          unit.Temperature `json:"dew"`
	Precip       float64          `json:"precip,omitempty"`
	PrecipProb   Percent          `json:"precipprob,omitempty"`
	PrecipCover  Percent          `json:"precipcover,omitempty"`
	Snow         float64          `json:"snow,omitempty"`
	SnowDepth    float64          `json:"snowdepth,omitempty"`
	PrecipType   []string         `json:"preciptype,omitempty"`
	WindGust     unit.WindSpeed   `json:"windgust,omitempty"`
	WindSpeed    unit.WindSpeed   `json:"windspeed,omitempty"`
	WindDir      float64          `json:"winddir,omitempty"`
	Pressure     float64          `json:"pressure,omitempty"`
	CloudCover   Percent          `json:"cloudcover"`
	Visibility   float64          `json:"visibility,omitempty"`
	UVIndex      unit.UVIndex     `json:"uvindex,omitempty"`
	SevereRisk   float64          `json:"severerisk,omitempty"`
	Condition    Condition        `json:"condition,omitempty"`
	Description  string           `json:"description,omitempty"`
	Icon         string           `json:"icon,omitempty"`
	Hours        []Hour           `json:"hours,omitempty"`
}

func (day Day) Time() time.Time {
	return time.Unix(day.DateEpoch, 0)
}

func (day Day) Until() time.Duration {
	return time.Until(day.Time().AddDate(0, 0, 1)).Truncate(24 * time.Hour)
}

func (w Day) Before(date time.Time) bool {
	year, month, day := date.Date()
	y, m, d := w.Time().Date()
	if year == y {
		if month == m {
			return day > d
		}
		return month > m
	}
	return year > y
}

func (day Day) IsExpired() bool {
	if day.Date != "" {
		return day.Before(time.Now())
	}
	return true
}

func (day Day) Weekday() string {
	return day.Time().Weekday().String()[:3]
}

func (day Day) DateInfo(condition bool) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Date: %s %s", day.Date, day.Weekday())
	if condition {
		fmt.Fprintf(&b, " (%s)", day.Condition)
	}
	return b.String()
}

func (day Day) DateInfoHTML() string {
	return fmt.Sprintf("%s %s %s", day.Date, day.Weekday(), day.Condition.Img(day.Icon))
}

func (day Day) Temperature() string {
	var b strings.Builder
	fmt.Fprintf(&b, "TempMax: %s, TempMin: %s, Temp: %s\n", day.TempMax, day.TempMin, day.Temp)
	fmt.Fprintf(&b, "FeelsLikeMax: %s, FeelsLikeMin: %s, FeelsLike: %s", day.FeelsLikeMax, day.FeelsLikeMin, day.FeelsLike)
	return b.String()
}

func (day Day) TemperatureHTML() string {
	var b strings.Builder
	fmt.Fprintf(&b, "<tr><td>TempMax:</td><td>%s</td>", day.TempMax)
	fmt.Fprintf(&b, "<td>TempMin:</td><td>%s</td>", day.TempMin)
	fmt.Fprintf(&b, "<td>Temp:</td><td>%s</td></tr>", day.Temp)
	fmt.Fprintf(&b, "<tr><td>FeelsLikeMax:</td><td>%s</td>", day.FeelsLikeMax)
	fmt.Fprintf(&b, "<td>FeelsLikeMin:</td><td>%s</td>", day.FeelsLikeMin)
	fmt.Fprintf(&b, "<td>FeelsLike:</td><td>%s</td></tr>", day.FeelsLike)
	return b.String()
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
		fmt.Fprintf(&b, "\nPrecipType: %s", strings.Join(day.PrecipType, ", "))
	}
	return b.String()
}

func (day Day) PrecipitationHTML(highlight ...int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Precip: %gmm, PrecipProb: %s, PrecipCover: %s\n", day.Precip, day.PrecipProb, day.PrecipCover)
	hours, precipHours := day.PrecipHours()
	if len(precipHours) == 0 {
		fmt.Fprint(&b, "PrecipHours: none")
	} else {
		fmt.Fprint(&b, "PrecipHours: ")
		for i, hour := range hours {
			if slices.Contains(highlight, hour) {
				fmt.Fprintf(&b, `<span style="color:red">%s</span>`, precipHours[i])
			} else {
				fmt.Fprint(&b, precipHours[i])
			}
			if i < len(hours)-1 {
				fmt.Fprint(&b, ", ")
			}
		}
	}
	if len(day.PrecipType) > 0 {
		fmt.Fprintf(&b, "\nPrecipType: %s", strings.Join(day.PrecipType, ", "))
	}
	return b.String()
}

func (day Day) PrecipHours() (hours []int, output []string) {
	for _, i := range day.Hours {
		if i.Precip > 0 {
			hours = append(hours, i.Hour())
			output = append(output, fmt.Sprintf("%02d(%gmm,%s)", i.Hour(), i.Precip, i.PrecipProb))
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
	fmt.Fprintf(&b, "WindGust: %s, WindSpeed: %s, WindDir: %g°\n", day.WindGust, day.WindSpeed, day.WindDir)
	fmt.Fprintf(&b, "CloudCover: %s, Visibility: %gkm, UVIndex: %s", day.CloudCover, day.Visibility, day.UVIndex)
	return b.String()
}

func (day Day) HTML() string {
	var b strings.Builder
	fmt.Fprint(&b, `<div style="display:list-item;margin-left:15px;list-style-type:disclosure-open">`, day.DateInfoHTML(), "</div>")
	fmt.Fprint(&b, "<table><tbody>")
	fmt.Fprint(&b, day.TemperatureHTML())
	fmt.Fprintf(&b, "<tr><td>Humidity:</td><td>%s</td>", day.Humidity)
	fmt.Fprintf(&b, "<td>Dew Point:</td><td>%s</td>", day.Dew)
	fmt.Fprintf(&b, "<td>Pressure:</td><td>%ghPa</td></tr>", day.Pressure)
	fmt.Fprintf(&b, "<tr><td>Precip:</td><td>%gmm</td>", day.Precip)
	fmt.Fprintf(&b, "<td>PrecipProb:</td><td>%s</td>", day.PrecipProb)
	fmt.Fprintf(&b, "<td>PrecipCover:</td><td>%s</td></tr>", day.PrecipCover)
	fmt.Fprintf(&b, "<tr><td>WindGust:</td><td>%s</td>", day.WindGust.HTML())
	fmt.Fprintf(&b, "<td>WindSpeed:</td><td>%s</td>", day.WindSpeed.HTML())
	fmt.Fprintf(&b, `<td>WindDir:</td><td>%g°</td></tr>`, day.WindDir)
	fmt.Fprintf(&b, "<tr><td>CloudCover:</td><td>%s</td>", day.CloudCover)
	fmt.Fprintf(&b, "<td>Visibility:</td><td>%gkm</td>", day.Visibility)
	fmt.Fprintf(&b, "<td>UVIndex:</td><td>%s</td></tr>", day.UVIndex.HTML())
	fmt.Fprint(&b, "</tbody></table>")
	return b.String()
}

type Hour struct {
	Time           string           `json:"time,omitempty"`
	TimeEpoch      int64            `json:"timeEpoch,omitempty"`
	Temp           unit.Temperature `json:"temp"`
	FeelsLike      unit.Temperature `json:"feelslike"`
	Humidity       Percent          `json:"humidity"`
	Dew            unit.Temperature `json:"dew"`
	Precip         float64          `json:"precip,omitempty"`
	PrecipProb     Percent          `json:"precipprob,omitempty"`
	Snow           float64          `json:"snow,omitempty"`
	SnowDepth      float64          `json:"snowdepth,omitempty"`
	PrecipType     []string         `json:"preciptype,omitempty"`
	WindGust       unit.WindSpeed   `json:"windgust,omitempty"`
	WindSpeed      unit.WindSpeed   `json:"windspeed,omitempty"`
	WindDir        float64          `json:"winddir,omitempty"`
	Pressure       float64          `json:"pressure,omitempty"`
	Visibility     float64          `json:"visibility,omitempty"`
	CloudCover     Percent          `json:"cloudcover"`
	SolarRadiation float64          `json:"solarradiation,omitempty"`
	SolarEnergy    float64          `json:"solarenergy,omitempty"`
	UVIndex        unit.UVIndex     `json:"uvindex,omitempty"`
	SevereRisk     float64          `json:"severerisk,omitempty"`
	Condition      Condition        `json:"condition,omitempty"`
	Icon           string           `json:"icon,omitempty"`
}

func (hour Hour) Hour() int {
	return time.Unix(hour.TimeEpoch, 0).Hour()
}

func (hour Hour) String() string {
	var format struct {
		Time       string       `json:"time"`
		Temp       string       `json:"temp"`
		FeelsLike  string       `json:"feelslike"`
		Humidity   Percent      `json:"humidity"`
		Dew        string       `json:"dew"`
		Precip     float64      `json:"precip"`
		PrecipProb Percent      `json:"precipprob"`
		WindGust   string       `json:"windgust"`
		WindSpeed  string       `json:"windspeed"`
		WindDir    float64      `json:"winddir"`
		Pressure   float64      `json:"pressure"`
		Visibility float64      `json:"visibility"`
		CloudCover Percent      `json:"cloudcover"`
		UVIndex    unit.UVIndex `json:"uvindex"`
		SevereRisk float64      `json:"severerisk"`
		Condition  Condition    `json:"condition"`
	}
	b, _ := json.Marshal(hour)
	json.Unmarshal(b, &format)
	b, _ = json.Marshal(format)
	return string(b)
}

type Percent float64

func (f Percent) Max(i Percent) Percent {
	return Percent(math.Max(float64(f), float64(i)))
}

func (f Percent) String() string {
	return fmt.Sprintf("%g%%", f)
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

func (s Condition) Img(icon string) string {
	return fmt.Sprintf(`<img style="height:20px" src=%q title=%q></img>`, icon, s)
}
