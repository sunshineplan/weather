package main

import (
	"cmp"
	"fmt"
	"html/template"
	"image/jpeg"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/utils/httpsvr"
	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/aqi"
	"github.com/sunshineplan/weather/maps"
	"github.com/sunshineplan/weather/unit"
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

	//if *debug {
	//	debug := router.Group("debug")
	//	debug.GET("/", gin.WrapF(pprof.Index))
	//	//debug.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	//	debug.GET("/profile", gin.WrapF(pprof.Profile))
	//	debug.GET("/symbol", gin.WrapF(pprof.Symbol))
	//	debug.POST("/symbol", gin.WrapF(pprof.Symbol))
	//	debug.GET("/trace", gin.WrapF(pprof.Trace))
	//	debug.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
	//	debug.GET("/block", gin.WrapH(pprof.Handler("block")))
	//	debug.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
	//	debug.GET("/heap", gin.WrapH(pprof.Handler("heap")))
	//	debug.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
	//	debug.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	//}

	router.GET("/img/:image", icon)
	for _, i := range animationDuration {
		d := strings.TrimSuffix(i.String(), "0m0s")
		router.GET("/"+d, func(c *gin.Context) {
			c.File(fmt.Sprintf("animation/%s", d) + ext)
		})
	}
	router.GET("/last", func(c *gin.Context) {
		b, err := os.ReadFile("last.webp")
		if err != nil {
			svc.Print(err)
			c.Status(500)
			return
		}
		c.Data(200, "image/webp", b)
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
		c.Data(200, "text/html; charset=utf-8", []byte(
			html.NewHTML().AppendChild(
				html.Head().AppendChild(
					html.Meta().Name("viewport").Attribute("content", "width=device-width"),
				),
				html.Body().Style("margin:0").
					Content(fullHTML(q, coords, current, days, avg, aqi, now, true, diff, "8px")+image),
			).HTML()),
		)
	})
	router.GET("/hourly", func(c *gin.Context) {
		var q string
		if q = c.Query("q"); q == "" {
			q = *query
		}
		forecasts, err := forecast.Forecast(q, 2)
		if err != nil {
			coords, err := getCoords(q, nil)
			if err != nil {
				svc.Print(err)
				c.Status(400)
				return
			}
			forecasts, err = forecast.ForecastByCoordinates(coords, 2)
			if err != nil {
				svc.Print(err)
				c.Status(500)
				return
			}
		}
		now := time.Now()
		var hours []weather.Hour
	Loop:
		for _, i := range forecasts {
			for _, i := range i.Hours {
				if d := i.TimeEpoch.Time().Sub(now); d > -2*time.Hour {
					if hours = append(hours, i); len(hours) == 48 {
						break Loop
					}
				}
			}
		}
		table := html.Table().Attribute("border", "1").Attribute("cellspacing", "0")
		th := []*html.TableCell{
			html.Th("Time").Colspan(2),
			html.Th("Temp./FL"),
			html.Th("RH"),
			html.Th("Pressure"),
			html.Th("Precip."),
			html.Th("Wind"),
			html.Th("Dir"),
		}
		table.AppendChild(html.Thead().AppendChild(html.Tr(th...)))
		tbody := html.Tbody()
		for _, hour := range hours {
			t := hour.TimeEpoch.Time()
			var dateContent any
			if date := t.Format("2006-01-02 15:04"); now.Truncate(time.Hour).Equal(t) {
				dateContent = html.Span().Style("color:red").Content(date)
			} else {
				dateContent = date
			}
			td := []*html.TableCell{
				html.Td(dateContent).Style("text-align:center;padding:0 5px"),
				html.Td(hour.Condition.Img(hour.Icon)),
				html.Td(hour.Temp.HTML() + " / " + hour.FeelsLike.HTML()).Style("text-align:center;padding:0 5px"),
				html.Td(hour.Humidity).Style("text-align:center;padding:0 5px"),
				html.Td(fmt.Sprintf("%ghPa", hour.Pressure)).Style("text-align:center;padding:0 5px"),
				html.Td(fmt.Sprintf("%gmm(%s)", hour.Precip, hour.PrecipProb)).Style("text-align:center;padding:0 5px"),
				html.Td(html.Span().Style("color:"+hour.WindSpeed.ForceColor()).Contentf("%sm/s", unit.FormatFloat64(hour.WindSpeed.MPS(), 1))),
				html.Td(html.Div().Style("display:flex;justify-content:center").Content(hour.WindDir)),
			}
			tbody.AppendChild(html.Tr(td...))
		}
		c.Data(200, "text/html; charset=utf-8", []byte(html.Div().Style("font-family:system-ui").AppendChild(table.AppendChild(tbody)).HTML()))
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
<head><meta name="viewport" content="width=device-width,initial-scale=1">
<title>{{.Title}}</title></head>
<body><h1>{{.Title}}</h1>
<pre>{{if not .Root}}<a href="..">..</a>
{{end}}{{range .Dirs}}
<a href="{{.}}">{{.}}</a>{{end}}</pre>
</body></html>`))
	storm := router.Group("storm")
	storm.GET("/", func(c *gin.Context) {
		years, err := os.ReadDir("storm")
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
		year := c.Param("year")
		files, err := os.ReadDir("storm/" + year)
		if err != nil {
			svc.Print(err)
			c.Status(500)
			return
		}
		var data = struct {
			Title string
			Root  bool
			Dirs  []string
		}{"Storm - " + year, false, nil}
		for _, i := range files {
			if i.IsDir() {
				data.Dirs = append(data.Dirs, i.Name())
			}
		}
		slices.SortStableFunc(data.Dirs, func(a, b string) int {
			id := func(s string) int {
				s, _, _ = strings.Cut(s, "-")
				id, _ := strconv.Atoi(s)
				return id
			}
			return cmp.Compare(id(a), id(b))
		})
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
