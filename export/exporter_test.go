package export

import (
	"site-monitor/monitor"
	"site-monitor/storage"
	"testing"
	"time"
)

// MockStorage implements storage.Storage interface for testing
type MockStorage struct {
	history []storage.HistoryEntry
	stats   map[string]storage.Stats
}

func (m *MockStorage) SaveResult(result monitor.Result) error {
	return nil
}

func (m *MockStorage) GetHistory(siteName string, since time.Time) ([]storage.HistoryEntry, error) {
	var filtered []storage.HistoryEntry
	for _, entry := range m.history {
		if entry.SiteName == siteName && entry.Timestamp.After(since) {
			filtered = append(filtered, entry)
		}
	}
	return filtered, nil
}

func (m *MockStorage) GetAllHistory(since time.Time) ([]storage.HistoryEntry, error) {
	var filtered []storage.HistoryEntry
	for _, entry := range m.history {
		if entry.Timestamp.After(since) {
			filtered = append(filtered, entry)
		}
	}
	return filtered, nil
}

func (m *MockStorage) GetStats(siteName string, since time.Time) (storage.Stats, error) {
	if stats, exists := m.stats[siteName]; exists {
		return stats, nil
	}
	return storage.Stats{}, nil
}

func (m *MockStorage) GetAllStats(since time.Time) (map[string]storage.Stats, error) {
	return m.stats, nil
}

func (m *MockStorage) Close() error {
	return nil
}

func (m *MockStorage) Init() error {
	return nil
}

func createMockStorage() *MockStorage {
	now := time.Now()

	return &MockStorage{
		history: []storage.HistoryEntry{
			{
				ID:        1,
				SiteName:  "Test Site 1",
				URL:       "https://test1.com",
				Status:    200,
				Duration:  100 * time.Millisecond,
				Success:   true,
				Timestamp: now.Add(-1 * time.Hour),
			},
			{
				ID:        2,
				SiteName:  "Test Site 1",
				URL:       "https://test1.com",
				Status:    500,
				Duration:  200 * time.Millisecond,
				Success:   false,
				Error:     "Internal Server Error",
				Timestamp: now.Add(-30 * time.Minute),
			},
			{
				ID:        3,
				SiteName:  "Test Site 2",
				URL:       "https://test2.com",
				Status:    200,
				Duration:  50 * time.Millisecond,
				Success:   true,
				Timestamp: now.Add(-15 * time.Minute),
			},
		},
		stats: map[string]storage.Stats{
			"Test Site 1": {
				SiteName:         "Test Site 1",
				TotalChecks:      100,
				SuccessfulChecks: 95,
				FailedChecks:     5,
				SuccessRate:      95.0,
				AvgResponseTime:  120 * time.Millisecond,
				MinResponseTime:  50 * time.Millisecond,
				MaxResponseTime:  300 * time.Millisecond,
				LastCheck:        now.Add(-5 * time.Minute),
				FirstCheck:       now.Add(-24 * time.Hour),
			},
			"Test Site 2": {
				SiteName:         "Test Site 2",
				TotalChecks:      50,
				SuccessfulChecks: 50,
				FailedChecks:     0,
				SuccessRate:      100.0,
				AvgResponseTime:  80 * time.Millisecond,
				MinResponseTime:  50 * time.Millisecond,
				MaxResponseTime:  120 * time.Millisecond,
				LastCheck:        now.Add(-2 * time.Minute),
				FirstCheck:       now.Add(-12 * time.Hour),
			},
		},
	}
}

