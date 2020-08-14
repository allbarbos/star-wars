package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckDependencies(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		hc := HealthCheck{
			Status: "ok",
			Dependencies: Dependencies{
				MongoDB: "ok",
			},
		}

		status := hc.CheckDependencies()

		assert.Equal(t, 200, status)
	})

	t.Run("when returns error", func(t *testing.T) {
		hc := HealthCheck{
			Status: "error",
			Dependencies: Dependencies{
				MongoDB: "error",
			},
		}

		status := hc.CheckDependencies()

		assert.Equal(t, 500, status)
	})
}
