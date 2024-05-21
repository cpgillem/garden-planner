package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
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
	Box *geometry.Box

	// Reference to formatter
	Formatter *Formatter

	OnSubmitted func(boxDelta geometry.Box)
}

func NewBoxEditor(box *geometry.Box, formatter *Formatter) *BoxEditor {
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
		container:   container.New(layout.NewFormLayout()),
		Formatter:   formatter,
		OnSubmitted: func(boxDelta geometry.Box) {},
	}

	boxEditor.XEntry.OnSubmitted = func(s string) {
		x, err := formatter.ToDimension(s)
		if err == nil {
			box.Location.X = float32(x.Float())
		}
		boxEditor.Refresh()
		boxEditor.OnSubmitted(geometry.NewBox(float32(x.Float()), 0, 0, 0))
	}

	boxEditor.YEntry.OnSubmitted = func(s string) {
		y, err := formatter.ToDimension(s)
		if err == nil {
			box.Location.Y = float32(y.Float())
		}
		boxEditor.Refresh()
		boxEditor.OnSubmitted(geometry.NewBox(0, float32(y.Float()), 0, 0))
	}

	boxEditor.WidthEntry.OnSubmitted = func(s string) {
		width, err := formatter.ToDimension(s)
		if err == nil {
			box.Location.X = float32(width.Float())
		}
		boxEditor.Refresh()
		boxEditor.OnSubmitted(geometry.NewBox(0, 0, float32(width.Float()), 0))
	}

	boxEditor.HeightEntry.OnSubmitted = func(s string) {
		height, err := formatter.ToDimension(s)
		if err == nil {
			box.Location.X = float32(height.Float())
		}
		boxEditor.Refresh()
		boxEditor.OnSubmitted(geometry.NewBox(0, 0, 0, float32(height.Float())))
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
	renderer := widget.NewSimpleRenderer(b.container)
	return renderer

}

func (b *BoxEditor) Refresh() {
	b.XEntry.SetText(b.Formatter.FormatDimension(b.Box.Location.X))
	b.YEntry.SetText(b.Formatter.FormatDimension(b.Box.Location.Y))
	b.WidthEntry.SetText(b.Formatter.FormatDimension(b.Box.Size.X))
	b.HeightEntry.SetText(b.Formatter.FormatDimension(b.Box.Size.Y))
}
