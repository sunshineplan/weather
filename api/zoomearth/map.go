package zoomearth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	"github.com/sunshineplan/chrome"
	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/option"
	"github.com/sunshineplan/weather/unit"
	"github.com/sunshineplan/weather/unit/coordinates"
)

var mapPath = map[weather.MapType]string{
	weather.Satellite:     "satellite",
	weather.Radar:         "radar",
	weather.Precipitation: "precipitation",
	weather.Wind:          "wind-speed",
	weather.Temperature:   "temperature",
	weather.Humidity:      "humidity",
	weather.DewPoint:      "dew-point",
	weather.Pressure:      "pressure",
}

var defaultMapOptions = MapOptions{
	zoom:     4,
	quality:  75,
	overlays: []string{"radar", "wind"},
}

var (
	_ option.Zoom     = MapOptions{}
	_ option.Quality  = MapOptions{}
	_ option.Overlays = MapOptions{}
)

func NewMapOptions(zoom float64, quality int, overlays []string) *MapOptions {
	return &MapOptions{zoom, quality, overlays}
}

func (o MapOptions) Zoom() float64 {
	return o.zoom
}

func (o MapOptions) Quality() int {
	return o.quality
}

func (o MapOptions) Overlays() []string {
	return o.overlays
}

func URL(path string, coords coordinates.Coordinates, zoom float64, overlays []string) string {
	if path == "" {
		path = mapPath[weather.Satellite]
	}
	if zoom == 0 {
		zoom = defaultMapOptions.zoom
	}
	url := fmt.Sprintf(
		"https://zoom.earth/maps/%s/#view=%g,%g,%sz", path, coords.Latitude(), coords.Longitude(), unit.FormatFloat64(zoom, 2),
	)
	if len(overlays) > 0 {
		url += "/overlays=" + strings.Join(overlays, ",")
	}
	return url
}

func Realtime(path string, coords coordinates.Coordinates, opt *MapOptions) (b []byte, err error) {
	if path == "" {
		path = mapPath[weather.Satellite]
	}
	o := defaultMapOptions
	if opt != nil {
		o.zoom = opt.zoom
		o.quality = opt.quality
		o.overlays = opt.overlays
	}
	if o.quality < 1 {
		o.quality = 1
	} else if o.quality > 100 {
		o.quality = 100
	}
	c := chrome.Headless().AddFlags(chromedp.WindowSize(600, 800))
	defer c.Close()
	ctx, cancel := context.WithTimeout(c, time.Minute)
	defer cancel()
	if err = chrome.EnableFetch(ctx, func(ev *fetch.EventRequestPaused) bool {
		return !strings.Contains(ev.Request.URL, "adsbygoogle")
	}); err != nil {
		return
	}
	notify := chrome.ListenEvent(ctx, "https://tiles.zoom.earth/times/geocolor.json", "GET", false)
	if err = chromedp.Run(ctx, chromedp.Navigate(URL(path, coords, o.zoom, o.overlays))); err != nil {
		return
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
		err = ctx.Err()
		return
	}
	ctx, cancel = context.WithTimeout(c, 5*time.Second)
	defer cancel()
	if err = chromedp.Run(ctx, chromedp.Click(".welcome .continue", chromedp.NodeVisible)); err == nil {
		if err = chromedp.Run(
			ctx,
			chromedp.EvaluateAsDevTools(`
$('.app-link.header').style.display='none'
$('.timeline .play').style.display='none'
$('.timeline .latest').style.display='none'
$('.scroll').style.display='none'
$('.timeline').style.top='calc(6px + env(safe-area-inset-top))'
$('.timeline').style.left='calc(50px + env(safe-area-inset-left))'
$('.timeline').style.right='calc(50px + env(safe-area-inset-right))'
$('.timeline').style.height='36px'
$('.timeline').style.width='150px'
$('.timeline').style.margin='0 auto'
$('span.day').innerText=new Date().toLocaleDateString('en-US',{month:'short',day:'numeric',})`, nil),
		); err != nil {
			return
		}
	} else {
		ctx, cancel = context.WithTimeout(c, 3*time.Second)
		defer cancel()
		if err = chromedp.Run(
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
			return
		}
	}
	ctx, cancel = context.WithTimeout(c, time.Minute)
	defer cancel()
	err = chromedp.Run(
		ctx,
		chromedp.EvaluateAsDevTools(`
$('button.title').style.display='none'
$('button.search').style.display='none'
$('.geolocation').style.display='none'
$('.group.overlays').style.display='none'
$('button.layers').style.display='none'
$('.notifications').style.display='none'`, nil),
		chromedp.Sleep(time.Second*2),
		chromedp.FullScreenshot(&b, o.quality),
	)
	return
}
