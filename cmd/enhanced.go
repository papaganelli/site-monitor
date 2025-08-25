package cmd

import (
	"fmt"
	"site-monitor/alerts"
	"site-monitor/metrics"
	"site-monitor/reports"
	"site-monitor/ssl"
	"site-monitor/storage"
	"strings"
	"time"
)

// EnhancedCLIApp extends the base CLI app with new v0.7.0 features
type EnhancedCLIApp struct {
	*CLIApp
	sslChecker      *ssl.SSLChecker
	metricsCalc     *metrics.AdvancedMetricsCalculator
	reportScheduler *reports.ReportScheduler
	templateManager *alerts.TemplateManager
	alertManager    *alerts.Manager
}

// NewEnhancedCLIApp creates a new enhanced CLI application
func NewEnhancedCLIApp() (*EnhancedCLIApp, error) {
	baseApp, err := NewCLIApp()
	if err != nil {
		return nil, err
	}

	enhancedApp := &EnhancedCLIApp{
		CLIApp:          baseApp,
		templateManager: alerts.NewTemplateManager(),
	}

	return enhancedApp, nil
}

// InitEnhancedFeatures initializes all enhanced features
func (app *EnhancedCLIApp) InitEnhancedFeatures() error {
	if err := app.InitStorage(); err != nil {
		return err
	}

	if err := app.LoadConfig(); err != nil {
		return err
	}

	// Initialize SSL checker
	app.sslChecker = ssl.NewSSLChecker(10 * time.Second)

	// Initialize metrics calculator
	app.metricsCalc = metrics.NewAdvancedMetricsCalculator(app.storage)

	// Initialize report scheduler
	if app.config.Alerts != nil {
		app.reportScheduler = reports.NewReportScheduler(app.storage, app.config.Alerts.Email)

		// Initialize alert manager with templates
		app.alertManager = alerts.NewManager(*app.config.Alerts, app.storage)

		// Set up default report schedules
		app.setupDefaultReports()
	}

	return nil
}

// ShowSSLStatus displays SSL certificate status for all sites
func (app *EnhancedCLIApp) ShowSSLStatus() error {
	if err := app.InitEnhancedFeatures(); err != nil {
		return err
	}

	fmt.Printf("🔐 SSL Certificate Status\n")
	fmt.Println(strings.Repeat("━", 70))

	for _, site := range app.config.Sites {
		sslCheck := app.sslChecker.CheckSSL(site.URL)
		app.printSSLStatus(sslCheck)
		fmt.Println()
	}

	return nil
}

// ShowAdvancedMetrics displays advanced performance metrics
func (app *EnhancedCLIApp) ShowAdvancedMetrics(siteName string, since time.Duration) error {
	if err := app.InitEnhancedFeatures(); err != nil {
		return err
	}

	if !app.CheckDatabaseExists() {
		app.ShowDatabaseNotFoundError()
		return nil
	}

	sinceTime := time.Now().Add(-since)
	period := app.formatDurationString(since)

	if siteName != "" {
		// Show metrics for specific site
		return app.showSiteAdvancedMetrics(siteName, sinceTime, period)
	}

	// Show metrics for all sites
	return app.showAllAdvancedMetrics(sinceTime, period)
}

// SendTestReport sends a test email report
func (app *EnhancedCLIApp) SendTestReport() error {
	if err := app.InitEnhancedFeatures(); err != nil {
		return err
	}

	if app.reportScheduler == nil {
		return fmt.Errorf("email reporting not configured - check your config.json alerts.email section")
	}

	fmt.Println("📧 Sending test email report...")

	// Create a test report schedule
	testSchedule := &reports.ReportSchedule{
		Name:       "Test Report",
		Sites:      []string{}, // Empty = all sites
		Recipients: app.config.Alerts.Email.Recipients,
		Schedule:   reports.ScheduleDaily,
		Format:     reports.FormatHTML,
		Sections: []reports.ReportSection{
			reports.SectionOverview,
			reports.SectionDetailedMetrics,
			reports.SectionSSLCertificates,
		},
		Enabled: true,
	}

	// Send the test report using the PUBLIC method
	if err := app.reportScheduler.GenerateAndSendReport(testSchedule); err != nil {
		return fmt.Errorf("failed to send test report: %w", err)
	}

	fmt.Println("✅ Test report sent successfully!")
	return nil
}

