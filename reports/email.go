package reports

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"site-monitor/config"
	"site-monitor/metrics"
	"site-monitor/ssl"
	"site-monitor/storage"
	"strings"
	"time"
)

// ReportScheduler manages periodic email reports
type ReportScheduler struct {
	storage     storage.Storage
	emailConfig config.EmailConfig
	sslChecker  *ssl.SSLChecker
	metricsCalc *metrics.AdvancedMetricsCalculator
	schedules   map[string]*ReportSchedule
	stopChannel chan struct{}
}

// ReportSchedule defines when and how to send reports
type ReportSchedule struct {
	Name       string          `json:"name"`
	Sites      []string        `json:"sites"` // Empty = all sites
	Recipients []string        `json:"recipients"`
	Schedule   ScheduleType    `json:"schedule"`
	Format     ReportFormat    `json:"format"`
	Sections   []ReportSection `json:"sections"`
	Enabled    bool            `json:"enabled"`
	LastSent   time.Time       `json:"last_sent"`
	NextDue    time.Time       `json:"next_due"`
	Template   string          `json:"template"`
}

// ScheduleType defines report frequency
type ScheduleType string

const (
	ScheduleDaily   ScheduleType = "daily"
	ScheduleWeekly  ScheduleType = "weekly"
	ScheduleMonthly ScheduleType = "monthly"
	ScheduleCustom  ScheduleType = "custom"
)

// ReportFormat defines output format
type ReportFormat string

const (
	FormatHTML ReportFormat = "html"
	FormatPDF  ReportFormat = "pdf"
	FormatCSV  ReportFormat = "csv"
)

// ReportSection defines what content to include
type ReportSection string

const (
	SectionOverview          ReportSection = "overview"
	SectionDetailedMetrics   ReportSection = "detailed_metrics"
	SectionSSLCertificates   ReportSection = "ssl_certificates"
	SectionAlertsSummary     ReportSection = "alerts_summary"
	SectionPerformanceTrends ReportSection = "performance_trends"
	SectionSLACompliance     ReportSection = "sla_compliance"
	SectionErrorAnalysis     ReportSection = "error_analysis"
	SectionRecommendations   ReportSection = "recommendations"
)

// ReportData contains all data for report generation
type ReportData struct {
	GeneratedAt time.Time `json:"generated_at"`
	Period      string    `json:"period"`
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`

	// Overview
	TotalSites    int     `json:"total_sites"`
	HealthySites  int     `json:"healthy_sites"`
	OverallUptime float64 `json:"overall_uptime"`
	TotalChecks   int64   `json:"total_checks"`

	// Site-specific data
	SiteMetrics map[string]*metrics.AdvancedMetrics `json:"site_metrics"`
	SSLChecks   map[string]*ssl.SSLCheck            `json:"ssl_checks"`

	// Aggregated insights
	TopPerformers   []SiteRanking    `json:"top_performers"`
	WorstPerformers []SiteRanking    `json:"worst_performers"`
	Recommendations []Recommendation `json:"recommendations"`

	// Executive summary
	ExecutiveSummary string `json:"executive_summary"`
}

// SiteRanking represents site performance ranking
type SiteRanking struct {
	SiteName        string        `json:"site_name"`
	UptimePercent   float64       `json:"uptime_percent"`
	AvgResponseTime time.Duration `json:"avg_response_time"`
	Score           float64       `json:"score"`
}

// Recommendation represents an actionable recommendation
type Recommendation struct {
	Type        RecommendationType `json:"type"`
	SiteName    string             `json:"site_name"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Priority    Priority           `json:"priority"`
	ActionItems []string           `json:"action_items"`
}

// RecommendationType defines types of recommendations
type RecommendationType string

const (
	RecommendationSSL         RecommendationType = "ssl"
	RecommendationPerformance RecommendationType = "performance"
	RecommendationReliability RecommendationType = "reliability"
	RecommendationSLA         RecommendationType = "sla"
)

// Priority defines recommendation priority
type Priority string

