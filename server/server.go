package main

import (
	"fmt"
	"html/template"
	"image/jpeg"
	"image/png"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/utils/httpsvr"
	"github.com/sunshineplan/weather/aqi"
	"github.com/sunshineplan/weather/maps"
	"github.com/sunshineplan/weather/unit/coordinates"
)

var server = httpsvr.New()

func runServer() error {
	svc.Println("Location:", *query)
	svc.Println("Coordinates:", location)
	svc.Println("AQI Type:", aqiType)
	svc.Println("AQI alert standard:", aqiStandard)
	if chatbot != nil {
		svc.Println("AI:", chatbot.LLMs())
		svc.Println("Model:", model)
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.TrustedPlatform = "X-Real-IP"
	server.Handler = router

	if *debug {
		debug := router.Group("debug")
		debug.GET("/", gin.WrapF(pprof.Index))
		//debug.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		debug.GET("/profile", gin.WrapF(pprof.Profile))
		debug.GET("/symbol", gin.WrapF(pprof.Symbol))
		debug.POST("/symbol", gin.WrapF(pprof.Symbol))
		debug.GET("/trace", gin.WrapF(pprof.Trace))
		debug.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		debug.GET("/block", gin.WrapH(pprof.Handler("block")))
		debug.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		debug.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		debug.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		debug.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}

	router.GET("/img/:image", icon)
	for _, i := range animationDuration {
		d := strings.TrimSuffix(i.String(), "0m0s")
		router.GET("/"+d, func(c *gin.Context) {
			c.File(fmt.Sprintf("animation/%s", d) + ext)
		})
	}
	router.GET("/last", func(c *gin.Context) {
		f, err := os.Open("last.png")
		if err != nil {
			svc.Print(err)
			c.Status(500)
			return
		}
		defer f.Close()
		img, err := png.Decode(f)
		if err != nil {
			svc.Print(err)
			c.Status(500)
			return
		}
		buf := bufPool.Get()
		defer func() {
			buf.Reset()
			bufPool.Put(buf)
		}()
		if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
			svc.Print(err)
			c.Status(500)
			return
		}
		c.Data(200, "image/jpeg", buf.Bytes())
	})
	router.GET("/map", func(c *gin.Context) {
		var q string
		if q = c.Query("q"); q == "" {
			q = *query
		}
		var z float64
		var err error
		if z, err = strconv.ParseFloat(c.Query("z"), 64); err != nil {
			z = *zoom
		}
		if q == *query {
			c.File("animation/24h" + ext)
		} else {
			coords, err := getCoords(q, nil)
			if err != nil {
				svc.Print(err)
				c.Status(400)
				return
			}
			_, img, err := mapAPI.Realtime(maps.Satellite, coords, mapOptions(z))
			if err != nil {
				svc.Print(err)
				c.Status(500)
				return
			}
			buf := bufPool.Get()
			defer func() {
				buf.Reset()
				bufPool.Put(buf)
			}()
			if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
				svc.Print(err)
				c.Status(500)
				return
			}
			c.Data(200, "image/jpeg", buf.Bytes())
		}
	})
	router.GET("/status", func(c *gin.Context) {
		now := time.Now()
		var q string
		if q = c.Query("q"); q == "" {
			q = *query
		}
		var n int
		if n, _ = strconv.Atoi(c.Query("n")); n < 1 {
			n = *days
		}
		var t aqi.Type
		var err error
		if err = t.UnmarshalText([]byte(c.Query("aqi"))); err != nil {
			t = aqiType
		}
		var diff, z float64
		if diff, err = strconv.ParseFloat(c.Query("diff"), 64); err != nil {
			diff = *difference
		}
		if z, err = strconv.ParseFloat(c.Query("z"), 64); err != nil {
			z = *zoom
		}
		current, days, avg, aqi, err := getAll(q, n, t, now, true)
		if err != nil {
			coords, err := getCoords(q, nil)
			if err != nil {
				svc.Print(err)
				c.Status(400)
				return
			}
			current, days, avg, aqi, err = getAllByCoordinates(coords, n, t, now, true)
			if err != nil {
				svc.Print(err)
				c.Status(500)
				return
			}
		}
		var coords coordinates.Coordinates
		var image html.HTML
		if q == *query {
			coords = location
			image = imageHTML(mapAPI.URL(maps.Satellite, time.Time{}, location, mapOptions(z)), "/6h")
		} else {
			coords, err = getCoords(q, nil)
			if err != nil {
				svc.Print(err)
				c.Status(400)
				return
			}
			image = imageHTML(mapAPI.URL(maps.Satellite, time.Time{}, coords, mapOptions(z)), "/map?q="+url.QueryEscape(q))
		}
		c.Data(200, "text/html", []byte(
			html.NewHTML().AppendChild(
				html.Head().AppendChild(
					html.Meta().Name("viewport").Attribute("content", "width=device-width"),
				),
				html.Body().Style("margin:0").
					Content(fullHTML(q, coords, current, days, avg, aqi, now, true, diff, "8px")+image),
			).HTML()),
		)
	})
	router.POST("/current", func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			q = c.ClientIP()
		}
		resp, err := realtime.Request("current.json", url.Values{"q": {q}})
		if err != nil {
			svc.Print(err)
			c.Status(500)
			return
		}
		c.JSON(200, resp)
	})
	router.POST("/history", func(c *gin.Context) {
		month := c.Query("month")
		if month == "" {
			month = time.Now().AddDate(0, -1, 0).Format("2006-01")
		}
		if _, err := time.Parse("2006-01", month); err != nil {
			svc.Print(err)
			c.Status(400)
			return
		}

		delete, _ := strconv.ParseBool(c.Query("delete"))

		res, err := export(month, delete)
		if err != nil {
			svc.Print(err)
			c.Status(500)
			return
		}
		c.String(200, res)
	})

	t := template.Must(template.New("").Parse(`<html>
<head><title>{{.Title}}</title></head>
<body><h1>{{.Title}}</h1>
<pre>{{if not .Root}}<a href="..">..</a>
{{end}}{{range .Dirs}}
<a href="{{.}}">{{.}}</a>{{end}}</pre>
</body></html>`))
	storm := router.Group("storm")
	storm.GET("/", func(c *gin.Context) {
		root, err := http.Dir("storm").Open(".")
		if err != nil {
			svc.Print(err)
			c.Status(500)
			return
		}
		years, err := root.Readdir(-1)
		if err != nil {
			svc.Print(err)
			c.Status(500)
			return
		}
		var data = struct {
			Title string
			Root  bool
			Dirs  []string
		}{"Storm", true, nil}
		for _, i := range years {
			data.Dirs = append(data.Dirs, i.Name())
		}
		if err := t.Execute(c.Writer, data); err != nil {
			svc.Print(err)
			c.Status(500)
		}
	})
	storm.GET("/:year/", func(c *gin.Context) {
		year, err := http.Dir("storm").Open(c.Param("year"))
		if err != nil {
			svc.Print(err)
			c.Status(500)
			return
		}
		files, err := year.Readdir(-1)
		if err != nil {
			svc.Print(err)
			c.Status(500)
			return
		}
		var data = struct {
			Title string
			Root  bool
			Dirs  []string
		}{"Storm - " + c.Param("year"), false, nil}
		for _, i := range files {
			if i.IsDir() {
				data.Dirs = append(data.Dirs, i.Name())
			}
		}
		if err := t.Execute(c.Writer, data); err != nil {
			svc.Print(err)
			c.Status(500)
		}
	})
	storm.GET("/:year/:id/", func(c *gin.Context) {
		res, err := filepath.Glob(filepath.Join("storm", c.Param("year"), c.Param("id")+ext))
		if err != nil {
			svc.Print(err)
			c.Status(500)
			return
		}
		if l := len(res); l == 0 {
			c.String(404, "404 page not found")
			return
		}
		c.File(res[0])
	})

	return server.Run()
}
