package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/cpgillem/garden-planner/geometry"
	"github.com/cpgillem/garden-planner/models"
	"golang.org/x/image/colornames"
)

type FeatureWidget struct {
	widget.BaseWidget

	// Internal data
	Feature *models.Feature

	// Events
	OnTapped        func()
	OnHandleDragged func(edge geometry.BoxEdge, e *fyne.DragEvent)
	OnHandleDragEnd func()
}

// Implement the Tappable interface to define click behavior.
func (fw *FeatureWidget) Tapped(e *fyne.PointEvent) {
	fw.OnTapped()
}

// Create a new widget representing a landscaping feature.
func NewFeatureWidget(feature *models.Feature) *FeatureWidget {
	featureWidget := FeatureWidget{
		Feature: feature,
	}

	featureWidget.ExtendBaseWidget(&featureWidget)

	return &featureWidget
}

func (featureWidget *FeatureWidget) CreateRenderer() fyne.WidgetRenderer {
	return newFeatureRenderer(featureWidget)
}

type featureRenderer struct {
	parent *FeatureWidget

	// Internal widgets
	Label        *widget.Label
	Border       *canvas.Rectangle
	TopHandle    *Handle
	BottomHandle *Handle
	LeftHandle   *Handle
	RightHandle  *Handle
}

func newFeatureRenderer(parent *FeatureWidget) featureRenderer {
	fr := featureRenderer{
		parent:       parent,
		Label:        widget.NewLabel(""),
		Border:       canvas.NewRectangle(colornames.Lawngreen),
		TopHandle:    NewHandle(),
		BottomHandle: NewHandle(),
		LeftHandle:   NewHandle(),
		RightHandle:  NewHandle(),
	}

	// Handle drag events.
	fr.TopHandle.OnDragged = func(e *fyne.DragEvent) {
		fr.parent.OnHandleDragged(geometry.TOP, e)
	}
	fr.BottomHandle.OnDragged = func(e *fyne.DragEvent) {
		fr.parent.OnHandleDragged(geometry.BOTTOM, e)
	}
	fr.LeftHandle.OnDragged = func(e *fyne.DragEvent) {
		fr.parent.OnHandleDragged(geometry.LEFT, e)
	}
	fr.RightHandle.OnDragged = func(e *fyne.DragEvent) {
		fr.parent.OnHandleDragged(geometry.RIGHT, e)
	}
	fr.TopHandle.OnDragEnd = fr.parent.OnHandleDragEnd
	fr.BottomHandle.OnDragEnd = fr.parent.OnHandleDragEnd
	fr.LeftHandle.OnDragEnd = fr.parent.OnHandleDragEnd
	fr.RightHandle.OnDragEnd = fr.parent.OnHandleDragEnd

	return fr
}

func (fr featureRenderer) Destroy() {

}

func (fr featureRenderer) Layout(size fyne.Size) {
	// Define size of rectangle.
	fr.Border.Resize(size)

	// Position handles.
	handleSize := fyne.NewSquareSize(10)
	fr.TopHandle.Resize(handleSize)
	fr.BottomHandle.Resize(handleSize)
	fr.LeftHandle.Resize(handleSize)
	fr.RightHandle.Resize(handleSize)

	fr.TopHandle.Move(fyne.NewPos(
		size.Width/2,
		0,
	))
	fr.BottomHandle.Move(fyne.NewPos(
		size.Width/2,
		size.Height,
	))
	fr.LeftHandle.Move(fyne.NewPos(
		0,
		size.Height/2,
	))
	fr.RightHandle.Move(fyne.NewPos(
		size.Width,
		size.Height/2,
	))

	// Center the rectangle.
	fr.Border.Move(fyne.NewPos(
		fr.LeftHandle.Size().Width/2,
		fr.TopHandle.Size().Height/2,
	))
}

func (fr featureRenderer) MinSize() fyne.Size {
	return fr.Border.MinSize().Add(fr.LeftHandle.MinSize())
}

func (fr featureRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		fr.Border,
		fr.TopHandle,
		fr.BottomHandle,
		fr.LeftHandle,
		fr.RightHandle,
		fr.Label,
	}
}

func (fr featureRenderer) Refresh() {
	fr.Border.Refresh()

	fr.Label.SetText(fr.parent.Feature.Name)
	fr.Label.Refresh()

	fr.TopHandle.Refresh()
	fr.BottomHandle.Refresh()
	fr.LeftHandle.Refresh()
	fr.RightHandle.Refresh()
}
