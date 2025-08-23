package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"site-monitor/config"
	"site-monitor/storage"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Dashboard represents the web dashboard server
type Dashboard struct {
	storage  storage.Storage
	config   *config.Config
	server   *http.Server
	clients  map[*websocket.Conn]bool
	upgrader websocket.Upgrader
}

// NewDashboard creates a new dashboard instance
func NewDashboard(storage storage.Storage, config *config.Config, port int) *Dashboard {
	dashboard := &Dashboard{
		storage: storage,
		config:  config,
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}

	// Create router
	router := mux.NewRouter()

	// Serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.HandlerFunc(dashboard.serveStatic)))

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/stats", dashboard.apiStats).Methods("GET")
	api.HandleFunc("/history", dashboard.apiHistory).Methods("GET")
	api.HandleFunc("/sites", dashboard.apiSites).Methods("GET")
	api.HandleFunc("/alerts", dashboard.apiAlerts).Methods("GET")
	api.HandleFunc("/overview", dashboard.apiOverview).Methods("GET")

	// WebSocket endpoint
	router.HandleFunc("/ws", dashboard.handleWebSocket)

	// Main dashboard page
	router.HandleFunc("/", dashboard.serveDashboard).Methods("GET")

	// Create HTTP server
	dashboard.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return dashboard
}

// Start starts the dashboard server
func (d *Dashboard) Start() error {
	log.Printf("🌐 Starting dashboard server on http://localhost%s", d.server.Addr)
	return d.server.ListenAndServe()
}

// Stop stops the dashboard server
func (d *Dashboard) Stop() error {
	return d.server.Close()
}

// serveDashboard serves the main dashboard HTML page
func (d *Dashboard) serveDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashboardHTML))
}

// serveStatic serves static files (CSS, JS, etc.)
func (d *Dashboard) serveStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/dashboard.css":
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte(dashboardCSS))
	case path == "/dashboard.js":
		w.Header().Set("Content-Type", "application/javascript")
		w.Write([]byte(dashboardJS))
	default:
		http.NotFound(w, r)
	}
}

// API Handlers

