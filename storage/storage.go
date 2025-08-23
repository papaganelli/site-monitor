package storage

import (
	"fmt"
	"site-monitor/monitor"
	"time"
)

// Storage defines the interface for persisting monitoring results
type Storage interface {
	// SaveResult stores a monitoring result
	SaveResult(result monitor.Result) error

	// GetHistory retrieves monitoring history for a specific site
	GetHistory(siteName string, since time.Time) ([]HistoryEntry, error)

	// GetAllHistory retrieves monitoring history for all sites
	GetAllHistory(since time.Time) ([]HistoryEntry, error)

	// GetStats calculates statistics for a specific site
	GetStats(siteName string, since time.Time) (Stats, error)

	// GetAllStats calculates statistics for all sites
	GetAllStats(since time.Time) (map[string]Stats, error)

	// Close closes the storage connection
	Close() error

	// Init initializes the storage (creates tables, migrations, etc.)
	Init() error
}

// HistoryEntry represents a stored monitoring result with metadata
type HistoryEntry struct {
	ID        int64         `json:"id"`
	SiteName  string        `json:"site_name"`
	URL       string        `json:"url"`
	Status    int           `json:"status_code"`
	Duration  time.Duration `json:"response_time_ms"`
	Success   bool          `json:"success"`
	Error     string        `json:"error_message,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
	CreatedAt time.Time     `json:"created_at"`
}

// Stats represents calculated statistics for a site
type Stats struct {
	SiteName         string        `json:"site_name"`
	TotalChecks      int64         `json:"total_checks"`
	SuccessfulChecks int64         `json:"successful_checks"`
	FailedChecks     int64         `json:"failed_checks"`
	SuccessRate      float64       `json:"success_rate_percent"`
	AvgResponseTime  time.Duration `json:"avg_response_time_ms"`
	MinResponseTime  time.Duration `json:"min_response_time_ms"`
	MaxResponseTime  time.Duration `json:"max_response_time_ms"`
	LastCheck        time.Time     `json:"last_check"`
	FirstCheck       time.Time     `json:"first_check"`
	Uptime           time.Duration `json:"uptime_duration"`
	Downtime         time.Duration `json:"downtime_duration"`
}

// String returns a formatted representation of the stats
func (s Stats) String() string {
	return fmt.Sprintf(
		"📊 %s: %d checks, %.1f%% uptime, avg: %v, last: %s",
		s.SiteName,
		s.TotalChecks,
		s.SuccessRate,
		s.AvgResponseTime,
		s.LastCheck.Format("15:04:05"),
	)
}
