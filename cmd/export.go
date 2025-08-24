package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"site-monitor/export"
	"strings"
	"time"
)

// ExportOptions contains options for the export command
type ExportCLIOptions struct {
	Format     string
	SiteName   string
	Since      time.Duration
	Until      *time.Time
	Limit      int
	OutputPath string
	Stats      bool
	Stdout     bool
}

// ShowExport handles the export command
func (app *CLIApp) ShowExport(opts ExportCLIOptions) error {
	if err := app.InitStorage(); err != nil {
		return err
	}
	defer app.Close()

	if !app.CheckDatabaseExists() {
		app.ShowDatabaseNotFoundError()
		return nil
	}

	// Parse and validate format
	format, err := app.parseExportFormat(opts.Format)
	if err != nil {
		return err
	}

	// Create exporter
	exporter := export.NewExporter(app.storage)

	// Build export options
	exportOpts := export.ExportOptions{
		Format:       format,
		SiteName:     opts.SiteName,
		Since:        opts.Since,
		Until:        opts.Until,
		Limit:        opts.Limit,
		IncludeStats: opts.Stats,
		OutputPath:   opts.OutputPath,
	}

	// Show export info
	app.showExportInfo(exportOpts)

	// Export data
	data, err := exporter.Export(exportOpts)
	if err != nil {
		return fmt.Errorf("export failed: %w", err)
	}

	// Get formatter
	formatter, err := export.GetFormatter(format)
	if err != nil {
		return fmt.Errorf("failed to get formatter: %w", err)
	}

	// Determine output destination
	var outputPath string
	if opts.Stdout {
		outputPath = "stdout"
	} else {
		outputPath = app.generateOutputPath(opts, formatter)
	}

	// Export to file or stdout
	if err := app.writeExport(data, formatter, outputPath); err != nil {
		return fmt.Errorf("failed to write export: %w", err)
	}

	// Show success message
	app.showExportSuccess(data, outputPath)

	return nil
}

// parseExportFormat parses and validates the export format
func (app *CLIApp) parseExportFormat(formatStr string) (export.ExportFormat, error) {
	if formatStr == "" {
		return export.FormatJSON, nil // Default format
	}

	formatStr = strings.ToLower(formatStr)
	switch formatStr {
	case "json":
		return export.FormatJSON, nil
	case "csv":
		return export.FormatCSV, nil
	case "html":
		return export.FormatHTML, nil
	default:
		return "", fmt.Errorf("unsupported format '%s'. Supported formats: json, csv, html", formatStr)
	}
}

// generateOutputPath generates an appropriate output file path
func (app *CLIApp) generateOutputPath(opts ExportCLIOptions, formatter export.Formatter) string {
	if opts.OutputPath != "" {
		// Ensure correct extension
		ext := formatter.FileExtension()
		if !strings.HasSuffix(opts.OutputPath, ext) {
			return opts.OutputPath + ext
		}
		return opts.OutputPath
	}

	// Generate automatic filename
	timestamp := time.Now().Format("20060102_150405")

	var filename string
	if opts.SiteName != "" {
		// Sanitize site name for filename
		safeName := strings.ReplaceAll(opts.SiteName, " ", "_")
		safeName = strings.ReplaceAll(safeName, "/", "_")
		filename = fmt.Sprintf("site-monitor_%s_%s", safeName, timestamp)
	} else {
		filename = fmt.Sprintf("site-monitor_export_%s", timestamp)
	}

	return filename + formatter.FileExtension()
}

