package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/sunshineplan/weather"
)

var (
	rainSnow     *weather.RainSnow
	tempRiseFall *weather.TempRiseFall
)

func report(t time.Time) {
	date := t.Format("01-02")
	days, err := forecast.Forecast(*query, *days)
	if err != nil {
		log.Print(err)
		return
	}
	if !strings.HasSuffix(days[0].Date, date) {
		log.Println("first forecast is not today:", days[0].Date)
		return
	}
	yesterday, err := history.History(*query, t.AddDate(0, 0, -1))
	if err != nil {
		log.Print(err)
		return
	}
	avg, err := average(date, 2)
	if err != nil {
		log.Print(err)
		return
	}
	runAlert(days, alertRainSnow)
	runAlert(append([]weather.Day{yesterday}, days...), alertTempRiseFall)
	today(days[0], yesterday, avg, t)
}

func daily(t time.Time) {
	date := t.Format("01-02")
	days, err := forecast.Forecast(*query, 1)
	if err != nil {
		log.Print(err)
		return
	}
	if !strings.HasSuffix(days[0].Date, date) {
		log.Println("first forecast is not today:", days[0].Date)
		return
	}
	yesterday, err := history.History(*query, t.AddDate(0, 0, -1))
	if err != nil {
		log.Print(err)
		return
	}
	avg, err := average(date, 2)
	if err != nil {
		log.Print(err)
		return
	}

	today(days[0], yesterday, avg, t)
}

func today(today, yesterday, avg weather.Day, t time.Time) {
	var body strings.Builder
	fmt.Fprintln(&body, today.String())
	fmt.Fprintln(&body)
	fmt.Fprintln(&body, "Compared with Yesterday")
	fmt.Fprintln(&body, weather.NewTempRiseFall(today, yesterday).DiffInfo())
	fmt.Fprintln(&body)
	fmt.Fprintln(&body, "Historical Average Temperature of", t.Format("01-02"))
	fmt.Fprintln(&body, avg.Temperature())
	fmt.Fprintln(&body, weather.NewTempRiseFall(today, avg).DiffInfo())
	fmt.Fprintln(&body)
	if rainSnow != nil {
		fmt.Fprintln(&body, "Recent Rain Snow Alert:")
		fmt.Fprintln(&body, rainSnow.String())
	} else {
		fmt.Fprintln(&body, "No Rain Snow Alert.")
	}
	fmt.Fprintln(&body)
	if tempRiseFall != nil {
		if tempRiseFall.IsRise() {
			fmt.Fprintln(&body, "Recent Temperature Rise Alert:")
		} else {
			fmt.Fprintln(&body, "Recent Temperature Fall Alert:")
		}
		fmt.Fprintln(&body, tempRiseFall.String())
	} else {
		fmt.Fprintln(&body, "No Temperature Alert.")
	}
	sendMail("[Weather]Daily Report"+timestamp(), body.String())
}

func alert(t time.Time) {
	days, err := forecast.Forecast(*query, *days)
	if err != nil {
		log.Print(err)
		return
	}
	yesterday, err := history.History(*query, t.AddDate(0, 0, -1))
	if err != nil {
		log.Print(err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		runAlert(days, alertRainSnow)
	}()
	go func() {
		defer wg.Done()
		runAlert(append([]weather.Day{yesterday}, days...), alertTempRiseFall)
	}()
	wg.Wait()
}

func runAlert(days []weather.Day, fn func([]weather.Day) (string, strings.Builder)) {
	if subject, body := fn(days); subject != "" {
		log.Print(subject)
		sendMail(subject, body.String())
	}
}

func alertRainSnow(days []weather.Day) (subject string, body strings.Builder) {
	if rainSnow != nil {
		if rainSnow.IsExpired() {
			rainSnow = nil
		}
	}

	if res, err := weather.WillRainSnow(days); err != nil {
		log.Print(err)
	} else if len(res) > 0 {
		var first weather.RainSnow
		for index, i := range res {
			if index == 0 {
				first = i
				if rainSnow == nil || rainSnow.Start().Date != i.Start().Date || rainSnow.Duration() != i.Duration() {
					subject = "[Weather]Rain Snow Alert - " + i.Start().Date + timestamp()
				}
			}
			fmt.Fprintln(&body, i.String())
			if index < len(res)-1 {
				fmt.Fprintln(&body)
			}
		}
		rainSnow = &first
	} else if rainSnow != nil {
		subject = "[Weather]Rain Snow Alert - Canceled" + timestamp()
		body.WriteString("No more rain snow")
		rainSnow = nil
	}
	return
}

func alertTempRiseFall(days []weather.Day) (subject string, body strings.Builder) {
	if tempRiseFall != nil {
		if tempRiseFall.IsExpired() {
			tempRiseFall = nil
		}
	}

	if res, err := weather.WillTempRiseFall(days, *difference); err != nil {
		log.Print(err)
	} else if len(res) > 0 {
		var first weather.TempRiseFall
		for index, i := range res {
			if index == 0 {
				first = i
				if tempRiseFall == nil || tempRiseFall.Day().Date != i.Day().Date {
					if i.IsRise() {
						subject = "[Weather]Temperature Rise Alert - " + i.Day().Date + timestamp()
					} else {
						subject = "[Weather]Temperature Fall Alert - " + i.Day().Date + timestamp()
					}
				}
			}
			fmt.Fprintln(&body, i.String())
			if index < len(res)-1 {
				fmt.Fprintln(&body)
			}
		}
		tempRiseFall = &first
	} else if tempRiseFall != nil {
		if tempRiseFall.IsRise() {
			subject = "[Weather]Temperature Rise Alert - Canceled" + timestamp()
			body.WriteString("No more temperature rise")
		} else {
			subject = "[Weather]Temperature Fall Alert - Canceled" + timestamp()
			body.WriteString("No more temperature fall")
		}
		tempRiseFall = nil
	}
	return
}
