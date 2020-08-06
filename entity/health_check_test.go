package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckDependencies( t *testing.T) {
	hc := HealthCheck{
		Status: "ok",
		Dependencies: Dependencies{
			MongoDB: "ok",
		},
	}

	status := hc.CheckDependencies()

	assert.Equal(t, 200, status)
}

func TestCheckDependencies_Error( t *testing.T) {
	hc := HealthCheck{
		Status: "error",
		Dependencies: Dependencies{
			MongoDB: "error",
		},
	}

	status := hc.CheckDependencies()

	assert.Equal(t, 500, status)
}