// writeExport writes the export data to file or stdout
func (app *CLIApp) writeExport(data *export.ExportData, formatter export.Formatter, outputPath string) error {
	if outputPath == "stdout" {
		return formatter.Format(data, os.Stdout)
	}

	// Create directory if needed
	dir := filepath.Dir(outputPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Write data
	return formatter.Format(data, file)
}

// showExportInfo displays information about the export operation
func (app *CLIApp) showExportInfo(opts export.ExportOptions) {
	fmt.Printf("📊 Exporting Site Monitor Data\n")
	fmt.Println(strings.Repeat("━", 40))

	// Format
	fmt.Printf("📄 Format: %s (%s)\n",
		strings.ToUpper(string(opts.Format)),
		export.FormatDescription(opts.Format))

	// Site filter
	if opts.SiteName != "" {
		fmt.Printf("🌐 Site: %s\n", opts.SiteName)
	} else {
		fmt.Printf("🌐 Sites: All monitored sites\n")
	}

	// Time range
	since := opts.Since
	if since == 0 {
		since = 24 * time.Hour // Default
	}
	fmt.Printf("⏱️  Time Range: Last %s\n", formatDuration(since))

	if opts.Until != nil {
		fmt.Printf("🕐 Until: %s\n", opts.Until.Format("2006-01-02 15:04:05"))
	}

	// Limit
	if opts.Limit > 0 {
		fmt.Printf("🔢 Limit: %d records\n", opts.Limit)
	}

	// Statistics
	if opts.IncludeStats {
		fmt.Printf("📈 Statistics: Included\n")
	}

	fmt.Println()
}

// showExportSuccess displays success information after export
func (app *CLIApp) showExportSuccess(data *export.ExportData, outputPath string) {
	fmt.Printf("✅ Export completed successfully!\n")
	fmt.Println(strings.Repeat("━", 40))

	if outputPath != "stdout" {
		fmt.Printf("📁 Output File: %s\n", outputPath)

		// Show file size
		if stat, err := os.Stat(outputPath); err == nil {
			fmt.Printf("📏 File Size: %s\n", formatBytes(stat.Size()))
		}
	}

	fmt.Printf("📊 Records Exported: %d\n", data.Metadata.TotalRecords)
	fmt.Printf("🌐 Sites Included: %d (%s)\n",
		len(data.Metadata.SitesIncluded),
		strings.Join(data.Metadata.SitesIncluded, ", "))

	fmt.Printf("⏰ Time Range: %s to %s\n",
		data.Metadata.TimeRange.From.Format("2006-01-02 15:04"),
		data.Metadata.TimeRange.To.Format("2006-01-02 15:04"))

	if data.Stats != nil {
		fmt.Printf("📈 Overall Uptime: %.1f%%\n", data.Stats.OverallUptime)
		if data.Stats.AvgResponseTime > 0 {
			fmt.Printf("⚡ Avg Response Time: %s\n", formatDuration(data.Stats.AvgResponseTime))
		}
	}

	fmt.Printf("🕐 Generated: %s\n", data.Metadata.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()

	// Show usage suggestions
	app.showExportSuggestions(data.Metadata.Format, outputPath)
}

// showExportSuggestions shows helpful suggestions for using the exported data
func (app *CLIApp) showExportSuggestions(format export.ExportFormat, outputPath string) {
	fmt.Printf("💡 Usage Suggestions:\n")

	switch format {
	case export.FormatJSON:
		fmt.Printf("   • Import into analysis tools or custom applications\n")
		fmt.Printf("   • Use with jq for command-line processing: jq '.' %s\n", outputPath)
		fmt.Printf("   • Load into Python/Node.js for data analysis\n")

	case export.FormatCSV:
		fmt.Printf("   • Open in Excel, Google Sheets, or LibreOffice Calc\n")
		fmt.Printf("   • Import into data analysis tools (R, Python pandas)\n")
		fmt.Printf("   • Use for creating custom reports and visualizations\n")

	case export.FormatHTML:
		if outputPath != "stdout" {
			fmt.Printf("   • Open in web browser: open %s\n", outputPath)
		}
		fmt.Printf("   • Share as a standalone report\n")
		fmt.Printf("   • Print or convert to PDF\n")
	}
	fmt.Println()
}

// ListExportFormats shows available export formats
func (app *CLIApp) ListExportFormats() {
	fmt.Printf("📊 Available Export Formats\n")
	fmt.Println(strings.Repeat("━", 40))

	formats := export.GetSupportedFormats()
	for _, format := range formats {
		fmt.Printf("• %-6s - %s\n", strings.ToUpper(string(format)), export.FormatDescription(format))
	}
	fmt.Println()

	fmt.Printf("💡 Usage Examples:\n")
	fmt.Printf("   site-monitor export --format json --output data.json\n")
	fmt.Printf("   site-monitor export --format csv --site \"My Site\" --since 7d\n")
	fmt.Printf("   site-monitor export --format html --stats --output report.html\n")
	fmt.Printf("   site-monitor export --format json --stdout | jq .\n")
	fmt.Println()
}

// Helper functions

// formatBytes formats byte size in human readable format
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
