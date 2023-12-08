package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/aqi"
	"github.com/sunshineplan/weather/storm"
	"github.com/sunshineplan/weather/unit"
)

var (
	location     coords
	rainSnow     []weather.RainSnow
	tempRiseFall []weather.TempRiseFall

	alertMutex sync.Mutex
	zoomMutex  sync.Mutex
)

func prepare(t time.Time) (forecasts []weather.Day, yesterday, avg weather.Day, aqiCurrent aqi.Current, err error) {
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
	if err != nil {
		return
	}
	aqiCurrent, err = aqiAPI.Realtime(aqi.China, *query)
	return
}

func report(t time.Time) {
	days, yesterday, avg, aqi, err := prepare(t)
	if err != nil {
		svc.Print(err)
		return
	}
	runAlert(days, alertRainSnow)
	runAlert(append([]weather.Day{yesterday}, days...), alertTempRiseFall)
	zoomEarth(t, true)
	sendMail(
		"[Weather]Daily Report"+timestamp(),
		today(days, yesterday, avg, aqi, t, ""),
		attachment("daily/daily-12h.gif"),
		true,
	)
}

func daily(t time.Time) {
	svc.Print("Start sending daily report...")
	days, yesterday, avg, aqi, err := prepare(t)
	if err != nil {
		svc.Print(err)
		return
	}
	sendMail(
		"[Weather]Daily Report"+timestamp(),
		today(days, yesterday, avg, aqi, t, ""),
		attachment("daily/daily-12h.gif"),
		true,
	)
}

func today(days []weather.Day, yesterday, avg weather.Day, aqi aqi.Current, t time.Time, src string) string {
	var b strings.Builder
	fmt.Fprint(&b, `<pre style="font-family:system-ui;margin:0">`)
	fmt.Fprint(&b, days[0].HTML())
	fmt.Fprintln(&b)
	fmt.Fprint(&b, aqiHTML(aqi))
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
	fmt.Fprint(&b, forecastHTML(days))

	alertMutex.Lock()
	defer alertMutex.Unlock()

	if len(rainSnow) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprint(&b, `<div style="display:list-item;margin-left:15px">`, "Recent Rain Snow Alert", "</div>")
		for index, i := range rainSnow {
			fmt.Fprint(&b, i.HTML(t))
			if index < len(rainSnow)-1 {
				fmt.Fprintln(&b)
			}
		}
	} else {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "No Rain Snow Alert.")
	}
	if len(tempRiseFall) > 0 {
		fmt.Fprintln(&b)
		fmt.Fprint(&b, `<div style="display:list-item;margin-left:15px">`, "Recent Temperature Alert", "</div>")
		for index, i := range tempRiseFall {
			fmt.Fprint(&b, i.HTML())
			if index < len(tempRiseFall)-1 {
				fmt.Fprintln(&b)
			}
		}
	} else {
		fmt.Fprintln(&b)
		fmt.Fprintln(&b, "No Temperature Alert.")
	}
	fmt.Fprint(&b, "\n</pre>")
	if src == "" {
		fmt.Fprintf(&b, "<a href=%q><img src='cid:attachment'></a>", location.url(*zoom))
	} else {
		fmt.Fprintf(&b, "<a href=%q><img src=%q></a>", location.url(*zoom), src)
	}
	return b.String()
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

	alertMutex.Lock()
	defer alertMutex.Unlock()
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		runAlert(days, alertRainSnow)
	}()
	go func() {
		defer wg.Done()
		runAlert(append([]weather.Day{yesterday}, days...), alertTempRiseFall)
	}()
	go func() {
		defer wg.Done()
		runAlert(nil, alertAQI)
	}()
	wg.Wait()
}

func runAlert(days []weather.Day, fn func([]weather.Day) (string, strings.Builder)) {
	if subject, body := fn(days); subject != "" {
		svc.Print(subject)
		sendMail(subject, `<pre style="font-family:system-ui;margin:0">`+body.String()+"</pre>", nil, false)
	}
}

func isRainSnow(now int, hours []weather.Hour) bool {
	for _, i := range hours {
		if hour := i.TimeEpoch.Time().Hour(); hour == now && i.Precip > 0 {
			return true
		}
	}
	return false
}

