package ui

import (
	"fyne.io/fyne/v2/widget"
	"github.com/bcicen/go-units"
)

// Extends the basic text box and performs functions to display/edit dimensions.
type DimensionEntry struct {
	widget.Entry

	// Internal State
	baseUnit           units.Unit
	value              units.Value
	dimensionFormatter *DimensionFormatter

	// Events
	OnValueChanged   func(val units.Value)
	OnDimensionError func(err error)
}

// Infers the base unit from the given unit.
func NewDimensionEntry(value units.Value, dimensionFormatter *DimensionFormatter) *DimensionEntry {
	dimensionEntry := &DimensionEntry{
		baseUnit:           value.Unit(),
		value:              value,
		dimensionFormatter: dimensionFormatter,
		OnValueChanged:     func(val units.Value) {},
		OnDimensionError:   func(err error) {},
	}
	dimensionEntry.ExtendBaseWidget(dimensionEntry)

	// When submitted, update the value and fire a new event.
	dimensionEntry.OnSubmitted = func(s string) {
		value, err := dimensionEntry.dimensionFormatter.ToDimension(s)
		if err != nil {
			// Revert the text content back to the existing value.
			dimensionEntry.Reset()

			// Fire an input error event, usually through showing a dialog.
			dimensionEntry.OnDimensionError(err)
			return
		}

		dimensionEntry.SetValue(value)
		dimensionEntry.OnValueChanged(value)
	}

	dimensionEntry.Reset()
	return dimensionEntry
}

// Sets the text of the entry widget to the current value. Does not fire an event.
func (e *DimensionEntry) Reset() {
	e.SetText(e.dimensionFormatter.FormatDimension(e.value))
}

func (e *DimensionEntry) SetValue(value units.Value) {
	e.value = value
	e.Reset()
}

// Parses the value out of the dimension entry's text. If it happens to be invalid, throws an error.
func (e *DimensionEntry) GetValue() (units.Value, error) {
	value, err := e.dimensionFormatter.ToDimensionBaseUnit(e.Text, e.baseUnit)
	if err != nil {
		z := units.NewValue(0, e.baseUnit)

		// Reset the value if the text happens to be incorrect.
		e.SetValue(z)
		return z, err
	}

	return value, nil
}

// Tries to parse the dimension entry's text, and returns 0 if impossible.
func (e *DimensionEntry) MustGetValue() units.Value {
	value, _ := e.GetValue()
	return value
}

func (e *DimensionEntry) Submitted() {

}
