package cmd

import (
	"fmt"
	"site-monitor/storage"
	"strings"
	"time"
)

// StatusOptions contains options for the status command
type StatusOptions struct {
	Watch    bool
	Interval time.Duration
}

// ShowStatus displays current status of all monitored sites
func (app *CLIApp) ShowStatus(opts StatusOptions) error {
	if err := app.InitStorage(); err != nil {
		return err
	}
	defer app.Close()

	if !app.CheckDatabaseExists() {
		app.ShowDatabaseNotFoundError()
		return nil
	}

	if opts.Watch {
		return app.watchStatus(opts.Interval)
	}

	return app.showCurrentStatus()
}

// showCurrentStatus shows a one-time status overview
func (app *CLIApp) showCurrentStatus() error {
	// Get recent stats (last 5 minutes to determine current status)
	since := time.Now().Add(-5 * time.Minute)
	allStats, err := app.storage.GetAllStats(since)
	if err != nil {
		return fmt.Errorf("failed to get current stats: %w", err)
	}

	if len(allStats) == 0 {
		fmt.Println("❌ No recent monitoring data found")
		fmt.Println()
		fmt.Println("Run 'site-monitor run' to start monitoring.")
		return nil
	}

	// Show header with timestamp
	fmt.Printf("🚀 Site Monitor Status - %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(strings.Repeat("━", 60))

	// Show each site's current status
	healthyCount := 0
	for _, stats := range allStats {
		if app.printSiteStatus(stats) {
			healthyCount++
		}
	}

	// Show summary
	fmt.Println(strings.Repeat("━", 60))
	totalSites := len(allStats)
	fmt.Printf("📊 Overall: %d/%d sites healthy", healthyCount, totalSites)

	if healthyCount == totalSites {
		fmt.Println(" ✅ All systems operational")
	} else if healthyCount == 0 {
		fmt.Println(" ❌ All systems down")
	} else {
		fmt.Println(" ⚠️ Some issues detected")
	}

	return nil
}

// watchStatus continuously shows status updates
func (app *CLIApp) watchStatus(interval time.Duration) error {
	fmt.Printf("👁️  Watching site status (refreshing every %s)\n", formatDuration(interval))
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Show initial status
	if err := app.showCurrentStatus(); err != nil {
		return err
	}

	// Refresh status at intervals
	for range ticker.C {
		// Clear screen (simple approach)
		fmt.Print("\033[2J\033[H")

		if err := app.showCurrentStatus(); err != nil {
			return err
		}
	}

	return nil
}

// printSiteStatus prints the current status of a single site
// Returns true if the site is healthy, false otherwise
func (app *CLIApp) printSiteStatus(stats storage.Stats) bool {
	// Determine current health based on recent success rate and last check
	isHealthy := true
	statusIcon := "✅"
	statusText := "HEALTHY"

	// Check success rate
	if stats.SuccessRate < 80.0 {
		isHealthy = false
		statusIcon = "❌"
		statusText = "DOWN"
	} else if stats.SuccessRate < 99.0 {
		isHealthy = false
		statusIcon = "⚠️"
		statusText = "DEGRADED"
	}

	// Check if last check was too long ago (more than 2x expected interval)
	timeSinceLastCheck := time.Since(stats.LastCheck)
	if timeSinceLastCheck > 5*time.Minute { // Assume max 5min interval
		isHealthy = false
		statusIcon = "🔄"
		statusText = "STALE"
	}

	// Site status line
	fmt.Printf("%s %-12s %s\n", statusIcon, statusText, stats.SiteName)

	// Additional details
	fmt.Printf("   📈 Recent Success: %.1f%% (%d/%d checks)\n",
		stats.SuccessRate, stats.SuccessfulChecks, stats.TotalChecks)

	if stats.SuccessfulChecks > 0 {
		fmt.Printf("   ⚡ Response Time: %v avg\n",
			stats.AvgResponseTime.Round(time.Millisecond))
	}

	fmt.Printf("   🕐 Last Check: %s ago\n",
		formatDuration(timeSinceLastCheck.Round(time.Second)))

	// Show recent failures if any
	if stats.FailedChecks > 0 {
		fmt.Printf("   💥 Recent Failures: %d\n", stats.FailedChecks)
	}

	// Show specific issues
	if !isHealthy {
		app.showStatusIssues(stats, timeSinceLastCheck)
	}

	fmt.Println() // Add spacing between sites
	return isHealthy
}

// showStatusIssues shows specific issues for unhealthy sites
func (app *CLIApp) showStatusIssues(stats storage.Stats, timeSinceLastCheck time.Duration) {
	fmt.Print("   🚨 Issues: ")

	var issues []string

	if stats.SuccessRate < 80.0 {
		issues = append(issues, "High failure rate")
	} else if stats.SuccessRate < 99.0 {
		issues = append(issues, "Some failures")
	}

	if timeSinceLastCheck > 5*time.Minute {
		issues = append(issues, "No recent checks")
	}

	if stats.SuccessfulChecks > 0 && stats.AvgResponseTime > 5*time.Second {
		issues = append(issues, "Slow response times")
	}

	fmt.Println(strings.Join(issues, ", "))
}
