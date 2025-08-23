package alerts

import (
	"fmt"
	"site-monitor/monitor"
	"time"
)

// AlertType represents the type of alert
type AlertType string

const (
	AlertTypeSiteDown     AlertType = "site_down"
	AlertTypeSiteUp       AlertType = "site_up"
	AlertTypeSlowResponse AlertType = "slow_response"
	AlertTypeLowUptime    AlertType = "low_uptime"
)

// AlertSeverity represents the severity level of an alert
type AlertSeverity string

const (
	SeverityInfo     AlertSeverity = "info"
	SeverityWarning  AlertSeverity = "warning"
	SeverityCritical AlertSeverity = "critical"
)

// Alert represents a monitoring alert
type Alert struct {
	ID         string        `json:"id"`
	Type       AlertType     `json:"type"`
	Severity   AlertSeverity `json:"severity"`
	SiteName   string        `json:"site_name"`
	SiteURL    string        `json:"site_url"`
	Message    string        `json:"message"`
	Details    string        `json:"details,omitempty"`
	Timestamp  time.Time     `json:"timestamp"`
	Resolved   bool          `json:"resolved"`
	ResolvedAt *time.Time    `json:"resolved_at,omitempty"`

	// Context information
	CurrentStatus    int           `json:"current_status,omitempty"`
	ResponseTime     time.Duration `json:"response_time,omitempty"`
	ConsecutiveFails int           `json:"consecutive_fails,omitempty"`
	UptimePercent    float64       `json:"uptime_percent,omitempty"`
	ErrorMessage     string        `json:"error_message,omitempty"`
}

// AlertChannel defines the interface for sending alerts
type AlertChannel interface {
	// Send sends an alert through this channel
	Send(alert Alert) error

	// Test verifies the channel configuration
	Test() error

	// Name returns the channel name for logging
	Name() string
}

// AlertState represents the current alert state for a site
type AlertState struct {
	SiteName         string    `json:"site_name"`
	IsDown           bool      `json:"is_down"`
	ConsecutiveFails int       `json:"consecutive_fails"`
	LastFailTime     time.Time `json:"last_fail_time,omitempty"`
	LastSuccessTime  time.Time `json:"last_success_time,omitempty"`
	LastAlertTime    time.Time `json:"last_alert_time,omitempty"`
	ActiveAlerts     []string  `json:"active_alerts"` // Alert IDs
}

// String returns a formatted string representation of the alert
func (a Alert) String() string {
	switch a.Type {
	case AlertTypeSiteDown:
		return fmt.Sprintf("🚨 SITE DOWN: %s is not responding", a.SiteName)
	case AlertTypeSiteUp:
		return fmt.Sprintf("✅ SITE RECOVERED: %s is back online", a.SiteName)
	case AlertTypeSlowResponse:
		return fmt.Sprintf("⚠️ SLOW RESPONSE: %s is responding slowly (%v)", a.SiteName, a.ResponseTime)
	case AlertTypeLowUptime:
		return fmt.Sprintf("📉 LOW UPTIME: %s uptime is %.1f%%", a.SiteName, a.UptimePercent)
	default:
		return fmt.Sprintf("🔔 ALERT: %s - %s", a.SiteName, a.Message)
	}
}

// IsRecoveryAlert returns true if this is a recovery/resolution alert
func (a Alert) IsRecoveryAlert() bool {
	return a.Type == AlertTypeSiteUp
}

// ShouldResolveAlert determines if an existing alert should be marked as resolved
func ShouldResolveAlert(currentResult monitor.Result, alertType AlertType) bool {
	switch alertType {
	case AlertTypeSiteDown:
		return currentResult.Success
	case AlertTypeSlowResponse:
		// This would need the threshold to determine, simplified for now
		return currentResult.Success && currentResult.Duration < 5*time.Second
	default:
		return false
	}
}
