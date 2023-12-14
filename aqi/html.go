package aqi

import (
	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/weather/unit"
)

func CurrentHTML(current Current) html.HTML {
	div := html.Div()
	div.AppendChild(
		html.Div().Style("display:list-item;margin-left:15px").
			Content(
				current.AQI().Type(),
				":",
				html.Span().Style("padding:0 1em;color:white;background-color:"+current.AQI().Level().Color()).
					Contentf("%d %s", current.AQI().Value(), current.AQI().Level()),
			),
	)
	table := html.Table()
	table.AppendChild(html.Tbody())
	var tr *html.Element
	for i, p := range current.Pollutants() {
		if tr == nil {
			tr = html.NewElement("tr")
		}
		tr.AppendChild(
			html.NewElement("td").Content(p.Kind(), ":"),
			html.NewElement("td").Style("color:"+p.Level().Color()).
				Contentf("%s %s", unit.FormatFloat64(p.Value(), 2), p.Unit()),
		)
		if (i+1)%3 == 0 {
			table.AppendChild(tr)
			tr = nil
		}
	}
	if tr != nil {
		table.AppendChild(tr)
	}
	div.AppendChild(table)
	return div.HTML()
}
