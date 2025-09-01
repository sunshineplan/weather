package zoomearth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/domstorage"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/chromedp"
	"github.com/sunshineplan/chrome"
	"github.com/sunshineplan/weather/maps"
	"github.com/sunshineplan/weather/option"
	"github.com/sunshineplan/weather/unit"
	"github.com/sunshineplan/weather/unit/coordinates"
)

var (
	DefaultTimeout    = 2 * time.Minute
	DefaultColorDepth = 5000
)

var mapPath = map[maps.MapType]string{
	maps.Satellite:     "satellite",
	maps.Radar:         "radar",
	maps.Precipitation: "precipitation",
	maps.Wind:          "wind-speed",
	maps.Temperature:   "temperature",
	maps.Humidity:      "humidity",
	maps.DewPoint:      "dew-point",
	maps.Pressure:      "pressure",
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

func URL(path string, t time.Time, coords coordinates.Coordinates, zoom float64, overlays []string) string {
	if path == "" {
		path = mapPath[maps.Satellite]
	}
	var date string
	if !t.IsZero() {
		date = "/date=" + t.UTC().Format("2006-01-02,15:04")
	}
	if zoom == 0 {
		zoom = defaultMapOptions.zoom
	}
	url := fmt.Sprintf(
		"%s/maps/%s/#view=%g,%g,%sz%s", root, path, coords.Latitude(), coords.Longitude(), unit.FormatFloat64(zoom, 2), date,
	)
	if len(overlays) > 0 {
		url += "/overlays=" + strings.Join(overlays, ",")
	}
	return url
}

func MapWithContext(ctx context.Context, path string, dt time.Time, coords coordinates.Coordinates, opt *MapOptions) (t time.Time, img image.Image, err error) {
	if path == "" {
		path = mapPath[maps.Satellite]
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
	nav, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()
	if err = chrome.EnableFetch(nav, func(ev *fetch.EventRequestPaused) bool {
		return !strings.Contains(ev.Request.URL, "adsbygoogle")
	}); err != nil {
		return
	}
	u, _ := url.Parse(root)
	chrome.SetCookies(nav, u, []*http.Cookie{{Name: "ze_language", Value: "en"}})
	if err = chromedp.Run(nav,
		chromedp.EmulateViewport(int64(o.width), int64(o.height)),
		chromedp.Navigate(root+"/assets/images/icon-100.jpg"),
	); err != nil {
		return
	}
	storageID := &domstorage.StorageID{StorageKey: domstorage.SerializedStorageKey(root + "/"), IsLocalStorage: true}
	for k, v := range map[string]string{
		"ze_distanceUnit": "metric",
		"ze_introsLayer":  "satellite",
		"ze_timeControl":  "timeline",
		"ze_timeFormat":   "hour24",
		"ze_timeZone":     "utc",
		"ze_welcome":      "false",
	} {
		if err = chrome.SetStorageItem(nav, storageID, k, v); err != nil {
			return
		}
	}
	geocolor := chrome.ListenEvent(nav, regexp.MustCompile(`https://tiles.zoom.earth/geocolor/.*\.jpg`), "GET", false)             // satellite
	icon := chrome.ListenEvent(nav, "https://tiles.zoom.earth/times/icon.json", "GET", false)                                      // wind
	windspeed := chrome.ListenEvent(nav, regexp.MustCompile(`https://tiles.zoom.earth/icon/v1/wind-speed/.*\.webp`), "GET", false) // wind
	var wg sync.WaitGroup
	wg.Go(func() { <-geocolor })
	wg.Go(func() { <-icon })
	wg.Go(func() { <-windspeed })
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	go chromedp.Run(nav, chromedp.Navigate(URL(path, dt, coords, o.zoom, o.overlays)))
	select {
	case <-done:
	case <-nav.Done():
		err = nav.Err()
		return
	}
	//rainviewerCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	//defer cancel()
	//rainviewer := chrome.ListenEvent(rainviewerCtx, regexp.MustCompile(`https://tilecache.rainviewer.com/.*\.webp`), "GET", false)
	//go chromedp.Run(rainviewerCtx, chromedp.ActionFunc(func(ctx context.Context) error {
	//	return input.DispatchKeyEvent(input.KeyDown).
	//		WithKey("r").
	//		WithCode("KeyR").
	//		WithWindowsVirtualKeyCode(82).
	//		Do(ctx)
	//}))
	//select {
	//case <-rainviewer:
	//case <-rainviewerCtx.Done():
	//}
	if err = chromedp.Run(ctx, chromedp.Evaluate("id=window.setTimeout(' ');for(i=1;i<id;i++)window.clearTimeout(i)", nil)); err != nil {
		return
	}
	click, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	chromedp.Run(click, chromedp.Click(".welcome .continue", chromedp.NodeVisible))
	var utcTime string
	if err = chromedp.Run(
		nav,
		chromedp.EvaluateAsDevTools(`
$$('nav.panel').forEach(i=>i.remove())
$$('.group').forEach(i=>i.remove())
$$('button').forEach(i=>i.remove())
$$('.notifications').forEach(i=>i.remove())
$$('.app-link').forEach(i=>i.remove())
$$('.scroll').forEach(i=>i.remove())
$$('.time-indicator').forEach(i=>i.remove())
$$('.hud').forEach(i=>i.remove())
$('.timeline').style.top='calc(6px + env(safe-area-inset-top))'
$('.timeline').style.left='calc(50px + env(safe-area-inset-left))'
$('.timeline').style.right='calc(50px + env(safe-area-inset-right))'
$('.timeline').style.height='36px'
$('.timeline').style.width='150px'
$('.timeline').style.margin='0 auto'`, nil),
		chromedp.Text("div.time-tooltip", &utcTime),
	); err != nil {
		return
	}
	var parseErr error
	if t, parseErr = time.Parse("Monday _2 January, 15:04MST", utcTime); parseErr != nil {
		t, parseErr = time.Parse("Mon _2 Jan, 15:04MST", utcTime)
	}
	if parseErr, _ = maps.ParseTimeError(parseErr); parseErr == nil {
		t = t.AddDate(time.Now().UTC().Year(), 0, 0).In(o.timezone)
		if err = chromedp.Run(nav, chromedp.EvaluateAsDevTools(
			fmt.Sprintf("$('.time-tooltip>.text').innerText='%s'", t.Format("Jan _2, 15:04")), nil)); err != nil {
			return
		}
	}
	for i := range 3 {
		if i == 0 {
			time.Sleep(3 * time.Second)
		} else {
			time.Sleep(10 * time.Second)
		}
		var b []byte
		if err = chromedp.Run(nav, chromedp.FullScreenshot(&b, 100)); err != nil {
			return
		}
		img, err = png.Decode(bytes.NewReader(b))
		if err != nil {
			return
		}
		if depth := colors(img); depth >= DefaultColorDepth {
			err = parseErr
			return
		} else {
			err = maps.InsufficientColor(depth)
		}
	}
	err = errors.Join(parseErr, err)
	return
}

func Map(path string, dt time.Time, coords coordinates.Coordinates, opt *MapOptions) (t time.Time, img image.Image, err error) {
	c := chrome.Headless()
	defer c.Close()
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("failed to get map: %v", e)
		}
	}()
	return MapWithContext(c, path, dt, coords, opt)
}

func Realtime(path string, coords coordinates.Coordinates, opt *MapOptions) (time.Time, image.Image, error) {
	return Map(path, time.Time{}, coords, opt)
}

func colors(img image.Image) int {
	m := make(map[color.Color]struct{})
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			c := img.At(x, y)
			m[c] = struct{}{}
		}
	}
	return len(m)
}
