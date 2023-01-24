package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/metadata"
	"github.com/sunshineplan/service"
	"github.com/sunshineplan/utils/flags"
	"github.com/sunshineplan/utils/mail"
)

var (
	meta   metadata.Server
	client mongodb.Client
	dialer mail.Dialer
)

var svc = service.Service{
	Name:     "Weather",
	Desc:     "weather api",
	Exec:     run,
	TestExec: test,
	Options: service.Options{
		Dependencies: []string{"Wants=network-online.target", "After=network.target"},
		Environment:  map[string]string{"GIN_MODE": "release"},
		ExcludeFiles: []string{"scripts/weather.conf"},
	},
}

var (
	query      = flag.String("query", "Shanghai", "weather query")
	difference = flag.Float64("difference", 5, "temperature difference")
	logPath    = flag.String("log", "", "Log Path")
)

func main() {
	self, err := os.Executable()
	if err != nil {
		log.Fatalln("Failed to get self path:", err)
	}

	flag.StringVar(&meta.Addr, "server", "", "Metadata Server Address")
	flag.StringVar(&meta.Header, "header", "", "Verify Header Header Name")
	flag.StringVar(&meta.Value, "value", "", "Verify Header Value")
	flag.StringVar(&server.Unix, "unix", "", "UNIX-domain Socket")
	flag.StringVar(&server.Host, "host", "0.0.0.0", "Server Host")
	flag.StringVar(&server.Port, "port", "12345", "Server Port")
	flag.StringVar(&svc.Options.UpdateURL, "update", "", "Update URL")
	flags.SetConfigFile(filepath.Join(filepath.Dir(self), "config.ini"))
	flags.Parse()

	if service.IsWindowsService() {
		svc.Run(false)
		return
	}

	switch flag.NArg() {
	case 0:
		run()
	case 1:
		cmd := flag.Arg(0)
		var ok bool
		if ok, err = svc.Command(cmd); !ok {
			log.Fatalln("Unknown argument:", cmd)
		}
	default:
		log.Fatalln("Unknown arguments:", strings.Join(flag.Args(), " "))
	}
	if err != nil {
		log.Fatalf("failed to %s: %v", flag.Arg(0), err)
	}
}
