package alerts

import (
	"fmt"
	"log"
	"site-monitor/config"
	"site-monitor/monitor"
	"site-monitor/storage"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Manager handles alert processing and routing
type Manager struct {
	config   config.AlertConfig
	storage  storage.Storage
	channels []AlertChannel
	states   map[string]*AlertState // Site name -> AlertState
	mu       sync.RWMutex
}

// NewManager creates a new alert manager
func NewManager(alertConfig config.AlertConfig, storage storage.Storage) *Manager {
	manager := &Manager{
		config:   alertConfig,
		storage:  storage,
		channels: make([]AlertChannel, 0),
		states:   make(map[string]*AlertState),
	}

	// Initialize alert channels based on configuration
	manager.initializeChannels()

	return manager
}

// initializeChannels sets up the configured alert channels
func (m *Manager) initializeChannels() {
	if m.config.Email.Enabled {
		emailChannel := NewEmailChannel(m.config.Email)
		m.channels = append(m.channels, emailChannel)
	}

	if m.config.Webhook.Enabled {
		webhookChannel := NewWebhookChannel(m.config.Webhook)
		m.channels = append(m.channels, webhookChannel)
	}

	log.Printf("📧 Initialized %d alert channels", len(m.channels))
}

// ProcessResult processes a monitoring result and generates alerts if needed
func (m *Manager) ProcessResult(result monitor.Result) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Get or create alert state for this site
	state := m.getOrCreateState(result.Name)

	// Update state based on the current result
	m.updateState(state, result)

	// Check for alert conditions
	alerts := m.checkAlertConditions(state, result)

	// Send any generated alerts
	for _, alert := range alerts {
		if err := m.sendAlert(alert); err != nil {
			log.Printf("❌ Failed to send alert: %v", err)
		} else {
			log.Printf("📧 Alert sent: %s", alert.String())
		}

		// Update state with alert information
		state.LastAlertTime = time.Now()
		if alert.Type == AlertTypeSiteDown {
			state.ActiveAlerts = append(state.ActiveAlerts, alert.ID)
		}
	}

	return nil
}

// getOrCreateState gets existing state or creates new one for a site
func (m *Manager) getOrCreateState(siteName string) *AlertState {
	if state, exists := m.states[siteName]; exists {
		return state
	}

	// Create new state
	state := &AlertState{
		SiteName:         siteName,
		IsDown:           false,
		ConsecutiveFails: 0,
		ActiveAlerts:     make([]string, 0),
	}
	m.states[siteName] = state
	return state
}

// updateState updates the alert state based on the current result
func (m *Manager) updateState(state *AlertState, result monitor.Result) {
	if result.Success {
		// Reset failure counter on success
		if state.ConsecutiveFails > 0 {
			state.ConsecutiveFails = 0
		}
		state.LastSuccessTime = result.Timestamp

		// Mark site as up if it was down
		if state.IsDown {
			state.IsDown = false
		}
	} else {
		// Increment failure counter
		state.ConsecutiveFails++
		state.LastFailTime = result.Timestamp

		// Mark site as down if threshold exceeded
		if state.ConsecutiveFails >= m.config.Thresholds.ConsecutiveFailures {
			state.IsDown = true
		}
	}
}

// checkAlertConditions checks if any alert conditions are met
func (m *Manager) checkAlertConditions(state *AlertState, result monitor.Result) []Alert {
	var alerts []Alert

	// Check if we should send alerts (respect cooldown)
	if m.shouldSkipDueToCooldown(state) {
		return alerts
	}

	// Check for site down alert
	if alert := m.checkSiteDownAlert(state, result); alert != nil {
		alerts = append(alerts, *alert)
	}

	// Check for site recovery alert
	if alert := m.checkSiteRecoveryAlert(state, result); alert != nil {
		alerts = append(alerts, *alert)
	}

	// Check for slow response alert
	if alert := m.checkSlowResponseAlert(state, result); alert != nil {
		alerts = append(alerts, *alert)
	}

	// Check for low uptime alert (less frequent check)
	if alert := m.checkLowUptimeAlert(state, result); alert != nil {
		alerts = append(alerts, *alert)
	}

	return alerts
}

// shouldSkipDueToCooldown checks if we should skip alerting due to cooldown period
func (m *Manager) shouldSkipDueToCooldown(state *AlertState) bool {
	if state.LastAlertTime.IsZero() {
		return false // No previous alert, can send
	}

	cooldown, err := m.config.Thresholds.GetAlertCooldown()
	if err != nil {
		log.Printf("⚠️ Invalid cooldown configuration: %v", err)
		return false // If config is invalid, don't skip
	}

	return time.Since(state.LastAlertTime) < cooldown
}

