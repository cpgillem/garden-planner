package controllers

import (
	"github.com/cpgillem/garden-planner/geometry"
	"github.com/cpgillem/garden-planner/models"
)

// Takes care of plan data in one place. Fires data-related events.
type PlanController struct {
	Plan *models.Plan

	selectedFeature models.FeatureID

	// Defines how to refresh UI code.
	OnFeatureSelected func(id models.FeatureID)
	OnFeatureAdded    func(id models.FeatureID)
	OnFeatureRemoved  func(id models.FeatureID)
}

func NewPlanController(plan *models.Plan) PlanController {
	return PlanController{
		Plan:              plan,
		OnFeatureSelected: func(id models.FeatureID) {},
		OnFeatureAdded:    func(id models.FeatureID) {},
		OnFeatureRemoved:  func(id models.FeatureID) {},
		selectedFeature:   -1,
	}
}

func (c *PlanController) MoveResizeFeature(id models.FeatureID, boxDelta *geometry.Box) {
	c.Plan.Features[id].Box.AddTo(boxDelta)
}

func (c *PlanController) SelectFeature(id models.FeatureID) {
	c.selectedFeature = id
	c.OnFeatureSelected(id)
}

func (c *PlanController) GetSelectedFeature() models.FeatureID {
	return c.selectedFeature
}

func (c *PlanController) HasSelection() bool {
	return c.Plan.Features[c.GetSelectedFeature()] != nil
}

func (c *PlanController) AddFeature(f models.Feature) {
	id := c.NewFeatureID()
	c.Plan.Features[id] = &f
	c.OnFeatureAdded(id)
}

func (c *PlanController) RemoveFeature(id models.FeatureID) {
	// c.Plan.Features[id] = nil
	c.OnFeatureRemoved(id)
}

func (c *PlanController) RemoveSelected() {
	if c.HasSelection() {
		c.RemoveFeature(c.GetSelectedFeature())
	}
}

func (c *PlanController) NewFeatureID() models.FeatureID {
	// Find max ID.
	max := models.FeatureID(-1)
	for i := range c.Plan.Features {
		if i > max {
			max = i
		}
	}
	return max + 1
}

func (c *PlanController) GetMaxName() string {
	// Find the maximum length of a feature.
	max := ""
	for i := range c.Plan.Features {
		n := c.Plan.Features[i].Name
		if len(n) > len(max) {
			max = n
		}
	}
	return max
}

func (c *PlanController) LenFeatures() int {
	return len(c.Plan.Features)
}

func (c *PlanController) HasFeature(id models.FeatureID) bool {
	return c.Plan.Features[id] != nil
}
