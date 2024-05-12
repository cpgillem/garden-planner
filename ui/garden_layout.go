package ui

import "fyne.io/fyne/v2"

type GardenLayout struct {
}

// Layout implements fyne.Layout.
func (g *GardenLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, containerSize.Height-g.MinSize(objects).Height)
	for _, o := range objects {
		size := o.MinSize()
		o.Resize(size)
		o.Move(pos)

		pos = pos.Add(fyne.NewPos(size.Width, size.Height))
	}
}

// MinSize implements fyne.Layout.
func (g *GardenLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()

		w += childSize.Width
		h += childSize.Height
	}
	return fyne.NewSize(w, h)
}

func NewGardenLayout() *GardenLayout {
	return &GardenLayout{}
}
