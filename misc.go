package weather

import (
	"fmt"
	"strings"
	"time"
)

func fmtDuration(d time.Duration) string {
	var b strings.Builder
	if d >= time.Hour*24 {
		fmt.Fprintf(&b, "%dd", d/24/time.Hour)
		d -= (d / time.Hour / 24) * (time.Hour * 24)
	}
	if d >= time.Hour {
		fmt.Fprintf(&b, "%dh", d/time.Hour)
		d -= d / time.Hour * time.Hour
	}
	if d >= time.Minute {
		fmt.Fprintf(&b, "%dm", d/time.Minute)
		d -= d / time.Minute * time.Minute
	}
	if s := d / time.Second; s > 0 {
		fmt.Fprintf(&b, "%ds", s)
	} else if s == 0 && b.Len() == 0 {
		fmt.Fprint(&b, "0s")
	}
	return b.String()
}
