package ui

import (
	"testing"

	"github.com/bcicen/go-units"
)

// A unit test in more ways than one.
func TestToDimension(t *testing.T) {
	// Happy path cases
	cases := []struct {
		in         string
		inBaseUnit units.Unit
		wantF      float64
		wantUnit   string
	}{
		{"1 in", units.Inch, 1, "inch"},
	}

	for _, c := range cases {
		formatter := NewFormatter(c.inBaseUnit)
		got, err := formatter.ToDimension(c.in)
		if err != nil {
			t.Errorf("ToDimension(%q): %q", c.in, err.Error())
		}

		if got.Float() != c.wantF || got.Unit().Name != c.wantUnit {
			t.Errorf("ToDimension(%q) == %f, %q; want %f, %q", c.in, got.Float(), got.Unit().Name, c.wantF, c.wantUnit)
		}
	}

	// Error cases
	errCases := []struct {
		in         string
		inBaseUnit units.Unit
		wantMsg    string
	}{
		{"1", units.Inch, "Dimension format: [quantity] [unit]."},
		{"A in", units.Inch, "Quantity must be a number."},
		{"2 horses", units.Inch, "Unrecognizable unit."},
	}

	for _, c := range errCases {
		formatter := NewFormatter(c.inBaseUnit)
		_, err := formatter.ToDimension(c.in)
		if err.Error() != c.wantMsg {
			t.Errorf("ToDimension(%q); got %q, want %q", c.in, err.Error(), c.wantMsg)
		}
	}
}
