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
	var b strings.Builder
	fmt.Fprint(&b, `<pre style="font-family:system-ui;margin:0">`)
	fmt.Fprint(&b, days[0].HTML())
	fmt.Fprintln(&b)
	fmt.Fprint(&b, `<div style="display:list-item;margin-left:15px">`, "Compared with Yesterday", "</div>")
	fmt.Fprint(&b, "<table><tbody>")
	fmt.Fprint(&b, yesterday.TemperatureHTML())
	fmt.Fprint(&b, weather.NewTempRiseFall(days[0], yesterday, 0).DiffInfoHTML())
	fmt.Fprint(&b, "</tbody></table>")
	fmt.Fprintln(&b)
	fmt.Fprint(&b, `<div style="display:list-item;margin-left:15px">`, "Historical Average Temperature of ", t.Format("01-02"), "</div>")
	fmt.Fprint(&b, "<table><tbody>")
	fmt.Fprint(&b, avg.TemperatureHTML())
	fmt.Fprint(&b, weather.NewTempRiseFall(days[0], avg, 0).DiffInfoHTML())
	fmt.Fprint(&b, "</tbody></table>")
	fmt.Fprintln(&b)
	fmt.Fprint(&b, `<div style="display:list-item;margin-left:15px">`, "Forecast", "</div>")
	fmt.Fprint(&b, table(days))
	if rainSnow != nil {
		fmt.Fprintln(&b)
		fmt.Fprint(&b, `<div style="display:list-item;margin-left:15px">`, "Recent Rain Snow Alert", "</div>")
		fmt.Fprint(&b, rainSnow.HTML())
	} else {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "No Rain Snow Alert.")
	}
	if tempRiseFall != nil {
		fmt.Fprintln(&b)
		if tempRiseFall.IsRise() {
			fmt.Fprint(&b, `<div style="display:list-item;margin-left:15px">`, "Recent Temperature Rise Alert", "</div>")
		} else {
			fmt.Fprint(&b, `<div style="display:list-item;margin-left:15px">`, "Recent Temperature Fall Alert", "</div>")
		}
		fmt.Fprint(&b, tempRiseFall.HTML())
	} else {
		fmt.Fprintln(&b)
		fmt.Fprint(&b, "No Temperature Alert.")
	}
	fmt.Fprint(&b, "</pre>")
	sendMail("[Weather]Daily Report"+timestamp(), b.String())
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
		sendMail(subject, `<pre style="font-family:system-ui;margin:0">`+body.String()+"</pre>")
	}
}

func isRainSnow(now int, hours []weather.Hour) bool {
	for _, i := range hours {
		if hour := i.Hour(); (hour == now || hour == now+1) && i.Precip > 0 {
			return true
		}
	}
	return false
}

func alertRainSnow(days []weather.Day) (subject string, b strings.Builder) {
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
				if rainSnow == nil ||
					rainSnow.Start().Date != i.Start().Date ||
					rainSnow.Duration() != i.Duration() ||
					(i.Start().Date == time.Now().Format("2006-01-02") && isRainSnow(time.Now().Hour(), i.Start().Hours)) {
					subject = "[Weather]Rain Snow Alert - " + i.Start().Date + timestamp()
				}
			}
			fmt.Fprint(&b, i.HTML())
			if index < len(res)-1 {
				fmt.Fprintln(&b)
			}
		}
		rainSnow = &first
	} else if rainSnow != nil {
		subject = "[Weather]Rain Snow Alert - Canceled" + timestamp()
		b.WriteString("No more rain snow")
		rainSnow = nil
	}
	return
}

func alertTempRiseFall(days []weather.Day) (subject string, b strings.Builder) {
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
			fmt.Fprint(&b, i.HTML())
			if index < len(res)-1 {
				fmt.Fprintln(&b)
			}
		}
		tempRiseFall = &first
	} else if tempRiseFall != nil {
		if tempRiseFall.IsRise() {
			subject = "[Weather]Temperature Rise Alert - Canceled" + timestamp()
			b.WriteString("No more temperature rise")
		} else {
			subject = "[Weather]Temperature Fall Alert - Canceled" + timestamp()
			b.WriteString("No more temperature fall")
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
	fmt.Fprint(&b, "<table border=1 cellspacing=0>")
	fmt.Fprint(&b, "<thead><tr><th colspan=2>Date</th><th>Max</th><th>Min</th><th>FLMax</th><th>FLMin</th><th>Rain%</th></tr></thead>")
	fmt.Fprint(&b, "<tbody>")
	for _, day := range days {
		fmt.Fprintf(&b, "<tr><td>%s</td>", day.DateInfo(false)[11:])
		fmt.Fprintf(&b, "<td>%s</td>", day.Condition.Img(day.Icon))
		fmt.Fprintf(&b, "<td>%s</td>", day.TempMax)
		fmt.Fprintf(&b, "<td>%s</td>", day.TempMin)
		fmt.Fprintf(&b, "<td>%s</td>", day.FeelsLikeMax)
		fmt.Fprintf(&b, "<td>%s</td>", day.FeelsLikeMin)
		fmt.Fprintf(&b, "<td>%s</td></tr>", day.PrecipProb)
	}
	fmt.Fprint(&b, "</tbody></table>")
	return b.String()
}