const (
	PriorityLow      Priority = "low"
	PriorityMedium   Priority = "medium"
	PriorityHigh     Priority = "high"
	PriorityCritical Priority = "critical"
)

// NewReportScheduler creates a new report scheduler
func NewReportScheduler(storage storage.Storage, emailConfig config.EmailConfig) *ReportScheduler {
	return &ReportScheduler{
		storage:     storage,
		emailConfig: emailConfig,
		sslChecker:  ssl.NewSSLChecker(10 * time.Second),
		metricsCalc: metrics.NewAdvancedMetricsCalculator(storage),
		schedules:   make(map[string]*ReportSchedule),
		stopChannel: make(chan struct{}),
	}
}

// AddSchedule adds a new report schedule
func (rs *ReportScheduler) AddSchedule(schedule *ReportSchedule) {
	rs.schedules[schedule.Name] = schedule
	rs.calculateNextDue(schedule)
}

// Start begins the report scheduler
func (rs *ReportScheduler) Start() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour) // Check every hour
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				rs.checkAndSendReports()
			case <-rs.stopChannel:
				return
			}
		}
	}()
}

// Stop stops the report scheduler
func (rs *ReportScheduler) Stop() {
	close(rs.stopChannel)
}

// checkAndSendReports checks if any reports are due and sends them
func (rs *ReportScheduler) checkAndSendReports() {
	now := time.Now()

	for _, schedule := range rs.schedules {
		if !schedule.Enabled {
			continue
		}

		if now.After(schedule.NextDue) {
			if err := rs.GenerateAndSendReport(schedule); err != nil {
				fmt.Printf("❌ Failed to send report '%s': %v\n", schedule.Name, err)
			} else {
				fmt.Printf("✅ Report '%s' sent successfully\n", schedule.Name)
				schedule.LastSent = now
				rs.calculateNextDue(schedule)
			}
		}
	}
}

// GenerateAndSendReport generates and sends a specific report (public method)
func (rs *ReportScheduler) GenerateAndSendReport(schedule *ReportSchedule) error {
	return rs.generateAndSendReport(schedule)
}

// generateAndSendReport generates and sends a specific report (private method)
func (rs *ReportScheduler) generateAndSendReport(schedule *ReportSchedule) error {
	// Determine time period based on schedule
	periodStart, periodEnd := rs.calculateReportPeriod(schedule.Schedule)

	// Generate report data
	reportData, err := rs.generateReportData(schedule, periodStart, periodEnd)
	if err != nil {
		return fmt.Errorf("failed to generate report data: %w", err)
	}

	// Generate report content based on format
	var content []byte
	var contentType string

	switch schedule.Format {
	case FormatHTML:
		content, err = rs.generateHTMLReport(reportData, schedule)
		contentType = "text/html"
	case FormatPDF:
		content, err = rs.generatePDFReport(reportData, schedule)
		contentType = "application/pdf"
	case FormatCSV:
		content, err = rs.generateCSVReport(reportData, schedule)
		contentType = "text/csv"
	default:
		content, err = rs.generateHTMLReport(reportData, schedule)
		contentType = "text/html"
	}

	if err != nil {
		return fmt.Errorf("failed to generate report content: %w", err)
	}

	// Send email
	subject := rs.generateSubject(schedule, reportData)
	return rs.sendEmailReport(schedule.Recipients, subject, content, contentType, schedule.Name)
}

