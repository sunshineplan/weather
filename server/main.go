package main

import (
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/metadata"
	"github.com/sunshineplan/service"
	"github.com/sunshineplan/utils/flags"
	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/api/weatherapi"
	"github.com/sunshineplan/weather/aqi"
	"github.com/sunshineplan/weather/storm"
)

var (
	svc    = service.New()
	meta   metadata.Server
	client mongodb.Client
	dialer mail.Dialer

	realtime *weatherapi.WeatherAPI
	forecast weather.API
	history  weather.API
	stormAPI storm.API
	aqiAPI   aqi.API
)

func init() {
	svc.Name = "Weather"
	svc.Desc = "weather service"
	svc.Exec = run
	svc.TestExec = test
	svc.Options = service.Options{
		Dependencies: []string{"Wants=network-online.target", "After=network.target"},
		Environment:  map[string]string{"GIN_MODE": "release"},
		ExcludeFiles: []string{"scripts/weather.conf"},
	}
	svc.RegisterCommand("report", "report", func(_ ...string) error {
		if err := initWeather(); err != nil {
			return err
		}
		report(time.Now())
		return nil
	}, 0, true)
	svc.RegisterCommand("record", "record", func(arg ...string) error {
		if err := initWeather(); err != nil {
			return err
		}
		date, err := time.Parse("2006-01-02", arg[0])
		if err != nil {
			return err
		}
		return record(date)
	}, 1, true)
}

var (
	query       = flag.String("query", "", "weather query")
	dailyReport = flag.String("daily", "7:00", "daily report time")
	start       = flag.String("start", "6:00", "alert start time")
	end         = flag.String("end", "22:00", "alert end time")
	days        = flag.Int("days", 15, "forecast days")
	difference  = flag.Float64("difference", 5, "temperature difference")
	aqiType     aqi.Type
	offset      = flag.Float64("offset", 0, "daily screenshot offset")
	zoom        = flag.Float64("zoom", 7, "daily screenshot zoom")
	quality     = flag.Int("quality", 95, "screenshot quality")
	radius      = flag.Float64("radius", 700, "storm affect radius (unit: km)")
	path        = flag.String("path", "storm", "storm screenshot save path")
	provider    = flag.String("provider", "visualcrossing", "weather provider")
	logPath     = flag.String("log", "", "Log file path")
)

func main() {
	self, err := os.Executable()
	if err != nil {
		svc.Fatalln("Failed to get self path:", err)
	}

	flag.TextVar(&aqiType, "aqi", aqi.US, "AQI Type")

	flag.StringVar(&meta.Addr, "server", "", "Metadata Server Address")
	flag.StringVar(&meta.Header, "header", "", "Verify Header Header Name")
	flag.StringVar(&meta.Value, "value", "", "Verify Header Value")
	flag.StringVar(&server.Unix, "unix", "", "UNIX-domain Socket")
	flag.StringVar(&server.Host, "host", "0.0.0.0", "Server Host")
	flag.StringVar(&server.Port, "port", "12345", "Server Port")
	flag.StringVar(&svc.Options.UpdateURL, "update", "", "Update URL")
	flag.StringVar(&svc.Options.PIDFile, "pid", "/var/run/weather.pid", "PID file path")
	flags.SetConfigFile(filepath.Join(filepath.Dir(self), "config.ini"))
	flags.Parse()

	if err := svc.ParseAndRun(flag.Args()); err != nil {
		svc.Fatal(err)
	}
}
