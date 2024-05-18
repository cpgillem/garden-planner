package ui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

type Formatter struct {
	window *fyne.Window
}

func NewFormatter(window *fyne.Window) *Formatter {
	return &Formatter{
		window: window,
	}
}

func (formatter *Formatter) ToInteger(s string) (int, error) {
	return strconv.Atoi(s)
}

// All-in-one function to convert to an integer and show a dialog if there is an error.
// If there is an error, the value will be nil.
func (formatter *Formatter) ToIntegerUI(s string) (int, bool) {
	i, err := formatter.ToInteger(s)
	if err != nil {
		formatter.IntegerErrorDialog()
		return 0, false
	}

	return i, true
}

func (formatter *Formatter) ToDecimal(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err
}

func (formatter *Formatter) ToDecimalUI(s string) (float32, bool) {
	f, err := formatter.ToDecimal(s)
	if err != nil {
		formatter.DecimalErrorDialog()
		return 0, false
	}

	return f, true
}

// TODO: This will become a more advanced type.
func (formatter *Formatter) ToDimension(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err
}

func (formatter *Formatter) ToDimensionUI(s string) (float32, bool) {
	f, err := formatter.ToDecimal(s)
	if err != nil {
		formatter.DimensionErrorDialog()
		return 0, false
	}

	return f, true
}

func (formatter *Formatter) FormatInteger(i int) string {
	return strconv.Itoa(i)
}

func (formatter *Formatter) FormatDecimal(f float32) string {
	return fmt.Sprintf("%.3f", f)
}

func (formatter *Formatter) FormatDimension(f float32) string {
	return fmt.Sprintf("%.3f", f)
}

func (formatter *Formatter) IntegerErrorDialog() {
	dialog.ShowInformation("Error", "Invalid integer.", *formatter.window)
}

func (formatter *Formatter) DecimalErrorDialog() {
	dialog.ShowInformation("Error", "Invalid decimal.", *formatter.window)
}

func (formatter *Formatter) DimensionErrorDialog() {
	dialog.ShowInformation("Error", "Invalid dimension.", *formatter.window)
}