// apiOverview returns overall system statistics
func (d *Dashboard) apiOverview(w http.ResponseWriter, r *http.Request) {
	since := time.Now().Add(-24 * time.Hour) // Last 24 hours by default

	// Parse since parameter if provided
	if sinceParam := r.URL.Query().Get("since"); sinceParam != "" {
		if parsed, err := time.ParseDuration(sinceParam); err == nil {
			since = time.Now().Add(-parsed)
		}
	}

	allStats, err := d.storage.GetAllStats(since)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Calculate overview metrics
	overview := OverviewResponse{
		TotalSites:       len(allStats),
		HealthySites:     0,
		TotalChecks:      0,
		SuccessfulChecks: 0,
		OverallUptime:    0,
		Sites:            make([]SiteOverview, 0),
		LastUpdate:       time.Now(),
	}

	for _, stats := range allStats {
		overview.TotalChecks += stats.TotalChecks
		overview.SuccessfulChecks += stats.SuccessfulChecks

		if stats.SuccessRate >= 99.0 {
			overview.HealthySites++
		}

		siteStatus := "healthy"
		if stats.SuccessRate < 80.0 {
			siteStatus = "down"
		} else if stats.SuccessRate < 99.0 {
			siteStatus = "degraded"
		}

		// Check if stale (no recent checks)
		if time.Since(stats.LastCheck) > 10*time.Minute {
			siteStatus = "stale"
		}

		overview.Sites = append(overview.Sites, SiteOverview{
			Name:         stats.SiteName,
			Status:       siteStatus,
			Uptime:       stats.SuccessRate,
			ResponseTime: stats.AvgResponseTime.Milliseconds(),
			LastCheck:    stats.LastCheck,
			TotalChecks:  stats.TotalChecks,
		})
	}

	// Calculate overall uptime
	if overview.TotalChecks > 0 {
		overview.OverallUptime = float64(overview.SuccessfulChecks) / float64(overview.TotalChecks) * 100
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(overview)
}

// apiStats returns statistics for sites
func (d *Dashboard) apiStats(w http.ResponseWriter, r *http.Request) {
	since := time.Now().Add(-24 * time.Hour)

	if sinceParam := r.URL.Query().Get("since"); sinceParam != "" {
		if parsed, err := time.ParseDuration(sinceParam); err == nil {
			since = time.Now().Add(-parsed)
		}
	}

	siteName := r.URL.Query().Get("site")

	if siteName != "" {
		// Get stats for specific site
		stats, err := d.storage.GetStats(siteName, since)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	} else {
		// Get stats for all sites
		allStats, err := d.storage.GetAllStats(since)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(allStats)
	}
}

// apiHistory returns monitoring history
func (d *Dashboard) apiHistory(w http.ResponseWriter, r *http.Request) {
	since := time.Now().Add(-24 * time.Hour)
	limit := 1000 // Default limit

	if sinceParam := r.URL.Query().Get("since"); sinceParam != "" {
		if parsed, err := time.ParseDuration(sinceParam); err == nil {
			since = time.Now().Add(-parsed)
		}
	}

	if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
		if parsed, err := strconv.Atoi(limitParam); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	siteName := r.URL.Query().Get("site")

	var history []storage.HistoryEntry
	var err error

	if siteName != "" {
		history, err = d.storage.GetHistory(siteName, since)
	} else {
		history, err = d.storage.GetAllHistory(since)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Apply limit
	if len(history) > limit {
		history = history[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// apiSites returns list of monitored sites
func (d *Dashboard) apiSites(w http.ResponseWriter, r *http.Request) {
	sites := make([]SiteInfo, len(d.config.Sites))
	for i, site := range d.config.Sites {
		sites[i] = SiteInfo{
			Name:     site.Name,
			URL:      site.URL,
			Interval: site.Interval,
			Timeout:  site.Timeout,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}

// apiAlerts returns alert configuration status
func (d *Dashboard) apiAlerts(w http.ResponseWriter, r *http.Request) {
	alertStatus := AlertStatus{
		EmailEnabled:   false,
		WebhookEnabled: false,
		TotalChannels:  0,
	}

	if d.config.Alerts != nil {
		alertStatus.EmailEnabled = d.config.Alerts.Email.Enabled
		alertStatus.WebhookEnabled = d.config.Alerts.Webhook.Enabled

		if alertStatus.EmailEnabled {
			alertStatus.TotalChannels++
		}
		if alertStatus.WebhookEnabled {
			alertStatus.TotalChannels++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(alertStatus)
}

// WebSocket handler for real-time updates
func (d *Dashboard) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := d.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Register client
	d.clients[conn] = true
	log.Printf("📡 WebSocket client connected (total: %d)", len(d.clients))

	// Send initial data
	d.sendOverviewUpdate(conn)

	// Handle client messages (keep connection alive)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			delete(d.clients, conn)
			log.Printf("📡 WebSocket client disconnected (remaining: %d)", len(d.clients))
			break
		}
	}
}

// sendOverviewUpdate sends overview data to a WebSocket client
func (d *Dashboard) sendOverviewUpdate(conn *websocket.Conn) {
	since := time.Now().Add(-24 * time.Hour)
	allStats, err := d.storage.GetAllStats(since)
	if err != nil {
		return
	}

	// Build overview
	overview := OverviewResponse{
		TotalSites:       len(allStats),
		HealthySites:     0,
		TotalChecks:      0,
		SuccessfulChecks: 0,
		Sites:            make([]SiteOverview, 0),
		LastUpdate:       time.Now(),
	}

	for _, stats := range allStats {
		overview.TotalChecks += stats.TotalChecks
		overview.SuccessfulChecks += stats.SuccessfulChecks

		if stats.SuccessRate >= 99.0 {
			overview.HealthySites++
		}

		siteStatus := "healthy"
		if stats.SuccessRate < 80.0 {
			siteStatus = "down"
		} else if stats.SuccessRate < 99.0 {
			siteStatus = "degraded"
		}

		overview.Sites = append(overview.Sites, SiteOverview{
			Name:         stats.SiteName,
			Status:       siteStatus,
			Uptime:       stats.SuccessRate,
			ResponseTime: stats.AvgResponseTime.Milliseconds(),
			LastCheck:    stats.LastCheck,
			TotalChecks:  stats.TotalChecks,
		})
	}

	if overview.TotalChecks > 0 {
		overview.OverallUptime = float64(overview.SuccessfulChecks) / float64(overview.TotalChecks) * 100
	}

	// Send update
	message := map[string]interface{}{
		"type": "overview_update",
		"data": overview,
	}

	conn.WriteJSON(message)
}

// BroadcastUpdate sends updates to all connected WebSocket clients
func (d *Dashboard) BroadcastUpdate() {
	for conn := range d.clients {
		d.sendOverviewUpdate(conn)
	}
}