func alertRainSnow(days []weather.Day) (subject string, b strings.Builder) {
	if len(rainSnow) > 0 {
		if rainSnow[0].IsExpired() {
			rainSnow = rainSnow[1:]
		}
	}

	if res, err := weather.WillRainSnow(days); err != nil {
		svc.Print(err)
	} else if n := len(res); n > 0 {
		for index, i := range res {
			now := time.Now()
			hour := now.Hour()
			if index == 0 {
				if start := i.Start(); len(rainSnow) == 0 ||
					rainSnow[0].Start().Date != start.Date ||
					rainSnow[0].Duration() != i.Duration() {
					subject = "[Weather]Rain Snow Alert - " + start.Date + timestamp()
				} else if start.Date == now.Format("2006-01-02") && isRainSnow(hour, start.Hours) {
					subject = "[Weather]Rain Snow Alert - Today" + timestamp()
					fmt.Fprintln(&b, start.DateInfoHTML())
					fmt.Fprintln(&b, start.PrecipitationHTML(hour))
					for index, n := 1, len(i.Days()); index < n && index < 3; index++ {
						fmt.Fprintln(&b)
						fmt.Fprintln(&b, i.Days()[index].DateInfoHTML())
						fmt.Fprintln(&b, i.Days()[index].PrecipitationHTML())
					}
					return
				}
			}
			fmt.Fprint(&b, i.HTML(now, hour))
			if index < len(res)-1 {
				fmt.Fprintln(&b)
			}
		}
		rainSnow = res
	} else if len(rainSnow) > 0 {
		subject = "[Weather]Rain Snow Alert - Canceled" + timestamp()
		b.WriteString("No more rain snow")
		rainSnow = nil
	}
	return
}

