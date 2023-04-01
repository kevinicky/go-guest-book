package entity

type Health struct {
	Status HealthComponent `json:"status"`
}

type HealthComponent struct {
	Server   string           `json:"server"`
	Database []HealthDatabase `json:"database"`
}

type HealthDatabase struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