// generateReportData compiles all data needed for the report
func (rs *ReportScheduler) generateReportData(schedule *ReportSchedule, periodStart, periodEnd time.Time) (*ReportData, error) {
	reportData := &ReportData{
		GeneratedAt: time.Now(),
		Period:      rs.formatPeriod(schedule.Schedule),
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		SiteMetrics: make(map[string]*metrics.AdvancedMetrics),
		SSLChecks:   make(map[string]*ssl.SSLCheck),
	}

	// Get list of sites to report on
	sites, err := rs.getSitesToReport(schedule)
	if err != nil {
		return nil, err
	}

	reportData.TotalSites = len(sites)
	var totalChecks int64
	var totalSuccessfulChecks int64
	var healthySites int

	// Generate metrics for each site
	for _, siteName := range sites {
		// Advanced metrics
		siteMetrics, err := rs.metricsCalc.CalculateAdvancedMetrics(siteName, periodStart, reportData.Period)
		if err != nil {
			fmt.Printf("⚠️ Failed to calculate metrics for %s: %v\n", siteName, err)
			continue
		}

		reportData.SiteMetrics[siteName] = siteMetrics

		totalChecks += siteMetrics.TotalChecks
		totalSuccessfulChecks += siteMetrics.SuccessfulChecks

		if siteMetrics.UptimePercent >= 99.0 {
			healthySites++
		}

		// SSL check if enabled
		if rs.shouldIncludeSection(schedule, SectionSSLCertificates) {
			// Get site URL (this would need to be stored or passed somehow)
			siteURL := rs.getSiteURL(siteName) // You'd need to implement this
			if siteURL != "" {
				sslCheck := rs.sslChecker.CheckSSL(siteURL)
				reportData.SSLChecks[siteName] = &sslCheck
			}
		}
	}

	reportData.HealthySites = healthySites
	reportData.TotalChecks = totalChecks

	if totalChecks > 0 {
		reportData.OverallUptime = float64(totalSuccessfulChecks) / float64(totalChecks) * 100
	}

	// Generate rankings and recommendations
	reportData.TopPerformers = rs.generateTopPerformers(reportData.SiteMetrics, 5)
	reportData.WorstPerformers = rs.generateWorstPerformers(reportData.SiteMetrics, 5)
	reportData.Recommendations = rs.generateRecommendations(reportData)
	reportData.ExecutiveSummary = rs.generateExecutiveSummary(reportData)

	return reportData, nil
}

// generateHTMLReport creates an HTML report
func (rs *ReportScheduler) generateHTMLReport(data *ReportData, schedule *ReportSchedule) ([]byte, error) {
	tmplContent := rs.getHTMLTemplate(schedule)

	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"formatDuration": formatDuration,
		"formatPercent":  formatPercent,
		"formatTime":     formatTime,
		"colorForUptime": colorForUptime,
		"priorityIcon":   priorityIcon,
	}).Parse(tmplContent)

	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// generatePDFReport creates a PDF report (placeholder - would need PDF library)
func (rs *ReportScheduler) generatePDFReport(data *ReportData, schedule *ReportSchedule) ([]byte, error) {
	// For now, generate HTML and note that PDF generation would require
	// additional libraries like wkhtmltopdf or a Go PDF library
	htmlContent, err := rs.generateHTMLReport(data, schedule)
	if err != nil {
		return nil, err
	}

	// TODO: Convert HTML to PDF using a library like:
	// - github.com/SebastiaanKlippert/go-wkhtmltopdf
	// - github.com/jung-kurt/gofpdf
	// - github.com/go-pdf/fpdf

	// For now, return HTML content
	return htmlContent, nil
}

// generateCSVReport creates a CSV report
func (rs *ReportScheduler) generateCSVReport(data *ReportData, schedule *ReportSchedule) ([]byte, error) {
	var buf bytes.Buffer

	// CSV Header
	buf.WriteString("Site Name,Uptime %,Avg Response Time (ms),P95 Response Time (ms),Total Checks,Failed Checks,MTTR (min),MTBF (hours),SSL Days Until Expiry\n")

	// Data rows
	for siteName, metrics := range data.SiteMetrics {
		sslDays := "N/A"
		if sslCheck, exists := data.SSLChecks[siteName]; exists && sslCheck.Valid {
			sslDays = fmt.Sprintf("%d", sslCheck.DaysUntilExpiry)
		}

		buf.WriteString(fmt.Sprintf("%s,%.2f,%.0f,%.0f,%d,%d,%.1f,%.1f,%s\n",
			siteName,
			metrics.UptimePercent,
			float64(metrics.P50.Milliseconds()),
			float64(metrics.P95.Milliseconds()),
			metrics.TotalChecks,
			metrics.FailedChecks,
			metrics.MTTR.Minutes(),
			metrics.MTBF.Hours(),
			sslDays,
		))
	}

	return buf.Bytes(), nil
}

