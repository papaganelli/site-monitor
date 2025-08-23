package alerts

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"site-monitor/config"
	"strconv"
	"strings"
	"time"
)

// EmailChannel implements AlertChannel for email notifications
type EmailChannel struct {
	config config.EmailConfig
}

// NewEmailChannel creates a new email alert channel
func NewEmailChannel(cfg config.EmailConfig) *EmailChannel {
	return &EmailChannel{
		config: cfg,
	}
}

// Name returns the channel name
func (e *EmailChannel) Name() string {
	return "Email"
}

// Send sends an alert via email
func (e *EmailChannel) Send(alert Alert) error {
	// Prepare email content
	subject := e.generateSubject(alert)
	body, err := e.generateBody(alert)
	if err != nil {
		return fmt.Errorf("failed to generate email body: %w", err)
	}

	// Send to all recipients
	for _, recipient := range e.config.Recipients {
		if err := e.sendEmail(recipient, subject, body); err != nil {
			return fmt.Errorf("failed to send email to %s: %w", recipient, err)
		}
	}

	return nil
}

// Test sends a test email to verify configuration
func (e *EmailChannel) Test() error {
	if !e.config.Enabled {
		return fmt.Errorf("email alerts are disabled")
	}

	if len(e.config.Recipients) == 0 {
		return fmt.Errorf("no email recipients configured")
	}

	// Create a test alert
	testAlert := Alert{
		ID:        "test-alert-" + strconv.FormatInt(time.Now().Unix(), 10),
		Type:      AlertTypeSiteDown,
		Severity:  SeverityInfo,
		SiteName:  "Test Site",
		SiteURL:   "https://example.com",
		Message:   "This is a test alert",
		Details:   "Testing email alert configuration",
		Timestamp: time.Now(),
	}

	// Send test email to first recipient only
	subject := "[TEST] " + e.generateSubject(testAlert)
	body, err := e.generateTestBody(testAlert)
	if err != nil {
		return fmt.Errorf("failed to generate test email body: %w", err)
	}

	return e.sendEmail(e.config.Recipients[0], subject, body)
}

// generateSubject creates the email subject line
func (e *EmailChannel) generateSubject(alert Alert) string {
	var prefix string
	switch alert.Severity {
	case SeverityCritical:
		prefix = "[CRITICAL]"
	case SeverityWarning:
		prefix = "[WARNING]"
	default:
		prefix = "[INFO]"
	}

	switch alert.Type {
	case AlertTypeSiteDown:
		return fmt.Sprintf("%s Site Monitor - %s is DOWN", prefix, alert.SiteName)
	case AlertTypeSiteUp:
		return fmt.Sprintf("%s Site Monitor - %s RECOVERED", prefix, alert.SiteName)
	case AlertTypeSlowResponse:
		return fmt.Sprintf("%s Site Monitor - %s SLOW RESPONSE", prefix, alert.SiteName)
	case AlertTypeLowUptime:
		return fmt.Sprintf("%s Site Monitor - %s LOW UPTIME", prefix, alert.SiteName)
	default:
		return fmt.Sprintf("%s Site Monitor - %s ALERT", prefix, alert.SiteName)
	}
}

