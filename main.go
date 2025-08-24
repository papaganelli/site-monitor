package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"site-monitor/cmd"
	"site-monitor/config"
	"site-monitor/monitor"
	"site-monitor/storage"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func main() {
	// Define command line flags
	var (
		showHelp    = flag.Bool("help", false, "Show help message")
		showVersion = flag.Bool("version", false, "Show version information")
	)

	// Parse command line arguments
	flag.Parse()
	args := flag.Args()

	// Show help
	if *showHelp || (len(args) == 0 && len(os.Args) == 1) {
		showUsage()
		return
	}

	// Show version
	if *showVersion {
		fmt.Println("Site Monitor v0.6.0")
		return
	}

	if len(args) == 0 {
		// Default behavior: run the monitor
		runMonitor()
		return
	}

	// Handle subcommands
	command := args[0]
	commandArgs := args[1:]

	app, err := cmd.NewCLIApp()
	if err != nil {
		log.Fatal("Failed to initialize CLI app:", err)
	}

	switch command {
	case "run", "monitor":
		runMonitor()
	case "stats":
		runStatsCommand(app, commandArgs)
	case "history":
		runHistoryCommand(app, commandArgs)
	case "status":
		runStatusCommand(app, commandArgs)
	case "dashboard", "web":
		runDashboardCommand(app, commandArgs)
	case "export":
		runExportCommand(app, commandArgs)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println()
		showUsage()
		os.Exit(1)
	}
}

// showUsage displays the help message
func showUsage() {
	fmt.Println("Site Monitor - Website monitoring tool")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  site-monitor [command] [options]")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("  run                     Start monitoring (default)")
	fmt.Println("  stats [options]         Show monitoring statistics")
	fmt.Println("  history [options]       Show monitoring history")
	fmt.Println("  status [options]        Show current status")
	fmt.Println("  dashboard [options]     Start web dashboard")
	fmt.Println("  export [options]        Export monitoring data")
	fmt.Println()
	fmt.Println("STATS OPTIONS:")
	fmt.Println("  --site <name>           Show stats for specific site")
	fmt.Println("  --since <duration>      Time period (e.g., 1h, 24h, 7d)")
	fmt.Println()
	fmt.Println("HISTORY OPTIONS:")
	fmt.Println("  --site <name>           Show history for specific site")
	fmt.Println("  --since <duration>      Time period (e.g., 1h, 24h)")
	fmt.Println("  --limit <number>        Limit number of entries")
	fmt.Println()
	fmt.Println("STATUS OPTIONS:")
	fmt.Println("  --watch                 Watch status with auto-refresh")
	fmt.Println("  --interval <duration>   Refresh interval (default: 30s)")
	fmt.Println()
	fmt.Println("DASHBOARD OPTIONS:")
	fmt.Println("  --port <number>         Web server port (default: 8080)")
	fmt.Println()
	fmt.Println("EXPORT OPTIONS:")
	fmt.Println("  --format <format>       Export format (json, csv, html)")
	fmt.Println("  --site <name>           Export data for specific site")
	fmt.Println("  --since <duration>      Time period to export")
	fmt.Println("  --output <file>         Output file path")
	fmt.Println("  --stats                 Include statistics")
	fmt.Println("  --stdout                Output to stdout")
	fmt.Println("  --list-formats          Show available formats")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  site-monitor run")
	fmt.Println("  site-monitor stats --since 24h")
	fmt.Println("  site-monitor stats --site \"My Site\"")
	fmt.Println("  site-monitor history --limit 50")
	fmt.Println("  site-monitor status --watch")
	fmt.Println("  site-monitor dashboard --port 3000")
	fmt.Println("  site-monitor export --format json --output data.json")
	fmt.Println("  site-monitor export --format csv --site \"My Site\" --since 7d")
	fmt.Println("  site-monitor export --format html --stats")
}

// runStatsCommand handles the stats subcommand
func runStatsCommand(app *cmd.CLIApp, args []string) {
	var siteName string
	var sinceStr string

	// Parse arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--site":
			if i+1 < len(args) {
				siteName = args[i+1]
				i++
			}
		case "--since":
			if i+1 < len(args) {
				sinceStr = args[i+1]
				i++
			}
		}
	}

	// Parse duration
	since := 24 * time.Hour // Default to 24 hours
	if sinceStr != "" {
		var err error
		since, err = parseDuration(sinceStr)
		if err != nil {
			log.Fatalf("Invalid duration '%s': %v", sinceStr, err)
		}
	}

	opts := cmd.StatsOptions{
		SiteName: siteName,
		Since:    since,
	}

	if err := app.ShowStats(opts); err != nil {
		log.Fatal(err)
	}
}

