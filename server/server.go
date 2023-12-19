package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/utils/httpsvr"
	"github.com/sunshineplan/utils/log"
	"github.com/sunshineplan/weather/aqi"
)

var server = httpsvr.New()

func runServer() error {
	if *logPath != "" {
		svc.Logger = log.New(*logPath, "", log.LstdFlags)
		gin.DefaultWriter = svc.Logger
		gin.DefaultErrorWriter = svc.Logger
	}
	svc.Println("Location:", *query)
	svc.Println("Coordinates:", location)

	router := gin.Default()
	router.Use(cors.Default())
	router.TrustedPlatform = "X-Real-IP"
	server.Handler = router

	router.GET("/img/:image", icon)
	router.GET("/storm/:storm", func(c *gin.Context) {
		storm := strings.ToLower(c.Param("storm"))
		c.File(filepath.Join(*path, storm, storm+".gif"))
	})
	router.GET("/24h", func(c *gin.Context) {
		c.File("daily/daily-24h.gif")
	})
	router.GET("/12h", func(c *gin.Context) {
		c.File("daily/daily-12h.gif")
	})
	router.GET("/6h", func(c *gin.Context) {
		c.File("daily/daily-6h.gif")
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
			c.File("daily/daily-24h.gif")
		} else {
			coords, err := getCoords(q)
			if err != nil {
				svc.Print(err)
				c.String(400, "")
				return
			}
			b, err := coords.screenshot(z, 95, 3)
			if err != nil {
				svc.Print(err)
				c.String(500, "")
				return
			}
			c.Data(200, "image/jpeg", b)
		}
	})
	router.GET("/status", func(c *gin.Context) {
		t := time.Now()
		var q string
		if q = c.Query("q"); q == "" {
			q = *query
		}
		var n int
		if n, _ = strconv.Atoi(c.Query("n")); n < 1 {
			n = *days
		}
		var diff, z float64
		var err error
		if diff, err = strconv.ParseFloat(c.Query("diff"), 64); err != nil {
			diff = *difference
		}
		if z, err = strconv.ParseFloat(c.Query("z"), 64); err != nil {
			z = *zoom
		}
		days, avg, aqi, err := getAll(q, n, aqi.China, t)
		if err != nil {
			svc.Print(err)
			c.String(500, "")
			return
		}
		var image html.HTML
		if q == *query {
			image = imageHTML(location.url(z), "/6h")
			q = fmt.Sprintf("%s(%s)", *query, location)
		} else {
			coords, err := getCoords(q)
			if err != nil {
				svc.Print(err)
				c.String(400, "")
				return
			}
			image = imageHTML(coords.url(z), "/map?q="+url.QueryEscape(q))
			q = fmt.Sprintf("%s(%s)", q, coords)
		}
		c.Data(200, "text/html", []byte(
			html.NewHTML().AppendChild(
				html.Head().AppendChild(
					html.Meta().Name("viewport").Attribute("content", "width=device-width"),
				),
				html.Body().Style("margin:0").
					Content(fullHTML(q, days, avg, aqi, t, diff)+image),
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
