package main

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/metadata"
	"github.com/sunshineplan/service"
	"github.com/sunshineplan/utils/flags"
	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/weatherapi"
)

var (
	svc    = service.New()
	meta   metadata.Server
	client mongodb.Client
	dialer mail.Dialer

	realtime *weatherapi.WeatherAPI
	forecast *weather.Weather
	history  weather.API
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
}

var (
	query       = flag.String("query", "", "weather query")
	dailyReport = flag.String("daily", "7:00", "daily report time")
	start       = flag.String("start", "6:00", "alert start time")
	end         = flag.String("end", "22:00", "alert end time")
	interval    = flag.Duration("interval", time.Hour, "alert interval")
	days        = flag.Int("days", 15, "forecast days")
	difference  = flag.Float64("difference", 5, "temperature difference")
	provider    = flag.String("provider", "visualcrossing", "weather provider")
	logPath     = flag.String("log", "", "Log file path")
)

func main() {
	self, err := os.Executable()
	if err != nil {
		svc.Fatalln("Failed to get self path:", err)
	}

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

	if service.IsWindowsService() {
		svc.Run()
		return
	}

	switch flag.NArg() {
	case 0:
		err = svc.Run()
	case 1:
		cmd := strings.ToLower(flag.Arg(0))
		var ok bool
		if ok, err = svc.Command(cmd); !ok {
			if cmd == "report" {
				if err := initWeather(); err != nil {
					svc.Fatal(err)
				}
				report(time.Now())
			} else {
				svc.Fatalln("Unknown argument:", cmd)
			}
		}
	default:
		svc.Fatalln("Unknown arguments:", strings.Join(flag.Args(), " "))
	}
	if err != nil {
		svc.Printf("failed to %s: %v", flag.Arg(0), err)
	}
}
