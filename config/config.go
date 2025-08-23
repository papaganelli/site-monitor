package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config represents the main configuration structure
type Config struct {
	Sites  []Site       `json:"sites"`
	Alerts *AlertConfig `json:"alerts,omitempty"`
}

// Site represents a single website to monitor
type Site struct {
	Name     string `json:"name"`     // Display name for the site
	URL      string `json:"url"`      // URL to monitor
	Interval string `json:"interval"` // How often to check (e.g., "30s", "5m")
	Timeout  string `json:"timeout"`  // HTTP request timeout
}

// AlertConfig represents the alert configuration
type AlertConfig struct {
	Email      EmailConfig     `json:"email"`
	Webhook    WebhookConfig   `json:"webhook"`
	Thresholds ThresholdConfig `json:"thresholds"`
}

// EmailConfig represents email alert configuration
type EmailConfig struct {
	Enabled    bool     `json:"enabled"`
	SMTPServer string   `json:"smtp_server"`
	SMTPPort   int      `json:"smtp_port"`
	Username   string   `json:"username"`
	Password   string   `json:"password"`
	From       string   `json:"from"`
	Recipients []string `json:"recipients"`
	UseTLS     bool     `json:"use_tls"`
}

// WebhookConfig represents webhook alert configuration
type WebhookConfig struct {
	Enabled    bool              `json:"enabled"`
	URL        string            `json:"url"`
	Format     string            `json:"format"` // slack, discord, teams, generic
	Headers    map[string]string `json:"headers,omitempty"`
	Timeout    string            `json:"timeout"`
	RetryCount int               `json:"retry_count"`
}

// ThresholdConfig represents alert threshold configuration
type ThresholdConfig struct {
	// Site down detection
	ConsecutiveFailures int `json:"consecutive_failures"`

	// Performance thresholds
	ResponseTimeThreshold string  `json:"response_time_threshold"` // e.g., "5s"
	UptimeThreshold       float64 `json:"uptime_threshold"`        // e.g., 95.0 for 95%

	// Time windows for calculations
	UptimeWindow      string `json:"uptime_window"`      // e.g., "24h"
	PerformanceWindow string `json:"performance_window"` // e.g., "1h"

	// Cooldown to prevent spam
	AlertCooldown string `json:"alert_cooldown"` // e.g., "5m"
}

// Load reads and parses configuration from JSON file
func Load(filename string) (*Config, error) {
	// Open the configuration file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode JSON into Config struct
	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GetInterval converts string interval to time.Duration
func (s *Site) GetInterval() (time.Duration, error) {
	return time.ParseDuration(s.Interval)
}

// GetTimeout converts string timeout to time.Duration
func (s *Site) GetTimeout() (time.Duration, error) {
	return time.ParseDuration(s.Timeout)
}

// Helper methods for ThresholdConfig

// GetResponseTimeThreshold parses and returns the response time threshold
func (tc ThresholdConfig) GetResponseTimeThreshold() (time.Duration, error) {
	return time.ParseDuration(tc.ResponseTimeThreshold)
}

// GetUptimeWindow parses and returns the uptime calculation window
func (tc ThresholdConfig) GetUptimeWindow() (time.Duration, error) {
	return time.ParseDuration(tc.UptimeWindow)
}

// GetPerformanceWindow parses and returns the performance calculation window
func (tc ThresholdConfig) GetPerformanceWindow() (time.Duration, error) {
	return time.ParseDuration(tc.PerformanceWindow)
}

// GetAlertCooldown parses and returns the alert cooldown duration
func (tc ThresholdConfig) GetAlertCooldown() (time.Duration, error) {
	return time.ParseDuration(tc.AlertCooldown)
}
