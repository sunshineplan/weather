package unit

import (
	"strconv"
	"strings"
)

func FormatFloat64[T ~float64](f T, prec int) string {
	s := strconv.FormatFloat(float64(f), 'f', prec, 64)
	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
	}
	return strings.TrimSuffix(s, ".")
}
