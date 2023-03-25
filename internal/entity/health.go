package entity

type Health struct {
	Status HealthComponent `json:"status"`
}

type HealthComponent struct {
	Server   string `json:"server"`
	Database string `json:"database"`
}
