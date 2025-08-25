package alerts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	textTemplate "text/template"
	"time"
)

// AlertTemplate defines a customizable alert template
type AlertTemplate struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	AlertType    AlertType              `json:"alert_type"`
	Channel      ChannelType            `json:"channel"`
	Subject      string                 `json:"subject"`
	Body         string                 `json:"body"`
	Variables    map[string]Variable    `json:"variables"`
	Conditions   []TemplateCondition    `json:"conditions"`
	Format       TemplateFormat         `json:"format"`
	Style        TemplateStyle          `json:"style"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	UsageCount   int64                  `json:"usage_count"`
	IsDefault    bool                   `json:"is_default"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

// ChannelType defines the target channel for the template
type ChannelType string

const (
	ChannelEmail   ChannelType = "email"
	ChannelSlack   ChannelType = "slack"
	ChannelDiscord ChannelType = "discord"
	ChannelTeams   ChannelType = "teams"
	ChannelWebhook ChannelType = "webhook"
	ChannelSMS     ChannelType = "sms"
)

// TemplateFormat defines the content format
type TemplateFormat string

const (
	FormatPlainText TemplateFormat = "plain"
	FormatHTML      TemplateFormat = "html"
	FormatMarkdown  TemplateFormat = "markdown"
	FormatJSON      TemplateFormat = "json"
)

// TemplateStyle defines visual styling options
type TemplateStyle struct {
	Theme         string            `json:"theme"` // dark, light, corporate, minimal
	Colors        map[string]string `json:"colors"`
	FontFamily    string            `json:"font_family"`
	IncludeLogo   bool              `json:"include_logo"`
	IncludeHeader bool              `json:"include_header"`
	IncludeFooter bool              `json:"include_footer"`
}

// Variable defines template variables and their properties
type Variable struct {
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	Type         VariableType `json:"type"`
	Required     bool         `json:"required"`
	DefaultValue interface{}  `json:"default_value"`
	Format       string       `json:"format"`  // For dates, numbers, etc.
	Options      []string     `json:"options"` // For enum-type variables
}

// VariableType defines the type of template variable
type VariableType string

const (
	VarTypeString   VariableType = "string"
	VarTypeNumber   VariableType = "number"
	VarTypeBoolean  VariableType = "boolean"
	VarTypeDateTime VariableType = "datetime"
	VarTypeDuration VariableType = "duration"
	VarTypeURL      VariableType = "url"
	VarTypeEnum     VariableType = "enum"
)

// TemplateCondition defines when to use a template
type TemplateCondition struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

// TemplateManager manages alert templates
type TemplateManager struct {
	templates map[string]*AlertTemplate
	defaults  map[AlertType]map[ChannelType]*AlertTemplate
}

// NewTemplateManager creates a new template manager
func NewTemplateManager() *TemplateManager {
	tm := &TemplateManager{
		templates: make(map[string]*AlertTemplate),
		defaults:  make(map[AlertType]map[ChannelType]*AlertTemplate),
	}

	// Initialize default templates
	tm.initializeDefaultTemplates()

	return tm
}

