package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/httpsvr"
	"github.com/sunshineplan/weather"
)

var server = httpsvr.New()

func runServer() {
	if *logPath != "" {
		f, err := os.OpenFile(*logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatalln("Failed to open log file:", err)
		}
		log.SetOutput(f)
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.TrustedPlatform = "X-Real-IP"
	server.Handler = router

	router.POST("/current", func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			q = c.ClientIP()
		}
		resp, err := weather.RealtimeWeather(q)
		if err != nil {
			log.Print(err)
			c.String(500, "")
			return
		}
		c.JSON(200, resp)
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
