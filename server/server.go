package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/httpsvr"
)

var server = httpsvr.New()

func runServer() {
	if *logPath != "" {
		f, err := os.OpenFile(*logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatalln("Failed to open log file:", err)
		}
		gin.DefaultWriter = f
		gin.DefaultErrorWriter = f
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
		resp, err := realtime.Request("current.json", fmt.Sprintf("q=%s", q))
		if err != nil {
			log.Print(err)
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
			log.Print(err)
			c.String(400, "")
			return
		}

		delete, _ := strconv.ParseBool("delete")

		buf, err := export(month, delete)
		if err != nil {
			log.Print(err)
			c.String(500, "")
			return
		}
		c.String(200, buf.String())
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