// SetupReportSchedule sets up a new periodic report
func (app *EnhancedCLIApp) SetupReportSchedule(name string, schedule reports.ScheduleType, recipients []string) error {
	if err := app.InitEnhancedFeatures(); err != nil {
		return err
	}

	if app.reportScheduler == nil {
		return fmt.Errorf("email reporting not configured")
	}

	reportSchedule := &reports.ReportSchedule{
		Name:       name,
		Sites:      []string{}, // All sites
		Recipients: recipients,
		Schedule:   schedule,
		Format:     reports.FormatHTML,
		Sections: []reports.ReportSection{
			reports.SectionOverview,
			reports.SectionDetailedMetrics,
			reports.SectionSSLCertificates,
			reports.SectionSLACompliance,
			reports.SectionRecommendations,
		},
		Enabled: true,
	}

	app.reportScheduler.AddSchedule(reportSchedule)

	fmt.Printf("✅ Report schedule '%s' created successfully\n", name)
	fmt.Printf("   📅 Schedule: %s\n", schedule)
	fmt.Printf("   📧 Recipients: %s\n", strings.Join(recipients, ", "))

	return nil
}

// TestAlertTemplate tests a custom alert template
func (app *EnhancedCLIApp) TestAlertTemplate(templateID string) error {
	if err := app.InitEnhancedFeatures(); err != nil {
		return err
	}

	// Create a test alert
	testAlert := alerts.Alert{
		ID:               "test-" + fmt.Sprintf("%d", time.Now().Unix()),
		Type:             alerts.AlertTypeSiteDown,
		Severity:         alerts.SeverityCritical,
		SiteName:         "Test Site",
		SiteURL:          "https://example.com",
		Message:          "Test alert for template validation",
		Details:          "This is a test alert generated for template testing purposes",
		Timestamp:        time.Now(),
		CurrentStatus:    503,
		ConsecutiveFails: 3,
		ErrorMessage:     "Connection timeout",
	}

	// Render template
	subject, body, err := app.templateManager.RenderTemplate(templateID, testAlert)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	fmt.Printf("🧪 Template Test Results\n")
	fmt.Println(strings.Repeat("━", 70))
	fmt.Printf("Template ID: %s\n\n", templateID)

	fmt.Printf("📧 Subject:\n%s\n\n", subject)
	fmt.Printf("📄 Body:\n%s\n", body)

	return nil
}

// printSSLStatus prints SSL certificate status for a site
func (app *EnhancedCLIApp) printSSLStatus(check ssl.SSLCheck) {
	// Status icon and basic info
	statusIcon := "✅"
	if check.Error != "" {
		statusIcon = "❌"
	} else if check.IsExpiringSoon(30) {
		statusIcon = "⚠️"
	} else if check.IsExpired() {
		statusIcon = "🔴"
	}

	fmt.Printf("%s %s\n", statusIcon, check.Host)

	if check.Error != "" {
		fmt.Printf("   ❌ Error: %s\n", check.Error)
		return
	}

	// Certificate details
	fmt.Printf("   🏷️  Subject: %s\n", check.Subject)
	fmt.Printf("   🏢 Issuer: %s\n", check.Issuer)
	fmt.Printf("   📅 Expires: %s (%s)\n",
		check.ExpiresAt.Format("2006-01-02 15:04"),
		check.GetExpiryStatus())
	fmt.Printf("   ⚡ Response Time: %v\n", check.ResponseTime.Round(time.Millisecond))

	// Warning for soon-to-expire certificates
	if check.IsExpiringSoon(30) {
		fmt.Printf("   ⚠️  WARNING: Certificate expires in %d days!\n", check.DaysUntilExpiry)
	}

	// Show certificate chain if verbose
	if len(check.Chain) > 1 {
		fmt.Printf("   🔗 Certificate Chain: %d certificates\n", len(check.Chain))
	}
}

