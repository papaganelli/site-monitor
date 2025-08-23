package cmd

import (
	"fmt"
	"site-monitor/web"
)

// DashboardOptions contains options for the dashboard command
type DashboardOptions struct {
	Port int
}

// ShowDashboard starts the web dashboard server
func (app *CLIApp) ShowDashboard(opts DashboardOptions) error {
	// Initialize storage
	if err := app.InitStorage(); err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	// Load configuration
	if err := app.LoadConfig(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if !app.CheckDatabaseExists() {
		fmt.Println("⚠️  Database not found!")
		fmt.Println()
		fmt.Println("The dashboard needs monitoring data to display.")
		fmt.Println("To start collecting data, run 'site-monitor run' in another terminal.")
		fmt.Println()
		fmt.Println("The dashboard will start anyway and update once data becomes available.")
		fmt.Println()
	}

	// Create dashboard instance
	dashboard := web.NewDashboard(app.storage, app.config, opts.Port)

	// Print startup information
	fmt.Printf("🌐 Starting Site Monitor Dashboard\n")
	fmt.Printf("📊 Port: %d\n", opts.Port)
	fmt.Printf("💾 Database: %s\n", app.dbPath)
	fmt.Printf("📋 Sites configured: %d\n", len(app.config.Sites))

	if app.config.Alerts != nil {
		alertChannels := 0
		if app.config.Alerts.Email.Enabled {
			alertChannels++
		}
		if app.config.Alerts.Webhook.Enabled {
			alertChannels++
		}
		fmt.Printf("🚨 Alert channels: %d\n", alertChannels)
	}

	fmt.Printf("\n🚀 Dashboard available at: http://localhost:%d\n", opts.Port)
	fmt.Printf("📱 WebSocket endpoint: ws://localhost:%d/ws\n", opts.Port)
	fmt.Printf("\n💡 Press Ctrl+C to stop the dashboard\n\n")

	// Start the dashboard server (blocking)
	if err := dashboard.Start(); err != nil {
		return fmt.Errorf("dashboard server error: %w", err)
	}

	return nil
}
