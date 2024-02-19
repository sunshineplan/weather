package wind

import "testing"

func TestDegree(t *testing.T) {
	for i, testcase := range []struct {
		Degree    Degree
		direction Direction
	}{
		{N, "N"},
		{difference, "NbE"},
		{NbE, "NbE"},
		{NbE + difference, "NNE"},
		{N - difference, "N"},
		{NbW + difference, "N"},
	} {
		if dir := testcase.Degree.Direction(); dir != testcase.direction {
			t.Errorf("#%d: expected %q; got %q", i, testcase.direction, dir)
		}
	}
}
