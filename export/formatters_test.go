package export

import (
	"bytes"
	"encoding/json"
	"site-monitor/storage"
	"strings"
	"testing"
	"time"
)

func createTestExportData() *ExportData {
	now := time.Now()

	return &ExportData{
		Metadata: ExportMetadata{
			GeneratedAt:   now,
			Format:        FormatJSON,
			TotalRecords:  2,
			SitesIncluded: []string{"Test Site 1", "Test Site 2"},
			TimeRange: TimeRange{
				From: now.Add(-1 * time.Hour),
				To:   now,
			},
		},
		Stats: &ExportStats{
			TotalSites:       2,
			TotalChecks:      2,
			SuccessfulChecks: 1,
			FailedChecks:     1,
			OverallUptime:    50.0,
			AvgResponseTime:  150 * time.Millisecond,
			MinResponseTime:  100 * time.Millisecond,
			MaxResponseTime:  200 * time.Millisecond,
			SiteStats: map[string]storage.Stats{
				"Test Site 1": {
					SiteName:         "Test Site 1",
					TotalChecks:      1,
					SuccessfulChecks: 1,
					SuccessRate:      100.0,
					AvgResponseTime:  100 * time.Millisecond,
				},
			},
		},
		History: []storage.HistoryEntry{
			{
				ID:        1,
				SiteName:  "Test Site 1",
				URL:       "https://test1.com",
				Status:    200,
				Duration:  100 * time.Millisecond,
				Success:   true,
				Timestamp: now.Add(-30 * time.Minute),
			},
			{
				ID:        2,
				SiteName:  "Test Site 2",
				URL:       "https://test2.com",
				Status:    500,
				Duration:  200 * time.Millisecond,
				Success:   false,
				Error:     "Internal Server Error",
				Timestamp: now.Add(-15 * time.Minute),
			},
		},
	}
}

