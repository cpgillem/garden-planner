package controllers

import (
	"github.com/cpgillem/garden-planner/geometry"
	"github.com/cpgillem/garden-planner/models"
)

// Takes care of plan data in one place. Fires data-related events.
type PlanController struct {
	Plan *models.Plan

	// Display data
	DisplayConfig *models.DisplayConfig

	// Defines how to refresh UI code.
	OnFeatureSelected func(id models.FeatureID)
}

func NewPlanController(plan *models.Plan, displayConfig *models.DisplayConfig) PlanController {
	return PlanController{
		Plan:              plan,
		DisplayConfig:     displayConfig,
		OnFeatureSelected: func(id models.FeatureID) {},
	}
}

func (c *PlanController) MoveResizeFeature(id models.FeatureID, boxDelta *geometry.AxisAlignedBoundingBox) {
	c.Plan.Features[id].Box.AddTo(boxDelta)
}

func (c *PlanController) SelectFeature(id models.FeatureID) {
	c.OnFeatureSelected(id)
}
