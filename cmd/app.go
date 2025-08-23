package cmd

import (
	"fmt"
	"os"
	"site-monitor/config"
	"site-monitor/storage"
)

// CLIApp represents the main CLI application
type CLIApp struct {
	storage storage.Storage
	config  *config.Config
	dbPath  string
}

// NewCLIApp creates a new CLI application instance
func NewCLIApp() (*CLIApp, error) {
	return &CLIApp{
		dbPath: "site-monitor.db",
	}, nil
}

// InitStorage initializes the storage connection
func (app *CLIApp) InitStorage() error {
	if app.storage != nil {
		return nil // Already initialized
	}

	db, err := storage.NewSQLiteStorage(app.dbPath)
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %w", err)
	}

	if err := db.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	app.storage = db
	return nil
}

// LoadConfig loads the configuration file
func (app *CLIApp) LoadConfig() error {
	cfg, err := config.Load("config.json")
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	app.config = cfg
	return nil
}

// Close closes the storage connection
func (app *CLIApp) Close() error {
	if app.storage != nil {
		return app.storage.Close()
	}
	return nil
}

// CheckDatabaseExists verifies if the database file exists
func (app *CLIApp) CheckDatabaseExists() bool {
	_, err := os.Stat(app.dbPath)
	return !os.IsNotExist(err)
}

// ShowDatabaseNotFoundError shows a helpful error when database is missing
func (app *CLIApp) ShowDatabaseNotFoundError() {
	fmt.Println("❌ Database not found!")
	fmt.Println()
	fmt.Println("It looks like you haven't run the monitor yet.")
	fmt.Println("To start collecting data, run:")
	fmt.Println()
	fmt.Println("  site-monitor run")
	fmt.Println()
	fmt.Println("Let it run for a few minutes to collect monitoring data,")
	fmt.Println("then try your command again.")
}