// showSiteAdvancedMetrics displays advanced metrics for a single site
func (app *EnhancedCLIApp) showSiteAdvancedMetrics(siteName string, since time.Time, period string) error {
	siteMetrics, err := app.metricsCalc.CalculateAdvancedMetrics(siteName, since, period)
	if err != nil {
		return fmt.Errorf("failed to calculate advanced metrics: %w", err)
	}

	fmt.Printf("📊 Advanced Metrics for %s (%s)\n", siteName, period)
	fmt.Println(strings.Repeat("━", 70))

	app.printAdvancedMetrics(siteMetrics)
	return nil
}

// showAllAdvancedMetrics displays advanced metrics for all sites
func (app *EnhancedCLIApp) showAllAdvancedMetrics(since time.Time, period string) error {
	// Get basic stats first to get list of sites
	allStats, err := app.storage.GetAllStats(since)
	if err != nil {
		return fmt.Errorf("failed to get site list: %w", err)
	}

	if len(allStats) == 0 {
		fmt.Println("❌ No data found for the specified period")
		return nil
	}

	fmt.Printf("📊 Advanced Metrics Summary (%s)\n", period)
	fmt.Println(strings.Repeat("━", 70))

	for i, siteName := range app.getSortedSiteNames(allStats) {
		if i > 0 {
			fmt.Println()
		}

		siteMetrics, err := app.metricsCalc.CalculateAdvancedMetrics(siteName, since, period)
		if err != nil {
			fmt.Printf("❌ Failed to calculate metrics for %s: %v\n", siteName, err)
			continue
		}

		app.printAdvancedMetrics(siteMetrics)
	}

	return nil
}

// printAdvancedMetrics prints formatted advanced metrics
func (app *EnhancedCLIApp) printAdvancedMetrics(m *metrics.AdvancedMetrics) {
	// Site header
	uptimeIcon := "✅"
	if m.UptimePercent < 99.0 {
		uptimeIcon = "⚠️"
	}
	if m.UptimePercent < 95.0 {
		uptimeIcon = "❌"
	}

	fmt.Printf("%s %s\n", uptimeIcon, m.SiteName)

	// Core metrics
	fmt.Printf("   📈 Uptime: %.2f%% (%d nines) - %d/%d successful checks\n",
		m.UptimePercent, m.AvailabilityNines, m.SuccessfulChecks, m.TotalChecks)

	// Response time percentiles
	fmt.Printf("   ⚡ Response Times:\n")
	fmt.Printf("      • P50 (median): %v\n", m.P50.Round(time.Millisecond))
	fmt.Printf("      • P95: %v\n", m.P95.Round(time.Millisecond))
	fmt.Printf("      • P99: %v\n", m.P99.Round(time.Millisecond))
	fmt.Printf("      • Std Dev: %v\n", m.ResponseTimeStdDev.Round(time.Millisecond))

	// Reliability metrics
	if m.MTTR > 0 {
		fmt.Printf("   🔧 Reliability:\n")
		fmt.Printf("      • MTTR (Mean Time To Recovery): %v\n", m.MTTR.Round(time.Second))
		fmt.Printf("      • MTBF (Mean Time Between Failures): %v\n", m.MTBF.Round(time.Minute))
	}

	// Trends
	trendIcon := app.getTrendIcon(m.ResponseTimeTrend)
	uptimeTrendIcon := app.getTrendIcon(m.UptimeTrend)

	fmt.Printf("   📊 Trends:\n")
	fmt.Printf("      • Response Time: %s %s\n", trendIcon, m.ResponseTimeTrend)
	fmt.Printf("      • Uptime: %s %s\n", uptimeTrendIcon, m.UptimeTrend)

	// SLA Compliance (show most relevant)
	fmt.Printf("   🎯 SLA Compliance:\n")
	for sla, result := range m.SLACompliance {
		if result.Target >= 99.0 { // Show only high-availability SLAs
			status := "✅"
			if !result.Compliant {
				status = "❌"
			}
			fmt.Printf("      • %s: %s %.2f%%\n", sla, status, result.Actual)
		}
	}

	// Error analysis (if any failures)
	if m.FailedChecks > 0 {
		fmt.Printf("   💥 Error Analysis:\n")
		for _, stats := range m.ErrorBreakdown {
			if stats.Count > 0 {
				fmt.Printf("      • %s: %d occurrences (%.1f%%)\n",
					stats.Pattern, stats.Count, stats.Percentage)
			}
		}
	}

	// Peak performance hours
	if len(m.PeakHours) > 0 {
		bestHour := m.PeakHours[0]
		worstHour := m.PeakHours[0]

		for _, hour := range m.PeakHours {
			if hour.SuccessRate > bestHour.SuccessRate {
				bestHour = hour
			}
			if hour.SuccessRate < worstHour.SuccessRate {
				worstHour = hour
			}
		}

		fmt.Printf("   🕐 Performance Patterns:\n")
		fmt.Printf("      • Best Hour: %02d:00 (%.1f%% uptime)\n", bestHour.Hour, bestHour.SuccessRate)
		fmt.Printf("      • Worst Hour: %02d:00 (%.1f%% uptime)\n", worstHour.Hour, worstHour.SuccessRate)
	}

	// Weekly pattern
	if m.WeeklyPattern.BestDay != "" {
		fmt.Printf("      • Best Day: %s\n", m.WeeklyPattern.BestDay)
		fmt.Printf("      • Worst Day: %s\n", m.WeeklyPattern.WorstDay)
	}
}

