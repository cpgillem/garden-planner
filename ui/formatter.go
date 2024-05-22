package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bcicen/go-units"
)

var AnyUnit units.Unit = units.NewUnit("Any", "")

type DimensionFormatter struct {
	fmtOptions units.FmtOptions
}

func NewFormatter() *DimensionFormatter {
	return &DimensionFormatter{
		fmtOptions: units.FmtOptions{
			Label:     true,
			Short:     true,
			Precision: 6,
		},
	}
}

func (formatter *DimensionFormatter) ToInteger(s string) (int, error) {
	return strconv.Atoi(s)
}

func (formatter *DimensionFormatter) ToDecimal(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err
}

// Parses a dimension with a quantity and a unit.
// Normally, the input would be a float, followed by a space, followed by a symbol.
// TODO: Accept units such as "
func (formatter *DimensionFormatter) ToDimension(s string) (units.Value, error) {
	zero := units.NewValue(0, AnyUnit)

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
	value := units.NewValue(float64(f), unit)

	return value, nil
}

// Useful if you would like to accept any dimension and convert it to a default unit.
func (formatter *DimensionFormatter) ToDimensionBaseUnit(s string, baseUnit units.Unit) (units.Value, error) {
	value, err := formatter.ToDimension(s)

	if err != nil {
		return value, err
	}

	if value.Unit().Name != baseUnit.Name {
		return value.MustConvert(baseUnit), nil
	}

	return value, nil
}

func (formatter *DimensionFormatter) FormatInteger(i int) string {
	return strconv.Itoa(i)
}

func (formatter *DimensionFormatter) FormatDecimal(f float32) string {
	return fmt.Sprintf("%.3f", f)
}

func (formatter *DimensionFormatter) FormatDimension(value units.Value) string {
	return value.Fmt(formatter.fmtOptions)
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