// initializeDefaultTemplates creates built-in templates
func (tm *TemplateManager) initializeDefaultTemplates() {
	// Initialize nested maps
	for _, alertType := range []AlertType{AlertTypeSiteDown, AlertTypeSiteUp, AlertTypeSlowResponse, AlertTypeLowUptime} {
		tm.defaults[alertType] = make(map[ChannelType]*AlertTemplate)
	}

	// Site Down Templates
	tm.addDefaultTemplate(&AlertTemplate{
		ID:        "default-site-down-email",
		Name:      "Site Down - Email",
		AlertType: AlertTypeSiteDown,
		Channel:   ChannelEmail,
		Format:    FormatHTML,
		Subject:   "🚨 CRITICAL: {{.SiteName}} is DOWN",
		Body: `<!DOCTYPE html>
<html>
<head><title>Site Down Alert</title></head>
<body style="font-family: Arial, sans-serif; margin: 20px; background-color: #f8f9fa;">
	<div style="background: linear-gradient(135deg, #dc3545, #c82333); color: white; padding: 20px; border-radius: 8px 8px 0 0;">
		<h1 style="margin: 0; font-size: 24px;">🚨 SITE DOWN ALERT</h1>
		<p style="margin: 5px 0 0 0; opacity: 0.9;">Critical monitoring alert</p>
	</div>
	
	<div style="background: white; padding: 20px; border-radius: 0 0 8px 8px; border: 1px solid #dee2e6;">
		<div style="background: #f8d7da; color: #721c24; padding: 15px; border-radius: 4px; margin-bottom: 20px;">
			<h2 style="margin: 0 0 10px 0;">{{.SiteName}} is not responding</h2>
			<p style="margin: 0;">The site has been down for {{.ConsecutiveFails}} consecutive checks.</p>
		</div>
		
		<table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
			<tr><td style="padding: 8px; font-weight: bold; width: 30%;">Site:</td><td style="padding: 8px;">{{.SiteName}}</td></tr>
			<tr><td style="padding: 8px; font-weight: bold;">URL:</td><td style="padding: 8px;"><a href="{{.SiteURL}}">{{.SiteURL}}</a></td></tr>
			<tr><td style="padding: 8px; font-weight: bold;">Time:</td><td style="padding: 8px;">{{.Timestamp | formatTime}}</td></tr>
			<tr><td style="padding: 8px; font-weight: bold;">Status Code:</td><td style="padding: 8px;">{{.CurrentStatus}}</td></tr>
			<tr><td style="padding: 8px; font-weight: bold;">Consecutive Failures:</td><td style="padding: 8px;">{{.ConsecutiveFails}}</td></tr>
			{{if .ErrorMessage}}<tr><td style="padding: 8px; font-weight: bold;">Error:</td><td style="padding: 8px; color: #dc3545; font-family: monospace;">{{.ErrorMessage}}</td></tr>{{end}}
		</table>
		
		<div style="background: #fff3cd; color: #856404; padding: 15px; border-radius: 4px;">
			<h4 style="margin: 0 0 10px 0;">🔧 Recommended Actions:</h4>
			<ul style="margin: 0; padding-left: 20px;">
				<li>Check server status and logs immediately</li>
				<li>Verify network connectivity</li>
				<li>Contact your hosting provider if needed</li>
				<li>Monitor for automatic recovery</li>
			</ul>
		</div>
	</div>
	
	<div style="text-align: center; margin-top: 20px; color: #6c757d; font-size: 12px;">
		Generated by Site Monitor | Alert ID: {{.ID}}
	</div>
</body>
</html>`,
		IsDefault: true,
		Variables: map[string]Variable{
			"SiteName":         {Name: "SiteName", Type: VarTypeString, Required: true},
			"SiteURL":          {Name: "SiteURL", Type: VarTypeURL, Required: true},
			"ConsecutiveFails": {Name: "ConsecutiveFails", Type: VarTypeNumber, Required: true},
			"CurrentStatus":    {Name: "CurrentStatus", Type: VarTypeNumber, Required: false},
			"ErrorMessage":     {Name: "ErrorMessage", Type: VarTypeString, Required: false},
			"Timestamp":        {Name: "Timestamp", Type: VarTypeDateTime, Required: true, Format: "2006-01-02 15:04:05"},
		},
	})

	tm.addDefaultTemplate(&AlertTemplate{
		ID:        "default-site-down-slack",
		Name:      "Site Down - Slack",
		AlertType: AlertTypeSiteDown,
		Channel:   ChannelSlack,
		Format:    FormatJSON,
		Subject:   "",
		Body: `{
	"text": "🚨 *SITE DOWN ALERT*",
	"attachments": [
		{
			"color": "danger",
			"title": "{{.SiteName}} is not responding",
			"text": "The site has been down for {{.ConsecutiveFails}} consecutive checks.",
			"fields": [
				{
					"title": "Site",
					"value": "{{.SiteName}}",
					"short": true
				},
				{
					"title": "Status Code", 
					"value": "{{.CurrentStatus}}",
					"short": true
				},
				{
					"title": "URL",
					"value": "<{{.SiteURL}}|{{.SiteURL}}>",
					"short": false
				}
				{{if .ErrorMessage}},
				{
					"title": "Error",
					"value": "` + "`{{.ErrorMessage}}`" + `",
					"short": false
				}
				{{end}}
			],
			"footer": "Site Monitor",
			"ts": {{.Timestamp | unixTime}}
		}
	]
}`,
		IsDefault: true,
	})

	// Site Recovery Templates
	tm.addDefaultTemplate(&AlertTemplate{
		ID:        "default-site-up-email",
		Name:      "Site Recovery - Email",
		AlertType: AlertTypeSiteUp,
		Channel:   ChannelEmail,
		Format:    FormatHTML,
		Subject:   "✅ RESOLVED: {{.SiteName}} is back online",
		Body: `<!DOCTYPE html>
<html>
<head><title>Site Recovery Alert</title></head>
<body style="font-family: Arial, sans-serif; margin: 20px; background-color: #f8f9fa;">
	<div style="background: linear-gradient(135deg, #28a745, #218838); color: white; padding: 20px; border-radius: 8px 8px 0 0;">
		<h1 style="margin: 0; font-size: 24px;">✅ SITE RECOVERED</h1>
		<p style="margin: 5px 0 0 0; opacity: 0.9;">Recovery notification</p>
	</div>
	
	<div style="background: white; padding: 20px; border-radius: 0 0 8px 8px; border: 1px solid #dee2e6;">
		<div style="background: #d4edda; color: #155724; padding: 15px; border-radius: 4px; margin-bottom: 20px;">
			<h2 style="margin: 0 0 10px 0;">{{.SiteName}} is back online! 🎉</h2>
			<p style="margin: 0;">The site is now responding normally with a {{.ResponseTime | formatDuration}} response time.</p>
		</div>
		
		<table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
			<tr><td style="padding: 8px; font-weight: bold; width: 30%;">Site:</td><td style="padding: 8px;">{{.SiteName}}</td></tr>
			<tr><td style="padding: 8px; font-weight: bold;">URL:</td><td style="padding: 8px;"><a href="{{.SiteURL}}">{{.SiteURL}}</a></td></tr>
			<tr><td style="padding: 8px; font-weight: bold;">Recovered At:</td><td style="padding: 8px;">{{.Timestamp | formatTime}}</td></tr>
			<tr><td style="padding: 8px; font-weight: bold;">Response Time:</td><td style="padding: 8px;">{{.ResponseTime | formatDuration}}</td></tr>
			<tr><td style="padding: 8px; font-weight: bold;">Status Code:</td><td style="padding: 8px; color: #28a745;">{{.CurrentStatus}}</td></tr>
		</table>
		
		<div style="background: #d1ecf1; color: #0c5460; padding: 15px; border-radius: 4px;">
			<p style="margin: 0;"><strong>✅ All systems are now operational.</strong> Continue monitoring for stability.</p>
		</div>
	</div>
</body>
</html>`,
		IsDefault: true,
	})

	// Slow Response Template
	tm.addDefaultTemplate(&AlertTemplate{
		ID:        "default-slow-response-email",
		Name:      "Slow Response - Email",
		AlertType: AlertTypeSlowResponse,
		Channel:   ChannelEmail,
		Format:    FormatHTML,
		Subject:   "⚠️ PERFORMANCE: {{.SiteName}} responding slowly",
		Body: `<!DOCTYPE html>
<html>
<head><title>Performance Alert</title></head>
<body style="font-family: Arial, sans-serif; margin: 20px; background-color: #f8f9fa;">
	<div style="background: linear-gradient(135deg, #ffc107, #e0a800); color: white; padding: 20px; border-radius: 8px 8px 0 0;">
		<h1 style="margin: 0; font-size: 24px;">⚠️ PERFORMANCE ALERT</h1>
		<p style="margin: 5px 0 0 0; opacity: 0.9;">Slow response detected</p>
	</div>
	
	<div style="background: white; padding: 20px; border-radius: 0 0 8px 8px; border: 1px solid #dee2e6;">
		<div style="background: #fff3cd; color: #856404; padding: 15px; border-radius: 4px; margin-bottom: 20px;">
			<h2 style="margin: 0 0 10px 0;">{{.SiteName}} is responding slowly</h2>
			<p style="margin: 0;">Response time of {{.ResponseTime | formatDuration}} exceeds the configured threshold.</p>
		</div>
		
		<table style="width: 100%; border-collapse: collapse; margin: 20px 0;">
			<tr><td style="padding: 8px; font-weight: bold; width: 30%;">Site:</td><td style="padding: 8px;">{{.SiteName}}</td></tr>
			<tr><td style="padding: 8px; font-weight: bold;">Response Time:</td><td style="padding: 8px; color: #ffc107; font-weight: bold;">{{.ResponseTime | formatDuration}}</td></tr>
			<tr><td style="padding: 8px; font-weight: bold;">Status Code:</td><td style="padding: 8px;">{{.CurrentStatus}}</td></tr>
			<tr><td style="padding: 8px; font-weight: bold;">Time:</td><td style="padding: 8px;">{{.Timestamp | formatTime}}</td></tr>
		</table>
		
		<div style="background: #e2e3e5; color: #383d41; padding: 15px; border-radius: 4px;">
			<h4 style="margin: 0 0 10px 0;">🔍 Potential Causes:</h4>
			<ul style="margin: 0; padding-left: 20px;">
				<li>High server load or CPU usage</li>
				<li>Database performance issues</li>
				<li>Network latency</li>
				<li>Large response payloads</li>
			</ul>
		</div>
	</div>
</body>
</html>`,
		IsDefault: true,
	})

	// Custom minimalist template
	tm.addDefaultTemplate(&AlertTemplate{
		ID:        "minimal-site-down-email",
		Name:      "Minimal Site Down - Email",
		AlertType: AlertTypeSiteDown,
		Channel:   ChannelEmail,
		Format:    FormatHTML,
		Subject:   "Site Down: {{.SiteName}}",
		Body: `<div style="font-family: 'Helvetica Neue', Arial, sans-serif; max-width: 600px; margin: 0 auto; background: #fff;">
	<div style="padding: 40px; text-align: center;">
		<h1 style="color: #e74c3c; font-weight: 300; font-size: 28px; margin-bottom: 20px;">Site Down</h1>
		<p style="font-size: 18px; color: #34495e; margin-bottom: 30px;">{{.SiteName}} is not responding</p>
		
		<div style="background: #ecf0f1; padding: 20px; border-radius: 4px; text-align: left;">
			<p><strong>URL:</strong> {{.SiteURL}}</p>
			<p><strong>Time:</strong> {{.Timestamp | formatTime}}</p>
			{{if .ErrorMessage}}<p><strong>Error:</strong> {{.ErrorMessage}}</p>{{end}}
		</div>
		
		<p style="color: #7f8c8d; font-size: 14px; margin-top: 30px;">Site Monitor Alert</p>
	</div>
</div>`,
		IsDefault: false,
		Style: TemplateStyle{
			Theme:       "minimal",
			FontFamily:  "Helvetica Neue",
			IncludeLogo: false,
		},
	})
}

