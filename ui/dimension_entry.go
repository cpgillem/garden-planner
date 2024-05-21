package ui

import (
	"fyne.io/fyne/widget"
)

// Extends the basic text box and performs functions to display/edit dimensions.
type DimensionEntry struct {
	widget.Entry
}

func NewUnitEntry() *DimensionEntry {
	dimensionEntry := &DimensionEntry{}
	dimensionEntry.ExtendBaseWidget(dimensionEntry)

	return dimensionEntry
}

func (e *DimensionEntry) TypedRune(r rune) {
	// Only allow numeric input.
	if (r >= '0' && r <= '9') || r == '.' {
		e.Entry.TypedRune(r)
	}
}

func (e *DimensionEntry) OnSubmitted() {

}
