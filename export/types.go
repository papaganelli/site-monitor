package export

import (
	"site-monitor/storage"
	"time"
)

// ExportFormat represents supported export formats
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
	FormatHTML ExportFormat = "html"
)

// ExportOptions contains configuration for data export
type ExportOptions struct {
	Format       ExportFormat  `json:"format"`
	SiteName     string        `json:"site_name,omitempty"`   // Filter by site (empty = all sites)
	Since        time.Duration `json:"since,omitempty"`       // Time period (e.g., 24h, 7d)
	Until        *time.Time    `json:"until,omitempty"`       // End time (default: now)
	Limit        int           `json:"limit,omitempty"`       // Max number of records
	IncludeStats bool          `json:"include_stats"`         // Include statistics summary
	OutputPath   string        `json:"output_path,omitempty"` // File path for CLI export
}

// ExportData represents the complete export dataset
type ExportData struct {
	Metadata ExportMetadata         `json:"metadata"`
	Stats    *ExportStats           `json:"stats,omitempty"`
	History  []storage.HistoryEntry `json:"history"`
}

// ExportMetadata contains information about the export
type ExportMetadata struct {
	GeneratedAt   time.Time     `json:"generated_at"`
	Format        ExportFormat  `json:"format"`
	TotalRecords  int           `json:"total_records"`
	SitesIncluded []string      `json:"sites_included"`
	TimeRange     TimeRange     `json:"time_range"`
	ExportOptions ExportOptions `json:"export_options"`
}

// TimeRange represents the time period covered by the export
type TimeRange struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// ExportStats contains aggregated statistics for the exported data
type ExportStats struct {
	TotalSites       int                      `json:"total_sites"`
	TotalChecks      int64                    `json:"total_checks"`
	SuccessfulChecks int64                    `json:"successful_checks"`
	FailedChecks     int64                    `json:"failed_checks"`
	OverallUptime    float64                  `json:"overall_uptime_percent"`
	SiteStats        map[string]storage.Stats `json:"site_stats"`

	// Response time statistics
	AvgResponseTime time.Duration `json:"avg_response_time"`
	MinResponseTime time.Duration `json:"min_response_time"`
	MaxResponseTime time.Duration `json:"max_response_time"`

	// Time distribution
	ChecksPerHour map[string]int `json:"checks_per_hour,omitempty"`
	ChecksPerDay  map[string]int `json:"checks_per_day,omitempty"`
}

// CSVRecord represents a single row for CSV export
type CSVRecord struct {
	Timestamp    string `csv:"timestamp"`
	SiteName     string `csv:"site_name"`
	URL          string `csv:"url"`
	Success      string `csv:"success"`
	StatusCode   string `csv:"status_code"`
	ResponseTime string `csv:"response_time_ms"`
	Error        string `csv:"error,omitempty"`
}

// HTMLExportData contains data formatted for HTML template
type HTMLExportData struct {
	ExportData
	FormattedHistory []FormattedHistoryEntry `json:"formatted_history"`
}

// FormattedHistoryEntry contains human-readable formatted data
type FormattedHistoryEntry struct {
	storage.HistoryEntry
	FormattedTimestamp    string `json:"formatted_timestamp"`
	FormattedResponseTime string `json:"formatted_response_time"`
	StatusIcon            string `json:"status_icon"`
	StatusClass           string `json:"status_class"`
}
