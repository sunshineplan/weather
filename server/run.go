package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/sunshineplan/ai"
	"github.com/sunshineplan/ai/client"
	"github.com/sunshineplan/database/mongodb/driver"
	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/utils/retry"
	"github.com/sunshineplan/utils/scheduler"
	"github.com/sunshineplan/weather/api/airmatters"
	"github.com/sunshineplan/weather/api/visualcrossing"
	"github.com/sunshineplan/weather/api/weatherapi"
	"github.com/sunshineplan/weather/api/zoomearth"
)

func initWeather() (err error) {
	if *query == "" {
		return errors.New("query is empty")
	}
	var res struct {
		WeatherAPI     string
		VisualCrossing string
		AirMatters     string
		Mongo          driver.Client
		Dialer         mail.Dialer
		Subscriber     mail.Receipts
		AI             ai.ClientConfig
	}
	if err = retry.Do(func() error {
		return meta.Get("weather", &res)
	}, 3, 20); err != nil {
		return
	}
	realtime = weatherapi.New(res.WeatherAPI)
	switch *provider {
	case "weatherapi":
		forecast = realtime
	default:
		forecast = visualcrossing.New(res.VisualCrossing)
	}
	location, err = getCoords(*query, forecast)
	if err != nil {
		return
	}
	history = forecast
	mapAPI = zoomearth.ZoomEarthAPI{}
	stormAPI = zoomearth.ZoomEarthAPI{}
	aqiAPI = airmatters.New(res.AirMatters)
	aqiStandard, err = getAQIStandard()
	if err != nil {
		return
	}
	db = &res.Mongo
	dialer = res.Dialer
	to = res.Subscriber

	if res.AI.LLMs != "" {
		chatbot, err = client.New(res.AI)
		if err != nil {
			svc.Error("Failed to connenct AI chatbot", "error", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if model, err = chatbot.Model(ctx); err != nil {
			svc.Error("Failed to get AI model", "error", err)
		}
	}

	return db.Connect()
}

func test() (err error) {
	e1 := initWeather()
	if e1 != nil {
		fmt.Println("Failed to initialize weather config:", e1)
	} else {
		db.Close()
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

func run() error {
	if err := initWeather(); err != nil {
		return err
	}
	defer db.Close()

	if _, _, _, _, err := getAll(*query, *days, aqiType, time.Now(), false); err != nil {
		return err
	}

	run := func() *scheduler.Scheduler {
		if *debug {
			svc.SetLevel(slog.LevelDebug)
			return scheduler.NewScheduler().WithDebug(slog.New(svc.Logger.LoggerHandler()))
		}
		return scheduler.NewScheduler()
	}
	run().At(scheduler.ScheduleFromString(*dailyReport)).Do(daily)
	run().At(scheduler.HourSchedule(9, 16, 23)).Do(func(t time.Time) { record(t.AddDate(0, 0, -3)) })
	run().At(scheduler.MinuteSchedule(0)).Do(alert)
	run().At(scheduler.MinuteSchedule(15, 45)).Do(alertStorm)
	run().At(scheduler.MinuteSchedule(10, 30, 50)).Do(updateSatellite)

	go alert(time.Now())

	return runServer()
}
