package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/cpgillem/garden-planner/models"
	"golang.org/x/image/colornames"
)

type FeatureWidget struct {
	widget.BaseWidget

	Label   *widget.Label
	Border  *canvas.Rectangle
	Feature *models.Feature
}

// Create a new widget representing a landscaping feature.
func NewFeatureWidget(feature *models.Feature) FeatureWidget {
	featureWidget := FeatureWidget{
		Label:   widget.NewLabelWithData(binding.BindString(&feature.Name)),
		Border:  canvas.NewRectangle(colornames.Lawngreen),
		Feature: feature,
	}

	featureWidget.Border.StrokeColor = colornames.Black
	featureWidget.Border.StrokeWidth = 1
	featureWidget.Label.Truncation = fyne.TextTruncateOff
	featureWidget.Label.Alignment = fyne.TextAlignCenter

	featureWidget.ExtendBaseWidget(&featureWidget)

	return featureWidget
}

func (featureWidget *FeatureWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewStack(featureWidget.Border, featureWidget.Label)
	return widget.NewSimpleRenderer(c)
}