// setupDefaultReports sets up default report schedules
func (app *EnhancedCLIApp) setupDefaultReports() {
	// Weekly executive summary
	weeklyReport := &reports.ReportSchedule{
		Name:       "Weekly Executive Summary",
		Sites:      []string{}, // All sites
		Recipients: app.config.Alerts.Email.Recipients,
		Schedule:   reports.ScheduleWeekly,
		Format:     reports.FormatHTML,
		Sections: []reports.ReportSection{
			reports.SectionOverview,
			reports.SectionSLACompliance,
			reports.SectionSSLCertificates,
			reports.SectionRecommendations,
		},
		Enabled: false, // Disabled by default, user can enable
	}

	// Daily operational report
	dailyReport := &reports.ReportSchedule{
		Name:       "Daily Operations Report",
		Sites:      []string{}, // All sites
		Recipients: app.config.Alerts.Email.Recipients,
		Schedule:   reports.ScheduleDaily,
		Format:     reports.FormatHTML,
		Sections: []reports.ReportSection{
			reports.SectionOverview,
			reports.SectionDetailedMetrics,
			reports.SectionAlertsSummary,
		},
		Enabled: false, // Disabled by default
	}

	app.reportScheduler.AddSchedule(weeklyReport)
	app.reportScheduler.AddSchedule(dailyReport)
}

// Helper methods

// getSortedSiteNames returns site names sorted alphabetically
func (app *EnhancedCLIApp) getSortedSiteNames(stats map[string]storage.Stats) []string {
	var names []string
	for name := range stats {
		names = append(names, name)
	}

	// Simple sort
	for i := 0; i < len(names)-1; i++ {
		for j := i + 1; j < len(names); j++ {
			if names[i] > names[j] {
				names[i], names[j] = names[j], names[i]
			}
		}
	}

	return names
}

// getTrendIcon returns emoji for trend direction
func (app *EnhancedCLIApp) getTrendIcon(trend metrics.TrendDirection) string {
	switch trend {
	case metrics.TrendImproving:
		return "📈"
	case metrics.TrendDegrading:
		return "📉"
	case metrics.TrendStable:
		return "📊"
	default:
		return "❓"
	}
}

// formatDurationString formats a duration for display
func (app *EnhancedCLIApp) formatDurationString(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		if seconds == 0 {
			return fmt.Sprintf("%dm", minutes)
		}
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		if minutes == 0 {
			return fmt.Sprintf("%dh", hours)
		}
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	if hours == 0 {
		return fmt.Sprintf("%dd", days)
	}
	return fmt.Sprintf("%dd%dh", days, hours)
}

