package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/utils/httpsvr"
	"github.com/sunshineplan/utils/log"
)

var server = httpsvr.New()

func runServer() error {
	if *logPath != "" {
		svc.Logger = log.New(*logPath, "", log.LstdFlags)
		gin.DefaultWriter = svc.Logger
		gin.DefaultErrorWriter = svc.Logger
	}
	svc.Println("Location:", *query)
	svc.Println("Coordinates:", coordinates)

	router := gin.Default()
	router.Use(cors.Default())
	router.TrustedPlatform = "X-Real-IP"
	server.Handler = router

	router.GET("/img/:image", icon)
	router.GET("/storm/:storm", func(c *gin.Context) {
		storm := strings.ToLower(c.Param("storm"))
		b, err := os.ReadFile(filepath.Join(*path, storm, storm+".gif"))
		if err != nil {
			svc.Print(err)
			if os.IsNotExist(err) {
				c.AbortWithStatus(404)
			} else {
				c.AbortWithStatus(500)
			}
			return
		}
		c.Data(200, "image/gif", b)
	})
	router.POST("/current", func(c *gin.Context) {
		q := c.Query("q")
		if q == "" {
			q = c.ClientIP()
		}
		resp, err := realtime.Request("current.json", fmt.Sprintf("q=%s", q))
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

		delete, _ := strconv.ParseBool("delete")

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
