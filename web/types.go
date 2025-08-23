package web

import "time"

// OverviewResponse represents the dashboard overview data
type OverviewResponse struct {
	TotalSites       int            `json:"total_sites"`
	HealthySites     int            `json:"healthy_sites"`
	TotalChecks      int64          `json:"total_checks"`
	SuccessfulChecks int64          `json:"successful_checks"`
	OverallUptime    float64        `json:"overall_uptime"`
	Sites            []SiteOverview `json:"sites"`
	LastUpdate       time.Time      `json:"last_update"`
}

// SiteOverview represents a site's overview data
type SiteOverview struct {
	Name         string    `json:"name"`
	Status       string    `json:"status"` // "healthy", "degraded", "down", "stale"
	Uptime       float64   `json:"uptime"`
	ResponseTime int64     `json:"response_time_ms"`
	LastCheck    time.Time `json:"last_check"`
	TotalChecks  int64     `json:"total_checks"`
}

// SiteInfo represents basic site configuration info
type SiteInfo struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	Interval string `json:"interval"`
	Timeout  string `json:"timeout"`
}

// AlertStatus represents alert configuration status
type AlertStatus struct {
	EmailEnabled   bool `json:"email_enabled"`
	WebhookEnabled bool `json:"webhook_enabled"`
	TotalChannels  int  `json:"total_channels"`
}

// ChartDataPoint represents a data point for charts
type ChartDataPoint struct {
	Timestamp    time.Time `json:"timestamp"`
	Value        float64   `json:"value"`
	ResponseTime int64     `json:"response_time,omitempty"`
	Success      bool      `json:"success"`
}

// TimeSeriesData represents time series data for charts
type TimeSeriesData struct {
	SiteName string           `json:"site_name"`
	Data     []ChartDataPoint `json:"data"`
}

// DashboardStats represents comprehensive dashboard statistics
type DashboardStats struct {
	Overview   OverviewResponse `json:"overview"`
	TimeSeries []TimeSeriesData `json:"time_series"`
	AlertInfo  AlertStatus      `json:"alert_info"`
}