// calculateReportPeriod determines the time period for the report
func (rs *ReportScheduler) calculateReportPeriod(schedule ScheduleType) (time.Time, time.Time) {
	now := time.Now()

	switch schedule {
	case ScheduleDaily:
		start := now.AddDate(0, 0, -1)
		return start, now
	case ScheduleWeekly:
		start := now.AddDate(0, 0, -7)
		return start, now
	case ScheduleMonthly:
		start := now.AddDate(0, -1, 0)
		return start, now
	default:
		start := now.AddDate(0, 0, -7)
		return start, now
	}
}

// calculateNextDue calculates when the next report is due
func (rs *ReportScheduler) calculateNextDue(schedule *ReportSchedule) {
	now := time.Now()

	switch schedule.Schedule {
	case ScheduleDaily:
		schedule.NextDue = now.AddDate(0, 0, 1)
	case ScheduleWeekly:
		schedule.NextDue = now.AddDate(0, 0, 7)
	case ScheduleMonthly:
		schedule.NextDue = now.AddDate(0, 1, 0)
	default:
		schedule.NextDue = now.AddDate(0, 0, 1)
	}

	// Set to specific time (e.g., 9:00 AM)
	schedule.NextDue = time.Date(
		schedule.NextDue.Year(),
		schedule.NextDue.Month(),
		schedule.NextDue.Day(),
		9, 0, 0, 0,
		schedule.NextDue.Location(),
	)
}

// sendEmailReport sends the report via email
func (rs *ReportScheduler) sendEmailReport(recipients []string, subject string, content []byte, contentType string, reportName string) error {
	// Parse SMTP server and port
	serverParts := strings.Split(rs.emailConfig.SMTPServer, ":")
	if len(serverParts) != 2 {
		return fmt.Errorf("invalid SMTP server format: %s", rs.emailConfig.SMTPServer)
	}

	server := serverParts[0]
	port := serverParts[1]

	// Set up authentication
	auth := smtp.PlainAuth("", rs.emailConfig.Username, rs.emailConfig.Password, server)

	// Prepare message
	from := rs.emailConfig.From
	if from == "" {
		from = rs.emailConfig.Username
	}

	for _, recipient := range recipients {
		msg := fmt.Sprintf("To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: %s; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
			recipient, from, subject, contentType, string(content))

		// Send the email
		addr := server + ":" + port
		err := smtp.SendMail(addr, auth, from, []string{recipient}, []byte(msg))
		if err != nil {
			return fmt.Errorf("failed to send email to %s: %w", recipient, err)
		}
	}

	return nil
}

// Helper functions for template processing
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

