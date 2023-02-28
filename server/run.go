package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sunshineplan/database/mongodb/api"
	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/utils/retry"
	"github.com/sunshineplan/utils/scheduler"
	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/visualcrossing"
	"github.com/sunshineplan/weather/weatherapi"
)

func initWeather() error {
	var res struct {
		WeatherAPI     string
		VisualCrossing string
		Mongo          api.Client
		Dialer         mail.Dialer
		Subscriber     []string
	}
	if err := retry.Do(func() error {
		return meta.Get("weather", &res)
	}, 3, 20); err != nil {
		return err
	}
	realtime = weatherapi.New(res.WeatherAPI)
	var api weather.API
	switch *provider {
	case "weatherapi":
		api = realtime
	default:
		api = visualcrossing.New(res.VisualCrossing)
	}
	forecast = weather.New(api)
	history = api
	client = &res.Mongo
	dialer = res.Dialer
	to = res.Subscriber

	return client.Connect()
}

func test() (err error) {
	e1 := initWeather()
	if e1 != nil {
		fmt.Println("Failed to initialize weather config:", e1)
	} else {
		client.Close()
	}

	_, e2 := realtime.Realtime("Shanghai")
	if e2 != nil {
		fmt.Println("Failed to fetch realtime weather:", e2)
	}

	_, e3 := history.History("Shanghai", time.Now().AddDate(0, 0, -1))
	if e3 != nil {
		fmt.Println("Failed to fetch history weather:", e2)
	}

	if e1 != nil || e2 != nil || e3 != nil {
		err = fmt.Errorf("test is failed")
	}

	return
}

func run() {
	if err := initWeather(); err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	at := scheduler.NewScheduler().At
	at(scheduler.ScheduleFromString(*dailyReport)).Do(daily)
	at(scheduler.HourSchedule(9, 23)).Do(func(t time.Time) { record(t.AddDate(0, 0, -1)) })
	at(scheduler.ClockSchedule(scheduler.ClockFromString(*start), scheduler.ClockFromString(*end), *interval)).Do(alert)

	runServer()
}