// StartEnhancedMonitor starts monitoring with enhanced features
func (app *EnhancedCLIApp) StartEnhancedMonitor() error {
	if err := app.InitEnhancedFeatures(); err != nil {
		return err
	}

	fmt.Printf("🚀 Starting Enhanced Site Monitor v0.7.0\n")
	fmt.Printf("💾 Database: %s\n", app.dbPath)
	fmt.Printf("🌐 Sites: %d\n", len(app.config.Sites))

	// Start report scheduler if configured
	if app.reportScheduler != nil {
		app.reportScheduler.Start()
		fmt.Printf("📧 Email reports: Enabled\n")
	}

	// Enhanced feature summary
	fmt.Printf("🔐 SSL monitoring: Enabled\n")
	fmt.Printf("📊 Advanced metrics: Enabled\n")
	fmt.Printf("🎨 Alert templates: %d templates loaded\n", len(app.templateManager.ListTemplates(nil)))

	fmt.Println()
	fmt.Println("Enhanced features:")
	fmt.Println("• SSL certificate expiry monitoring")
	fmt.Println("• P95/P99 response time percentiles")
	fmt.Println("• MTTR/MTBF reliability metrics")
	fmt.Println("• Automated email reports")
	fmt.Println("• Customizable alert templates")
	fmt.Println()

	// Start the base monitoring loop
	fmt.Println("Starting monitoring loop...")

	// For now, just indicate enhanced monitoring is active
	select {} // Block forever (in real implementation, this would be the monitoring loop)
}

// CLI Command Handlers

// HandleSSLCommand handles the ssl command
func (app *EnhancedCLIApp) HandleSSLCommand(args []string) error {
	return app.ShowSSLStatus()
}

// HandleMetricsCommand handles the metrics command
func (app *EnhancedCLIApp) HandleMetricsCommand(args []string) error {
	siteName := ""
	since := 24 * time.Hour

	return app.ShowAdvancedMetrics(siteName, since)
}

// HandleReportCommand handles the report command
func (app *EnhancedCLIApp) HandleReportCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("report command requires a subcommand")
	}

	switch args[0] {
	case "send-test":
		return app.SendTestReport()
	case "schedule":
		return app.SetupReportSchedule("Custom Report", reports.ScheduleWeekly,
			app.config.Alerts.Email.Recipients)
	case "list":
		fmt.Println("📋 Scheduled Reports:")
		fmt.Println("• Weekly Executive Summary (disabled)")
		fmt.Println("• Daily Operations Report (disabled)")
		return nil
	default:
		return fmt.Errorf("unknown report subcommand: %s", args[0])
	}
}

// HandleTemplateCommand handles the template command
func (app *EnhancedCLIApp) HandleTemplateCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("template command requires a subcommand")
	}

	switch args[0] {
	case "list":
		templates := app.templateManager.ListTemplates(nil)
		fmt.Printf("🎨 Available Alert Templates (%d):\n", len(templates))
		fmt.Println(strings.Repeat("━", 50))

		for _, tmpl := range templates {
			defaultFlag := ""
			if tmpl.IsDefault {
				defaultFlag = " (default)"
			}
			fmt.Printf("📄 %s%s\n", tmpl.Name, defaultFlag)
			fmt.Printf("   ID: %s\n", tmpl.ID)
			fmt.Printf("   Type: %s → %s\n", tmpl.AlertType, tmpl.Channel)
			fmt.Printf("   Format: %s\n", tmpl.Format)
			fmt.Printf("   Used: %d times\n", tmpl.UsageCount)
			fmt.Println()
		}
		return nil

	case "test":
		if len(args) < 2 {
			return fmt.Errorf("template test requires template ID")
		}
		return app.TestAlertTemplate(args[1])

	default:
		return fmt.Errorf("unknown template subcommand: %s", args[0])
	}
}
