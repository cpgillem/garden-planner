package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bcicen/go-units"
)

type Formatter struct {
	BaseUnit units.Unit
}

func NewFormatter(baseUnit units.Unit) *Formatter {
	return &Formatter{
		BaseUnit: baseUnit,
	}
}

func (formatter *Formatter) ToInteger(s string) (int, error) {
	return strconv.Atoi(s)
}

func (formatter *Formatter) ToDecimal(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err
}

// Parses a dimension with a quantity and a unit.
// Normally, the input would be a float, followed by a space, followed by a symbol.
// TODO: Accept units such as "
func (formatter *Formatter) ToDimension(s string) (units.Value, error) {
	zero := units.NewValue(0, formatter.BaseUnit)

	// Check format.
	firstSpace := strings.Index(strings.TrimSpace(s), " ")
	if firstSpace < 0 {
		return zero, NewDimensionError(s, "Dimension format: [quantity] [unit].")
	}

	// Parse out quantity and unit string.
	qty := strings.TrimSpace(s[:firstSpace])
	unitStr := strings.TrimSpace(s[firstSpace+1:])

	// Parse quantity to float.
	f, err := strconv.ParseFloat(qty, 32)
	if err != nil {
		return zero, NewDimensionError(s, "Quantity must be a number.")
	}

	// Parse unit string.
	unit, err := units.Find(unitStr)
	if err != nil {
		return zero, NewDimensionError(s, "Unrecognizable unit.")
	}

	// Create value.
	dim := units.NewValue(float64(f), unit)

	// Convert to base unit if necessary.
	converted := dim.MustConvert(formatter.BaseUnit)

	return converted, nil
}

func (formatter *Formatter) FormatInteger(i int) string {
	return strconv.Itoa(i)
}

func (formatter *Formatter) FormatDecimal(f float32) string {
	return fmt.Sprintf("%.3f", f)
}

func (formatter *Formatter) FormatDimension(f float32) string {
	val := units.NewValue(float64(f), formatter.BaseUnit)
	return val.Fmt(units.FmtOptions{
		Label:     true,
		Short:     true,
		Precision: 6,
	})
}

type DimensionError struct {
	input string
	msg   string
}

func NewDimensionError(input string, msg string) DimensionError {
	return DimensionError{
		input: input,
		msg:   msg,
	}
}

func (de DimensionError) Error() string {
	return de.msg
}
