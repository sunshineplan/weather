package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sunshineplan/ai/prompt"
	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/aqi"
	"github.com/sunshineplan/weather/storm"
	"github.com/sunshineplan/weather/unit"
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
	_, days, avg, aqi, err := getAll(*query, *days, aqiType, t, false)
	if err != nil {
		svc.Print(err)
		return
	}
	runAlert(days[1:], alertRainSnow)
	runAlert(days, alertTempRiseFall)
	zoomEarth(t, true)
	sendMail(
		"[Weather]Daily Report"+timestamp(),
		fullHTML(*query, location, weather.Current{}, days, avg, aqi, t, *difference, "0")+
			html.Br().HTML()+
			imageHTML(location.url(*zoom), "cid:attachment"),
		mail.TextHTML,
		attachment("daily/daily-12h.gif"),
		true,
	)
	if chatbot != nil {
		res, err := aiReport(*query, days, t)
		if err != nil {
			svc.Print(err)
			return
		}
		sendMail("[Weather]Daily AI Report"+timestamp(), res, mail.TextPlain, nil, true)
	}
}

func daily(t time.Time) {
	svc.Print("Start sending daily report...")
	_, days, avg, aqi, err := getAll(*query, *days, aqiType, t, false)
	if err != nil {
		svc.Print(err)
		return
	}
	go sendMail(
		"[Weather]Daily Report"+timestamp(),
		fullHTML(*query, location, weather.Current{}, days, avg, aqi, t, *difference, "0")+
			html.Br().HTML()+
			imageHTML(location.url(*zoom), "cid:attachment"),
		mail.TextHTML,
		attachment("daily/daily-12h.gif"),
		true,
	)
	if chatbot != nil {
		svc.Print("Start sending daily AI report...")
		res, err := aiReport(*query, days, t)
		if err != nil {
			svc.Print(err)
			return
		}
		go sendMail("[Weather]Daily AI Report"+timestamp(), res, mail.TextPlain, nil, true)
	}
}

func fullHTML(
	q string, location *coords,
	current weather.Current, days []weather.Day, avg weather.Day, currentAQI aqi.Current,
	t time.Time, diff float64, margin string,
) html.HTML {
	div := html.Div().Style("font-family:system-ui;margin:" + margin)
	div.AppendChild(
		html.Span().Style("display:list-item;list-style:circle;margin-left:1em").AppendChild(
			html.Span().Style("font-size:1.5em").Contentf("Weather of %s", cases.Title(language.English).String(q)),
			html.Span().Style("display:inline-block;font-family:monospace;font-size:.75em;text-align:right;margin-left:1em").
				AppendChild(
					html.Span().Style("display:block").Content(location.Latitude()),
					html.Span().Style("display:block").Content(location.Longitude()),
				),
		),
	)
	if current.Datetime != "" {
		now := html.Div().AppendContent(
			html.Span().Style("display:list-item;margin-left:15px").
				Contentf("Current(%s %s)", t.Format("2006-01-02"), t.Weekday().String()[:3]),
			current,
		)
		if currentAQI != nil {
			now.AppendContent(aqi.CurrentHTML(currentAQI))
		}
		div.AppendChild(now.AppendContent(html.Br()))
	}
	div.AppendContent(
		html.Div().AppendContent(
			html.Span().Style("display:list-item;margin-left:15px").Content("Today"),
			days[1],
		),
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
	div.AppendContent(forecastHTML(days[1:], current.Datetime != ""))

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
			html.Div().Content("No Rain Snow Alert."),
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
			html.Div().Content("No Temperature Alert."),
		)
	}
	return div.HTML()
}

func aiReport(query string, days []weather.Day, t time.Time) (string, error) {
	var s []string
	for _, i := range days {
		day := i.String()
		_, precipHours := i.PrecipHours()
		if len(precipHours) != 0 {
			day += "\nPrecipHours: " + strings.Join(precipHours, ", ")
		}
		if len(i.PrecipType) > 0 {
			day += "\nPrecipType: " + strings.Join(i.PrecipType, ", ")
		}
		s = append(s, day)
	}
	q := prompt.New(fmt.Sprintf(`Today's date is %s and location is %s.
Based on the provided weather data, please:
1. Compare today's weather with yesterday's.
2. Analyze temperature trends.
3. Analyze precipitation trends, including precipitation hour, amount, probability, and coverage.
4. Generate a weather forecast.
The output result language is required to be the location language.`, t.Format("2006-01-02"), query))
	c, _, err := q.Execute(chatbot, s, "")
	if err != nil {
		return "", err
	}
	res := <-c
	if res.Error != nil {
		return "", res.Error
	}
	return res.Result[0], nil
}

func alert(t time.Time) {
	svc.Print("Start alerting...")
	_, days, err := getWeather(*query, *days, t, false)
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
		sendMail(subject, html.Div().Style("font-family:system-ui;margin:0").AppendChild(body).HTML(), mail.TextHTML, nil, false)
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

func forecastHTML(days []weather.Day, dir bool) html.HTML {
	div := html.Div()
	div.AppendChild(html.Span().Style("display:list-item;margin-left:15px").Content("Forecast"))
	table := html.Table().Attribute("border", "1").Attribute("cellspacing", "0")
	th := []*html.TableCell{
		html.Th("Date").Colspan(2),
		html.Th("Temp"),
		html.Th("FeelsLike"),
		html.Th("Rain"),
		html.Th("Wind"),
	}
	if dir {
		th = append(th, html.Th("Dir"))
	}
	table.AppendChild(html.Thead().AppendChild(html.Tr(th...)))
	tbody := html.Tbody()
	for _, day := range days {
		td := []*html.TableCell{
			html.Td(html.Background().AppendChild(
				html.Span().Content(day.DateEpoch.Time().Format("01-02")),
				html.Span().Content(day.DateEpoch.Weekday()),
			)).Style("display:grid;text-align:center;font-size:.9em"),
			html.Td(day.Condition.Img(day.Icon)),
			html.Td(html.HTML(strings.ReplaceAll(string(day.TempMax.HTML()), "°C", "")) + " / " + day.TempMin.HTML()),
			html.Td(html.HTML(strings.ReplaceAll(string(day.FeelsLikeMax.HTML()), "°C", "")) + " / " + day.FeelsLikeMin.HTML()),
			html.Td(html.Background().AppendChild(
				html.Span().Contentf("%gmm", day.Precip),
				html.Span().Content(day.PrecipProb),
			)).Style("display:grid;text-align:center;font-size:.9em"),
			html.Td(html.Span().Style("color:"+day.WindSpeed.ForceColor()).Contentf("%sm/s", unit.FormatFloat64(day.WindSpeed.MPS(), 1))),
		}
		if dir {
			td = append(td, html.Td(html.Div().Style("display:flex;justify-content:center").Content(day.WindDir)))
		}
		tbody.AppendChild(html.Tr(td...))
	}
	return div.AppendChild(table.AppendChild(tbody)).HTML()
}

func zoomEarth(t time.Time, isReport bool) {
	if !isReport {
		go func() {
			svc.Print("Start saving satellite map...")
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
			mail.TextHTML,
			attachments,
			true,
		)
	}
}

var aqiStandard int

func alertAQI(_ []weather.Day) (subject string, body *html.Element) {
	current, err := aqiAPI.Realtime(aqiType, *query)
	if err != nil {
		svc.Print(err)
		return
	}
	if index := current.AQI(); index.Value() >= aqiStandard {
		subject = "[Weather]Air Quality Alert - " + index.Level().String() + timestamp()
		body = html.Background().Content(aqi.CurrentHTML(current))
	}
	return
}
