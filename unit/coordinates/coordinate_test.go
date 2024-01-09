package coordinates

import "testing"

func TestDMS(t *testing.T) {
	for _, testcase := range []struct {
		f      float64
		lat    bool
		expect string
	}{
		{1, true, "1°00′00″N"},
		{1.1, false, "1°06′00″E"},
		{-1.11, true, "1°06′36″S"},
		{-1.111, false, "1°06′40″W"}, //1°06′39.6″
		{121.41, false, "121°24′36″E"},
	} {
		dms := newDMS(testcase.f)
		if res := dms.str(testcase.lat); res != testcase.expect {
			t.Errorf("expected %q; got %q", testcase.expect, res)
		}
	}
}