// checkSiteDownAlert checks for site down conditions
func (m *Manager) checkSiteDownAlert(state *AlertState, result monitor.Result) *Alert {
	// Only alert if site just went down (threshold reached) and wasn't already down
	if state.ConsecutiveFails == m.config.Thresholds.ConsecutiveFailures && !state.IsDown {
		return &Alert{
			ID:               uuid.New().String(),
			Type:             AlertTypeSiteDown,
			Severity:         SeverityCritical,
			SiteName:         result.Name,
			SiteURL:          result.URL,
			Message:          fmt.Sprintf("Site %s is down", result.Name),
			Details:          fmt.Sprintf("Failed %d consecutive checks. Last error: %s", state.ConsecutiveFails, result.Error),
			Timestamp:        time.Now(),
			CurrentStatus:    result.Status,
			ConsecutiveFails: state.ConsecutiveFails,
			ErrorMessage:     result.Error,
		}
	}
	return nil
}

// checkSiteRecoveryAlert checks for site recovery conditions
func (m *Manager) checkSiteRecoveryAlert(state *AlertState, result monitor.Result) *Alert {
	// Alert if site just recovered (was down and now successful)
	if state.IsDown && result.Success && len(state.ActiveAlerts) > 0 {
		// Clear active alerts
		state.ActiveAlerts = make([]string, 0)

		return &Alert{
			ID:            uuid.New().String(),
			Type:          AlertTypeSiteUp,
			Severity:      SeverityInfo,
			SiteName:      result.Name,
			SiteURL:       result.URL,
			Message:       fmt.Sprintf("Site %s has recovered", result.Name),
			Details:       fmt.Sprintf("Site is responding normally. Response time: %v", result.Duration),
			Timestamp:     time.Now(),
			CurrentStatus: result.Status,
			ResponseTime:  result.Duration,
			Resolved:      true,
		}
	}
	return nil
}

// checkSlowResponseAlert checks for slow response conditions
func (m *Manager) checkSlowResponseAlert(state *AlertState, result monitor.Result) *Alert {
	if !result.Success {
		return nil // Don't alert on slow response if site is down
	}

	threshold, err := m.config.Thresholds.GetResponseTimeThreshold()
	if err != nil {
		return nil // Invalid configuration
	}

	if result.Duration > threshold {
		return &Alert{
			ID:            uuid.New().String(),
			Type:          AlertTypeSlowResponse,
			Severity:      SeverityWarning,
			SiteName:      result.Name,
			SiteURL:       result.URL,
			Message:       fmt.Sprintf("Site %s is responding slowly", result.Name),
			Details:       fmt.Sprintf("Response time %v exceeds threshold of %v", result.Duration, threshold),
			Timestamp:     time.Now(),
			CurrentStatus: result.Status,
			ResponseTime:  result.Duration,
		}
	}
	return nil
}

// checkLowUptimeAlert checks for low uptime conditions (less frequent)
func (m *Manager) checkLowUptimeAlert(state *AlertState, result monitor.Result) *Alert {
	// Only check uptime periodically (e.g., every 10 minutes)
	if time.Since(state.LastAlertTime) < 10*time.Minute {
		return nil
	}

	// Get uptime stats for the configured window
	window, err := m.config.Thresholds.GetUptimeWindow()
	if err != nil {
		return nil
	}

	since := time.Now().Add(-window)
	stats, err := m.storage.GetStats(result.Name, since)
	if err != nil || stats.TotalChecks < 10 { // Need minimum data
		return nil
	}

	if stats.SuccessRate < m.config.Thresholds.UptimeThreshold {
		return &Alert{
			ID:            uuid.New().String(),
			Type:          AlertTypeLowUptime,
			Severity:      SeverityWarning,
			SiteName:      result.Name,
			SiteURL:       result.URL,
			Message:       fmt.Sprintf("Site %s has low uptime", result.Name),
			Details:       fmt.Sprintf("Uptime %.1f%% is below threshold of %.1f%% over the last %v", stats.SuccessRate, m.config.Thresholds.UptimeThreshold, window),
			Timestamp:     time.Now(),
			UptimePercent: stats.SuccessRate,
		}
	}
	return nil
}

// sendAlert sends an alert through all configured channels
func (m *Manager) sendAlert(alert Alert) error {
	if len(m.channels) == 0 {
		log.Printf("⚠️ No alert channels configured, alert not sent: %s", alert.String())
		return nil
	}

	var errors []error
	for _, channel := range m.channels {
		if err := channel.Send(alert); err != nil {
			errors = append(errors, fmt.Errorf("channel %s: %w", channel.Name(), err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to send through some channels: %v", errors)
	}

	return nil
}

// TestChannels tests all configured alert channels
func (m *Manager) TestChannels() error {
	if len(m.channels) == 0 {
		return fmt.Errorf("no alert channels configured")
	}

	var errors []error
	for _, channel := range m.channels {
		log.Printf("🧪 Testing %s channel...", channel.Name())
		if err := channel.Test(); err != nil {
			errors = append(errors, fmt.Errorf("channel %s failed test: %w", channel.Name(), err))
		} else {
			log.Printf("✅ %s channel test successful", channel.Name())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("some channels failed tests: %v", errors)
	}

	log.Printf("✅ All %d alert channels tested successfully", len(m.channels))
	return nil
}

// GetAlertStates returns current alert states for all sites
func (m *Manager) GetAlertStates() map[string]AlertState {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy to avoid race conditions
	states := make(map[string]AlertState)
	for siteName, state := range m.states {
		states[siteName] = *state
	}
	return states
}
