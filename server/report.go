package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/aqi"
	"github.com/sunshineplan/weather/storm"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	location     *coords
	rainSnow     []weather.RainSnow
	tempRiseFall []weather.TempRiseFall

	alertMutex sync.Mutex
	zoomMutex  sync.Mutex
)

func report(t time.Time) {
	days, avg, aqi, err := getAll(*query, *days, aqi.China, t)
	if err != nil {
		svc.Print(err)
		return
	}
	runAlert(days[1:], alertRainSnow)
	runAlert(days, alertTempRiseFall)
	zoomEarth(t, true)
	sendMail(
		"[Weather]Daily Report"+timestamp(),
		fullHTML(fmt.Sprintf("%s(%s)", *query, location), days, avg, aqi, t, *difference)+
			html.Br().HTML()+
			imageHTML(location.url(*zoom), "cid:attachment"),
		attachment("daily/daily-12h.gif"),
		true,
	)
}

func daily(t time.Time) {
	svc.Print("Start sending daily report...")
	days, avg, aqi, err := getAll(*query, *days, aqi.China, t)
	if err != nil {
		svc.Print(err)
		return
	}
	sendMail(
		"[Weather]Daily Report"+timestamp(),
		fullHTML(fmt.Sprintf("%s(%s)", *query, location), days, avg, aqi, t, *difference)+imageHTML(location.url(*zoom), "cid:attachment"),
		attachment("daily/daily-12h.gif"),
		true,
	)
}

func fullHTML(q string, days []weather.Day, avg weather.Day, currentAQI aqi.Current, t time.Time, diff float64) html.HTML {
	div := html.Div().Style("font-family:system-ui;margin:0")
	div.AppendContent(
		html.Span().Style("display:list-item;list-style:circle;margin-left:1em;font-size:1.5em").
			Contentf("Weather of %s", cases.Title(language.English).String(q)),
		days[1],
		html.Br(),
		aqi.CurrentHTML(currentAQI),
		html.Br(),
		html.Div().AppendChild(
			html.Span().Style("display:list-item;margin-left:15px").Content("Compared with Yesterday"),
			html.Table().AppendChild(
				html.Tbody().Content(
					days[0].TemperatureHTML(),
					weather.NewTempRiseFall(days[1], days[0], 0).DiffInfoHTML(),
				),
			),
		),
		html.Br(),
	)
	if avg.Date != "" {
		div.AppendChild(
			html.Div().AppendChild(
				html.Span().Style("display:list-item;margin-left:15px").
					Contentf("Historical Average Temperature of %s", t.Format("01-02")),
				html.Table().AppendChild(
					html.Tbody().Content(
						avg.TemperatureHTML(),
						weather.NewTempRiseFall(days[1], avg, 0).DiffInfoHTML(),
					),
				),
			),
			html.Br(),
		)
	}
	div.AppendContent(forecastHTML(days[1:]))

	if rainSnow := weather.WillRainSnow(days[1:]); len(rainSnow) > 0 {
		res := html.Div()
		for index, i := range rainSnow {
			res.AppendContent(i.HTML(t))
			if index < len(rainSnow)-1 {
				res.AppendChild(html.Br())
			}
		}
		div.AppendChild(
			html.Br(),
			html.Div().AppendChild(
				html.Span().Style("display:list-item;margin-left:15px").Content("Recent Rain Snow Alert"),
				html.Div().AppendChild(res),
			),
		)
	} else {
		div.AppendContent(
			html.Br(),
			"No Rain Snow Alert.",
		)
	}
	if tempRiseFall := weather.WillTempRiseFall(days, diff); len(tempRiseFall) > 0 {
		res := html.Div()
		for index, i := range tempRiseFall {
			res.AppendContent(i)
			if index < len(tempRiseFall)-1 {
				res.AppendChild(html.Br())
			}
		}
		div.AppendChild(
			html.Br(),
			html.Div().AppendChild(
				html.Span().Style("display:list-item;margin-left:15px").Content("Recent Temperature Alert"),
				html.Div().AppendChild(res),
			),
		)
	} else {
		div.AppendContent(
			html.Br(),
			"No Temperature Alert.",
		)
	}
	return div.HTML()
}

func alert(t time.Time) {
	svc.Print("Start alerting...")
	days, err := getWeather(*query, *days, t)
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
		runAlert(days[1:], alertRainSnow)
	}()
	go func() {
		defer wg.Done()
		runAlert(days, alertTempRiseFall)
	}()
	go func() {
		defer wg.Done()
		runAlert(nil, alertAQI)
	}()
	wg.Wait()
}

