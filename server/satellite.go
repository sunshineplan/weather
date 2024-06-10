package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/kettek/apng"
	"github.com/sunshineplan/weather/api/zoomearth"
	"github.com/sunshineplan/weather/maps"
	"github.com/sunshineplan/weather/storm"
	"github.com/sunshineplan/weather/unit/coordinates"
)

var (
	format      = "200601021504"
	shortFormat = "01021504"
	width       = 600
	height      = 800
	timezone    = time.FixedZone("CST", 8*60*60)

	animationDuration = []time.Duration{
		6 * time.Hour,
		12 * time.Hour,
		18 * time.Hour,
		24 * time.Hour,
	}
	keep = int(slices.Max(animationDuration)/time.Hour) * 6
)

func mapOptions(zoom float64) *zoomearth.MapOptions {
	return zoomearth.NewMapOptions().
		SetSize(width, height).
		SetZoom(zoom).
		SetOverlays([]string{"radar", "wind"}).
		SetTimeZone(timezone)
}

var satelliteMutex sync.Mutex

func satellite(t time.Time, coords coordinates.Coordinates, path, format string, opt any) error {
	satelliteMutex.Lock()
	defer satelliteMutex.Unlock()
	time.Sleep(time.Second)
	t, img, err := mapAPI.Map(maps.Satellite, t, coords, opt)
	if err != nil {
		if errors.Is(err, maps.ErrInsufficientColor) {
			svc.Print(err)
			return nil
		} else {
			return err
		}
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
	file := filepath.Join(path, t.Format(format)+".png")
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
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

var enc = apng.Encoder{CompressionLevel: apng.BestCompression}

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
	var step int
	if d != 0 {
		step = int(math.Logb(float64(d / time.Hour)))
	} else if step = int(math.Round(math.Log(1+float64(len(res))))) - 2; step <= 0 {
		step = 1
	}
	slices.Reverse(res)
	var imgs []string
	for i, img := range res {
		if i%step == 0 {
			imgs = append(imgs, img)
		}
	}
	slices.Reverse(imgs)
	apngImg, n := apng.APNG{}, len(imgs)
	var delay int
	if d != 0 {
		delay = 40
	} else if delay = 6000 / n; delay > 40 {
		delay = 40
	}
	for i, img := range imgs {
		f, err := os.Open(img)
		if err != nil {
			return err
		}
		if img, _, err := image.Decode(f); err != nil {
			svc.Print(err)
		} else {
			if i != n-1 {
				apngImg.Frames = append(apngImg.Frames, apng.Frame{Image: img, DelayNumerator: uint16(delay)})
			} else {
				apngImg.Frames = append(apngImg.Frames, apng.Frame{Image: img, DelayNumerator: 300})
			}
		}
		f.Close()
	}
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	f, err := os.Create(output + ".png")
	if err != nil {
		return err
	}
	defer f.Close()
	return enc.Encode(f, apngImg)
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
		}
	}
	for _, d := range animationDuration {
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
				}
			}
		}
		if err := animation(dir+"/*", dir, 0, shortFormat, false); err != nil {
			svc.Print(err)
		}
	}
}
