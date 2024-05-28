package main

import (
	"bytes"
	"image/jpeg"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/utils/httpsvr"
	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/aqi"
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

	router.GET("/img/:image", icon)
	router.GET("/storm/:storm", func(c *gin.Context) {
		storm := strings.ToLower(c.Param("storm"))
		c.File(filepath.Join(*path, storm, storm+".png"))
	})
	router.GET("/24h", func(c *gin.Context) {
		c.File("daily/daily-24h.png")
	})
	router.GET("/12h", func(c *gin.Context) {
		c.File("daily/daily-12h.png")
	})
	router.GET("/6h", func(c *gin.Context) {
		c.File("daily/daily-6h.png")
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
			c.File("daily/daily-24h.png")
		} else {
			coords, err := getCoords(q, nil)
			if err != nil {
				svc.Print(err)
				c.String(400, "")
				return
			}
			_, img, err := mapAPI.Realtime(weather.Satellite, coords, mapOptions(z))
			if err != nil {
				svc.Print(err)
				c.String(500, "")
				return
			}
			var buf bytes.Buffer
			if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: *quality}); err != nil {
				svc.Print(err)
				c.String(500, "")
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
				c.String(400, "")
				return
			}
			current, days, avg, aqi, err = getAllByCoordinates(coords, n, t, now, true)
			if err != nil {
				svc.Print(err)
				c.String(500, "")
				return
			}
		}
		var coords coordinates.Coordinates
		var image html.HTML
		if q == *query {
			coords = location
			image = imageHTML(mapAPI.URL(weather.Satellite, time.Time{}, location, mapOptions(z)), "/6h")
		} else {
			coords, err = getCoords(q, nil)
			if err != nil {
				svc.Print(err)
				c.String(400, "")
				return
			}
			image = imageHTML(mapAPI.URL(weather.Satellite, time.Time{}, coords, mapOptions(z)), "/map?q="+url.QueryEscape(q))
		}
		c.Data(200, "text/html", []byte(
			html.NewHTML().AppendChild(
				html.Head().AppendChild(
					html.Meta().Name("viewport").Attribute("content", "width=device-width"),
				),
				html.Body().Style("margin:0").
					Content(fullHTML(q, coords, current, days, avg, aqi, now, diff, "8px")+image),
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
			c.String(500, "")
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
			c.String(400, "")
			return
		}

		delete, _ := strconv.ParseBool(c.Query("delete"))

		res, err := export(month, delete)
		if err != nil {
			svc.Print(err)
			c.String(500, "")
			return
		}
		c.String(200, res)
	})

	return server.Run()
}
