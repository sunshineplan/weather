package zoomearth

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chromedp/cdproto/domstorage"
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
	width:    800,
	height:   600,
	zoom:     4,
	overlays: []string{"radar", "wind"},
	timezone: time.Local,
}

var (
	_ option.Size     = new(MapOptions)
	_ option.Zoom     = new(MapOptions)
	_ option.Overlays = new(MapOptions)
	_ option.TimeZone = new(MapOptions)
)

func NewMapOptions() *MapOptions {
	return new(MapOptions)
}

func (o MapOptions) Size() (int, int)         { return o.width, o.height }
func (o MapOptions) Zoom() float64            { return o.zoom }
func (o MapOptions) Overlays() []string       { return o.overlays }
func (o MapOptions) TimeZone() *time.Location { return o.timezone }

func (o *MapOptions) SetSize(width int, height int) *MapOptions {
	o.width = width
	o.height = height
	return o
}
func (o *MapOptions) SetZoom(zoom float64) *MapOptions {
	o.zoom = zoom
	return o
}
func (o *MapOptions) SetOverlays(overlays []string) *MapOptions {
	o.overlays = overlays
	return o
}
func (o *MapOptions) SetTimeZone(timezone *time.Location) *MapOptions {
	o.timezone = timezone
	return o
}

func URL(path string, coords coordinates.Coordinates, zoom float64, overlays []string) string {
	if path == "" {
		path = mapPath[weather.Satellite]
	}
	if zoom == 0 {
		zoom = defaultMapOptions.zoom
	}
	url := fmt.Sprintf(
		"%s/maps/%s/#view=%g,%g,%sz", root, path, coords.Latitude(), coords.Longitude(), unit.FormatFloat64(zoom, 2),
	)
	if len(overlays) > 0 {
		url += "/overlays=" + strings.Join(overlays, ",")
	}
	return url
}

func Realtime(path string, coords coordinates.Coordinates, opt *MapOptions) (t time.Time, img image.Image, err error) {
	if path == "" {
		path = mapPath[weather.Satellite]
	}
	o := defaultMapOptions
	if opt != nil {
		if opt.width > 0 {
			o.width = opt.width
		}
		if opt.height > 0 {
			o.height = opt.height
		}
		if opt.zoom > 0 {
			o.zoom = opt.zoom
		}
		o.overlays = opt.overlays
		if opt.timezone != nil {
			o.timezone = opt.timezone
		}
	}
	c := chrome.Headless().AddFlags(chromedp.WindowSize(o.width, o.height))
	defer c.Close()
	ctx, cancel := context.WithTimeout(c, time.Minute)
	defer cancel()
	if err = chrome.EnableFetch(ctx, func(ev *fetch.EventRequestPaused) bool {
		return !strings.Contains(ev.Request.URL, "adsbygoogle")
	}); err != nil {
		return
	}
	u, _ := url.Parse(root)
	c.SetCookies(u, []*http.Cookie{{Name: "ze_language", Value: "en"}})
	if err = chromedp.Run(ctx, chromedp.Navigate(root+"/assets/images/icon-100.jpg")); err != nil {
		return
	}
	storageID := &domstorage.StorageID{StorageKey: domstorage.SerializedStorageKey(root + "/"), IsLocalStorage: true}
	for k, v := range map[string]string{
		"ze_timeZone":    "utc",
		"ze_timeFormat":  "hour24",
		"ze_timeControl": "timeline",
		"ze_welcome":     "false",
	} {
		if err = c.SetStorageItem(storageID, k, v); err != nil {
			return
		}
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
	chromedp.Run(ctx,
		chromedp.Click(".welcome .continue", chromedp.NodeVisible),
		chromedp.Evaluate("$('.app-link.header').style.display='none'", nil),
	)
	ctx, cancel = context.WithTimeout(c, time.Minute)
	defer cancel()
	var utcTime string
	if err = chromedp.Run(
		ctx,
		chromedp.EvaluateAsDevTools(`
$('.timeline .play').style.display='none'
$('.timeline .latest').style.display='none'
$('.scroll').style.display='none'
$('.timeline').style.top='calc(6px + env(safe-area-inset-top))'
$('.timeline').style.left='calc(50px + env(safe-area-inset-left))'
$('.timeline').style.right='calc(50px + env(safe-area-inset-right))'
$('.timeline').style.height='36px'
$('.timeline').style.width='150px'
$('.timeline').style.margin='0 auto'
$('button.title').style.display='none'
$('button.search').style.display='none'
$('.geolocation').style.display='none'
$('.group.overlays').style.display='none'
$('button.layers').style.display='none'
$('.notifications').style.display='none'`, nil),
		chromedp.Text("div.time-tooltip", &utcTime),
	); err != nil {
		return
	}
	if t, err = time.Parse("Monday _2 Jan, 15:04MST", utcTime); err != nil {
		if t, err = time.Parse("Mon _2 Jan, 15:04MST", utcTime); err != nil {
			return
		}
	}
	t = t.In(o.timezone)
	var b []byte
	if err = chromedp.Run(
		ctx,
		chromedp.EvaluateAsDevTools(fmt.Sprintf("$('.time-tooltip>.text').innerText='%s'", t.Format("Jan _2, 15:04")), nil),
		chromedp.Sleep(time.Second*2),
		chromedp.FullScreenshot(&b, 100),
	); err != nil {
		return
	}
	img, err = png.Decode(bytes.NewReader(b))
	return
}
