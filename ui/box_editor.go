package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/bcicen/go-units"
	"github.com/cpgillem/garden-planner/geometry"
)

type BoxEditor struct {
	widget.BaseWidget

	// Internal widgets
	XLabel      *widget.Label
	YLabel      *widget.Label
	WidthLabel  *widget.Label
	HeightLabel *widget.Label

	XEntry      *DimensionEntry
	YEntry      *DimensionEntry
	WidthEntry  *DimensionEntry
	HeightEntry *DimensionEntry

	// Container
	container *fyne.Container

	// Reference to formatter
	Formatter *DimensionFormatter

	// Events
	OnSubmitted func(newBox geometry.Box)
}

func NewBoxEditor(initialBox geometry.Box, baseUnit units.Unit, formatter *DimensionFormatter) *BoxEditor {
	boxEditor := &BoxEditor{
		XLabel:      widget.NewLabel("X"),
		YLabel:      widget.NewLabel("Y"),
		WidthLabel:  widget.NewLabel("Width"),
		HeightLabel: widget.NewLabel("Height"),
		XEntry:      NewDimensionEntry(units.NewValue(float64(initialBox.GetX()), baseUnit), formatter),
		YEntry:      NewDimensionEntry(units.NewValue(float64(initialBox.GetY()), baseUnit), formatter),
		WidthEntry:  NewDimensionEntry(units.NewValue(float64(initialBox.GetWidth()), baseUnit), formatter),
		HeightEntry: NewDimensionEntry(units.NewValue(float64(initialBox.GetHeight()), baseUnit), formatter),
		container:   container.New(layout.NewFormLayout()),
		Formatter:   formatter,
		OnSubmitted: func(newBox geometry.Box) {},
	}

	boxEditor.XEntry.OnValueChanged = func(val units.Value) {
		boxEditor.UpdateBox()
	}

	boxEditor.YEntry.OnValueChanged = func(val units.Value) {
		boxEditor.UpdateBox()
	}

	boxEditor.WidthEntry.OnValueChanged = func(val units.Value) {
		boxEditor.UpdateBox()
	}

	boxEditor.HeightEntry.OnValueChanged = func(val units.Value) {
		boxEditor.UpdateBox()
	}

	boxEditor.container.Add(boxEditor.XLabel)
	boxEditor.container.Add(boxEditor.XEntry)
	boxEditor.container.Add(boxEditor.YLabel)
	boxEditor.container.Add(boxEditor.YEntry)
	boxEditor.container.Add(boxEditor.WidthLabel)
	boxEditor.container.Add(boxEditor.WidthEntry)
	boxEditor.container.Add(boxEditor.HeightLabel)
	boxEditor.container.Add(boxEditor.HeightEntry)

	boxEditor.ExtendBaseWidget(boxEditor)

	return boxEditor
}

// Called when one of the entries is successfully submitted.
func (b *BoxEditor) UpdateBox() {
	newBox := geometry.NewBox(
		float32(b.XEntry.GetValue().Float()),
		float32(b.YEntry.GetValue().Float()),
		float32(b.WidthEntry.GetValue().Float()),
		float32(b.HeightEntry.GetValue().Float()),
	)
	b.OnSubmitted(newBox)
}

func (b *BoxEditor) SetBox(box geometry.Box) {
	b.XEntry.SetValue(units.NewValue(float64(box.GetX()), b.XEntry.baseUnit))
	b.YEntry.SetValue(units.NewValue(float64(box.GetX()), b.YEntry.baseUnit))
	b.WidthEntry.SetValue(units.NewValue(float64(box.GetWidth()), b.WidthEntry.baseUnit))
	b.HeightEntry.SetValue(units.NewValue(float64(box.GetHeight()), b.HeightEntry.baseUnit))
}

func (b *BoxEditor) CreateRenderer() fyne.WidgetRenderer {
	renderer := widget.NewSimpleRenderer(b.container)
	return renderer

}
