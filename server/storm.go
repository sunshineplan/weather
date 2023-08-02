package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	"github.com/sunshineplan/chrome"
	"github.com/sunshineplan/weather/storm"
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

func (coords Coordinates) screenshot(zoom float64, quality int, clock bool, retry int) (b []byte, err error) {
	c := chrome.Headless().AddFlags(chromedp.WindowSize(600, 800))
	defer c.Close()
	if err = c.EnableFetch(func(ev *fetch.EventRequestPaused) bool {
		return !strings.Contains(ev.Request.URL, "adsbygoogle")
	}); err != nil {
		return
	}
	done := c.ListenEvent(chrome.URLContains("notifications"), "GET", false)
	if err = c.Run(chromedp.Navigate(coords.url(zoom))); err != nil {
		return
	}
	select {
	case <-done:
	case <-time.After(time.Minute):
		if retry = retry - 1; retry == 0 {
			return nil, fmt.Errorf("timeout")
		}
		time.Sleep(3 * time.Minute)
		return coords.screenshot(zoom, quality, clock, retry)
	}
	if !clock {
		if err = c.Run(chromedp.EvaluateAsDevTools("$('div.panel.clock').style.display='none'", nil)); err != nil {
			return
		}
	}
	err = c.Run(
		//chromedp.EvaluateAsDevTools("$('nav.panel.layers').style.display='none'", nil),
		chromedp.EvaluateAsDevTools("$('div.layers').style.display='none'", nil),
		chromedp.EvaluateAsDevTools("$('aside.notifications').style.display='none'", nil),
		chromedp.Sleep(time.Second),
		chromedp.FullScreenshot(&b, quality),
	)
	return
}

func willAffect(storm storm.Data, coords Coordinates, radius float64) bool {
	if !storm.Active || storm.Cone == nil {
		return false
	}
	for _, i := range storm.Track {
		if Coordinates(i.Coordinates).inArea(coords, radius) {
			return true
		}
	}
	return false
}
