package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/cpgillem/garden-planner/models"
)

type FeatureWidget struct {
	widget.BaseWidget

	Label *widget.Label

	Feature *models.Feature
}

func NewFeatureWidget(feature *models.Feature) *FeatureWidget {
	featureWidget := &FeatureWidget{
		Label: widget.NewLabel(feature.Name),
	}

	featureWidget.Label.Truncation = fyne.TextTruncateEllipsis
	featureWidget.ExtendBaseWidget(featureWidget)
	return featureWidget
}

func (featureWidget *FeatureWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewCenter(featureWidget.Label)
	return widget.NewSimpleRenderer(c)
}
