package airmatters

import (
	"cmp"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/sunshineplan/weather/aqi"
	"github.com/sunshineplan/weather/unit"
)

var typeMap = map[aqi.Type]string{
	aqi.Australia:  "aqi_au",
	aqi.Canada:     "aqi_ca",
	aqi.China:      "aqi_cn",
	aqi.Europe:     "caqi_eu",
	aqi.India:      "naqi_in",
	aqi.Netherland: "aqi_nl",
	aqi.UK:         "daqi_uk",
	aqi.US:         "aqi_us",
}

var kindMap = map[string]aqi.Kind{
	"co":   aqi.CO,
	"no2":  aqi.NO2,
	"o3":   aqi.O3,
	"pm25": aqi.PM2Dot5,
	"pm10": aqi.PM10,
	"so2":  aqi.SO2,
}

func unixTime(s string) unit.UnixTime {
	t, _ := time.Parse("2006-01-02 15:04:05", s)
	return unit.UnixTime(t.Unix())
}

var _ aqi.Current = Current{}

func (i Current) Unix() unit.UnixTime {
	return unixTime(i.Time)
}

func (i Current) Date() string {
	return i.Unix().Date()
}

func (i Current) AQI() aqi.AQI {
	for _, item := range i.Items {
		if item.Type == "index" && item.Kind == "aqi" {
			value, _ := strconv.Atoi(item.Value)
			return aqi.NewAQI(i.AQIType, value, aqi.NewLevel(item.Level, item.Color))
		}
	}
	return nil
}

func parseFloat(s string) float64 {
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return v
	}
	if regexp.MustCompile(`\d+~\d+`).MatchString(s) {
		v := strings.Split(s, "~")
		v1, _ := strconv.ParseFloat(v[0], 64)
		v2, _ := strconv.ParseFloat(v[1], 64)
		return (v1 + v2) / 2
	}
	return 0
}

func (i Current) Pollutants() (pollutants []aqi.Pollutant) {
	for _, i := range i.Items {
		if i.Type == "pollutant" {
			if v, ok := kindMap[i.Kind]; ok {
				pollutants = append(pollutants, aqi.NewPollutant(v, i.Unit, parseFloat(i.Value), aqi.NewLevel(i.Level, i.Color)))
			}
		}
	}
	slices.SortFunc(pollutants, func(a, b aqi.Pollutant) int { return cmp.Compare(a.Kind(), b.Kind()) })
	return
}

var _ aqi.Day = Item{}

func (i Item) Unix() unit.UnixTime {
	return unixTime(i.Time)
}

func (i Item) Date() string {
	return i.Unix().Date()
}

func (i Item) AQI() aqi.AQI {
	return aqi.NewAQI(i.AQIType, int(parseFloat(i.Value)), aqi.NewLevel(i.Level, i.Color))
}