// addDefaultTemplate adds a template to both templates and defaults
func (tm *TemplateManager) addDefaultTemplate(template *AlertTemplate) {
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()

	tm.templates[template.ID] = template
	tm.defaults[template.AlertType][template.Channel] = template
}

// GetTemplate retrieves a template by ID
func (tm *TemplateManager) GetTemplate(id string) (*AlertTemplate, bool) {
	template, exists := tm.templates[id]
	return template, exists
}

// GetDefaultTemplate gets the default template for an alert type and channel
func (tm *TemplateManager) GetDefaultTemplate(alertType AlertType, channel ChannelType) (*AlertTemplate, bool) {
	if channelTemplates, exists := tm.defaults[alertType]; exists {
		if template, exists := channelTemplates[channel]; exists {
			return template, true
		}
	}
	return nil, false
}

// AddTemplate adds a new custom template
func (tm *TemplateManager) AddTemplate(template *AlertTemplate) error {
	// Validate template
	if err := tm.validateTemplate(template); err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	template.ID = tm.generateTemplateID(template)
	template.CreatedAt = time.Now()
	template.UpdatedAt = time.Now()
	template.UsageCount = 0

	tm.templates[template.ID] = template
	return nil
}

// UpdateTemplate updates an existing template
func (tm *TemplateManager) UpdateTemplate(id string, updates *AlertTemplate) error {
	template, exists := tm.templates[id]
	if !exists {
		return fmt.Errorf("template with ID %s not found", id)
	}

	// Don't allow updating default templates
	if template.IsDefault {
		return fmt.Errorf("cannot update default template")
	}

	// Validate updates
	if err := tm.validateTemplate(updates); err != nil {
		return fmt.Errorf("template validation failed: %w", err)
	}

	// Update fields
	updates.ID = id
	updates.CreatedAt = template.CreatedAt
	updates.UpdatedAt = time.Now()
	updates.UsageCount = template.UsageCount

	tm.templates[id] = updates
	return nil
}

