package zoomearth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sunshineplan/weather/storm"
	"github.com/sunshineplan/weather/unit/coordinates"
)

var _ storm.Storm = StormID("")

func (id StormID) URL() string {
	return fmt.Sprintf("%s/storms/%s/", root, id)
}

type stormData struct {
	ID     StormID
	Season string
	Name   string
	Title  string
	Active bool
	Type   string
	Place  string
	Cone   []coordinates.LongLat
	Track  []track
	JA     int
}

var _ storm.Track = track{}

type track struct {
	D time.Time           `json:"date"`
	C coordinates.LongLat `json:"coordinates"`
	F bool                `json:"forecast"`
}

func (t track) Date() time.Time                      { return t.D }
func (t track) Coordinates() coordinates.Coordinates { return t.C }
func (t track) Forecast() bool                       { return t.F }

func (id StormID) Data() (storm.Data, error) {
	resp, err := http.Get(fmt.Sprint(root, "/data/storms/?id=", id))
	if err != nil {
		return storm.Data{}, err
	}
	defer resp.Body.Close()
	var data stormData
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return storm.Data{}, err
	}
	return data.Convert(), nil
}

func GetStorms(t time.Time) (storms []storm.Storm, err error) {
	t = t.UTC().Truncate(6 * time.Hour)
	resp, err := http.Get(root + "/data/storms/?date=" + t.Format("2006-01-02T15:04Z"))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var res struct {
		Storms []StormID
		Error  string
	}
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return
	}
	if err := res.Error; err != "" {
		return nil, errors.New(err)
	}
	for _, i := range res.Storms {
		storms = append(storms, i)
	}
	return
}
