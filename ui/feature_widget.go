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

	// Internal data
	FeatureID models.FeatureID

	// Controller Reference
	Controller *controllers.PlanController

	// Events
	OnDragged       func(e *fyne.DragEvent)
	OnDragEnd       func()
	OnHandleDragged func(edge geometry.BoxEdge, e *fyne.DragEvent)
	OnHandleDragEnd func(edge geometry.BoxEdge)
}

// Implement the Tappable interface to define click behavior.
func (fw *FeatureWidget) Tapped(e *fyne.PointEvent) {
	fw.Controller.SelectFeature(fw.FeatureID)
	fw.Refresh()
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
	featureWidget := FeatureWidget{
		FeatureID:       id,
		Controller:      controller,
		OnDragEnd:       func() {},
		OnDragged:       func(e *fyne.DragEvent) {},
		OnHandleDragged: func(edge geometry.BoxEdge, e *fyne.DragEvent) {},
		OnHandleDragEnd: func(edge geometry.BoxEdge) {},
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
		fr.parent.HandleDragged(geometry.TOP, e)
	}
	fr.BottomHandle.OnDragged = func(e *fyne.DragEvent) {
		fr.parent.HandleDragged(geometry.BOTTOM, e)
	}
	fr.LeftHandle.OnDragged = func(e *fyne.DragEvent) {
		fr.parent.HandleDragged(geometry.LEFT, e)
	}
	fr.RightHandle.OnDragged = func(e *fyne.DragEvent) {
		fr.parent.HandleDragged(geometry.RIGHT, e)
	}

	fr.TopHandle.OnDragEnd = func() {
		fr.parent.HandleDragEnd(geometry.TOP)
	}
	fr.BottomHandle.OnDragEnd = func() {
		fr.parent.HandleDragEnd(geometry.BOTTOM)
	}
	fr.LeftHandle.OnDragEnd = func() {
		fr.parent.HandleDragEnd(geometry.LEFT)
	}
	fr.RightHandle.OnDragEnd = func() {
		fr.parent.HandleDragEnd(geometry.RIGHT)
	}

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

	fr.Label.SetText(fr.parent.Controller.Plan.Features[fr.parent.FeatureID].Name)

	fr.TopHandle.Refresh()
	fr.BottomHandle.Refresh()
	fr.LeftHandle.Refresh()
	fr.RightHandle.Refresh()

	// fr.Layout(fr.MinSize())
}
