package aqi

import (
	"fmt"
	"strings"

	"github.com/sunshineplan/weather/unit"
)

func CurrentHTML(current Current) string {
	var b strings.Builder
	fmt.Fprint(&b, "<div>")
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
	fmt.Fprint(&b, "</tbody></table></div>")
	return b.String()
}
