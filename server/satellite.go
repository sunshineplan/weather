package main

import (
	"fmt"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/sunshineplan/utils/pool"
	"github.com/sunshineplan/weather/api/zoomearth"
	"github.com/sunshineplan/weather/maps"
	"github.com/sunshineplan/weather/storm"
	"github.com/sunshineplan/weather/unit/coordinates"
)

const (
	format = "200601021504"
	width  = 600
	height = 800
)

var (
	satelliteMutex sync.Mutex
	pngEncoder     = png.Encoder{
		CompressionLevel: png.BestSpeed,
		BufferPool:       pool.New[png.EncoderBuffer](),
	}
	timezone          = time.FixedZone("CST", 8*60*60)
	animationDuration = []time.Duration{
		6 * time.Hour,
		12 * time.Hour,
		24 * time.Hour,
		36 * time.Hour,
	}
	keep        = int(slices.Max(animationDuration)/time.Hour) * 6
	stormMinute = []int{0, 30}
)

func mapOptions(zoom float64) *zoomearth.MapOptions {
	return zoomearth.NewMapOptions().
		SetSize(width, height).
		SetZoom(zoom).
		SetOverlays([]string{"radar", "wind"}).
		SetTimeZone(timezone)
}

func satellite(t time.Time, coords coordinates.Coordinates, path string, opt any) error {
	satelliteMutex.Lock()
	defer satelliteMutex.Unlock()
	time.Sleep(time.Second)
	t, img, err := mapAPI.Map(maps.Satellite, t, coords, opt)
	if err != nil {
		return err
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
	return pngEncoder.Encode(f, img)
}

func getTimes(path string) (ts []time.Time) {
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
	for i := time.Duration(1); i <= 2*time.Hour/(10*time.Minute); i++ {
		if t := last.Add(-i * 10 * time.Minute); slices.IndexFunc(res, func(s string) bool {
			return strings.HasSuffix(s, t.Format(format)+".png")
		}) == -1 {
			ts = append(ts, t)
		}
	}
	return
}

func getImages(path string, d time.Duration, daily bool) (imgs []string, err error) {
	res, err := filepath.Glob(path)
	if err != nil {
		return
	}
	if daily {
		for ; len(res) > keep; res = res[1:] {
			if err := os.Remove(res[0]); err != nil {
				svc.Print(err)
			}
		}
	} else {
		res = slices.DeleteFunc(res, func(i string) bool {
			file := filepath.Base(i)
			if index := strings.LastIndex(file, "."); index != -1 {
				file = file[:index]
			}
			t, err := time.ParseInLocation(format, file, timezone)
			if err != nil {
				svc.Print(err)
				return false
			}
			if !slices.Contains(stormMinute, t.Minute()) {
				if err := os.Remove(i); err != nil {
					svc.Print(err)
				}
				return true
			}
			return false
		})
	}
	if daily {
		now := time.Now().Truncate(10 * time.Minute).In(timezone)
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
			return now.Sub(t) > d
		})
	}
	var step int
	if daily {
		step = int(math.Logb(float64(d / time.Hour)))
	} else if step = int(math.Round(float64(len(res)) / 30)); step == 0 {
		step = 1
	}
	slices.Reverse(res)
	for i, img := range res {
		if i%step == 0 {
			imgs = append(imgs, img)
		}
	}
	slices.Reverse(imgs)
	return
}

func animation(path, output string, d time.Duration, daily bool) error {
	imgs, err := getImages(path, d, daily)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}
	return encodeAnimation(output, imgs)
}

func updateSatellite(_ time.Time) {
	svc.Print("Start saving satellite map...")
	zoomMutex.Lock()
	defer zoomMutex.Unlock()
	if err := satellite(time.Time{}, location, "daily", mapOptions(*zoom)); err != nil {
		svc.Print(err)
		return
	}
	for _, t := range getTimes("daily") {
		if err := satellite(t, location, "daily", mapOptions(*zoom)); err != nil {
			svc.Print(err)
		}
	}
	for _, d := range animationDuration {
		if err := animation("daily/*", "animation/"+strings.TrimSuffix(d.String(), "0m0s"), d, true); err != nil {
			svc.Print(err)
		}
	}
}

func updateStorm(storms []storm.Data) {
	for _, i := range storms {
		dir := filepath.Join(*path, i.Season, fmt.Sprintf("%d-%s", i.No, i.ID))
		if err := satellite(time.Time{}, i.Coordinates(time.Now()), dir, mapOptions(*stormZoom)); err != nil {
			svc.Print(err)
			continue
		}
		for _, t := range getTimes(dir) {
			if coords := i.Coordinates(t); coords != nil && slices.Contains(stormMinute, t.Minute()) {
				if err := satellite(t, coords, dir, mapOptions(*stormZoom)); err != nil {
					svc.Print(err)
				}
			}
		}
		if err := animation(dir+"/*", dir, 0, false); err != nil {
			svc.Print(err)
		}
	}
}