func alertTempRiseFall(days []weather.Day) (subject string, b strings.Builder) {
	if len(tempRiseFall) > 0 {
		if tempRiseFall[0].IsExpired() {
			tempRiseFall = tempRiseFall[1:]
		}
	}

	if res, err := weather.WillTempRiseFall(days, *difference); err != nil {
		svc.Print(err)
	} else if len(res) > 0 {
		for index, i := range res {
			if index == 0 {
				if len(tempRiseFall) == 0 ||
					tempRiseFall[0].Day().Date != i.Day().Date ||
					(tempRiseFall[0].Day().Date == i.Day().Date && tempRiseFall[0].IsRise() != i.IsRise()) {
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
		tempRiseFall = res
	} else if len(tempRiseFall) > 0 {
		if tempRiseFall[0].IsRise() {
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

func aqiHTML(current aqi.Current) string {
	var b strings.Builder
	fmt.Fprintf(&b,
		`<div style="display:list-item;margin-left:15px">%s:<span style="padding:0 1em;color:white;background-color:%s">%d %s</span></div>`,
		current.AQI().Type(), current.AQI().Level().Color(), current.AQI().Value(), current.AQI().Level(),
	)
	fmt.Fprint(&b, "<table><tbody>")
	for i, p := range current.Pollutants() {
		if i%3 == 0 {
			if i != 0 {
				fmt.Fprint(&b, "</tr>")
			}
			fmt.Fprint(&b, "<tr>")
		}
		fmt.Fprintf(&b, `<td>%s:</td><td style="color:%s">%s %s</td>`,
			p.Kind().HTML(), p.Level().Color(), unit.FormatFloat64(p.Value(), 2), p.Unit())
	}
	fmt.Fprint(&b, "</tr>")
	fmt.Fprint(&b, "</tbody></table>")
	return b.String()
}

func forecastHTML(days []weather.Day) string {
	if len(days) > 10 {
		days = days[:10]
	}
	var b strings.Builder
	fmt.Fprint(&b, `<div style="display:list-item;margin-left:15px">Forecast</div>`)
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

func zoomEarth(t time.Time, isReport bool) {
	if !isReport {
		go func() {
			zoomMutex.Lock()
			defer zoomMutex.Unlock()
			b, err := location.offset(0, *offset).screenshot(*zoom, *quality, 5)
			if err != nil {
				svc.Print(err)
				return
			}
			file := fmt.Sprintf("daily/%s.jpg", t.Format("200601021504"))
			if err := os.MkdirAll("daily", 0755); err != nil {
				svc.Print(err)
				return
			}
			if err := os.WriteFile(file, b, 0644); err != nil {
				svc.Print(err)
				return
			}
			res, err := filepath.Glob("daily/*.jpg")
			if err != nil {
				svc.Print(err)
				return
			}
			for ; len(res) > 48; res = res[1:] {
				if err := os.Remove(res[0]); err != nil {
					svc.Print(err)
					return
				}
			}
			if err := jpg2gif("daily/*.jpg", "daily/daily-24h.gif", 48); err != nil {
				svc.Print(err)
			}
			if err := jpg2gif("daily/*.jpg", "daily/daily-12h.gif", 24); err != nil {
				svc.Print(err)
			}
			if err := jpg2gif("daily/*.jpg", "daily/daily-6h.gif", 12); err != nil {
				svc.Print(err)
			}
		}()
	}
	storms, err := stormAPI.GetStorms(t)
	if err != nil {
		svc.Print(err)
		return
	}
	var found, alert []storm.Data
	for _, i := range storms {
		storm, err := i.Data()
		if err != nil {
			svc.Print(err)
			continue
		}
		if affect, future := willAffect(storm, location, *radius); affect {
			found = append(found, storm)
			if future {
				alert = append(alert, storm)
				svc.Printf("Alerting storm %s(%s)", storm.ID, storm.Coordinates)
			} else {
				svc.Printf("Recording storm %s(%s)", storm.ID, storm.Coordinates)
			}
		}
	}
	if len(found) == 0 {
		return
	}
	if !isReport {
		for _, i := range found {
			b, err := coords{i.Coordinates}.screenshot(5.4, *quality, 3)
			if err != nil {
				svc.Print(err)
				return
			}
			dir := fmt.Sprintf("%s/%s", *path, i.ID)
			file := fmt.Sprintf("%s/%s.jpg", dir, time.Now().Format("20060102-1504"))
			if err := os.MkdirAll(dir, 0755); err != nil {
				svc.Print(err)
				continue
			}
			if err := os.WriteFile(file, b, 0644); err != nil {
				svc.Print(err)
				continue
			}
			if err := jpg2gif(dir+"/*.jpg", fmt.Sprintf("%s/%s.gif", dir, i.ID), 0); err != nil {
				svc.Print(err)
			}
		}
	}
	if len(alert) == 0 {
		return
	}
	var affectStorms, bodys []string
	var attachments []*mail.Attachment
	for i, storm := range alert {
		affectStorms = append(affectStorms, storm.Name)
		bodys = append(bodys, fmt.Sprintf("%s - %s<a href=%q><img src='cid:map%d'></a>", storm.Title, storm.Place, storm.URL, i))
		b, err := os.ReadFile(fmt.Sprintf("%s/%s/%[2]s.gif", *path, storm.ID))
		if err != nil {
			svc.Print(err)
			return
		}
		attachments = append(attachments, &mail.Attachment{
			Filename:  fmt.Sprintf("image%d.gif", i),
			Bytes:     b,
			ContentID: fmt.Sprintf("map%d", i),
		})
	}
	if hour := t.Hour(); isReport || ((hour == 6 || hour == 12 || hour == 21) && t.Minute() < 30) {
		sendMail(
			fmt.Sprintf("[Weather]Storm Alert - %s%s", strings.Join(affectStorms, "|"), timestamp()),
			strings.Join(bodys, "\n"),
			attachments,
			true,
		)
	}
}

func alertAQI(_ []weather.Day) (subject string, b strings.Builder) {
	current, err := aqiAPI.Realtime(aqi.China, *query)
	if err != nil {
		svc.Print(err)
		return
	}
	if level := current.AQI().Level().String(); level != "Excellent" && level != "Good" {
		subject = "[Weather]Air Quality Alert - " + level + timestamp()
		b.WriteString(aqiHTML(current))
	}
	return
}