// runHistoryCommand handles the history subcommand
func runHistoryCommand(app *cmd.CLIApp, args []string) {
	var siteName string
	var sinceStr string
	var limitStr string

	// Parse arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--site":
			if i+1 < len(args) {
				siteName = args[i+1]
				i++
			}
		case "--since":
			if i+1 < len(args) {
				sinceStr = args[i+1]
				i++
			}
		case "--limit":
			if i+1 < len(args) {
				limitStr = args[i+1]
				i++
			}
		}
	}

	// Parse duration
	since := 24 * time.Hour // Default to 24 hours
	if sinceStr != "" {
		var err error
		since, err = parseDuration(sinceStr)
		if err != nil {
			log.Fatalf("Invalid duration '%s': %v", sinceStr, err)
		}
	}

	// Parse limit
	limit := 0
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			log.Fatalf("Invalid limit '%s': %v", limitStr, err)
		}
	}

	opts := cmd.HistoryOptions{
		SiteName: siteName,
		Since:    since,
		Limit:    limit,
	}

	if err := app.ShowHistory(opts); err != nil {
		log.Fatal(err)
	}
}

// runStatusCommand handles the status subcommand
func runStatusCommand(app *cmd.CLIApp, args []string) {
	var watch bool
	var intervalStr string

	// Parse arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--watch":
			watch = true
		case "--interval":
			if i+1 < len(args) {
				intervalStr = args[i+1]
				i++
			}
		}
	}

	// Parse interval
	interval := 30 * time.Second // Default to 30 seconds
	if intervalStr != "" {
		var err error
		interval, err = parseDuration(intervalStr)
		if err != nil {
			log.Fatalf("Invalid interval '%s': %v", intervalStr, err)
		}
	}

	opts := cmd.StatusOptions{
		Watch:    watch,
		Interval: interval,
	}

	if err := app.ShowStatus(opts); err != nil {
		log.Fatal(err)
	}
}

// runDashboardCommand handles the dashboard subcommand
func runDashboardCommand(app *cmd.CLIApp, args []string) {
	var portStr string

	// Parse arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--port":
			if i+1 < len(args) {
				portStr = args[i+1]
				i++
			}
		}
	}

	// Parse port
	port := 8080 // Default port
	if portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil || port < 1 || port > 65535 {
			log.Fatalf("Invalid port '%s': must be between 1 and 65535", portStr)
		}
	}

	opts := cmd.DashboardOptions{
		Port: port,
	}

	if err := app.ShowDashboard(opts); err != nil {
		log.Fatal(err)
	}
}

// runExportCommand handles the export subcommand
func runExportCommand(app *cmd.CLIApp, args []string) {
	var format string
	var siteName string
	var sinceStr string
	var untilStr string
	var limitStr string
	var outputPath string
	var stats bool
	var stdout bool
	var listFormats bool

	// Parse arguments
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--format", "-f":
			if i+1 < len(args) {
				format = args[i+1]
				i++
			}
		case "--site", "-s":
			if i+1 < len(args) {
				siteName = args[i+1]
				i++
			}
		case "--since":
			if i+1 < len(args) {
				sinceStr = args[i+1]
				i++
			}
		case "--until":
			if i+1 < len(args) {
				untilStr = args[i+1]
				i++
			}
		case "--limit", "-l":
			if i+1 < len(args) {
				limitStr = args[i+1]
				i++
			}
		case "--output", "-o":
			if i+1 < len(args) {
				outputPath = args[i+1]
				i++
			}
		case "--stats":
			stats = true
		case "--stdout":
			stdout = true
		case "--list-formats":
			listFormats = true
		case "--help", "-h":
			showExportHelp()
			return
		}
	}

	// Show format list if requested
	if listFormats {
		app.ListExportFormats()
		return
	}

	// Parse duration
	since := 24 * time.Hour // Default to 24 hours
	if sinceStr != "" {
		var err error
		since, err = parseDuration(sinceStr)
		if err != nil {
			log.Fatalf("Invalid duration '%s': %v", sinceStr, err)
		}
	}

	// Parse until time
	var until *time.Time
	if untilStr != "" {
		parsedTime, err := time.Parse("2006-01-02 15:04:05", untilStr)
		if err != nil {
			// Try alternative format
			parsedTime, err = time.Parse("2006-01-02", untilStr)
			if err != nil {
				log.Fatalf("Invalid until time '%s': use format '2006-01-02 15:04:05' or '2006-01-02'", untilStr)
			}
		}
		until = &parsedTime
	}

	// Parse limit
	limit := 0
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			log.Fatalf("Invalid limit '%s': must be a positive number", limitStr)
		}
	}

	opts := cmd.ExportCLIOptions{
		Format:     format,
		SiteName:   siteName,
		Since:      since,
		Until:      until,
		Limit:      limit,
		OutputPath: outputPath,
		Stats:      stats,
		Stdout:     stdout,
	}

	if err := app.ShowExport(opts); err != nil {
		log.Fatal(err)
	}
}

