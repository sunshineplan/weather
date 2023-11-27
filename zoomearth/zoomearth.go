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

var _ storm.API = ZoomEarthAPI{}

func (ZoomEarthAPI) GetStorms(t time.Time) (storms []storm.Storm, err error) {
	t = t.UTC().Truncate(6 * time.Hour)
	resp, err := http.Get("https://zoom.earth/data/storms/?date=" + t.Format("2006-01-02T15:04Z"))
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

var _ storm.Storm = StormID("")

func (id StormID) URL() string {
	return fmt.Sprintf("https://zoom.earth/storms/%s/", id)
}

func (id StormID) Data() (storm.Data, error) {
	resp, err := http.Get(fmt.Sprint("https://zoom.earth/data/storms/?id=", id))
	if err != nil {
		return storm.Data{}, err
	}
	defer resp.Body.Close()
	var data StormData
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return storm.Data{}, err
	}
	data.calcCoordinates()
	return data.Convert(), nil
}

func calcCoordinates(a, b Track, t time.Time) coordinates.Coordinates {
	start, end, n := a.Date.Unix(), b.Date.Unix(), t.Unix()
	if n < start {
		return a.Coordinates
	} else if n > end {
		return b.Coordinates
	}
	rate := float64(n-start) / float64(end-start)
	return coordinates.LongLat{
		a.Coordinates[0] + (b.Coordinates[0]-a.Coordinates[0])*rate,
		a.Coordinates[1] + (b.Coordinates[1]-a.Coordinates[1])*rate,
	}
}

func (data *StormData) calcCoordinates() {
	var a, b Track
	for _, i := range data.Track {
		if !i.Forecast {
			a = i
		} else {
			b = i
			break
		}
	}
	if b.Coordinates == [2]float64{} {
		data.Coordinates = a.Coordinates
	} else {
		data.Coordinates = calcCoordinates(a, b, time.Now())
	}
}
