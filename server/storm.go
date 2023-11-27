package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	"github.com/sunshineplan/chrome"
	"github.com/sunshineplan/weather/storm"
	"github.com/sunshineplan/weather/unit/coordinates"
)

type coords struct{ coordinates.Coordinates }

func (coords coords) inArea(c coords, radius float64) bool {
	return coordinates.Distance(coords, c) <= radius
}

func (c coords) offset(x, y float64) coords {
	return coords{coordinates.New(float64(c.Latitude())+x, float64(c.Longitude())+y)}
}

func (coords coords) url(zoom float64) string {
	return fmt.Sprintf(
		"https://zoom.earth/maps/satellite/#view=%g,%g,%.2fz/overlays=radar,wind", coords.Latitude(), coords.Longitude(), zoom,
	)
}

func (coords coords) screenshot(zoom float64, quality int, retry int) (b []byte, err error) {
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

func willAffect(storm storm.Data, coordinates coords, radius float64) (affect, future bool) {
	if !storm.Active {
		return
	}
	for _, i := range storm.Track {
		if (coords{i.Coordinates()}).inArea(coordinates, radius) {
			affect = true
			if i.Forecast() {
				future = true
				break
			}
		}
	}
	return
}