// generateBody creates the email body using HTML template
func (e *EmailChannel) generateBody(alert Alert) (string, error) {
	tmpl := template.Must(template.New("email").Parse(emailTemplate))

	var buf bytes.Buffer
	data := struct {
		Alert             Alert
		FormattedTime     string
		SeverityColor     string
		TypeIcon          string
		IsRecovery        bool
		FormattedDuration string
	}{
		Alert:             alert,
		FormattedTime:     alert.Timestamp.Format("2006-01-02 15:04:05 MST"),
		SeverityColor:     getSeverityColor(alert.Severity),
		TypeIcon:          getTypeIcon(alert.Type),
		IsRecovery:        alert.IsRecoveryAlert(),
		FormattedDuration: formatDuration(alert.ResponseTime),
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// generateTestBody creates a test email body
func (e *EmailChannel) generateTestBody(alert Alert) (string, error) {
	testTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Site Monitor Test Alert</title>
</head>
<body style="font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5;">
    <div style="background-color: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
        <h2 style="color: #2196F3; margin-top: 0;">🧪 Site Monitor Test Alert</h2>
        
        <p>This is a test alert to verify your email configuration is working correctly.</p>
        
        <div style="background-color: #E3F2FD; padding: 15px; border-radius: 4px; margin: 20px 0;">
            <h3 style="margin-top: 0; color: #1976D2;">Configuration Test Results:</h3>
            <ul>
                <li><strong>SMTP Server:</strong> {{.SMTPServer}}</li>
                <li><strong>From Address:</strong> {{.From}}</li>
                <li><strong>Recipients:</strong> {{.Recipients}}</li>
                <li><strong>TLS Enabled:</strong> {{.UseTLS}}</li>
            </ul>
        </div>
        
        <p><strong>Test sent at:</strong> {{.FormattedTime}}</p>
        
        <hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
        <p style="color: #666; font-size: 0.9em;">
            If you received this email, your Site Monitor alert configuration is working correctly!<br>
            You can now expect to receive alerts when your monitored sites experience issues.
        </p>
    </div>
</body>
</html>`

	tmpl := template.Must(template.New("testEmail").Parse(testTemplate))

	var buf bytes.Buffer
	data := struct {
		SMTPServer    string
		From          string
		Recipients    string
		UseTLS        bool
		FormattedTime string
	}{
		SMTPServer:    e.config.SMTPServer,
		From:          e.config.From,
		Recipients:    strings.Join(e.config.Recipients, ", "),
		UseTLS:        e.config.UseTLS,
		FormattedTime: time.Now().Format("2006-01-02 15:04:05 MST"),
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// sendEmail sends an email using SMTP
func (e *EmailChannel) sendEmail(recipient, subject, body string) error {
	// Parse SMTP server and port
	serverParts := strings.Split(e.config.SMTPServer, ":")
	if len(serverParts) != 2 {
		return fmt.Errorf("invalid SMTP server format: %s", e.config.SMTPServer)
	}

	server := serverParts[0]
	port := serverParts[1]

	// Set up authentication
	auth := smtp.PlainAuth("", e.config.Username, e.config.Password, server)

	// Prepare message
	from := e.config.From
	if from == "" {
		from = e.config.Username // Use username as from if not specified
	}

	msg := fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=UTF-8\r\n"+
		"\r\n"+
		"%s",
		recipient, from, subject, body)

	// Send the email
	addr := server + ":" + port
	return smtp.SendMail(addr, auth, from, []string{recipient}, []byte(msg))
}

// getSeverityColor returns the HTML color for the severity
func getSeverityColor(severity AlertSeverity) string {
	switch severity {
	case SeverityCritical:
		return "#F44336" // Red
	case SeverityWarning:
		return "#FF9800" // Orange
	default:
		return "#2196F3" // Blue
	}
}

// getTypeIcon returns an emoji icon for the alert type
func getTypeIcon(alertType AlertType) string {
	switch alertType {
	case AlertTypeSiteDown:
		return "🚨"
	case AlertTypeSiteUp:
		return "✅"
	case AlertTypeSlowResponse:
		return "⚠️"
	case AlertTypeLowUptime:
		return "📉"
	default:
		return "🔔"
	}
}

// formatDuration formats a duration for display
func formatDuration(d time.Duration) string {
	if d == 0 {
		return "N/A"
	}

	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}

	return d.Round(time.Millisecond).String()
}

// emailTemplate is the HTML template for alert emails
const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Site Monitor Alert</title>
</head>
<body style="font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5;">
    <div style="background-color: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
        <h2 style="color: {{.SeverityColor}}; margin-top: 0;">{{.TypeIcon}} Site Monitor Alert</h2>
        
        <div style="background-color: {{if .IsRecovery}}#E8F5E8{{else}}#FFEBEE{{end}}; padding: 15px; border-radius: 4px; margin: 20px 0;">
            <h3 style="margin-top: 0; color: {{.SeverityColor}};">{{.Alert.Message}}</h3>
            {{if .Alert.Details}}<p>{{.Alert.Details}}</p>{{end}}
        </div>
        
        <table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
            <tr>
                <td style="padding: 8px; border-bottom: 1px solid #eee; font-weight: bold; width: 120px;">Site:</td>
                <td style="padding: 8px; border-bottom: 1px solid #eee;">{{.Alert.SiteName}}</td>
            </tr>
            <tr>
                <td style="padding: 8px; border-bottom: 1px solid #eee; font-weight: bold;">URL:</td>
                <td style="padding: 8px; border-bottom: 1px solid #eee;"><a href="{{.Alert.SiteURL}}" style="color: #2196F3;">{{.Alert.SiteURL}}</a></td>
            </tr>
            <tr>
                <td style="padding: 8px; border-bottom: 1px solid #eee; font-weight: bold;">Time:</td>
                <td style="padding: 8px; border-bottom: 1px solid #eee;">{{.FormattedTime}}</td>
            </tr>
            <tr>
                <td style="padding: 8px; border-bottom: 1px solid #eee; font-weight: bold;">Severity:</td>
                <td style="padding: 8px; border-bottom: 1px solid #eee;"><span style="color: {{.SeverityColor}}; font-weight: bold;">{{.Alert.Severity}}</span></td>
            </tr>
            {{if .Alert.CurrentStatus}}
            <tr>
                <td style="padding: 8px; border-bottom: 1px solid #eee; font-weight: bold;">HTTP Status:</td>
                <td style="padding: 8px; border-bottom: 1px solid #eee;">{{.Alert.CurrentStatus}}</td>
            </tr>
            {{end}}
            {{if .Alert.ResponseTime}}
            <tr>
                <td style="padding: 8px; border-bottom: 1px solid #eee; font-weight: bold;">Response Time:</td>
                <td style="padding: 8px; border-bottom: 1px solid #eee;">{{.FormattedDuration}}</td>
            </tr>
            {{end}}
            {{if .Alert.ConsecutiveFails}}
            <tr>
                <td style="padding: 8px; border-bottom: 1px solid #eee; font-weight: bold;">Consecutive Fails:</td>
                <td style="padding: 8px; border-bottom: 1px solid #eee;">{{.Alert.ConsecutiveFails}}</td>
            </tr>
            {{end}}
            {{if .Alert.UptimePercent}}
            <tr>
                <td style="padding: 8px; border-bottom: 1px solid #eee; font-weight: bold;">Uptime:</td>
                <td style="padding: 8px; border-bottom: 1px solid #eee;">{{printf "%.1f%%" .Alert.UptimePercent}}</td>
            </tr>
            {{end}}
            {{if .Alert.ErrorMessage}}
            <tr>
                <td style="padding: 8px; border-bottom: 1px solid #eee; font-weight: bold;">Error:</td>
                <td style="padding: 8px; border-bottom: 1px solid #eee; color: #F44336; font-family: monospace;">{{.Alert.ErrorMessage}}</td>
            </tr>
            {{end}}
        </table>
        
        {{if not .IsRecovery}}
        <div style="background-color: #FFF3E0; padding: 15px; border-radius: 4px; margin: 20px 0;">
            <h4 style="margin-top: 0; color: #F57C00;">🔧 Recommended Actions:</h4>
            <ul style="margin-bottom: 0;">
                <li>Check the site manually: <a href="{{.Alert.SiteURL}}" style="color: #2196F3;">{{.Alert.SiteURL}}</a></li>
                <li>Verify server status and logs</li>
                <li>Check network connectivity</li>
                {{if eq .Alert.Type "slow_response"}}
                <li>Monitor server performance and database queries</li>
                <li>Check for high CPU or memory usage</li>
                {{end}}
            </ul>
        </div>
        {{end}}
        
        <hr style="border: none; border-top: 1px solid #eee; margin: 20px 0;">
        <p style="color: #666; font-size: 0.9em;">
            This alert was generated by Site Monitor.<br>
            Alert ID: {{.Alert.ID}}
        </p>
    </div>
</body>
</html>`