// DeleteTemplate removes a custom template
func (tm *TemplateManager) DeleteTemplate(id string) error {
	template, exists := tm.templates[id]
	if !exists {
		return fmt.Errorf("template with ID %s not found", id)
	}

	if template.IsDefault {
		return fmt.Errorf("cannot delete default template")
	}

	delete(tm.templates, id)
	return nil
}

// RenderTemplate renders a template with alert data
func (tm *TemplateManager) RenderTemplate(templateID string, alert Alert) (string, string, error) {
	template, exists := tm.templates[templateID]
	if !exists {
		return "", "", fmt.Errorf("template %s not found", templateID)
	}

	// Increment usage count
	template.UsageCount++

	// Prepare template data
	data := tm.prepareTemplateData(alert, template)

	// Render subject
	subject, err := tm.renderString(template.Subject, data, template.Format)
	if err != nil {
		return "", "", fmt.Errorf("failed to render subject: %w", err)
	}

	// Render body
	body, err := tm.renderString(template.Body, data, template.Format)
	if err != nil {
		return "", "", fmt.Errorf("failed to render body: %w", err)
	}

	return subject, body, nil
}

// titleCase capitalizes the first letter of each word (replacement for deprecated strings.Title)
func titleCase(s string) string {
	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, " ")
}