func TestJSONFormatter(t *testing.T) {
	formatter := &JSONFormatter{}
	data := createTestExportData()

	var buf bytes.Buffer
	err := formatter.Format(data, &buf)
	if err != nil {
		t.Fatalf("JSON formatting failed: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("Generated JSON is invalid: %v", err)
	}

	// Verify required fields are present
	if _, exists := result["metadata"]; !exists {
		t.Error("Missing metadata field in JSON output")
	}

	if _, exists := result["stats"]; !exists {
		t.Error("Missing stats field in JSON output")
	}

	if _, exists := result["history"]; !exists {
		t.Error("Missing history field in JSON output")
	}

	// Verify content type and extension
	if formatter.ContentType() != "application/json" {
		t.Errorf("Expected content type application/json, got %s", formatter.ContentType())
	}

	if formatter.FileExtension() != ".json" {
		t.Errorf("Expected file extension .json, got %s", formatter.FileExtension())
	}
}

func TestCSVFormatter(t *testing.T) {
	formatter := &CSVFormatter{}
	data := createTestExportData()

	var buf bytes.Buffer
	err := formatter.Format(data, &buf)
	if err != nil {
		t.Fatalf("CSV formatting failed: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	// Verify header line
	if len(lines) < 1 {
		t.Fatal("CSV output is empty")
	}

	headerExpected := "timestamp,site_name,url,success,status_code,response_time_ms,error"
	if lines[0] != headerExpected {
		t.Errorf("Expected header: %s\nGot: %s", headerExpected, lines[0])
	}

	// Verify we have the right number of data lines (header + 2 records)
	if len(lines) != 3 {
		t.Errorf("Expected 3 lines (header + 2 records), got %d", len(lines))
	}

	// Verify first data line contains expected values
	firstDataLine := lines[1]
	expectedValues := []string{
		"Test Site 1",
		"https://test1.com",
		"true",
		"200",
		"100.00", // response time in ms
	}

	for _, expected := range expectedValues {
		if !strings.Contains(firstDataLine, expected) {
			t.Errorf("Expected CSV line to contain '%s', got: %s", expected, firstDataLine)
		}
	}

	// Verify content type and extension
	if formatter.ContentType() != "text/csv" {
		t.Errorf("Expected content type text/csv, got %s", formatter.ContentType())
	}

	if formatter.FileExtension() != ".csv" {
		t.Errorf("Expected file extension .csv, got %s", formatter.FileExtension())
	}
}

func TestHTMLFormatter(t *testing.T) {
	formatter := &HTMLFormatter{}
	data := createTestExportData()

	var buf bytes.Buffer
	err := formatter.Format(data, &buf)
	if err != nil {
		t.Fatalf("HTML formatting failed: %v", err)
	}

	output := buf.String()

	// Verify it's HTML
	if !strings.Contains(output, "<!DOCTYPE html>") {
		t.Error("Output doesn't appear to be HTML")
	}

	if !strings.Contains(output, "<html") {
		t.Error("Missing HTML tag")
	}

	// Verify title is present
	if !strings.Contains(output, "Site Monitor Export Report") {
		t.Error("Missing report title")
	}

	// Verify site data is included
	if !strings.Contains(output, "Test Site 1") {
		t.Error("Missing Test Site 1 in output")
	}

	if !strings.Contains(output, "Test Site 2") {
		t.Error("Missing Test Site 2 in output")
	}

	// Verify statistics are included
	if !strings.Contains(output, "50.0%") { // Overall uptime
		t.Error("Missing overall uptime in output")
	}

	// Verify status indicators
	if !strings.Contains(output, "✅") { // Success icon
		t.Error("Missing success status icon")
	}

	if !strings.Contains(output, "❌") { // Error icon
		t.Error("Missing error status icon")
	}

	// Verify error message is included
	if !strings.Contains(output, "Internal Server Error") {
		t.Error("Missing error message in output")
	}

	// Verify CSS is included
	if !strings.Contains(output, "<style>") {
		t.Error("Missing CSS styles")
	}

	// Verify content type and extension
	if formatter.ContentType() != "text/html" {
		t.Errorf("Expected content type text/html, got %s", formatter.ContentType())
	}

	if formatter.FileExtension() != ".html" {
		t.Errorf("Expected file extension .html, got %s", formatter.FileExtension())
	}
}

func TestGetFormatter(t *testing.T) {
	tests := []struct {
		format      ExportFormat
		expectError bool
		expectType  interface{}
	}{
		{FormatJSON, false, &JSONFormatter{}},
		{FormatCSV, false, &CSVFormatter{}},
		{FormatHTML, false, &HTMLFormatter{}},
		{ExportFormat("invalid"), true, nil},
	}

	for _, test := range tests {
		formatter, err := GetFormatter(test.format)

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for format %s", test.format)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for format %s: %v", test.format, err)
			continue
		}

		// Check type
		switch test.format {
		case FormatJSON:
			if _, ok := formatter.(*JSONFormatter); !ok {
				t.Errorf("Expected JSONFormatter for JSON format")
			}
		case FormatCSV:
			if _, ok := formatter.(*CSVFormatter); !ok {
				t.Errorf("Expected CSVFormatter for CSV format")
			}
		case FormatHTML:
			if _, ok := formatter.(*HTMLFormatter); !ok {
				t.Errorf("Expected HTMLFormatter for HTML format")
			}
		}
	}
}

func TestTemplateHelperFunctions(t *testing.T) {
	// Test formatDuration
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{100 * time.Millisecond, "100ms"},
		{1500 * time.Millisecond, "1.5s"},
		{2 * time.Second, "2s"},
	}

	for _, test := range tests {
		result := formatDuration(test.duration)
		if result != test.expected {
			t.Errorf("formatDuration(%v) = %s, expected %s", test.duration, result, test.expected)
		}
	}

	// Test formatTime
	testTime := time.Date(2024, 1, 15, 14, 30, 45, 0, time.UTC)
	expected := "2024-01-15 14:30:45"
	result := formatTime(testTime)
	if result != expected {
		t.Errorf("formatTime() = %s, expected %s", result, expected)
	}

	// Test statusIcon
	if statusIcon(true) != "✅" {
		t.Errorf("statusIcon(true) should return ✅")
	}
	if statusIcon(false) != "❌" {
		t.Errorf("statusIcon(false) should return ❌")
	}

	// Test statusClass
	if statusClass(true) != "success" {
		t.Errorf("statusClass(true) should return 'success'")
	}
	if statusClass(false) != "error" {
		t.Errorf("statusClass(false) should return 'error'")
	}

	// Test formatPercent
	if formatPercent(95.567) != "95.6%" {
		t.Errorf("formatPercent(95.567) should return '95.6%%'")
	}
}

func TestHTMLFormatterFormatHistoryForHTML(t *testing.T) {
	formatter := &HTMLFormatter{}

	now := time.Now()
	history := []storage.HistoryEntry{
		{
			SiteName:  "Test Site",
			Success:   true,
			Duration:  100 * time.Millisecond,
			Timestamp: now,
		},
		{
			SiteName:  "Test Site",
			Success:   false,
			Duration:  200 * time.Millisecond,
			Timestamp: now.Add(-1 * time.Hour),
		},
	}

	formatted := formatter.formatHistoryForHTML(history)

	if len(formatted) != 2 {
		t.Errorf("Expected 2 formatted entries, got %d", len(formatted))
	}

	// Check first entry (success)
	if formatted[0].StatusIcon != "✅" {
		t.Errorf("Expected success icon ✅, got %s", formatted[0].StatusIcon)
	}

	if formatted[0].StatusClass != "success" {
		t.Errorf("Expected success class, got %s", formatted[0].StatusClass)
	}

	// Check second entry (failure)
	if formatted[1].StatusIcon != "❌" {
		t.Errorf("Expected error icon ❌, got %s", formatted[1].StatusIcon)
	}

	if formatted[1].StatusClass != "error" {
		t.Errorf("Expected error class, got %s", formatted[1].StatusClass)
	}

	// Verify formatted timestamp is not empty
	if formatted[0].FormattedTimestamp == "" {
		t.Error("FormattedTimestamp should not be empty")
	}

	// Verify formatted response time is not empty
	if formatted[0].FormattedResponseTime == "" {
		t.Error("FormattedResponseTime should not be empty")
	}
}
