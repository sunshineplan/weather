package main

import (
	"errors"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/kettek/apng"
	"github.com/sunshineplan/weather/api/zoomearth"
	"github.com/sunshineplan/weather/maps"
	"github.com/sunshineplan/weather/storm"
	"github.com/sunshineplan/weather/unit/coordinates"
)

var (
	keep        = 432
	format      = "200601021504"
	shortFormat = "01021504"
	timezone    = time.FixedZone("CST", 8*60*60)
)

func mapOptions(zoom float64) *zoomearth.MapOptions {
	return zoomearth.NewMapOptions().
		SetSize(600, 800).
		SetZoom(zoom).
		SetOverlays([]string{"radar", "wind"}).
		SetTimeZone(timezone)
}

func satellite(t time.Time, coords coordinates.Coordinates, path, format string, opt any) (err error) {
	t, img, err := mapAPI.Map(maps.Satellite, t, coords, opt)
	if err != nil {
		if errors.Is(err, maps.ErrInsufficientColor) {
			svc.Print(err)
		} else {
			return
		}
	}
	if err = os.MkdirAll(path, 0755); err != nil {
		return
	}
	file := filepath.Join(path, t.Format(format)+".png")
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	if err = png.Encode(f, img); err != nil {
		return
	}
	return
}

func getTimes(path, format string) (ts []time.Time) {
	res, err := filepath.Glob(path + "/*.png")
	if err != nil {
		panic(err)
	}
	if len(res) == 0 {
		return
	}
	last, err := time.ParseInLocation(format, strings.TrimSuffix(filepath.Base(res[len(res)-1]), ".png"), timezone)
	if err != nil {
		panic(err)
	}
	for i := time.Duration(1); i <= time.Hour/(10*time.Minute); i++ {
		if t := last.Add(-i * 10 * time.Minute); slices.IndexFunc(res, func(s string) bool {
			return strings.HasSuffix(s, t.Format(format)+".png")
		}) == -1 {
			ts = append(ts, t)
		}
	}
	return
}

func animation(path, output string, d time.Duration, format string, remove bool) error {
	res, err := filepath.Glob(path)
	if err != nil {
		return err
	}
	if remove {
		for ; len(res) > keep; res = res[1:] {
			if err := os.Remove(res[0]); err != nil {
				svc.Print(err)
			}
		}
	}
	if d != 0 {
		now := time.Now().In(timezone)
		res = slices.DeleteFunc(res, func(i string) bool {
			file := filepath.Base(i)
			if index := strings.LastIndex(file, "."); index != -1 {
				file = file[:index]
			}
			t, err := time.ParseInLocation(format, file, timezone)
			if err != nil {
				svc.Print(err)
				return true
			}
			return now.Sub(t) >= d
		})
	}
	slices.Reverse(res)
	var step int
	if d != 0 {
		step = int(math.Logb(float64(d / time.Hour)))
	} else if step = int(math.Round(math.Log(1+float64(len(res))))) - 2; step <= 0 {
		step = 1
	}
	var imgs []image.Image
	for i, name := range res {
		if i%step == 0 {
			f, err := os.Open(name)
			if err != nil {
				return err
			}
			defer f.Close()
			img, _, err := image.Decode(f)
			if err != nil {
				return err
			}
			imgs = append(imgs, img)
		}
	}
	slices.Reverse(imgs)
	gifImg, apngImg, n := new(gif.GIF), apng.APNG{}, len(imgs)
	var delay int
	if d != 0 {
		delay = 40
	} else if delay = 6000 / n; delay > 40 {
		delay = 40
	}
	for i, img := range imgs {
		p := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.Draw(p, p.Rect, img, image.Point{}, draw.Over)
		gifImg.Image = append(gifImg.Image, p)
		if i != n-1 {
			gifImg.Delay = append(gifImg.Delay, delay)
			apngImg.Frames = append(apngImg.Frames, apng.Frame{Image: img, DelayNumerator: uint16(delay)})
		} else {
			gifImg.Delay = append(gifImg.Delay, 300)
			apngImg.Frames = append(apngImg.Frames, apng.Frame{Image: img, DelayNumerator: 300})
		}
	}
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	f, err := os.Create(output + ".gif")
	if err != nil {
		return err
	}
	defer f.Close()
	if err := gif.EncodeAll(f, gifImg); err != nil {
		return err
	}
	f, err = os.Create(output + ".png")
	if err != nil {
		return err
	}
	defer f.Close()
	return apng.Encode(f, apngImg)
}

func updateDaily() {
	svc.Print("Start saving satellite map...")
	zoomMutex.Lock()
	defer zoomMutex.Unlock()
	if err := satellite(time.Time{}, location, "daily", format, mapOptions(*zoom)); err != nil {
		svc.Print(err)
		return
	}
	for _, t := range getTimes("daily", format) {
		if err := satellite(t, location, "daily", format, mapOptions(*zoom)); err != nil {
			svc.Print(err)
			continue
		}
	}
	for _, d := range []time.Duration{72, 48, 24, 12, 6} {
		d = d * time.Hour
		if err := animation("daily/*", "animation/"+strings.TrimSuffix(d.String(), "0m0s"), d, format, true); err != nil {
			svc.Print(err)
		}
	}
}

func updateStorm(storms []storm.Data) {
	for _, i := range storms {
		dir := filepath.Join(*path, i.Season, fmt.Sprintf("%d-%s", i.No, i.ID))
		if err := satellite(time.Time{}, i.Coordinates(time.Now()), dir, shortFormat, mapOptions(*stormZoom)); err != nil {
			svc.Print(err)
			continue
		}
		for _, t := range getTimes(dir, shortFormat) {
			if coords := i.Coordinates(t); coords != nil {
				if err := satellite(t, coords, dir, shortFormat, mapOptions(*stormZoom)); err != nil {
					svc.Print(err)
					continue
				}
			}
		}
		if err := animation(dir+"/*", dir, 0, shortFormat, false); err != nil {
			log.Print(err)
		}
	}
}
