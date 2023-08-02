package storm

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Storm string

func GetStorms(t time.Time) ([]Storm, error) {
	t = t.UTC().Truncate(6 * time.Hour)
	resp, err := http.Get("https://zoom.earth/data/storms/?date=" + t.Format("2006-01-02T15:04Z"))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var res struct {
		Storms []Storm
		Error  string
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	if err := res.Error; err != "" {
		return nil, errors.New(err)
	}
	return res.Storms, nil
}

func (storm Storm) URL() string {
	return fmt.Sprintf("https://zoom.earth/storms/%s/", storm)
}

type Data struct {
	ID     Storm
	Name   string
	Title  string
	Active bool
	Type   string
	Cone   [][2]float64
	Track  []Track
	JA     int

	Coordinates [2]float64
}

type Track struct {
	Date        time.Time
	Coordinates [2]float64
	Forecast    bool
}

func calcCoordinates(a, b Track, t time.Time) [2]float64 {
	start, end, n := a.Date.Unix(), b.Date.Unix(), t.Unix()
	if n < start {
		return a.Coordinates
	} else if n > end {
		return b.Coordinates
	}
	rate := float64(n-start) / float64(end-start)
	return [2]float64{
		a.Coordinates[0] + (b.Coordinates[0]-a.Coordinates[0])*rate,
		a.Coordinates[1] + (b.Coordinates[1]-a.Coordinates[1])*rate,
	}
}

func (data *Data) calcCoordinates() {
	var a, b Track
	for _, i := range data.Track {
		if !i.Forecast {
			a = i
		} else {
			b = i
			break
		}
	}
	if b.Coordinates == [2]float64{0, 0} {
		data.Coordinates = a.Coordinates
	} else {
		data.Coordinates = calcCoordinates(a, b, time.Now())
	}
}

func (s Storm) Data() (data Data, err error) {
	resp, err := http.Get(fmt.Sprint("https://zoom.earth/data/storms/?id=", s))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return
	}
	data.calcCoordinates()
	return
}