// showExportHelp displays help for the export command
func showExportHelp() {
	fmt.Println("Site Monitor - Export Command Help")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  site-monitor export [options]")
	fmt.Println()
	fmt.Println("OPTIONS:")
	fmt.Println("  --format, -f <format>   Export format (json, csv, html) [default: json]")
	fmt.Println("  --site, -s <name>       Export data for specific site only")
	fmt.Println("  --since <duration>      Time period to export (e.g., 1h, 24h, 7d) [default: 24h]")
	fmt.Println("  --until <time>          End time (format: '2006-01-02 15:04:05' or '2006-01-02')")
	fmt.Println("  --limit, -l <number>    Maximum number of records to export")
	fmt.Println("  --output, -o <file>     Output file path [default: auto-generated]")
	fmt.Println("  --stats                 Include statistical summary in export")
	fmt.Println("  --stdout                Output to stdout instead of file")
	fmt.Println("  --list-formats          Show available export formats")
	fmt.Println("  --help, -h              Show this help message")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  site-monitor export --format json --output data.json")
	fmt.Println("  site-monitor export --format csv --site \"My Site\" --since 7d")
	fmt.Println("  site-monitor export --format html --stats --output report.html")
	fmt.Println("  site-monitor export --format json --stdout | jq .")
	fmt.Println("  site-monitor export --list-formats")
	fmt.Println("  site-monitor export --since 1h --until \"2024-01-01 12:00:00\"")
	fmt.Println()
}

// runMonitor runs the original monitoring daemon
func runMonitor() {
	// Load configuration from file
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize storage
	db, err := storage.NewSQLiteStorage("site-monitor.db")
	if err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}
	defer db.Close()

	// Initialize database tables
	if err := db.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	fmt.Printf("🚀 Starting monitoring for %d sites\n", len(cfg.Sites))
	fmt.Printf("💾 Database initialized: site-monitor.db\n")

	var wg sync.WaitGroup

	// Channel to handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start monitoring each site in a separate goroutine
	for _, site := range cfg.Sites {
		wg.Add(1)

		go func(s config.Site) {
			defer wg.Done()

			interval, err := s.GetInterval()
			if err != nil {
				log.Printf("Invalid interval for %s: %v", s.Name, err)
				return
			}

			timeout, err := s.GetTimeout()
			if err != nil {
				log.Printf("Invalid timeout for %s: %v", s.Name, err)
				return
			}

			m := monitor.New(s.URL, interval)
			m.SetName(s.Name)
			m.SetTimeout(timeout)
			m.SetStorage(db) // Attach storage to monitor

			fmt.Printf("📍 Starting %s (%s) - checking every %s\n",
				s.Name, s.URL, s.Interval)

			if err := m.Start(); err != nil {
				log.Printf("Error monitoring %s: %v", s.Name, err)
			}
		}(site)
	}

	// Handle graceful shutdown
	go func() {
		<-sigChan
		fmt.Println("\n🛑 Received shutdown signal, stopping monitors...")
		// Note: In a real implementation, we would send a stop signal to monitors
		// For now, the program will exit and goroutines will be terminated
		os.Exit(0)
	}()

	wg.Wait()
}

// parseDuration parses duration strings like "1h", "30m", "24h", "7d"
func parseDuration(s string) (time.Duration, error) {
	// Handle days (Go's time.ParseDuration doesn't support 'd')
	if strings.HasSuffix(s, "d") {
		daysStr := strings.TrimSuffix(s, "d")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return 0, fmt.Errorf("invalid number of days: %s", daysStr)
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}

	// Use standard time.ParseDuration for other units
	return time.ParseDuration(s)
}
