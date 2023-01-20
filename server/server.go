package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
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

	router := httprouter.New()
	server.Handler = router

	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Access-Control-Allow-Methods", r.Method)
		header.Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusNoContent)
	})

	router.POST("/current", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		q := r.FormValue("q")
		if q == "" {
			q = getClientIP(r)
		}
		resp, err := weather.RealtimeWeather(q)
		if err != nil {
			log.Print(err)
			w.WriteHeader(500)
			return
		}
		b, _ := json.Marshal(resp)
		w.Write(b)
	})

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func getClientIP(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
	if clientIP == "" {
		clientIP = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	}
	if clientIP != "" {
		return clientIP
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}
	return ""
}
