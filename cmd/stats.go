package cmd

import (
	"fmt"
	"site-monitor/storage"
	"strings"
	"time"
)

// StatsOptions contains options for the stats command
type StatsOptions struct {
	SiteName string
	Since    time.Duration
}

// ShowStats displays statistics for monitored sites
func (app *CLIApp) ShowStats(opts StatsOptions) error {
	if err := app.InitStorage(); err != nil {
		return err
	}
	defer app.Close()

	if !app.CheckDatabaseExists() {
		app.ShowDatabaseNotFoundError()
		return nil
	}

	since := time.Now().Add(-opts.Since)

	// Show header
	fmt.Printf("📊 Monitoring Statistics")
	if opts.Since > 0 {
		fmt.Printf(" (Last %s)", formatDuration(opts.Since))
	}
	fmt.Println()
	fmt.Println(strings.Repeat("━", 50))

	if opts.SiteName != "" {
		// Show stats for specific site
		return app.showSiteStats(opts.SiteName, since)
	}

	// Show stats for all sites
	return app.showAllStats(since)
}

// showSiteStats displays statistics for a specific site
func (app *CLIApp) showSiteStats(siteName string, since time.Time) error {
	stats, err := app.storage.GetStats(siteName, since)
	if err != nil {
		return fmt.Errorf("failed to get stats for %s: %w", siteName, err)
	}

	if stats.TotalChecks == 0 {
		fmt.Printf("❌ No data found for site: %s\n", siteName)
		return nil
	}

	app.printSiteStats(stats)
	return nil
}

// showAllStats displays statistics for all sites
func (app *CLIApp) showAllStats(since time.Time) error {
	allStats, err := app.storage.GetAllStats(since)
	if err != nil {
		return fmt.Errorf("failed to get all stats: %w", err)
	}

	if len(allStats) == 0 {
		fmt.Println("❌ No monitoring data found")
		fmt.Println()
		fmt.Println("Run 'site-monitor run' to start collecting data.")
		return nil
	}

	// Sort sites by name for consistent output
	var sites []string
	for siteName := range allStats {
		sites = append(sites, siteName)
	}

	for i, siteName := range sites {
		if i > 0 {
			fmt.Println() // Add spacing between sites
		}
		app.printSiteStats(allStats[siteName])
	}

	// Show summary
	fmt.Println()
	fmt.Println(strings.Repeat("━", 50))
	app.printSummary(allStats)

	return nil
}

// printSiteStats prints formatted statistics for a single site
func (app *CLIApp) printSiteStats(stats storage.Stats) {
	// Site name and uptime
	uptimeIcon := "✅"
	if stats.SuccessRate < 99.0 {
		uptimeIcon = "⚠️"
	}
	if stats.SuccessRate < 95.0 {
		uptimeIcon = "❌"
	}

	fmt.Printf("%s %s\n", uptimeIcon, stats.SiteName)

	// Uptime percentage
	fmt.Printf("   📈 Uptime: %.1f%% (%d/%d checks)\n",
		stats.SuccessRate, stats.SuccessfulChecks, stats.TotalChecks)

	// Response time stats
	if stats.SuccessfulChecks > 0 {
		fmt.Printf("   ⚡ Response: %v avg (min: %v, max: %v)\n",
			stats.AvgResponseTime.Round(time.Millisecond),
			stats.MinResponseTime.Round(time.Millisecond),
			stats.MaxResponseTime.Round(time.Millisecond))
	}

	// Last check timing
	if !stats.LastCheck.IsZero() {
		lastCheckAgo := time.Since(stats.LastCheck).Round(time.Second)
		fmt.Printf("   🕐 Last Check: %s ago\n", formatDuration(lastCheckAgo))
	}

	// Failed checks (if any)
	if stats.FailedChecks > 0 {
		fmt.Printf("   💥 Failed Checks: %d\n", stats.FailedChecks)
	}

	// Total monitoring duration
	if !stats.FirstCheck.IsZero() && !stats.LastCheck.IsZero() {
		monitoringDuration := stats.LastCheck.Sub(stats.FirstCheck)
		if monitoringDuration > time.Minute {
			fmt.Printf("   📅 Monitoring Duration: %s\n", formatDuration(monitoringDuration))
		}
	}
}

// printSummary prints an overall summary for all sites
func (app *CLIApp) printSummary(allStats map[string]storage.Stats) {
	var totalChecks int64
	var totalSuccessful int64
	var healthySites int
	totalSites := len(allStats)

	for _, stats := range allStats {
		totalChecks += stats.TotalChecks
		totalSuccessful += stats.SuccessfulChecks
		if stats.SuccessRate >= 99.0 {
			healthySites++
		}
	}

	overallSuccess := float64(0)
	if totalChecks > 0 {
		overallSuccess = float64(totalSuccessful) / float64(totalChecks) * 100
	}

	fmt.Printf("📋 Summary: %d sites monitored\n", totalSites)
	fmt.Printf("🎯 Overall Uptime: %.1f%% (%d/%d checks)\n",
		overallSuccess, totalSuccessful, totalChecks)
	fmt.Printf("💚 Healthy Sites: %d/%d (≥99%% uptime)\n",
		healthySites, totalSites)
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
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
