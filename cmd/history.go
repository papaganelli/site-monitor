package cmd

import (
	"fmt"
	"site-monitor/storage"
	"strings"
	"time"
)

// HistoryOptions contains options for the history command
type HistoryOptions struct {
	SiteName string
	Since    time.Duration
	Limit    int
}

// ShowHistory displays monitoring history
func (app *CLIApp) ShowHistory(opts HistoryOptions) error {
	if err := app.InitStorage(); err != nil {
		return err
	}
	defer app.Close()

	if !app.CheckDatabaseExists() {
		app.ShowDatabaseNotFoundError()
		return nil
	}

	since := time.Now().Add(-opts.Since)

	// Get history entries
	var entries []storage.HistoryEntry
	var err error

	if opts.SiteName != "" {
		entries, err = app.storage.GetHistory(opts.SiteName, since)
	} else {
		entries, err = app.storage.GetAllHistory(since)
	}

	if err != nil {
		return fmt.Errorf("failed to get history: %w", err)
	}

	// Limit results if specified
	if opts.Limit > 0 && len(entries) > opts.Limit {
		entries = entries[:opts.Limit]
	}

	// Show header
	fmt.Printf("📋 Monitoring History")
	if opts.SiteName != "" {
		fmt.Printf(" - %s", opts.SiteName)
	}
	if opts.Since > 0 {
		fmt.Printf(" (Last %s)", formatDuration(opts.Since))
	}
	if opts.Limit > 0 {
		fmt.Printf(" - Limited to %d entries", opts.Limit)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("━", 70))

	if len(entries) == 0 {
		fmt.Println("❌ No history entries found")
		return nil
	}

	// Group entries by site for better readability (when showing all sites)
	if opts.SiteName == "" {
		app.showHistoryGroupedBySite(entries)
	} else {
		app.showHistoryList(entries)
	}

	// Show summary
	fmt.Println(strings.Repeat("━", 70))
	app.showHistorySummary(entries, opts)

	return nil
}

// showHistoryGroupedBySite shows history entries grouped by site
func (app *CLIApp) showHistoryGroupedBySite(entries []storage.HistoryEntry) {
	// Group entries by site name
	siteGroups := make(map[string][]storage.HistoryEntry)
	for _, entry := range entries {
		siteGroups[entry.SiteName] = append(siteGroups[entry.SiteName], entry)
	}

	isFirst := true
	for siteName, siteEntries := range siteGroups {
		if !isFirst {
			fmt.Println() // Add spacing between sites
		}
		isFirst = false

		fmt.Printf("🌐 %s (%d entries)\n", siteName, len(siteEntries))
		fmt.Println(strings.Repeat("─", 50))

		// Show only the most recent entries for each site to avoid clutter
		limit := 5
		if len(siteEntries) > limit {
			fmt.Printf("   Showing %d most recent entries (of %d total)\n", limit, len(siteEntries))
			siteEntries = siteEntries[:limit]
		}

		for _, entry := range siteEntries {
			app.printHistoryEntry(entry, false) // false = don't show site name
		}
	}
}

// showHistoryList shows history entries as a simple list
func (app *CLIApp) showHistoryList(entries []storage.HistoryEntry) {
	for _, entry := range entries {
		app.printHistoryEntry(entry, true) // true = show site name
	}
}

// printHistoryEntry prints a single history entry
func (app *CLIApp) printHistoryEntry(entry storage.HistoryEntry, showSiteName bool) {
	// Status icon
	statusIcon := "✅"
	statusText := "OK"
	if !entry.Success {
		statusIcon = "❌"
		statusText = "FAIL"
	}

	// Format timestamp
	timestamp := entry.Timestamp.Format("15:04:05")

	// Site name (if requested)
	siteInfo := ""
	if showSiteName {
		siteInfo = fmt.Sprintf(" (%s)", entry.SiteName)
	}

	// Basic info line
	fmt.Printf("   [%s] %s %s%s - %d - %v",
		timestamp,
		statusIcon,
		statusText,
		siteInfo,
		entry.Status,
		entry.Duration.Round(time.Millisecond))

	// Error message (if any)
	if entry.Error != "" {
		fmt.Printf(" - %s", entry.Error)
	}

	fmt.Println()
}

// showHistorySummary shows a summary of the history entries
func (app *CLIApp) showHistorySummary(entries []storage.HistoryEntry, opts HistoryOptions) {
	if len(entries) == 0 {
		return
	}

	successCount := 0
	var totalDuration time.Duration
	var minDuration, maxDuration time.Duration
	sites := make(map[string]bool)

	for i, entry := range entries {
		if entry.Success {
			successCount++
			if i == 0 || entry.Duration < minDuration {
				minDuration = entry.Duration
			}
			if entry.Duration > maxDuration {
				maxDuration = entry.Duration
			}
			totalDuration += entry.Duration
		}
		sites[entry.SiteName] = true
	}

	successRate := float64(successCount) / float64(len(entries)) * 100
	avgDuration := time.Duration(0)
	if successCount > 0 {
		avgDuration = totalDuration / time.Duration(successCount)
	}

	fmt.Printf("📊 Summary: %d entries", len(entries))
	if len(sites) > 1 {
		fmt.Printf(" from %d sites", len(sites))
	}
	fmt.Println()

	fmt.Printf("✅ Success Rate: %.1f%% (%d/%d)\n", successRate, successCount, len(entries))

	if successCount > 0 {
		fmt.Printf("⚡ Response Times: %v avg (min: %v, max: %v)\n",
			avgDuration.Round(time.Millisecond),
			minDuration.Round(time.Millisecond),
			maxDuration.Round(time.Millisecond))
	}

	// Time span
	if len(entries) > 1 {
		timeSpan := entries[0].Timestamp.Sub(entries[len(entries)-1].Timestamp)
		fmt.Printf("⏱️  Time Span: %s\n", formatDuration(timeSpan))
	}
}
