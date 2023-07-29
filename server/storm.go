package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	"github.com/sunshineplan/chrome"
	"github.com/sunshineplan/gohttp"
)

type Coordinates [2]float64

func (coords Coordinates) String() string {
	return fmt.Sprintf("%.1f,%.1f", coords[1], coords[0])
}

func (coords Coordinates) distance(c Coordinates) float64 {
	return math.Sqrt(math.Pow(coords[0]-c[0], 2) + math.Pow(coords[1]-c[1], 2))
}

func (coords Coordinates) inArea(c Coordinates, radius float64) bool {
	return coords.distance(c) <= radius
}

func (coords Coordinates) offset(x, y float64) Coordinates {
	return Coordinates{coords[0] + x, coords[1] + y}
}

func (coords Coordinates) url(zoom float64) string {
	return fmt.Sprintf("https://zoom.earth/maps/satellite/#view=%s,%.2fz/overlays=radar,wind", coords, zoom)
}

func (coords Coordinates) screenshot(zoom float64, quality int) (b []byte, err error) {
	c := chrome.Headless().AddFlags(chromedp.WindowSize(600, 800))
	defer c.Close()
	if err = c.EnableFetch(func(ev *fetch.EventRequestPaused) bool {
		return !strings.Contains(ev.Request.URL, "adsbygoogle")
	}); err != nil {
		return
	}
	done := c.ListenEvent(chrome.URLContains("manifest"), "GET", false)
	if err = c.Run(chromedp.Navigate(coords.url(zoom))); err != nil {
		return
	}
	select {
	case <-done:
	case <-time.After(time.Minute):
		svc.Print("storm screenshot timeout")
	}
	err = c.Run(
		//chromedp.EvaluateAsDevTools("$('nav.panel.layers').style.display='none'", nil),
		chromedp.EvaluateAsDevTools("$('div.panel.clock').style.display='none'", nil),
		chromedp.EvaluateAsDevTools("$('aside.notifications').style.display='none'", nil),
		chromedp.Sleep(time.Second),
		chromedp.FullScreenshot(&b, quality),
	)
	return
}

type Storm string

func getStorms(t time.Time) ([]Storm, error) {
	t = t.UTC().Truncate(6 * time.Hour)
	resp, err := gohttp.Get("https://zoom.earth/data/storms/?date="+t.Format("2006-01-02T15:04Z"), nil)
	if err != nil {
		return nil, err
	}
	var res struct {
		Storms []Storm
		Error  string
	}
	if err := resp.JSON(&res); err != nil {
		return nil, err
	}
	if err := res.Error; err != "" {
		return nil, errors.New(err)
	}
	return res.Storms, nil
}

func (s Storm) data() (*gohttp.Response, error) {
	return gohttp.Get(fmt.Sprint("https://zoom.earth/data/storms/?id=", s), nil)
}

func (s Storm) willAffect(coords Coordinates, radius float64) bool {
	resp, err := s.data()
	if err != nil {
		svc.Print(err)
		return false
	}
	var res struct {
		Active bool
		Cone   []Coordinates
		Track  []struct {
			Coordinates Coordinates
		}
	}
	if err := resp.JSON(&res); err != nil {
		svc.Println(err, resp)
		return false
	}
	if !res.Active || res.Cone == nil {
		return false
	}
	for _, i := range res.Track {
		if i.Coordinates.inArea(coords, radius) {
			return true
		}
	}
	return false
}
