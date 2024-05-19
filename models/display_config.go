package models

type DisplayConfig struct {
	Scale float32
}

func NewDisplayConfig() DisplayConfig {
	return DisplayConfig{
		Scale: 2,
	}
}
