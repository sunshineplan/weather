package unit

import (
	"encoding/json"
	"fmt"

	"github.com/sunshineplan/utils/html"
)

var _ json.Unmarshaler = new(UVIndex)

type UVIndex int

func (i *UVIndex) UnmarshalJSON(b []byte) error {
	var f float64
	if err := json.Unmarshal(b, &f); err != nil {
		return err
	}
	*i = UVIndex(f)
	return nil
}

func (i UVIndex) Color() string {
	if i <= 2 {
		return "#8cd600" // Green PMS 375
	} else if i <= 5 {
		return "#f9e814" // Yellow PMS 102
	} else if i <= 7 {
		return "#f77f00" // Orange PMS 151
	} else if i <= 10 {
		return "#ef2b2d" // Red PMS 032
	}
	return "#9663c4" // Purple PMS 265
}

func (i UVIndex) Risk() string {
	if i <= 2 {
		return "Low"
	} else if i <= 5 {
		return "Moderate"
	} else if i <= 7 {
		return "High"
	} else if i <= 10 {
		return "Very High"
	}
	return "Extreme"
}

func (i UVIndex) String() string {
	return fmt.Sprintf("%d(%s)", i, i.Risk())
}

func (i UVIndex) HTML() html.HTML {
	return html.Span().Style("color:" + i.Color()).Content(i.String()).HTML()
}
