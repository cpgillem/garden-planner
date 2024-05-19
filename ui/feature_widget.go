package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/cpgillem/garden-planner/controllers"
	"github.com/cpgillem/garden-planner/geometry"
	"github.com/cpgillem/garden-planner/models"
	"golang.org/x/image/colornames"
)

type FeatureWidget struct {
	widget.BaseWidget

	// Internal widgets
	Label        *widget.Label
	Border       *canvas.Rectangle
	TopHandle    *Handle
	BottomHandle *Handle
	LeftHandle   *Handle
	RightHandle  *Handle

	// Internal data
	FeatureID models.FeatureID
	selected  bool

	// Controller Reference
	Controller *controllers.PlanController

	// Events
	OnDragged       func(e *fyne.DragEvent)
	OnDragEnd       func()
	OnHandleDragged func(edge geometry.BoxEdge, e *fyne.DragEvent)
	OnHandleDragEnd func(edge geometry.BoxEdge)
	OnTapped        func()
}

// Implement the Tappable interface to define click behavior.
func (fw *FeatureWidget) Tapped(e *fyne.PointEvent) {
	fw.Controller.SelectFeature(fw.FeatureID)
	fw.OnTapped()
}
func (fw *FeatureWidget) Dragged(e *fyne.DragEvent) {
	boxDelta := geometry.NewBoxWithValues(
		e.Dragged.DX/fw.Controller.DisplayConfig.Scale,
		e.Dragged.DY/fw.Controller.DisplayConfig.Scale,
		0,
		0,
	)
	fw.Controller.MoveResizeFeature(fw.FeatureID, &boxDelta)
	fw.OnDragged(e)
}
func (fw *FeatureWidget) DragEnd() {
	fw.OnDragEnd()
}

func (fw *FeatureWidget) HandleDragged(edge geometry.BoxEdge, e *fyne.DragEvent) {
	dx := e.Dragged.DX / fw.Controller.DisplayConfig.Scale
	dy := e.Dragged.DY / fw.Controller.DisplayConfig.Scale
	dbox := geometry.NewBox()

	// Handle edge cases (lol)
	switch edge {
	case geometry.TOP:
		dbox.Location.Y = dy
		dbox.Size.Y = -dy
	case geometry.BOTTOM:
		dbox.Size.Y = dy
	case geometry.LEFT:
		dbox.Location.X = dx
		dbox.Size.X = -dx
	case geometry.RIGHT:
		dbox.Size.X = dx
	}

	// Add box delta.
	fw.Controller.MoveResizeFeature(fw.FeatureID, &dbox)
	fw.OnHandleDragged(edge, e)
}

func (fw *FeatureWidget) HandleDragEnd(edge geometry.BoxEdge) {
	fw.OnHandleDragEnd(edge)
}

// Create a new widget representing a landscaping feature.
func NewFeatureWidget(id models.FeatureID, controller *controllers.PlanController) *FeatureWidget {
	fw := FeatureWidget{
		FeatureID:       id,
		Controller:      controller,
		selected:        false,
		OnDragEnd:       func() {},
		OnDragged:       func(e *fyne.DragEvent) {},
		OnHandleDragged: func(edge geometry.BoxEdge, e *fyne.DragEvent) {},
		OnHandleDragEnd: func(edge geometry.BoxEdge) {},
		OnTapped:        func() {},
		Label:           widget.NewLabel(""),
		Border:          canvas.NewRectangle(colornames.Lawngreen),
		TopHandle:       NewHandle(),
		BottomHandle:    NewHandle(),
		LeftHandle:      NewHandle(),
		RightHandle:     NewHandle(),
	}

	// Handle drag events.
	fw.TopHandle.OnDragged = func(e *fyne.DragEvent) {
		fw.HandleDragged(geometry.TOP, e)
	}
	fw.BottomHandle.OnDragged = func(e *fyne.DragEvent) {
		fw.HandleDragged(geometry.BOTTOM, e)
	}
	fw.LeftHandle.OnDragged = func(e *fyne.DragEvent) {
		fw.HandleDragged(geometry.LEFT, e)
	}
	fw.RightHandle.OnDragged = func(e *fyne.DragEvent) {
		fw.HandleDragged(geometry.RIGHT, e)
	}

	fw.TopHandle.OnDragEnd = func() {
		fw.HandleDragEnd(geometry.TOP)
	}
	fw.BottomHandle.OnDragEnd = func() {
		fw.HandleDragEnd(geometry.BOTTOM)
	}
	fw.LeftHandle.OnDragEnd = func() {
		fw.HandleDragEnd(geometry.LEFT)
	}
	fw.RightHandle.OnDragEnd = func() {
		fw.HandleDragEnd(geometry.RIGHT)
	}

	fw.ExtendBaseWidget(&fw)

	return &fw
}

func (featureWidget *FeatureWidget) CreateRenderer() fyne.WidgetRenderer {
	return newFeatureRenderer(featureWidget)
}

func (fw *FeatureWidget) Select() {
	fw.selected = true
	fw.Border.StrokeColor = colornames.Black
	fw.Border.StrokeWidth = 1
}

func (fw *FeatureWidget) Deselect() {
	fw.selected = false
	fw.Border.StrokeWidth = 0
}

func (fw *FeatureWidget) IsSelected() bool {
	return fw.selected
}

type featureRenderer struct {
	parent *FeatureWidget
}

func newFeatureRenderer(parent *FeatureWidget) featureRenderer {
	fr := featureRenderer{
		parent: parent,
	}

	return fr
}

func (fr featureRenderer) Destroy() {

}

func (fr featureRenderer) Layout(size fyne.Size) {
	// TODO: Layout can have more objects depending on properties feature contains, e.g. row spacing.

	// Define size of rectangle.
	fr.parent.Border.Resize(size)

	// Position handles.
	handleSize := fyne.NewSquareSize(10)
	fr.parent.TopHandle.Resize(handleSize)
	fr.parent.BottomHandle.Resize(handleSize)
	fr.parent.LeftHandle.Resize(handleSize)
	fr.parent.RightHandle.Resize(handleSize)

	fr.parent.TopHandle.Move(fyne.NewPos(
		size.Width/2,
		0,
	))
	fr.parent.BottomHandle.Move(fyne.NewPos(
		size.Width/2,
		size.Height,
	))
	fr.parent.LeftHandle.Move(fyne.NewPos(
		0,
		size.Height/2,
	))
	fr.parent.RightHandle.Move(fyne.NewPos(
		size.Width,
		size.Height/2,
	))

	// Center the rectangle.
	fr.parent.Border.Move(fyne.NewPos(
		fr.parent.LeftHandle.Size().Width/2,
		fr.parent.TopHandle.Size().Height/2,
	))
}

func (fr featureRenderer) MinSize() fyne.Size {
	return fr.parent.Border.MinSize().Add(fr.parent.LeftHandle.MinSize())
}

func (fr featureRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		fr.parent.Border,
		fr.parent.TopHandle,
		fr.parent.BottomHandle,
		fr.parent.LeftHandle,
		fr.parent.RightHandle,
		fr.parent.Label,
	}
}

func (fr featureRenderer) Refresh() {
	fr.parent.Border.Refresh()

	fr.parent.Label.SetText(fr.parent.Controller.Plan.Features[fr.parent.FeatureID].Name)

	fr.parent.TopHandle.Refresh()
	fr.parent.BottomHandle.Refresh()
	fr.parent.LeftHandle.Refresh()
	fr.parent.RightHandle.Refresh()

	// fr.Layout(fr.MinSize())
}
