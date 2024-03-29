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
	ctx, cancel = context.WithTimeout(c, 5*time.Second)
	defer cancel()
	if err := chromedp.Run(ctx, chromedp.Click(".intro-app>button", chromedp.NodeVisible)); err == nil {
		if err := chromedp.Run(
			ctx,
			chromedp.EvaluateAsDevTools(`
$('.app-link.header').style.display='none'
$('.timeline .play').style.display='none'
$('.timeline .latest').style.display='none'
$('.scroll').style.display='none'
$('.time-indicator').style.display='none'
$('.timeline').style.top='calc(6px + env(safe-area-inset-top))'
$('.timeline').style.left='calc(50px + env(safe-area-inset-left))'
$('.timeline').style.right='calc(50px + env(safe-area-inset-right))'
$('.timeline').style.height='36px'
$('.timeline').style.width='150px'
$('.timeline').style.margin='0 auto'
$('span.day').innerText=new Date().toLocaleDateString('en-US',{month:'short',day:'numeric',})`, nil),
		); err != nil {
			panic(err)
		}
	} else {
		ctx, cancel = context.WithTimeout(c, 3*time.Second)
		defer cancel()
		if err := chromedp.Run(
			ctx,
			chromedp.EvaluateAsDevTools(`
$$('.up,.down').forEach(e=>e.remove())
$$('div .text').forEach(e=>{e.style.top='18px'})
$('.clock .play').style.display='none'
$('.clock .latest').style.display='none'
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
		); err != nil {
			panic(err)
		}
	}
	ctx, cancel = context.WithTimeout(c, time.Minute)
	defer cancel()
	if err = chromedp.Run(
		ctx,
		chromedp.EvaluateAsDevTools(`
$('button.title').style.display='none'
$('button.search').style.display='none'
$('.geolocation').style.display='none'
$('.group.overlays').style.display='none'
$('button.layers').style.display='none'
$('.notifications').style.display='none'`, nil),
		chromedp.Sleep(time.Second*2),
		chromedp.FullScreenshot(&b, quality),
	); err != nil {
		panic(err)
	}
	if len(b) <= 30*1024 {
		panic("bad screenshot")
	}
	return
}

func willAffect(storm storm.Data, coords *coords, radius float64) (affect, future bool) {
	if !storm.Active {
		return
	}
	for _, i := range storm.Track {
		if coordinates.Distance(i.Coordinates(), coords) <= radius {
			affect = true
			if i.Forecast() {
				future = true
				break
			}
		}
	}
	return
}
