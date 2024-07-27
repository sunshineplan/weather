package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sunshineplan/ai/prompt"
	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/aqi"
	"github.com/sunshineplan/weather/maps"
	"github.com/sunshineplan/weather/storm"
	"github.com/sunshineplan/weather/unit"
	"github.com/sunshineplan/weather/unit/coordinates"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	location     coordinates.Coordinates
	rainSnow     []weather.RainSnow
	tempRiseFall []weather.TempRiseFall

	alertMutex sync.Mutex
	zoomMutex  sync.Mutex

	isReport bool
)

func report(t time.Time) {
	isReport = true
	_, days, avg, aqi, err := getAll(*query, *days, aqiType, t, false)
	if err != nil {
		svc.Print(err)
		return
	}
	runAlert(days[1:], alertRainSnow)
	runAlert(days, alertTempRiseFall)
	alertStorm(t)
	sendMail(
		"[Weather]Daily Report"+timestamp(),
		fullHTML(*query, location, weather.Current{}, days, avg, aqi, t, true, *difference, "0")+
			html.Br().HTML()+
			imageHTML(mapAPI.URL(maps.Satellite, time.Time{}, location, mapOptions(*zoom)), "cid:attachment"),
		mail.TextHTML,
		attach6hGIF(),
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
		fullHTML(*query, location, weather.Current{}, days, avg, aqi, t, false, *difference, "0")+
			html.Br().HTML()+
			imageHTML(mapAPI.URL(maps.Satellite, time.Time{}, location, mapOptions(*zoom)), "cid:attachment"),
		mail.TextHTML,
		attach6hGIF(),
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
	q string, location coordinates.Coordinates,
	current weather.Current, days []weather.Day, avg weather.Day, currentAQI aqi.Current,
	t time.Time, highlight bool, diff float64, margin string,
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
			if highlight {
				res.AppendContent(i.HTML(t, t.Hour()))
			} else {
				res.AppendContent(i.HTML(t))
			}
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

func runAlert(days []weather.Day, fn func([]weather.Day) (string, *html.Element, []*mail.Attachment)) {
	if subject, body, attach := fn(days); subject != "" {
		svc.Print(subject)
		body := html.Div().Style("font-family:system-ui;margin:0").AppendChild(body).HTML()
		if attach != nil {
			body = body +
				html.Br().HTML() +
				imageHTML(mapAPI.URL(maps.Satellite, time.Time{}, location, mapOptions(*zoom)), "cid:attachment")
		}
		sendMail(
			subject,
			body,
			mail.TextHTML,
			attach,
			false,
		)
	}
}

func isRainSnow(hour int, hours []weather.Hour) bool {
	return slices.IndexFunc(hours, func(h weather.Hour) bool {
		return h.TimeEpoch.Time().Hour() == hour && h.Precip > 0
	}) >= 0
}

func alertRainSnow(days []weather.Day) (subject string, body *html.Element, attch []*mail.Attachment) {
	if len(rainSnow) > 0 {
		if rainSnow[0].IsExpired() {
			rainSnow = rainSnow[1:]
		}
	}

	body = html.Background()
	if res := weather.WillRainSnow(days); len(res) > 0 {
		now := time.Now()
		hour := now.Hour()
		for index, i := range res {
			if index == 0 {
				start := i.Start()
				var isRainNow bool
				if start.Date == now.Format("2006-01-02") && isRainSnow(hour, start.Hours) {
					isRainNow = true
				}
				if len(rainSnow) == 0 ||
					rainSnow[0].Start().Date != start.Date ||
					rainSnow[0].Duration() != i.Duration() {
					subject = "[Weather]Rain Snow Alert - " + start.Date + timestamp()
					if isRainNow {
						attch = attachLast()
					}
				} else if isRainNow {
					subject = "[Weather]Rain Snow Alert - Today" + timestamp()
					body.AppendContent(start.DateInfoHTML(), start.PrecipitationHTML(hour))
					for index, n := 1, len(i.Days()); index < n && index < 3; index++ {
						body.AppendContent(
							html.Br(),
							i.Days()[index].DateInfoHTML(),
							i.Days()[index].PrecipitationHTML(),
						)
					}
					attch = attachLast()
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

func alertTempRiseFall(days []weather.Day) (subject string, body *html.Element, _ []*mail.Attachment) {
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

var aqiStandard int

func alertAQI(_ []weather.Day) (subject string, body *html.Element, _ []*mail.Attachment) {
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

func alertStorm(t time.Time) {
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
		if storm.Season == "" || storm.No == 0 {
			continue
		}
		if affect, future := storm.Affect(location, *radius); affect {
			found = append(found, storm)
			coords := storm.Coordinates(t)
			if coords == nil {
				coords = storm.Track[len(storm.Track)-1].Coordinates()
			}
			if future {
				alert = append(alert, storm)
				svc.Printf("Alerting storm %s(%s)", storm.ID, coords)
			} else {
				svc.Printf("Recording storm %s(%s)", storm.ID, coords)
			}
		}
	}
	if len(found) == 0 {
		return
	}
	if !isReport {
		updateStorm(found)
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
		if attachment := attachStorm(i, storm); attachment != nil {
			attachments = append(attachments, attachment)
		}
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
