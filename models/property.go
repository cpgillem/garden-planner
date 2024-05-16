package models

type Property struct {
	Name         string `json:"name"`
	DisplayName  string `json:"display_name"`
	Description  string `json:"description"`
	PropertyType string `json:"property_type"`
}
