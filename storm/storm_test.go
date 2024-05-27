package storm

import (
	"testing"
	"time"

	"github.com/sunshineplan/weather/unit/coordinates"
)

var _ Track = track{}

type track struct {
	t time.Time
	c coordinates.Coordinates
}

func (t track) Date() time.Time                      { return t.t }
func (t track) Coordinates() coordinates.Coordinates { return t.c }
func (t track) Forecast() bool                       { return false }

func TestTrack(t *testing.T) {
	data := Data{
		Track: []Track{
			track{time.Date(2000, 1, 1, 12, 0, 0, 0, time.Local), coordinates.New(2, 2)},
			track{time.Date(2000, 1, 1, 13, 0, 0, 0, time.Local), coordinates.New(4, 4)},
		},
	}
	for i, testcase := range []struct {
		t      time.Time
		expect coordinates.Coordinates
	}{
		{time.Date(2000, 1, 1, 11, 0, 0, 0, time.Local), nil},
		{time.Date(2000, 1, 1, 12, 0, 0, 0, time.Local), coordinates.New(2, 2)},
		{time.Date(2000, 1, 1, 12, 30, 0, 0, time.Local), coordinates.New(3, 3)},
		{time.Date(2000, 1, 1, 13, 0, 0, 0, time.Local), coordinates.New(4, 4)},
		{time.Date(2000, 1, 1, 13, 30, 0, 0, time.Local), nil},
	} {
		if res := data.Coordinates(testcase.t); res != testcase.expect {
			t.Errorf("%d expected %q; got %q", i, testcase.expect, res)
		}
	}
}
