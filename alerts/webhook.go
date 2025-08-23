package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"site-monitor/config"
	"strconv"
	"time"
)

// WebhookChannel implements AlertChannel for webhook notifications
type WebhookChannel struct {
	config config.WebhookConfig
	client *http.Client
}

// NewWebhookChannel creates a new webhook alert channel
func NewWebhookChannel(cfg config.WebhookConfig) *WebhookChannel {
	timeout := 30 * time.Second
	if cfg.Timeout != "" {
		if parsed, err := time.ParseDuration(cfg.Timeout); err == nil {
			timeout = parsed
		}
	}

	return &WebhookChannel{
		config: cfg,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// Name returns the channel name
func (w *WebhookChannel) Name() string {
	return fmt.Sprintf("Webhook (%s)", w.config.Format)
}

// Send sends an alert via webhook
func (w *WebhookChannel) Send(alert Alert) error {
	payload, err := w.generatePayload(alert)
	if err != nil {
		return fmt.Errorf("failed to generate webhook payload: %w", err)
	}

	// Retry logic
	maxRetries := w.config.RetryCount
	if maxRetries <= 0 {
		maxRetries = 1
	}

	var lastError error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := w.sendRequest(payload); err != nil {
			lastError = err
			if attempt < maxRetries {
				// Wait before retry (exponential backoff)
				waitTime := time.Duration(attempt*attempt) * time.Second
				time.Sleep(waitTime)
				continue
			}
		} else {
			return nil // Success
		}
	}

	return fmt.Errorf("failed to send webhook after %d attempts: %w", maxRetries, lastError)
}

// Test sends a test webhook to verify configuration
func (w *WebhookChannel) Test() error {
	if !w.config.Enabled {
		return fmt.Errorf("webhook alerts are disabled")
	}

	if w.config.URL == "" {
		return fmt.Errorf("webhook URL not configured")
	}

	// Create a test alert
	testAlert := Alert{
		ID:        "test-alert-" + strconv.FormatInt(time.Now().Unix(), 10),
		Type:      AlertTypeSiteDown,
		Severity:  SeverityInfo,
		SiteName:  "Test Site",
		SiteURL:   "https://example.com",
		Message:   "This is a test alert",
		Details:   "Testing webhook alert configuration",
		Timestamp: time.Now(),
	}

	return w.Send(testAlert)
}

// generatePayload creates the webhook payload based on the configured format
func (w *WebhookChannel) generatePayload(alert Alert) ([]byte, error) {
	switch w.config.Format {
	case "slack":
		return w.generateSlackPayload(alert)
	case "discord":
		return w.generateDiscordPayload(alert)
	case "teams":
		return w.generateTeamsPayload(alert)
	default:
		return w.generateGenericPayload(alert)
	}
}

// generateSlackPayload creates a Slack-formatted payload
func (w *WebhookChannel) generateSlackPayload(alert Alert) ([]byte, error) {
	color := "good"
	if alert.Severity == SeverityCritical {
		color = "danger"
	} else if alert.Severity == SeverityWarning {
		color = "warning"
	}

	payload := SlackPayload{
		Text: alert.String(),
		Attachments: []SlackAttachment{
			{
				Color:  color,
				Title:  alert.Message,
				Text:   alert.Details,
				Fields: w.buildSlackFields(alert),
				Footer: "Site Monitor",
				Ts:     alert.Timestamp.Unix(),
			},
		},
	}

	return json.Marshal(payload)
}

// generateDiscordPayload creates a Discord-formatted payload
func (w *WebhookChannel) generateDiscordPayload(alert Alert) ([]byte, error) {
	color := 3447003 // Blue
	if alert.Severity == SeverityCritical {
		color = 15158332 // Red
	} else if alert.Severity == SeverityWarning {
		color = 15105570 // Orange
	}

	embed := DiscordEmbed{
		Title:       alert.Message,
		Description: alert.Details,
		Color:       color,
		Fields:      w.buildDiscordFields(alert),
		Footer: DiscordFooter{
			Text: "Site Monitor",
		},
		Timestamp: alert.Timestamp.Format(time.RFC3339),
	}

	payload := DiscordPayload{
		Content: alert.String(),
		Embeds:  []DiscordEmbed{embed},
	}

	return json.Marshal(payload)
}

// generateTeamsPayload creates a Microsoft Teams-formatted payload
func (w *WebhookChannel) generateTeamsPayload(alert Alert) ([]byte, error) {
	themeColor := "0078D4" // Blue
	if alert.Severity == SeverityCritical {
		themeColor = "D13438" // Red
	} else if alert.Severity == SeverityWarning {
		themeColor = "FF8C00" // Orange
	}

	payload := TeamsPayload{
		Type:       "MessageCard",
		Context:    "https://schema.org/extensions",
		Summary:    alert.Message,
		ThemeColor: themeColor,
		Sections: []TeamsSection{
			{
				ActivityTitle:    alert.Message,
				ActivitySubtitle: fmt.Sprintf("Site: %s", alert.SiteName),
				Text:             alert.Details,
				Facts:            w.buildTeamsFacts(alert),
			},
		},
		PotentialAction: []TeamsAction{
			{
				Type: "OpenUri",
				Name: "View Site",
				Targets: []TeamsTarget{
					{
						OS:  "default",
						URI: alert.SiteURL,
					},
				},
			},
		},
	}

	return json.Marshal(payload)
}

// generateGenericPayload creates a generic JSON payload
func (w *WebhookChannel) generateGenericPayload(alert Alert) ([]byte, error) {
	payload := GenericPayload{
		Alert:     alert,
		Message:   alert.String(),
		Timestamp: alert.Timestamp.Format(time.RFC3339),
	}

	return json.Marshal(payload)
}

// sendRequest sends the HTTP request to the webhook URL
func (w *WebhookChannel) sendRequest(payload []byte) error {
	req, err := http.NewRequest("POST", w.config.URL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set default headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SiteMonitor/1.0")

	// Add custom headers
	for key, value := range w.config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// buildSlackFields builds Slack attachment fields from alert data
func (w *WebhookChannel) buildSlackFields(alert Alert) []SlackField {
	var fields []SlackField

	fields = append(fields, SlackField{
		Title: "Site",
		Value: alert.SiteName,
		Short: true,
	})

	fields = append(fields, SlackField{
		Title: "Severity",
		Value: string(alert.Severity),
		Short: true,
	})

	if alert.CurrentStatus > 0 {
		fields = append(fields, SlackField{
			Title: "HTTP Status",
			Value: strconv.Itoa(alert.CurrentStatus),
			Short: true,
		})
	}

	if alert.ResponseTime > 0 {
		fields = append(fields, SlackField{
			Title: "Response Time",
			Value: formatDuration(alert.ResponseTime),
			Short: true,
		})
	}

	if alert.ConsecutiveFails > 0 {
		fields = append(fields, SlackField{
			Title: "Consecutive Failures",
			Value: strconv.Itoa(alert.ConsecutiveFails),
			Short: true,
		})
	}

	if alert.UptimePercent > 0 {
		fields = append(fields, SlackField{
			Title: "Uptime",
			Value: fmt.Sprintf("%.1f%%", alert.UptimePercent),
			Short: true,
		})
	}

	return fields
}

// buildDiscordFields builds Discord embed fields from alert data
func (w *WebhookChannel) buildDiscordFields(alert Alert) []DiscordField {
	var fields []DiscordField

	fields = append(fields, DiscordField{
		Name:   "Site",
		Value:  alert.SiteName,
		Inline: true,
	})

	fields = append(fields, DiscordField{
		Name:   "Severity",
		Value:  string(alert.Severity),
		Inline: true,
	})

	if alert.CurrentStatus > 0 {
		fields = append(fields, DiscordField{
			Name:   "HTTP Status",
			Value:  strconv.Itoa(alert.CurrentStatus),
			Inline: true,
		})
	}

	if alert.ResponseTime > 0 {
		fields = append(fields, DiscordField{
			Name:   "Response Time",
			Value:  formatDuration(alert.ResponseTime),
			Inline: true,
		})
	}

	return fields
}

// buildTeamsFacts builds Teams message facts from alert data
func (w *WebhookChannel) buildTeamsFacts(alert Alert) []TeamsFact {
	var facts []TeamsFact

	facts = append(facts, TeamsFact{
		Name:  "Site",
		Value: alert.SiteName,
	})

	facts = append(facts, TeamsFact{
		Name:  "URL",
		Value: alert.SiteURL,
	})

	facts = append(facts, TeamsFact{
		Name:  "Severity",
		Value: string(alert.Severity),
	})

	if alert.CurrentStatus > 0 {
		facts = append(facts, TeamsFact{
			Name:  "HTTP Status",
			Value: strconv.Itoa(alert.CurrentStatus),
		})
	}

	if alert.ResponseTime > 0 {
		facts = append(facts, TeamsFact{
			Name:  "Response Time",
			Value: formatDuration(alert.ResponseTime),
		})
	}

	return facts
}

// Payload structures for different webhook formats

// SlackPayload represents a Slack webhook payload
type SlackPayload struct {
	Text        string            `json:"text"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Color  string       `json:"color"`
	Title  string       `json:"title"`
	Text   string       `json:"text"`
	Fields []SlackField `json:"fields"`
	Footer string       `json:"footer"`
	Ts     int64        `json:"ts"`
}

type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// DiscordPayload represents a Discord webhook payload
type DiscordPayload struct {
	Content string         `json:"content"`
	Embeds  []DiscordEmbed `json:"embeds"`
}

type DiscordEmbed struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Color       int            `json:"color"`
	Fields      []DiscordField `json:"fields"`
	Footer      DiscordFooter  `json:"footer"`
	Timestamp   string         `json:"timestamp"`
}

type DiscordField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type DiscordFooter struct {
	Text string `json:"text"`
}

// TeamsPayload represents a Microsoft Teams webhook payload
type TeamsPayload struct {
	Type            string         `json:"@type"`
	Context         string         `json:"@context"`
	Summary         string         `json:"summary"`
	ThemeColor      string         `json:"themeColor"`
	Sections        []TeamsSection `json:"sections"`
	PotentialAction []TeamsAction  `json:"potentialAction"`
}

type TeamsSection struct {
	ActivityTitle    string      `json:"activityTitle"`
	ActivitySubtitle string      `json:"activitySubtitle"`
	Text             string      `json:"text"`
	Facts            []TeamsFact `json:"facts"`
}

type TeamsFact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TeamsAction struct {
	Type    string        `json:"@type"`
	Name    string        `json:"name"`
	Targets []TeamsTarget `json:"targets"`
}

type TeamsTarget struct {
	OS  string `json:"os"`
	URI string `json:"uri"`
}

// GenericPayload represents a generic webhook payload
type GenericPayload struct {
	Alert     Alert  `json:"alert"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
}