func TestExporter_Export_JSON(t *testing.T) {
	mockStorage := createMockStorage()
	exporter := NewExporter(mockStorage)

	opts := ExportOptions{
		Format:       FormatJSON,
		Since:        2 * time.Hour,
		IncludeStats: true,
	}

	data, err := exporter.Export(opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify metadata
	if data.Metadata.Format != FormatJSON {
		t.Errorf("Expected format %s, got %s", FormatJSON, data.Metadata.Format)
	}

	if data.Metadata.TotalRecords != 3 {
		t.Errorf("Expected 3 records, got %d", data.Metadata.TotalRecords)
	}

	if len(data.Metadata.SitesIncluded) != 2 {
		t.Errorf("Expected 2 sites, got %d", len(data.Metadata.SitesIncluded))
	}

	// Verify stats are included
	if data.Stats == nil {
		t.Error("Expected stats to be included")
	} else {
		if data.Stats.TotalChecks != 3 {
			t.Errorf("Expected 3 total checks in stats, got %d", data.Stats.TotalChecks)
		}
		if data.Stats.SuccessfulChecks != 2 {
			t.Errorf("Expected 2 successful checks, got %d", data.Stats.SuccessfulChecks)
		}
	}

	// Verify history
	if len(data.History) != 3 {
		t.Errorf("Expected 3 history entries, got %d", len(data.History))
	}
}

func TestExporter_Export_WithSiteFilter(t *testing.T) {
	mockStorage := createMockStorage()
	exporter := NewExporter(mockStorage)

	opts := ExportOptions{
		Format:   FormatJSON,
		SiteName: "Test Site 1",
		Since:    2 * time.Hour,
	}

	data, err := exporter.Export(opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify only Test Site 1 entries are included
	if data.Metadata.TotalRecords != 2 {
		t.Errorf("Expected 2 records for Test Site 1, got %d", data.Metadata.TotalRecords)
	}

	if len(data.Metadata.SitesIncluded) != 1 {
		t.Errorf("Expected 1 site, got %d", len(data.Metadata.SitesIncluded))
	}

	if data.Metadata.SitesIncluded[0] != "Test Site 1" {
		t.Errorf("Expected site 'Test Site 1', got '%s'", data.Metadata.SitesIncluded[0])
	}

	// Verify all entries are for Test Site 1
	for _, entry := range data.History {
		if entry.SiteName != "Test Site 1" {
			t.Errorf("Expected all entries to be for Test Site 1, found entry for %s", entry.SiteName)
		}
	}
}

func TestExporter_Export_WithLimit(t *testing.T) {
	mockStorage := createMockStorage()
	exporter := NewExporter(mockStorage)

	opts := ExportOptions{
		Format: FormatJSON,
		Since:  2 * time.Hour,
		Limit:  2,
	}

	data, err := exporter.Export(opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify limit is applied
	if data.Metadata.TotalRecords != 2 {
		t.Errorf("Expected 2 records due to limit, got %d", data.Metadata.TotalRecords)
	}

	if len(data.History) != 2 {
		t.Errorf("Expected 2 history entries, got %d", len(data.History))
	}
}

func TestExporter_Export_InvalidOptions(t *testing.T) {
	mockStorage := createMockStorage()
	exporter := NewExporter(mockStorage)

	// Test invalid format
	opts := ExportOptions{
		Format: ExportFormat("invalid"),
		Since:  1 * time.Hour,
	}

	_, err := exporter.Export(opts)
	if err == nil {
		t.Error("Expected error for invalid format")
	}

	// Test negative limit
	opts = ExportOptions{
		Format: FormatJSON,
		Since:  1 * time.Hour,
		Limit:  -1,
	}

	_, err = exporter.Export(opts)
	if err == nil {
		t.Error("Expected error for negative limit")
	}

	// Test negative since duration
	opts = ExportOptions{
		Format: FormatJSON,
		Since:  -1 * time.Hour,
	}

	_, err = exporter.Export(opts)
	if err == nil {
		t.Error("Expected error for negative since duration")
	}
}

func TestExporter_TimeRangeFiltering(t *testing.T) {
	mockStorage := createMockStorage()
	exporter := NewExporter(mockStorage)

	now := time.Now()
	until := now.Add(-20 * time.Minute) // Only entries older than 20 minutes

	opts := ExportOptions{
		Format: FormatJSON,
		Since:  2 * time.Hour,
		Until:  &until,
	}

	data, err := exporter.Export(opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Should only get entries older than 20 minutes (2 entries)
	if data.Metadata.TotalRecords != 2 {
		t.Errorf("Expected 2 records with until filter, got %d", data.Metadata.TotalRecords)
	}

	// Verify all entries are before the until time
	for _, entry := range data.History {
		if entry.Timestamp.After(until) {
			t.Errorf("Found entry after until time: %v > %v", entry.Timestamp, until)
		}
	}
}

func TestExporter_StatsGeneration(t *testing.T) {
	mockStorage := createMockStorage()
	exporter := NewExporter(mockStorage)

	opts := ExportOptions{
		Format:       FormatJSON,
		Since:        2 * time.Hour,
		IncludeStats: true,
	}

	data, err := exporter.Export(opts)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	if data.Stats == nil {
		t.Fatal("Expected stats to be generated")
	}

	// Verify calculated stats
	if data.Stats.TotalSites != 2 {
		t.Errorf("Expected 2 sites in stats, got %d", data.Stats.TotalSites)
	}

	if data.Stats.TotalChecks != 3 {
		t.Errorf("Expected 3 total checks, got %d", data.Stats.TotalChecks)
	}

	if data.Stats.SuccessfulChecks != 2 {
		t.Errorf("Expected 2 successful checks, got %d", data.Stats.SuccessfulChecks)
	}

	if data.Stats.FailedChecks != 1 {
		t.Errorf("Expected 1 failed check, got %d", data.Stats.FailedChecks)
	}

	// Check uptime calculation
	expectedUptime := float64(2) / float64(3) * 100
	if abs(data.Stats.OverallUptime-expectedUptime) > 0.1 {
		t.Errorf("Expected uptime %.1f%%, got %.1f%%", expectedUptime, data.Stats.OverallUptime)
	}

	// Verify time distributions
	if len(data.Stats.ChecksPerHour) == 0 {
		t.Error("Expected checks per hour data")
	}

	if len(data.Stats.ChecksPerDay) == 0 {
		t.Error("Expected checks per day data")
	}

	// Verify site stats are included
	if len(data.Stats.SiteStats) != 2 {
		t.Errorf("Expected site stats for 2 sites, got %d", len(data.Stats.SiteStats))
	}
}

func TestGetSupportedFormats(t *testing.T) {
	formats := GetSupportedFormats()

	expectedFormats := []ExportFormat{FormatJSON, FormatCSV, FormatHTML}
	if len(formats) != len(expectedFormats) {
		t.Errorf("Expected %d formats, got %d", len(expectedFormats), len(formats))
	}

	for _, expected := range expectedFormats {
		found := false
		for _, format := range formats {
			if format == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected format %s not found", expected)
		}
	}
}

func TestFormatDescription(t *testing.T) {
	tests := []struct {
		format   ExportFormat
		contains string
	}{
		{FormatJSON, "JSON"},
		{FormatCSV, "CSV"},
		{FormatHTML, "HTML"},
		{ExportFormat("invalid"), "Unknown"},
	}

	for _, test := range tests {
		desc := FormatDescription(test.format)
		if !contains(desc, test.contains) {
			t.Errorf("Expected description for %s to contain '%s', got '%s'",
				test.format, test.contains, desc)
		}
	}
}

// Helper function for floating point comparison
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				indexOfSubstring(s, substr) >= 0))
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
