package monitor

import (
	"fmt"
	"net/http"
	"time"
)

type Monitor struct {
	Name     string // Display name for the monitor
	URL      string
	Interval time.Duration
	client   *http.Client
}

// New creates a new monitor instance
func New(url string, interval time.Duration) *Monitor {
	return &Monitor{
		Name:     url, // Default name is URL
		URL:      url,
		Interval: interval,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SetName sets a custom name for the monitor
func (m *Monitor) SetName(name string) {
	m.Name = name
}

// SetTimeout sets custom timeout for HTTP requests
func (m *Monitor) SetTimeout(timeout time.Duration) {
	m.client.Timeout = timeout
}

// Start begins the monitoring loop
func (m *Monitor) Start() error {
	ticker := time.NewTicker(m.Interval)
	defer ticker.Stop()

	// First check immediately
	result := m.check()
	fmt.Println(result)

	// Use for range instead of for { select {} }
	for range ticker.C {
		result := m.check()
		fmt.Println(result)
	}

	return nil // This will never be reached, but satisfies the function signature
}

// check performs a single HTTP check
func (m *Monitor) check() Result {
	start := time.Now()

	resp, err := m.client.Get(m.URL)
	duration := time.Since(start)

	result := Result{
		Name:      m.Name, // Include monitor name
		URL:       m.URL,
		Duration:  duration,
		Timestamp: time.Now(),
	}

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.Status = resp.StatusCode
	result.Success = resp.StatusCode >= 200 && resp.StatusCode < 400

	return result
}
