package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/executor"
)

var imgCache sync.Map

func icon(c *gin.Context) {
	file := strings.ToLower(c.Param("image"))
	if !strings.HasSuffix(file, ".png") {
		c.AbortWithStatus(404)
		return
	}
	icon := strings.TrimSuffix(file, ".png")
	if b, ok := imgCache.Load(icon); ok {
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
		imgCache.Store(icon, b)
		c.Data(200, "image/png", b)
	}
}
