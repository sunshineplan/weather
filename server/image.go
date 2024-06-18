package main

import (
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/executor"
	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/utils/pool"
)

var (
	iconCache sync.Map
	gifPool   = pool.New[gif.GIF]()
)

func icon(c *gin.Context) {
	file := strings.ToLower(c.Param("image"))
	if !strings.HasSuffix(file, ".png") {
		c.AbortWithStatus(404)
		return
	}
	icon := strings.TrimSuffix(file, ".png")
	if b, ok := iconCache.Load(icon); ok {
		c.Data(200, "image/png", b.([]byte))
		return
	}
	v, err := executor.ExecuteSerial(
		[]string{
			"https://cdn.jsdelivr.net/gh/visualcrossing/WeatherIcons@main/PNG/2nd Set - Color/%s.png",
			"https://fastly.jsdelivr.net/gh/visualcrossing/WeatherIcons@main/PNG/2nd Set - Color/%s.png",
			"https://raw.githubusercontent.com/visualcrossing/WeatherIcons/main/PNG/2nd Set - Color/%s.png",
		},
		func(url string) (any, error) {
			resp, err := http.Get(fmt.Sprintf(url, icon))
			if err != nil {
				return nil, err
			}
			if status := resp.StatusCode; status != 200 && status != 404 {
				return nil, fmt.Errorf("no StatusOK response: %d", status)
			}
			return resp, nil
		},
	)
	if err != nil {
		svc.Print(err)
		c.AbortWithStatus(500)
		return
	}

	resp := v.(*http.Response)
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		c.AbortWithStatus(404)
		return
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatus(500)
	} else {
		iconCache.Store(icon, b)
		c.Data(200, "image/png", b)
	}
}

func imageHTML(href, src string) html.HTML {
	return html.A().Href(href).AppendChild(html.Img().Src(src)).HTML()
}

func encodeGIF(w io.Writer, imgs []string) error {
	gifImg := gifPool.Get()
	defer func() {
		gifImg.Image = gifImg.Image[:0]
		gifImg.Delay = gifImg.Delay[:0]
		gifPool.Put(gifImg)
	}()
	for i, img := range imgs {
		f, err := os.Open(img)
		if err != nil {
			return err
		}
		if img, _, err := image.Decode(f); err != nil {
			svc.Print(err)
		} else {
			p := image.NewPaletted(img.Bounds(), palette.Plan9)
			draw.Draw(p, p.Rect, img, image.Point{}, draw.Over)
			gifImg.Image = append(gifImg.Image, p)
			if i != len(imgs)-1 {
				gifImg.Delay = append(gifImg.Delay, 40)
			} else {
				gifImg.Delay = append(gifImg.Delay, 300)
			}
		}
		f.Close()
	}
	return gif.EncodeAll(w, gifImg)
}