// renderString renders a template string with data
func (tm *TemplateManager) renderString(templateStr string, data interface{}, format TemplateFormat) (string, error) {
	funcMap := template.FuncMap{
		"formatTime":     formatTemplateTime,
		"formatDuration": formatTemplateDuration,
		"unixTime":       unixTime,
		"upper":          strings.ToUpper,
		"lower":          strings.ToLower,
		"title":          titleCase, // Use custom titleCase function
		"join":           strings.Join,
		"replace":        strings.Replace,
		"contains":       strings.Contains,
		"now":            time.Now,
		"add":            func(a, b int) int { return a + b },
		"sub":            func(a, b int) int { return a - b },
		"mul":            func(a, b int) int { return a * b },
		"div":            func(a, b int) int { return a / b },
	}

	switch format {
	case FormatHTML:
		tmpl, err := template.New("alert").Funcs(funcMap).Parse(templateStr)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}
		return buf.String(), nil

	case FormatPlainText, FormatMarkdown:
		tmpl, err := textTemplate.New("alert").Funcs(textTemplate.FuncMap(funcMap)).Parse(templateStr)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}
		return buf.String(), nil

	case FormatJSON:
		// For JSON templates, parse as text template but validate JSON after rendering
		tmpl, err := textTemplate.New("alert").Funcs(textTemplate.FuncMap(funcMap)).Parse(templateStr)
		if err != nil {
			return "", err
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}

		// Validate JSON
		var js json.RawMessage
		if err := json.Unmarshal(buf.Bytes(), &js); err != nil {
			return "", fmt.Errorf("rendered template is not valid JSON: %w", err)
		}

		return buf.String(), nil

	default:
		return "", fmt.Errorf("unsupported template format: %s", format)
	}
}

