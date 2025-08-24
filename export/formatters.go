package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"site-monitor/storage"
	"strconv"
	"time"
)

// Formatter interface for different export formats
type Formatter interface {
	Format(data *ExportData, writer io.Writer) error
	ContentType() string
	FileExtension() string
}

// JSONFormatter handles JSON export
type JSONFormatter struct{}

func (f *JSONFormatter) Format(data *ExportData, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ") // Pretty print
	return encoder.Encode(data)
}

func (f *JSONFormatter) ContentType() string {
	return "application/json"
}

func (f *JSONFormatter) FileExtension() string {
	return ".json"
}

// CSVFormatter handles CSV export
type CSVFormatter struct{}

func (f *CSVFormatter) Format(data *ExportData, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	header := []string{
		"timestamp",
		"site_name",
		"url",
		"success",
		"status_code",
		"response_time_ms",
		"error",
	}
	if err := csvWriter.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, entry := range data.History {
		record := []string{
			entry.Timestamp.Format(time.RFC3339),
			entry.SiteName,
			entry.URL,
			strconv.FormatBool(entry.Success),
			strconv.Itoa(entry.Status),
			fmt.Sprintf("%.2f", float64(entry.Duration.Nanoseconds())/1000000), // Convert to ms
			entry.Error,
		}

		if err := csvWriter.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	return nil
}

func (f *CSVFormatter) ContentType() string {
	return "text/csv"
}

func (f *CSVFormatter) FileExtension() string {
	return ".csv"
}

// HTMLFormatter handles HTML report export
type HTMLFormatter struct{}

func (f *HTMLFormatter) Format(data *ExportData, writer io.Writer) error {
	// Convert to HTML-friendly format
	htmlData := &HTMLExportData{
		ExportData:       *data,
		FormattedHistory: f.formatHistoryForHTML(data.History),
	}

	tmpl := template.Must(template.New("report").Funcs(template.FuncMap{
		"formatDuration": formatDuration,
		"formatTime":     formatTime,
		"statusIcon":     statusIcon,
		"statusClass":    statusClass,
		"formatPercent":  formatPercent,
		"add":            func(a, b int) int { return a + b },
	}).Parse(htmlReportTemplate))

	return tmpl.Execute(writer, htmlData)
}

func (f *HTMLFormatter) ContentType() string {
	return "text/html"
}

func (f *HTMLFormatter) FileExtension() string {
	return ".html"
}

func (f *HTMLFormatter) formatHistoryForHTML(history []storage.HistoryEntry) []FormattedHistoryEntry {
	formatted := make([]FormattedHistoryEntry, len(history))

	for i, entry := range history {
		formatted[i] = FormattedHistoryEntry{
			HistoryEntry:          entry,
			FormattedTimestamp:    entry.Timestamp.Format("2006-01-02 15:04:05"),
			FormattedResponseTime: formatDuration(entry.Duration),
			StatusIcon:            statusIcon(entry.Success),
			StatusClass:           statusClass(entry.Success),
		}
	}

	return formatted
}

// GetFormatter returns appropriate formatter for the given format
func GetFormatter(format ExportFormat) (Formatter, error) {
	switch format {
	case FormatJSON:
		return &JSONFormatter{}, nil
	case FormatCSV:
		return &CSVFormatter{}, nil
	case FormatHTML:
		return &HTMLFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// Template helper functions
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return d.Round(time.Millisecond).String()
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func statusIcon(success bool) string {
	if success {
		return "✅"
	}
	return "❌"
}

func statusClass(success bool) string {
	if success {
		return "success"
	}
	return "error"
}

func formatPercent(value float64) string {
	return fmt.Sprintf("%.1f%%", value)
}

// HTML report template
const htmlReportTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Site Monitor Export Report</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; 
            margin: 0; 
            padding: 20px; 
            background: #f8fafc; 
            color: #1e293b;
        }
        .container { 
            max-width: 1200px; 
            margin: 0 auto; 
            background: white; 
            border-radius: 8px; 
            box-shadow: 0 1px 3px rgba(0,0,0,0.1);
        }
        .header { 
            background: #2563eb; 
            color: white; 
            padding: 30px; 
            border-radius: 8px 8px 0 0; 
        }
        .header h1 { 
            margin: 0; 
            font-size: 2rem; 
        }
        .metadata { 
            font-size: 0.9rem; 
            opacity: 0.9; 
            margin-top: 10px; 
        }
        .content { 
            padding: 30px; 
        }
        .stats-grid { 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); 
            gap: 20px; 
            margin-bottom: 30px; 
        }
        .stat-card { 
            background: #f1f5f9; 
            padding: 20px; 
            border-radius: 6px; 
            text-align: center; 
        }
        .stat-value { 
            font-size: 2rem; 
            font-weight: bold; 
            color: #2563eb; 
        }
        .stat-label { 
            color: #64748b; 
            font-size: 0.9rem; 
            margin-top: 5px; 
        }
        .section { 
            margin-bottom: 30px; 
        }
        .section h2 { 
            border-bottom: 2px solid #e2e8f0; 
            padding-bottom: 10px; 
            color: #1e293b; 
        }
        table { 
            width: 100%; 
            border-collapse: collapse; 
            margin-top: 15px; 
        }
        th, td { 
            padding: 12px; 
            text-align: left; 
            border-bottom: 1px solid #e2e8f0; 
        }
        th { 
            background: #f8fafc; 
            font-weight: 600; 
            color: #475569; 
        }
        .success { color: #059669; }
        .error { color: #dc2626; }
        .site-stats { 
            display: grid; 
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); 
            gap: 20px; 
        }
        .site-card { 
            border: 1px solid #e2e8f0; 
            border-radius: 6px; 
            padding: 20px; 
        }
        .site-card h3 { 
            margin: 0 0 15px 0; 
            color: #1e293b; 
        }
        .history-table { 
            max-height: 400px; 
            overflow-y: auto; 
        }
        @media (max-width: 768px) {
            .stats-grid { grid-template-columns: 1fr; }
            .site-stats { grid-template-columns: 1fr; }
            table { font-size: 0.9rem; }
            th, td { padding: 8px; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Site Monitor Export Report</h1>
            <div class="metadata">
                Generated on {{ formatTime .Metadata.GeneratedAt }} • 
                {{ .Metadata.TotalRecords }} records • 
                {{ len .Metadata.SitesIncluded }} sites •
                {{ formatTime .Metadata.TimeRange.From }} to {{ formatTime .Metadata.TimeRange.To }}
            </div>
        </div>

        <div class="content">
            {{if .Stats}}
            <div class="section">
                <h2>📊 Overview</h2>
                <div class="stats-grid">
                    <div class="stat-card">
                        <div class="stat-value">{{ .Stats.TotalSites }}</div>
                        <div class="stat-label">Total Sites</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-value">{{ .Stats.TotalChecks }}</div>
                        <div class="stat-label">Total Checks</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-value">{{ formatPercent .Stats.OverallUptime }}</div>
                        <div class="stat-label">Overall Uptime</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-value">{{ formatDuration .Stats.AvgResponseTime }}</div>
                        <div class="stat-label">Avg Response Time</div>
                    </div>
                </div>
            </div>

            {{if .Stats.SiteStats}}
            <div class="section">
                <h2>🌐 Site Statistics</h2>
                <div class="site-stats">
                    {{range $siteName, $stats := .Stats.SiteStats}}
                    <div class="site-card">
                        <h3>{{ $siteName }}</h3>
                        <table>
                            <tr><td>Uptime</td><td>{{ formatPercent $stats.SuccessRate }}</td></tr>
                            <tr><td>Total Checks</td><td>{{ $stats.TotalChecks }}</td></tr>
                            <tr><td>Successful</td><td>{{ $stats.SuccessfulChecks }}</td></tr>
                            <tr><td>Failed</td><td>{{ $stats.FailedChecks }}</td></tr>
                            <tr><td>Avg Response</td><td>{{ formatDuration $stats.AvgResponseTime }}</td></tr>
                            <tr><td>Last Check</td><td>{{ formatTime $stats.LastCheck }}</td></tr>
                        </table>
                    </div>
                    {{end}}
                </div>
            </div>
            {{end}}
            {{end}}

            <div class="section">
                <h2>📋 History ({{ len .History }} records)</h2>
                <div class="history-table">
                    <table>
                        <thead>
                            <tr>
                                <th>Time</th>
                                <th>Site</th>
                                <th>Status</th>
                                <th>Response Time</th>
                                <th>Error</th>
                            </tr>
                        </thead>
                        <tbody>
                            {{range .FormattedHistory}}
                            <tr>
                                <td>{{ .FormattedTimestamp }}</td>
                                <td>{{ .SiteName }}</td>
                                <td class="{{ .StatusClass }}">{{ .StatusIcon }} 
                                    {{if .Success}}OK ({{ .Status }}){{else}}ERROR{{if ne .Status 0}} ({{ .Status }}){{end}}{{end}}
                                </td>
                                <td>{{ .FormattedResponseTime }}</td>
                                <td class="error">{{ .Error }}</td>
                            </tr>
                            {{end}}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>
</body>
</html>`
