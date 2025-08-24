package export

import (
	"fmt"
	"site-monitor/storage"
	"sort"
	"time"
)

// Exporter handles data export in various formats
type Exporter struct {
	storage storage.Storage
}

// NewExporter creates a new data exporter
func NewExporter(storage storage.Storage) *Exporter {
	return &Exporter{
		storage: storage,
	}
}

// Export exports monitoring data according to the specified options
func (e *Exporter) Export(opts ExportOptions) (*ExportData, error) {
	// Validate options
	if err := e.validateOptions(opts); err != nil {
		return nil, fmt.Errorf("invalid export options: %w", err)
	}

	// Calculate time range
	timeRange := e.calculateTimeRange(opts)

	// Fetch history data
	var history []storage.HistoryEntry
	var err error

	if opts.SiteName != "" {
		history, err = e.storage.GetHistory(opts.SiteName, timeRange.From)
	} else {
		history, err = e.storage.GetAllHistory(timeRange.From)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch history: %w", err)
	}

	// Filter by end time if specified
	if opts.Until != nil {
		history = e.filterByEndTime(history, *opts.Until)
	}

	// Apply limit
	if opts.Limit > 0 && len(history) > opts.Limit {
		history = history[:opts.Limit]
	}

	// Build site list
	sitesIncluded := e.getSitesFromHistory(history)

	// Generate statistics if requested
	var exportStats *ExportStats
	if opts.IncludeStats {
		exportStats = e.generateStats(history, timeRange, sitesIncluded)
	}

	// Create export data
	exportData := &ExportData{
		Metadata: ExportMetadata{
			GeneratedAt:   time.Now(),
			Format:        opts.Format,
			TotalRecords:  len(history),
			SitesIncluded: sitesIncluded,
			TimeRange:     timeRange,
			ExportOptions: opts,
		},
		Stats:   exportStats,
		History: history,
	}

	return exportData, nil
}

// validateOptions validates the export options
func (e *Exporter) validateOptions(opts ExportOptions) error {
	// Validate format
	switch opts.Format {
	case FormatJSON, FormatCSV, FormatHTML:
		// Valid formats
	default:
		return fmt.Errorf("unsupported format: %s (supported: json, csv, html)", opts.Format)
	}

	// Validate limit
	if opts.Limit < 0 {
		return fmt.Errorf("limit cannot be negative")
	}

	// Validate time range
	if opts.Since < 0 {
		return fmt.Errorf("since duration cannot be negative")
	}

	return nil
}

// calculateTimeRange determines the time range for the export
func (e *Exporter) calculateTimeRange(opts ExportOptions) TimeRange {
	now := time.Now()

	// Default to last 24 hours
	since := 24 * time.Hour
	if opts.Since > 0 {
		since = opts.Since
	}

	from := now.Add(-since)
	to := now
	if opts.Until != nil {
		to = *opts.Until
	}

	return TimeRange{
		From: from,
		To:   to,
	}
}

// filterByEndTime filters history entries by end time
func (e *Exporter) filterByEndTime(history []storage.HistoryEntry, endTime time.Time) []storage.HistoryEntry {
	var filtered []storage.HistoryEntry
	for _, entry := range history {
		if entry.Timestamp.Before(endTime) || entry.Timestamp.Equal(endTime) {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// getSitesFromHistory extracts unique site names from history
func (e *Exporter) getSitesFromHistory(history []storage.HistoryEntry) []string {
	siteSet := make(map[string]bool)
	for _, entry := range history {
		siteSet[entry.SiteName] = true
	}

	var sites []string
	for site := range siteSet {
		sites = append(sites, site)
	}

	// Sort for consistent output
	sort.Strings(sites)
	return sites
}

// generateStats calculates comprehensive statistics for the export
func (e *Exporter) generateStats(history []storage.HistoryEntry, timeRange TimeRange, sites []string) *ExportStats {
	stats := &ExportStats{
		TotalSites:    len(sites),
		TotalChecks:   int64(len(history)),
		SiteStats:     make(map[string]storage.Stats),
		ChecksPerHour: make(map[string]int),
		ChecksPerDay:  make(map[string]int),
	}

	if len(history) == 0 {
		return stats
	}

	// Calculate basic metrics
	var totalResponseTime time.Duration
	var successCount int64
	var minResponseTime, maxResponseTime time.Duration
	first := true

	// Track time distributions
	for _, entry := range history {
		if entry.Success {
			successCount++
			totalResponseTime += entry.Duration

			if first || entry.Duration < minResponseTime {
				minResponseTime = entry.Duration
				first = false
			}
			if entry.Duration > maxResponseTime {
				maxResponseTime = entry.Duration
			}
		}

		// Hour distribution
		hourKey := entry.Timestamp.Format("2006-01-02 15")
		stats.ChecksPerHour[hourKey]++

		// Day distribution
		dayKey := entry.Timestamp.Format("2006-01-02")
		stats.ChecksPerDay[dayKey]++
	}

	stats.SuccessfulChecks = successCount
	stats.FailedChecks = stats.TotalChecks - successCount

	// Calculate overall uptime
	if stats.TotalChecks > 0 {
		stats.OverallUptime = float64(stats.SuccessfulChecks) / float64(stats.TotalChecks) * 100
	}

	// Calculate response time stats
	if successCount > 0 {
		stats.AvgResponseTime = totalResponseTime / time.Duration(successCount)
		stats.MinResponseTime = minResponseTime
		stats.MaxResponseTime = maxResponseTime
	}

	// Generate per-site statistics
	for _, site := range sites {
		if siteStats, err := e.storage.GetStats(site, timeRange.From); err == nil {
			stats.SiteStats[site] = siteStats
		}
	}

	return stats
}

// GetSupportedFormats returns list of supported export formats
func GetSupportedFormats() []ExportFormat {
	return []ExportFormat{FormatJSON, FormatCSV, FormatHTML}
}

// FormatDescription returns human-readable format descriptions
func FormatDescription(format ExportFormat) string {
	switch format {
	case FormatJSON:
		return "JSON - Machine-readable structured data"
	case FormatCSV:
		return "CSV - Spreadsheet compatible comma-separated values"
	case FormatHTML:
		return "HTML - Human-readable web page report"
	default:
		return "Unknown format"
	}
}