func formatPercent(f float64) string {
	return fmt.Sprintf("%.2f%%", f)
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func colorForUptime(uptime float64) string {
	if uptime >= 99.5 {
		return "#10B981" // Green
	} else if uptime >= 95.0 {
		return "#F59E0B" // Yellow
	}
	return "#EF4444" // Red
}

func priorityIcon(priority Priority) string {
	switch priority {
	case PriorityCritical:
		return "🔴"
	case PriorityHigh:
		return "🟡"
	case PriorityMedium:
		return "🔵"
	case PriorityLow:
		return "🟢"
	default:
		return "⚪"
	}
}

// getSitesToReport returns the list of sites to report on
func (rs *ReportScheduler) getSitesToReport(schedule *ReportSchedule) ([]string, error) {
	if len(schedule.Sites) > 0 {
		return schedule.Sites, nil
	}

	// Get all sites from storage (simplified approach)
	// In a real implementation, you'd get this from config or a site registry
	allStats, err := rs.storage.GetAllStats(time.Now().Add(-24 * time.Hour))
	if err != nil {
		return nil, err
	}

	var sites []string
	for siteName := range allStats {
		sites = append(sites, siteName)
	}

	return sites, nil
}

// getSiteURL returns URL for a site name (placeholder - should be implemented properly)
func (rs *ReportScheduler) getSiteURL(siteName string) string {
	// This is a placeholder. In a real implementation, you would:
	// 1. Store site URLs in the database
	// 2. Pass the config to ReportScheduler
	// 3. Or have a site registry service
	return "" // Return empty for now, SSL checks will be skipped
}

// shouldIncludeSection checks if a section should be included
func (rs *ReportScheduler) shouldIncludeSection(schedule *ReportSchedule, section ReportSection) bool {
	for _, s := range schedule.Sections {
		if s == section {
			return true
		}
	}
	return false
}

// formatPeriod formats schedule type to readable string
func (rs *ReportScheduler) formatPeriod(schedule ScheduleType) string {
	switch schedule {
	case ScheduleDaily:
		return "Daily"
	case ScheduleWeekly:
		return "Weekly"
	case ScheduleMonthly:
		return "Monthly"
	default:
		return "Custom"
	}
}

// generateSubject creates email subject for the report
func (rs *ReportScheduler) generateSubject(schedule *ReportSchedule, data *ReportData) string {
	return fmt.Sprintf("📊 %s - %s Report", schedule.Name, data.Period)
}

// generateTopPerformers generates top performing sites ranking
func (rs *ReportScheduler) generateTopPerformers(siteMetrics map[string]*metrics.AdvancedMetrics, limit int) []SiteRanking {
	var rankings []SiteRanking

	for _, metrics := range siteMetrics {
		score := metrics.UptimePercent*0.7 + (1.0-float64(metrics.P95.Milliseconds())/1000.0)*0.3
		rankings = append(rankings, SiteRanking{
			SiteName:        metrics.SiteName,
			UptimePercent:   metrics.UptimePercent,
			AvgResponseTime: metrics.P50, // Use P50 as average
			Score:           score,
		})
	}

	// Sort by score (simplified bubble sort)
	for i := 0; i < len(rankings)-1; i++ {
		for j := i + 1; j < len(rankings); j++ {
			if rankings[i].Score < rankings[j].Score {
				rankings[i], rankings[j] = rankings[j], rankings[i]
			}
		}
	}

	if len(rankings) > limit {
		rankings = rankings[:limit]
	}

	return rankings
}

// generateWorstPerformers generates worst performing sites ranking
func (rs *ReportScheduler) generateWorstPerformers(siteMetrics map[string]*metrics.AdvancedMetrics, limit int) []SiteRanking {
	rankings := rs.generateTopPerformers(siteMetrics, len(siteMetrics))

	// Reverse order for worst performers
	for i := len(rankings)/2 - 1; i >= 0; i-- {
		opp := len(rankings) - 1 - i
		rankings[i], rankings[opp] = rankings[opp], rankings[i]
	}

	if len(rankings) > limit {
		rankings = rankings[:limit]
	}

	return rankings
}

// generateRecommendations generates actionable recommendations
func (rs *ReportScheduler) generateRecommendations(data *ReportData) []Recommendation {
	var recommendations []Recommendation

	for siteName, sslCheck := range data.SSLChecks {
		if sslCheck.IsExpiringSoon(30) {
			recommendations = append(recommendations, Recommendation{
				Type:        RecommendationSSL,
				SiteName:    siteName,
				Title:       "SSL Certificate Expiring Soon",
				Description: fmt.Sprintf("Certificate expires in %d days", sslCheck.DaysUntilExpiry),
				Priority:    PriorityHigh,
				ActionItems: []string{
					"Renew SSL certificate",
					"Update certificate in server configuration",
					"Verify certificate chain after renewal",
				},
			})
		}
	}

	for siteName, metrics := range data.SiteMetrics {
		if metrics.UptimePercent < 99.0 {
			priority := PriorityMedium
			if metrics.UptimePercent < 95.0 {
				priority = PriorityCritical
			}

			recommendations = append(recommendations, Recommendation{
				Type:        RecommendationReliability,
				SiteName:    siteName,
				Title:       "Low Uptime Detected",
				Description: fmt.Sprintf("Uptime %.1f%% is below target", metrics.UptimePercent),
				Priority:    priority,
				ActionItems: []string{
					"Investigate root cause of downtime",
					"Review server logs and metrics",
					"Consider implementing redundancy",
				},
			})
		}

		if metrics.P95 > 2*time.Second {
			recommendations = append(recommendations, Recommendation{
				Type:        RecommendationPerformance,
				SiteName:    siteName,
				Title:       "Slow Response Times",
				Description: fmt.Sprintf("P95 response time %v exceeds 2 seconds", metrics.P95),
				Priority:    PriorityMedium,
				ActionItems: []string{
					"Optimize database queries",
					"Review application performance",
					"Consider implementing caching",
				},
			})
		}
	}

	return recommendations
}

// generateExecutiveSummary creates an executive summary
func (rs *ReportScheduler) generateExecutiveSummary(data *ReportData) string {
	healthyPercent := float64(data.HealthySites) / float64(data.TotalSites) * 100

	summary := fmt.Sprintf("Monitoring %d sites with %.1f%% overall uptime. ",
		data.TotalSites, data.OverallUptime)

	if healthyPercent >= 90 {
		summary += "✅ System performance is excellent with all key services operating normally."
	} else if healthyPercent >= 70 {
		summary += "⚠️ Some performance issues detected. Recommend immediate attention to underperforming services."
	} else {
		summary += "🚨 Critical performance issues detected. Immediate action required."
	}

	return summary
}

// getHTMLTemplate returns the HTML template for reports
func (rs *ReportScheduler) getHTMLTemplate(schedule *ReportSchedule) string {
	// Return a basic HTML template
	return `<!DOCTYPE html>
<html>
<head>
    <title>{{.GeneratedAt | formatTime}} - Site Monitor Report</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .header { background: #2563eb; color: white; padding: 20px; text-align: center; }
        .summary { background: #f8fafc; padding: 15px; margin: 20px 0; }
        .metric { display: inline-block; margin: 10px; padding: 10px; background: white; border: 1px solid #e2e8f0; }
        .site { margin: 15px 0; padding: 10px; border-left: 4px solid #10b981; }
        .recommendation { margin: 10px 0; padding: 10px; background: #fff3cd; border-left: 4px solid #f59e0b; }
    </style>
</head>
<body>
    <div class="header">
        <h1>📊 Site Monitor Report</h1>
        <p>{{.Period}} - Generated {{.GeneratedAt | formatTime}}</p>
    </div>
    
    <div class="summary">
        <h2>Executive Summary</h2>
        <p>{{.ExecutiveSummary}}</p>
        
        <div class="metric">
            <strong>{{.TotalSites}}</strong><br>Total Sites
        </div>
        <div class="metric">
            <strong>{{.HealthySites}}</strong><br>Healthy Sites
        </div>
        <div class="metric">
            <strong>{{.OverallUptime | formatPercent}}</strong><br>Overall Uptime
        </div>
        <div class="metric">
            <strong>{{.TotalChecks}}</strong><br>Total Checks
        </div>
    </div>
    
    <h2>Site Performance</h2>
    {{range $siteName, $metrics := .SiteMetrics}}
    <div class="site">
        <h3>{{$siteName}}</h3>
        <p>Uptime: {{$metrics.UptimePercent | formatPercent}} | P95: {{$metrics.P95 | formatDuration}}</p>
    </div>
    {{end}}
    
    {{if .Recommendations}}
    <h2>Recommendations</h2>
    {{range .Recommendations}}
    <div class="recommendation">
        <h4>{{.Title}} - {{.SiteName}}</h4>
        <p>{{.Description}}</p>
    </div>
    {{end}}
    {{end}}
    
    <hr>
    <p><small>Generated by Site Monitor v0.6.0</small></p>
</body>
</html>`
}
