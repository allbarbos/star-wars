package entity

import "net/http"

// HealthCheck entity
type HealthCheck struct {
	Status       string       `json:"status"`
	Dependencies Dependencies `json:"dependencies"`
}

// CheckDependencies externals
func (h *HealthCheck) CheckDependencies() int {
	if h.Dependencies.MongoDB == "ok" {
		return http.StatusOK
	}
	h.Status = h.Dependencies.MongoDB
	return http.StatusInternalServerError
}

// Dependencies entity
type Dependencies struct {
	MongoDB string `json:"mongoDb"`
}
