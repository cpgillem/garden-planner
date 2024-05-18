package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

type Handle struct {
	widget.BaseWidget

	// Drawing
	Circle canvas.Circle

	// Events
	OnDragged func(d *fyne.DragEvent)
	OnDragEnd func()
}

func NewHandle() *Handle {
	h := Handle{
		Circle:    *canvas.NewCircle(colornames.White),
		OnDragged: func(e *fyne.DragEvent) {},
		OnDragEnd: func() {},
	}

	h.Circle.StrokeColor = colornames.Black
	h.Circle.StrokeWidth = 1

	h.ExtendBaseWidget(&h)
	return &h
}

func (h *Handle) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewStack(&h.Circle)
	return widget.NewSimpleRenderer(c)
}

func (h *Handle) Dragged(e *fyne.DragEvent) {
	h.OnDragged(e)
}

func (h *Handle) DragEnd() {
	h.OnDragEnd()
}
