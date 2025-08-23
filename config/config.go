package config

import (
    "encoding/json"
    "os"
    "time"
)

// Config represents the main configuration structure
type Config struct {
    Sites []Site `json:"sites"`
}

// Site represents a single website to monitor
type Site struct {
    Name     string `json:"name"`     // Display name for the site
    URL      string `json:"url"`      // URL to monitor
    Interval string `json:"interval"` // How often to check (e.g., "30s", "5m")
    Timeout  string `json:"timeout"`  // HTTP request timeout
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