// prepareTemplateData converts alert to template-friendly data structure
func (tm *TemplateManager) prepareTemplateData(alert Alert, tmpl *AlertTemplate) map[string]interface{} {
	data := map[string]interface{}{
		"ID":               alert.ID,
		"Type":             string(alert.Type),
		"Severity":         string(alert.Severity),
		"SiteName":         alert.SiteName,
		"SiteURL":          alert.SiteURL,
		"Message":          alert.Message,
		"Details":          alert.Details,
		"Timestamp":        alert.Timestamp,
		"Resolved":         alert.Resolved,
		"CurrentStatus":    alert.CurrentStatus,
		"ResponseTime":     alert.ResponseTime,
		"ConsecutiveFails": alert.ConsecutiveFails,
		"UptimePercent":    alert.UptimePercent,
		"ErrorMessage":     alert.ErrorMessage,
		"ResolvedAt":       alert.ResolvedAt,
	}

	// Add custom fields if present
	for key, value := range tmpl.CustomFields {
		data[key] = value
	}

	return data
}

// validateTemplate validates template structure and content
func (tm *TemplateManager) validateTemplate(tmpl *AlertTemplate) error {
	if tmpl.Name == "" {
		return fmt.Errorf("template name is required")
	}

	if tmpl.Subject == "" && tmpl.Channel == ChannelEmail {
		return fmt.Errorf("subject is required for email templates")
	}

	if tmpl.Body == "" {
		return fmt.Errorf("template body is required")
	}

	// Validate template syntax by parsing
	switch tmpl.Format {
	case FormatHTML:
		_, err := template.New("test").Parse(tmpl.Body)
		if err != nil {
			return fmt.Errorf("invalid HTML template syntax: %w", err)
		}
	case FormatPlainText, FormatMarkdown:
		_, err := textTemplate.New("test").Parse(tmpl.Body)
		if err != nil {
			return fmt.Errorf("invalid text template syntax: %w", err)
		}
	case FormatJSON:
		// Parse as text template first
		_, err := textTemplate.New("test").Parse(tmpl.Body)
		if err != nil {
			return fmt.Errorf("invalid JSON template syntax: %w", err)
		}
	}

	return nil
}

// generateTemplateID generates a unique ID for the template
func (tm *TemplateManager) generateTemplateID(tmpl *AlertTemplate) string {
	base := strings.ToLower(strings.ReplaceAll(tmpl.Name, " ", "-"))
	base = strings.ReplaceAll(base, "_", "-")

	// Add suffix if ID already exists
	id := base
	counter := 1
	for {
		if _, exists := tm.templates[id]; !exists {
			break
		}
		id = fmt.Sprintf("%s-%d", base, counter)
		counter++
	}

	return id
}

// ListTemplates returns all templates, optionally filtered
func (tm *TemplateManager) ListTemplates(filters map[string]interface{}) []*AlertTemplate {
	var result []*AlertTemplate

	for _, template := range tm.templates {
		include := true

		// Apply filters
		if alertType, ok := filters["alert_type"]; ok {
			if template.AlertType != AlertType(alertType.(string)) {
				include = false
			}
		}

		if channel, ok := filters["channel"]; ok {
			if template.Channel != ChannelType(channel.(string)) {
				include = false
			}
		}

		if isDefault, ok := filters["is_default"]; ok {
			if template.IsDefault != isDefault.(bool) {
				include = false
			}
		}

		if include {
			result = append(result, template)
		}
	}

	return result
}

// ExportTemplate exports a template as JSON
func (tm *TemplateManager) ExportTemplate(id string) ([]byte, error) {
	template, exists := tm.templates[id]
	if !exists {
		return nil, fmt.Errorf("template %s not found", id)
	}

	return json.MarshalIndent(template, "", "  ")
}

// ImportTemplate imports a template from JSON
func (tm *TemplateManager) ImportTemplate(data []byte) error {
	var template AlertTemplate
	if err := json.Unmarshal(data, &template); err != nil {
		return fmt.Errorf("failed to parse template JSON: %w", err)
	}

	// Reset metadata for imported template
	template.ID = ""
	template.IsDefault = false
	template.UsageCount = 0

	return tm.AddTemplate(&template)
}

// Template helper functions
func formatTemplateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func formatTemplateDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1fm", d.Minutes())
	}
	return fmt.Sprintf("%.1fh", d.Hours())
}

func unixTime(t time.Time) int64 {
	return t.Unix()
}
