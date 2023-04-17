package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sunshineplan/weather"
)

var (
	rainSnow     *weather.RainSnow
	tempRiseFall *weather.TempRiseFall
)

func prepare(t time.Time) (forecasts []weather.Day, yesterday, avg weather.Day, err error) {
	date := t.Format("01-02")
	forecasts, err = forecast.Forecast(*query, *days)
	if err != nil {
		return
	}
	if !strings.HasSuffix(forecasts[0].Date, date) {
		err = fmt.Errorf("first forecast is not today: %s", forecasts[0].Date)
		return
	}
	yesterday, err = history.History(*query, t.AddDate(0, 0, -1))
	if err != nil {
		return
	}
	avg, err = average(date, 2)
	return
}

func report(t time.Time) {
	days, yesterday, avg, err := prepare(t)
	if err != nil {
		svc.Print(err)
		return
	}
	runAlert(days, alertRainSnow)
	runAlert(append([]weather.Day{yesterday}, days...), alertTempRiseFall)
	today(days, yesterday, avg, t)
}

func daily(t time.Time) {
	svc.Print("Start sending daily report...")
	days, yesterday, avg, err := prepare(t)
	if err != nil {
		svc.Print(err)
		return
	}
	today(days, yesterday, avg, t)
}

func today(days []weather.Day, yesterday, avg weather.Day, t time.Time) {
	var body strings.Builder
	fmt.Fprintln(&body, `<pre style="font-family:system-ui">`)
	fmt.Fprintln(&body, days[0])
	fmt.Fprintln(&body)
	fmt.Fprintln(&body, "Compared with Yesterday")
	fmt.Fprintln(&body, weather.NewTempRiseFall(days[0], yesterday).DiffInfo())
	fmt.Fprintln(&body)
	fmt.Fprintln(&body, "Historical Average Temperature of", t.Format("01-02"))
	fmt.Fprintln(&body, avg.Temperature())
	fmt.Fprintln(&body, weather.NewTempRiseFall(days[0], avg).DiffInfo())
	fmt.Fprintln(&body)
	fmt.Fprintln(&body, "Forecast:")
	fmt.Fprintln(&body, table(days))
	if rainSnow != nil {
		fmt.Fprintln(&body, "Recent Rain Snow Alert:")
		fmt.Fprintln(&body, rainSnow)
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
		fmt.Fprintln(&body, tempRiseFall)
	} else {
		fmt.Fprintln(&body, "No Temperature Alert.")
	}
	fmt.Fprintln(&body, "</pre>")
	sendMail("[Weather]Daily Report"+timestamp(), body.String())
}

func alert(t time.Time) {
	svc.Print("Start alerting...")
	days, err := forecast.Forecast(*query, *days)
	if err != nil {
		svc.Print(err)
		return
	}
	yesterday, err := history.History(*query, t.AddDate(0, 0, -1))
	if err != nil {
		svc.Print(err)
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
		svc.Print(subject)
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
		svc.Print(err)
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
		svc.Print(err)
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

func table(days []weather.Day) string {
	if len(days) > 7 {
		days = days[:7]
	}
	var b strings.Builder
	fmt.Fprintln(&b, `<table border="1" cellspacing="0">`)
	fmt.Fprintln(&b, "<thead>")
	fmt.Fprintln(&b, "<tr>")
	fmt.Fprintln(&b, "<th>Date</th>")
	fmt.Fprintln(&b, "<th>Max</th>")
	fmt.Fprintln(&b, "<th>Min</th>")
	fmt.Fprintln(&b, "<th>FLMax</th>")
	fmt.Fprintln(&b, "<th>FLMin</th>")
	fmt.Fprintln(&b, "<th>Rain%</th>")
	fmt.Fprintln(&b, "<th>Condition</th>")
	fmt.Fprintln(&b, "</tr>")
	fmt.Fprintln(&b, "</thead>")
	fmt.Fprintln(&b, "<tbody>")
	for _, day := range days {
		fmt.Fprintln(&b, "<tr>")
		fmt.Fprintf(&b, "<td>%s</td>\n", day.DateInfo(false)[11:])
		fmt.Fprintf(&b, "<td>%s</td>\n", day.TempMax)
		fmt.Fprintf(&b, "<td>%s</td>\n", day.TempMin)
		fmt.Fprintf(&b, "<td>%s</td>\n", day.FeelsLikeMax)
		fmt.Fprintf(&b, "<td>%s</td>\n", day.FeelsLikeMin)
		fmt.Fprintf(&b, "<td>%s</td>\n", day.PrecipProb)
		fmt.Fprintf(&b, "<td>%s</td>\n", day.Condition)
		fmt.Fprintln(&b, "</tr>")
	}
	fmt.Fprintln(&b, "</tbody>")
	fmt.Fprint(&b, "</table>")
	return b.String()
}
