package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/HugoSmits86/nativewebp"
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
	satelliteMutex    sync.Mutex
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

func last() error {
	satelliteMutex.Lock()
	defer satelliteMutex.Unlock()
	time.Sleep(time.Second)
	_, img, err := mapAPI.Map(maps.Satellite, time.Time{}, location, mapOptions(*lastZoom))
	if err != nil {
		return err
	}
	f, err := os.Create("last.webp")
	if err != nil {
		return err
	}
	defer f.Close()
	return nativewebp.Encode(f, img, nil)
}

func satellite(t time.Time, coords coordinates.Coordinates, path string, opt any) (satelliteTime time.Time, err error) {
	satelliteMutex.Lock()
	defer satelliteMutex.Unlock()
	time.Sleep(time.Second)
	satelliteTime, img, err := mapAPI.Map(maps.Satellite, t, coords, opt)
	if err != nil {
		return
	}
	if err = os.MkdirAll(path, 0755); err != nil {
		return
	}
	file := filepath.Join(path, satelliteTime.Format(format)+".webp")
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return
	}
	defer f.Close()
	err = nativewebp.Encode(f, img, nil)
	return
}

func getTimes(t time.Time, s []string) (ts []time.Time) {
	for i := time.Duration(1); i <= 2*time.Hour/(10*time.Minute); i++ {
		if t := t.Add(-i * 10 * time.Minute); slices.IndexFunc(s, func(s string) bool {
			return strings.HasSuffix(s, t.Format(format)+".webp")
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
	t, err := satellite(time.Time{}, location, "daily", mapOptions(*zoom))
	if err != nil {
		svc.Print(err)
		return
	}
	res, err := filepath.Glob("daily/*.webp")
	if err != nil {
		svc.Print(err)
		return
	}
	for _, t := range getTimes(t, res) {
		if _, err := satellite(t, location, "daily", mapOptions(*zoom)); err != nil {
			svc.Print(err)
		}
	}
	for _, d := range animationDuration {
		if err := animation("daily/*", "animation/"+strings.TrimSuffix(d.String(), "0m0s"), d, true); err != nil {
			svc.Print(err)
		}
	}
	if err := last(); err != nil {
		svc.Print(err)
	}
}

func updateStorm(storms []storm.Data) {
	now := time.Now().In(timezone)
	for _, i := range storms {
		var last time.Time
		var err error
		dir := filepath.Join(*path, i.Season, fmt.Sprintf("%d-%s", i.No, i.ID))
		if coords := i.Coordinates(now); coords != nil {
			if last, err = satellite(time.Time{}, coords, dir, mapOptions(*stormZoom)); err != nil {
				svc.Print(err)
				continue
			}
		}
		res, err := filepath.Glob(dir + "/*.webp")
		if err != nil {
			svc.Print(err)
			continue
		}
		if last.IsZero() {
			last = now.Truncate(10 * time.Minute).Add(-time.Hour)
		}
		for _, t := range getTimes(last, res) {
			if coords := i.Coordinates(t); coords != nil && slices.Contains(stormMinute, t.Minute()) {
				if _, err := satellite(t, coords, dir, mapOptions(*stormZoom)); err != nil {
					svc.Print(err)
				}
			}
		}
		if err := animation(dir+"/*", dir, 0, false); err != nil {
			svc.Print(err)
		}
	}
}