func runAlert(days []weather.Day, fn func([]weather.Day) (string, *html.Element)) {
	if subject, body := fn(days); subject != "" {
		svc.Print(subject)
		sendMail(subject, html.Div().Style("font-family:system-ui;margin:0").AppendChild(body).HTML(), nil, false)
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

func alertRainSnow(days []weather.Day) (subject string, body *html.Element) {
	if len(rainSnow) > 0 {
		if rainSnow[0].IsExpired() {
			rainSnow = rainSnow[1:]
		}
	}

	body = html.Background()
	if res := weather.WillRainSnow(days); len(res) > 0 {
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
					body.AppendContent(start.DateInfoHTML(), start.PrecipitationHTML(hour))
					for index, n := 1, len(i.Days()); index < n && index < 3; index++ {
						body.AppendContent(
							html.Br(),
							i.Days()[index].DateInfoHTML(),
							i.Days()[index].PrecipitationHTML(),
						)
					}
					return
				}
			}
			body.AppendContent(i.HTML(now, hour))
			if index < len(res)-1 {
				body.AppendChild(html.Br())
			}
		}
		rainSnow = res
	} else if len(rainSnow) > 0 {
		subject = "[Weather]Rain Snow Alert - Canceled" + timestamp()
		body.Content("No more rain snow")
		rainSnow = nil
	}
	return
}

func alertTempRiseFall(days []weather.Day) (subject string, body *html.Element) {
	if len(tempRiseFall) > 0 {
		if tempRiseFall[0].IsExpired() {
			tempRiseFall = tempRiseFall[1:]
		}
	}

	body = html.Background()
	if res := weather.WillTempRiseFall(days, *difference); len(res) > 0 {
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
			body.AppendContent(i)
			if index < len(res)-1 {
				body.AppendChild(html.Br())
			}
		}
		tempRiseFall = res
	} else if len(tempRiseFall) > 0 {
		if tempRiseFall[0].IsRise() {
			subject = "[Weather]Temperature Rise Alert - Canceled" + timestamp()
			body.Content("No more temperature rise")
		} else {
			subject = "[Weather]Temperature Fall Alert - Canceled" + timestamp()
			body.Content("No more temperature fall")
		}
		tempRiseFall = nil
	}
	return
}

func forecastHTML(days []weather.Day) html.HTML {
	if len(days) > 10 {
		days = days[:10]
	}
	div := html.Div()
	div.AppendChild(html.Span().Style("display:list-item;margin-left:15px").Content("Forecast"))
	table := html.Table().Attribute("border", "1").Attribute("cellspacing", "0")
	table.AppendChild(
		html.Thead().AppendChild(
			html.Tr(
				html.Th("Date").Colspan(2),
				html.Th("Max"),
				html.Th("Min"),
				html.Th("FLMax"),
				html.Th("FLMin"),
				html.Th("Rain%"),
			),
		),
	)
	tbody := html.Tbody()
	for _, day := range days {
		tbody.AppendChild(
			html.Tr(
				html.Td(day.DateInfo(false)[11:]),
				html.Td(day.Condition.Img(day.Icon)),
				html.Td(day.TempMax),
				html.Td(day.TempMin),
				html.Td(day.FeelsLikeMax),
				html.Td(day.FeelsLikeMin),
				html.Td(day.PrecipProb),
			))
	}
	table.AppendChild(tbody)
	div.AppendChild(table)
	return div.HTML()
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
		bodys = append(bodys, html.Background().
			Contentf("%s - %s", storm.Title, storm.Place).
			AppendChild(html.A().Href(storm.URL).AppendChild(html.Img().Src("cid:map"+strconv.Itoa(i)))).String(),
		)
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
			html.HTML(strings.Join(bodys, "\n")),
			attachments,
			true,
		)
	}
}

func alertAQI(_ []weather.Day) (subject string, body *html.Element) {
	current, err := aqiAPI.Realtime(aqi.China, *query)
	if err != nil {
		svc.Print(err)
		return
	}
	if level := current.AQI().Level().String(); level != "Excellent" && level != "Good" {
		subject = "[Weather]Air Quality Alert - " + level + timestamp()
		body = html.Background().Content(aqi.CurrentHTML(current))
	}
	return
}
