package main

import (
	"context"
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

func (coords Coordinates) screenshot(zoom float64, quality int, retry int) (b []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			svc.Print(e)
			if retry--; retry == 0 {
				err = fmt.Errorf("screenshot failed")
			} else {
				time.Sleep(time.Minute)
				b, err = coords.screenshot(zoom, quality, retry)
			}
		}
	}()
	c := chrome.Headless().AddFlags(chromedp.WindowSize(600, 800))
	defer c.Close()
	ctx, cancel := context.WithTimeout(c, time.Minute)
	defer cancel()
	if err = chrome.EnableFetch(ctx, func(ev *fetch.EventRequestPaused) bool {
		return !strings.Contains(ev.Request.URL, "adsbygoogle")
	}); err != nil {
		panic(err)
	}
	notify := chrome.ListenEvent(ctx, "https://tiles.zoom.earth/times/geocolor.json", "GET", false)
	if err = chromedp.Run(ctx, chromedp.Navigate(coords.url(zoom))); err != nil {
		panic(err)
	}
	done := make(chan struct{})
	go func() {
		var n int
		for range notify {
			n++
			if n == 4 {
				close(done)
				return
			}
		}
	}()
	select {
	case <-done:
	case <-ctx.Done():
		panic("screenshot timeout, wait for retry")
	}
	if err = chromedp.Run(
		ctx,
		chromedp.EvaluateAsDevTools(`
$('button.title').style.display='none'
$('button.search').style.display='none'
$('.geolocation').style.display='none'
$('.group.overlays').style.display='none'
$('.group.layers').style.display='none'
$('.notifications').style.display='none'

$$('.up,.down').forEach(e=>e.remove())
$$('div .text').forEach(e=>{e.style.top='18px'})
$('.play').style.display='none'
$('.latest').style.display='none'
$('.clock').style.top='22px'
$('.clock').style.height='50px'
$('.clock').style.width='180px'
$('.clock').style.marginLeft='-90px'
$('div.date').style.left='16px'
$('.hour').style.left='70px'
$('.colon').style.left='108px'
$('.colon').style.animation='none'
$('.minute').style.left='110px'
$('.am-pm').style.left='146px'`, nil),
		chromedp.Sleep(time.Second),
		chromedp.FullScreenshot(&b, quality),
	); err != nil {
		panic(err)
	}
	if len(b) <= 30*1024 {
		panic("bad screenshot")
	}
	return
}

func willAffect(storm storm.Data, coords Coordinates, radius float64) (affect, future bool) {
	if !storm.Active {
		return
	}
	for _, i := range storm.Track {
		if Coordinates(i.Coordinates()).inArea(coords, radius) {
			affect = true
			if i.Forecast() {
				future = true
				break
			}
		}
	}
	return
}
