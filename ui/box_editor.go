package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/cpgillem/garden-planner/geometry"
)

type BoxEditor struct {
	widget.BaseWidget

	// Internal widgets
	XLabel      *widget.Label
	YLabel      *widget.Label
	WidthLabel  *widget.Label
	HeightLabel *widget.Label

	XEntry      *widget.Entry
	YEntry      *widget.Entry
	WidthEntry  *widget.Entry
	HeightEntry *widget.Entry

	// Container
	container *fyne.Container

	// Reference to box
	Box *geometry.AxisAlignedBoundingBox

	// Reference to formatter
	Formatter *Formatter

	// Events
	OnUpdate func()
}

func NewBoxEditor(box *geometry.AxisAlignedBoundingBox, formatter *Formatter) *BoxEditor {
	boxEditor := &BoxEditor{
		Box:         box,
		XLabel:      widget.NewLabel("X"),
		YLabel:      widget.NewLabel("Y"),
		WidthLabel:  widget.NewLabel("Width"),
		HeightLabel: widget.NewLabel("Height"),
		XEntry:      widget.NewEntry(),
		YEntry:      widget.NewEntry(),
		WidthEntry:  widget.NewEntry(),
		HeightEntry: widget.NewEntry(),
		container:   container.NewGridWithColumns(2),
		Formatter:   formatter,
	}

	boxEditor.XEntry.OnSubmitted = func(s string) {
		if x, ok := formatter.ToDimensionUI(s); ok {
			box.Location.X = x
		}
		boxEditor.Refresh()
	}

	boxEditor.YEntry.OnSubmitted = func(s string) {
		if y, ok := formatter.ToDimensionUI(s); ok {
			box.Location.Y = y
		}
		boxEditor.Refresh()
	}

	boxEditor.WidthEntry.OnSubmitted = func(s string) {
		if width, ok := formatter.ToDimensionUI(s); ok {
			box.Size.X = width
		}
		boxEditor.Refresh()
	}

	boxEditor.HeightEntry.OnSubmitted = func(s string) {
		if height, ok := formatter.ToDimensionUI(s); ok {
			box.Size.Y = height
		}
		boxEditor.Refresh()
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

func (b *BoxEditor) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(b.container)
}

func (b *BoxEditor) Refresh() {
	b.XEntry.SetText(b.Formatter.FormatDimension(b.Box.Location.X))
	b.YEntry.SetText(b.Formatter.FormatDimension(b.Box.Location.Y))
	b.WidthEntry.SetText(b.Formatter.FormatDimension(b.Box.Size.X))
	b.HeightEntry.SetText(b.Formatter.FormatDimension(b.Box.Size.Y))

	b.OnUpdate()
}