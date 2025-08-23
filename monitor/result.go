package monitor

import (
	"fmt"
	"time"
)

type Result struct {
	Name      string        `json:"name"` // Monitor name
	URL       string        `json:"url"`
	Status    int           `json:"status"`
	Duration  time.Duration `json:"duration"`
	Timestamp time.Time     `json:"timestamp"`
	Success   bool          `json:"success"`
	Error     string        `json:"error,omitempty"`
}

// String returns a formatted string representation of the result
func (r Result) String() string {
	status := "✅ OK"
	if !r.Success {
		status = "❌ ERROR" // ← Corrigé (pas de := car on réassigne)
	}

	return fmt.Sprintf("[%s] %s (%s) - Status: %d - Duration: %v",
		r.Timestamp.Format("15:04:05"),
		status,
		r.Name,
		r.Status,
		r.Duration)
}
