package web

// dashboardHTML contains the main dashboard HTML template
const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Site Monitor Dashboard</title>
    <link rel="stylesheet" href="/static/dashboard.css">
    <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/4.4.0/chart.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/date-fns/2.30.0/index.min.js"></script>
</head>
<body>
    <div class="dashboard">
        <!-- Header -->
        <header class="header">
            <div class="header-content">
                <div class="logo">
                    <h1>🚀 Site Monitor</h1>
                    <span class="version">v0.5.0</span>
                </div>
                <div class="header-actions">
                    <div class="status-indicator" id="connection-status">
                        <div class="status-dot"></div>
                        <span>Connected</span>
                    </div>
                    <button class="btn-refresh" onclick="refreshData()">
                        <span class="refresh-icon">🔄</span>
                        Refresh
                    </button>
                </div>
            </div>
        </header>

        <!-- Main Content -->
        <main class="main-content">
            <!-- Overview Cards -->
            <section class="overview-section">
                <div class="overview-grid">
                    <div class="overview-card total-sites">
                        <div class="card-icon">🌐</div>
                        <div class="card-content">
                            <div class="card-value" id="total-sites">-</div>
                            <div class="card-label">Total Sites</div>
                        </div>
                    </div>
                    
                    <div class="overview-card healthy-sites">
                        <div class="card-icon">✅</div>
                        <div class="card-content">
                            <div class="card-value" id="healthy-sites">-</div>
                            <div class="card-label">Healthy Sites</div>
                        </div>
                    </div>
                    
                    <div class="overview-card uptime">
                        <div class="card-icon">📈</div>
                        <div class="card-content">
                            <div class="card-value" id="overall-uptime">-.-%</div>
                            <div class="card-label">Overall Uptime</div>
                        </div>
                    </div>
                    
                    <div class="overview-card total-checks">
                        <div class="card-icon">🔍</div>
                        <div class="card-content">
                            <div class="card-value" id="total-checks">-</div>
                            <div class="card-label">Total Checks</div>
                        </div>
                    </div>
                </div>
            </section>

            <!-- Sites Grid -->
            <section class="sites-section">
                <h2 class="section-title">Sites Status</h2>
                <div class="sites-grid" id="sites-grid">
                    <!-- Sites will be populated by JavaScript -->
                </div>
            </section>

            <!-- Charts Section -->
            <section class="charts-section">
                <div class="charts-grid">
                    <div class="chart-container">
                        <h3 class="chart-title">Response Time Trends (Last 24h)</h3>
                        <canvas id="response-time-chart"></canvas>
                    </div>
                    
                    <div class="chart-container">
                        <h3 class="chart-title">Uptime Distribution</h3>
                        <canvas id="uptime-chart"></canvas>
                    </div>
                </div>
            </section>

            <!-- Recent Activity -->
            <section class="activity-section">
                <h2 class="section-title">Recent Activity</h2>
                <div class="activity-container">
                    <div class="activity-list" id="activity-list">
                        <div class="activity-loading">Loading recent activity...</div>
                    </div>
                </div>
            </section>
        </main>
    </div>

    <!-- Loading Overlay -->
    <div class="loading-overlay" id="loading-overlay">
        <div class="loading-spinner"></div>
        <div class="loading-text">Loading Dashboard...</div>
    </div>

    <!-- Toast Notifications -->
    <div class="toast-container" id="toast-container"></div>

    <script src="/static/dashboard.js"></script>
</body>
</html>`

// dashboardCSS contains the modern CSS styles
const dashboardCSS = `/* Modern Dashboard Styles */
:root {
    --primary-color: #2563eb;
    --primary-dark: #1d4ed8;
    --secondary-color: #64748b;
    --success-color: #059669;
    --warning-color: #d97706;
    --error-color: #dc2626;
    --background-primary: #ffffff;
    --background-secondary: #f8fafc;
    --background-dark: #0f172a;
    --text-primary: #1e293b;
    --text-secondary: #64748b;
    --border-color: #e2e8f0;
    --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
    --shadow-md: 0 4px 6px -1px rgb(0 0 0 / 0.1);
    --shadow-lg: 0 10px 15px -3px rgb(0 0 0 / 0.1);
    --radius-sm: 0.375rem;
    --radius-md: 0.5rem;
    --radius-lg: 0.75rem;
}

@media (prefers-color-scheme: dark) {
    :root {
        --background-primary: #1e293b;
        --background-secondary: #0f172a;
        --text-primary: #f1f5f9;
        --text-secondary: #94a3b8;
        --border-color: #334155;
    }
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

body {
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: var(--background-secondary);
    color: var(--text-primary);
    line-height: 1.6;
    overflow-x: hidden;
}

/* Dashboard Layout */
.dashboard {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
}

/* Header */
.header {
    background: var(--background-primary);
    border-bottom: 1px solid var(--border-color);
    box-shadow: var(--shadow-sm);
    position: sticky;
    top: 0;
    z-index: 100;
}

.header-content {
    max-width: 1400px;
    margin: 0 auto;
    padding: 1rem 2rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.logo {
    display: flex;
    align-items: center;
    gap: 0.75rem;
}

.logo h1 {
    font-size: 1.5rem;
    font-weight: 700;
    background: linear-gradient(135deg, var(--primary-color), var(--primary-dark));
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
}

.version {
    background: var(--primary-color);
    color: white;
    padding: 0.25rem 0.5rem;
    border-radius: var(--radius-sm);
    font-size: 0.75rem;
    font-weight: 500;
}

.header-actions {
    display: flex;
    align-items: center;
    gap: 1rem;
}

.status-indicator {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--text-secondary);
    font-size: 0.875rem;
}

.status-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: var(--success-color);
    animation: pulse 2s infinite;
}

@keyframes pulse {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.5; }
}

.btn-refresh {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 1rem;
    background: var(--primary-color);
    color: white;
    border: none;
    border-radius: var(--radius-md);
    cursor: pointer;
    font-weight: 500;
    transition: all 0.2s;
}

.btn-refresh:hover {
    background: var(--primary-dark);
    transform: translateY(-1px);
}

.refresh-icon {
    transition: transform 0.3s;
}

.btn-refresh:active .refresh-icon {
    transform: rotate(180deg);
}

/* Main Content */
.main-content {
    flex: 1;
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
    width: 100%;
}

.section-title {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 1.5rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
}

/* Overview Cards */
.overview-section {
    margin-bottom: 3rem;
}

.overview-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
}

.overview-card {
    background: var(--background-primary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    padding: 1.5rem;
    display: flex;
    align-items: center;
    gap: 1rem;
    box-shadow: var(--shadow-sm);
    transition: all 0.2s;
    position: relative;
    overflow: hidden;
}

.overview-card:hover {
    box-shadow: var(--shadow-md);
    transform: translateY(-2px);
}

.overview-card::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 3px;
    background: var(--primary-color);
}

.overview-card.healthy-sites::before {
    background: var(--success-color);
}

.overview-card.uptime::before {
    background: linear-gradient(90deg, var(--success-color), var(--primary-color));
}

.card-icon {
    font-size: 2rem;
    opacity: 0.8;
}

.card-content {
    flex: 1;
}

.card-value {
    font-size: 2rem;
    font-weight: 700;
    color: var(--text-primary);
    line-height: 1;
}

.card-label {
    color: var(--text-secondary);
    font-size: 0.875rem;
    font-weight: 500;
    margin-top: 0.25rem;
}

/* Sites Grid */
.sites-section {
    margin-bottom: 3rem;
}

.sites-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
    gap: 1.5rem;
}

.site-card {
    background: var(--background-primary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    padding: 1.5rem;
    box-shadow: var(--shadow-sm);
    transition: all 0.2s;
    position: relative;
}

.site-card:hover {
    box-shadow: var(--shadow-md);
}

.site-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1rem;
}

.site-name {
    font-weight: 600;
    color: var(--text-primary);
    font-size: 1.1rem;
}

.site-status {
    padding: 0.25rem 0.75rem;
    border-radius: var(--radius-sm);
    font-size: 0.75rem;
    font-weight: 600;
    text-transform: uppercase;
}

.site-status.healthy {
    background: #dcfce7;
    color: var(--success-color);
}

.site-status.degraded {
    background: #fef3c7;
    color: var(--warning-color);
}

.site-status.down {
    background: #fecaca;
    color: var(--error-color);
}

.site-status.stale {
    background: #f1f5f9;
    color: var(--secondary-color);
}

.site-metrics {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
}

.metric {
    text-align: center;
}

.metric-value {
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-primary);
}

.metric-label {
    color: var(--text-secondary);
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-top: 0.25rem;
}

.site-chart {
    height: 40px;
    margin-top: 1rem;
    background: linear-gradient(90deg, var(--success-color), var(--primary-color));
    border-radius: var(--radius-sm);
    opacity: 0.1;
}

/* Charts */
.charts-section {
    margin-bottom: 3rem;
}

.charts-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
    gap: 2rem;
}

.chart-container {
    background: var(--background-primary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    padding: 1.5rem;
    box-shadow: var(--shadow-sm);
}

.chart-title {
    font-size: 1rem;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 1rem;
    text-align: center;
}

.chart-container canvas {
    max-height: 300px;
}

/* Activity Section */
.activity-section {
    margin-bottom: 2rem;
}

.activity-container {
    background: var(--background-primary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-lg);
    overflow: hidden;
    box-shadow: var(--shadow-sm);
}

.activity-list {
    max-height: 400px;
    overflow-y: auto;
}

.activity-item {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--border-color);
    transition: background 0.2s;
}

.activity-item:hover {
    background: var(--background-secondary);
}

.activity-item:last-child {
    border-bottom: none;
}

.activity-icon {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    flex-shrink: 0;
}

.activity-icon.success {
    background: var(--success-color);
}

.activity-icon.error {
    background: var(--error-color);
}

.activity-content {
    flex: 1;
}

.activity-message {
    color: var(--text-primary);
    font-weight: 500;
}

.activity-details {
    color: var(--text-secondary);
    font-size: 0.875rem;
    margin-top: 0.25rem;
}

.activity-time {
    color: var(--text-secondary);
    font-size: 0.75rem;
    flex-shrink: 0;
}

.activity-loading {
    padding: 2rem;
    text-align: center;
    color: var(--text-secondary);
}

/* Loading Overlay */
.loading-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(255, 255, 255, 0.9);
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    backdrop-filter: blur(4px);
}

.loading-spinner {
    width: 40px;
    height: 40px;
    border: 4px solid var(--border-color);
    border-left: 4px solid var(--primary-color);
    border-radius: 50%;
    animation: spin 1s linear infinite;
}

@keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
}

.loading-text {
    margin-top: 1rem;
    color: var(--text-secondary);
    font-weight: 500;
}

.loading-overlay.hidden {
    display: none;
}

/* Toast Notifications */
.toast-container {
    position: fixed;
    top: 1rem;
    right: 1rem;
    z-index: 1100;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
}

.toast {
    background: var(--background-primary);
    border: 1px solid var(--border-color);
    border-radius: var(--radius-md);
    padding: 1rem 1.5rem;
    box-shadow: var(--shadow-lg);
    display: flex;
    align-items: center;
    gap: 0.75rem;
    min-width: 300px;
    animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
    from {
        transform: translateX(100%);
        opacity: 0;
    }
    to {
        transform: translateX(0);
        opacity: 1;
    }
}

.toast.success {
    border-left: 4px solid var(--success-color);
}

.toast.error {
    border-left: 4px solid var(--error-color);
}

.toast.warning {
    border-left: 4px solid var(--warning-color);
}

.toast-icon {
    font-size: 1.25rem;
}

.toast-message {
    flex: 1;
    font-weight: 500;
}

.toast-close {
    background: none;
    border: none;
    color: var(--text-secondary);
    cursor: pointer;
    padding: 0;
    font-size: 1.25rem;
    line-height: 1;
}

.toast-close:hover {
    color: var(--text-primary);
}

/* Responsive Design */
@media (max-width: 768px) {
    .header-content {
        padding: 1rem;
        flex-direction: column;
        gap: 1rem;
    }
    
    .main-content {
        padding: 1rem;
    }
    
    .overview-grid {
        grid-template-columns: 1fr;
    }
    
    .sites-grid {
        grid-template-columns: 1fr;
    }
    
    .charts-grid {
        grid-template-columns: 1fr;
    }
    
    .site-metrics {
        grid-template-columns: 1fr;
        gap: 0.5rem;
    }
}

@media (max-width: 480px) {
    .logo h1 {
        font-size: 1.25rem;
    }
    
    .overview-card {
        padding: 1rem;
    }
    
    .card-value {
        font-size: 1.5rem;
    }
    
    .site-card {
        padding: 1rem;
    }
    
    .toast {
        min-width: 280px;
        margin: 0 1rem;
    }
}

/* Animations and Transitions */
.fade-in {
    animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
}

.scale-in {
    animation: scaleIn 0.3s ease-out;
}

@keyframes scaleIn {
    from { transform: scale(0.95); opacity: 0; }
    to { transform: scale(1); opacity: 1; }
}

/* Custom Scrollbar */
::-webkit-scrollbar {
    width: 6px;
    height: 6px;
}

::-webkit-scrollbar-track {
    background: var(--background-secondary);
}

::-webkit-scrollbar-thumb {
    background: var(--border-color);
    border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
    background: var(--text-secondary);
}`

// dashboardJS contains the JavaScript code for dashboard functionality
const dashboardJS = `
// Dashboard JavaScript
class SiteMonitorDashboard {
    constructor() {
        this.ws = null;
        this.charts = {};
        this.lastUpdate = null;
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        
        this.init();
    }
    
    async init() {
        await this.loadInitialData();
        this.initWebSocket();
        this.initCharts();
        this.startPeriodicUpdates();
        this.hideLoadingOverlay();
    }
    
    async loadInitialData() {
        try {
            const [overview, history] = await Promise.all([
                fetch('/api/overview').then(r => r.json()),
                fetch('/api/history?limit=100').then(r => r.json())
            ]);
            
            this.updateOverview(overview);
            this.updateSitesGrid(overview.sites);
            this.updateActivityFeed(history);
            
        } catch (error) {
            console.error('Failed to load initial data:', error);
            this.showToast('Failed to load dashboard data', 'error');
        }
    }
    
    initWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = protocol + '//' + window.location.host + '/ws';
        
        this.ws = new WebSocket(wsUrl);
        
        this.ws.onopen = () => {
            console.log('WebSocket connected');
            this.updateConnectionStatus(true);
            this.reconnectAttempts = 0;
        };
        
        this.ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            this.handleWebSocketMessage(message);
        };
        
        this.ws.onclose = () => {
            console.log('WebSocket disconnected');
            this.updateConnectionStatus(false);
            this.attemptReconnect();
        };
        
        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.updateConnectionStatus(false);
        };
    }
    
    attemptReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            const delay = Math.pow(2, this.reconnectAttempts) * 1000; // Exponential backoff
            
            console.log('Attempting to reconnect in ${delay}ms (attempt ${this.reconnectAttempts})');
            
            setTimeout(() => {
                this.initWebSocket();
            }, delay);
        } else {
            this.showToast('Connection lost. Please refresh the page.', 'error');
        }
    }
    
    handleWebSocketMessage(message) {
        switch (message.type) {
            case series_update':
                this.updateOverview(message.data);
                this.updateSitesGrid(message.data.sites);
                this.updateCharts();
                break;
            default:
                console.log('Unknown message type:', message.type);
        }
    }
    
    updateConnectionStatus(connected) {
        const statusIndicator = document.getElementById('connection-status');
        const statusDot = statusIndicator.querySelector('.status-dot');
        const statusText = statusIndicator.querySelector('span');
        
        if (connected) {
            statusDot.style.background = '#059669';
            statusText.textContent = 'Connected';
        } else {
            statusDot.style.background = '#dc2626';
            statusText.textContent = 'Disconnected';
        }
    }
    
    updateOverview(data) {
        document.getElementById('total-sites').textContent = data.total_sites;
        document.getElementById('healthy-sites').textContent = data.healthy_sites;
        document.getElementById('overall-uptime').textContent = data.overall_uptime.toFixed(1) + '%';
        document.getElementById('total-checks').textContent = this.formatNumber(data.total_checks);
        
        this.lastUpdate = new Date(data.last_update);
    }
    
    updateSitesGrid(sites) {
        const grid = document.getElementById('sites-grid');
        grid.innerHTML = '';
        
        sites.forEach(site => {
            const siteCard = this.createSiteCard(site);
            grid.appendChild(siteCard);
        });
    }
    
    createSiteCard(site) {
        const card = document.createElement('div');
        card.className = 'site-card fade-in';
        
        const lastCheckTime = this.formatTimeAgo(new Date(site.last_check));
        
        card.innerHTML = '
            <div class="site-header">
                <div class="site-name">${this.escapeHtml(site.name)}</div>
                <div class="site-status ${site.status}">${site.status}</div>
            </div>
            <div class="site-metrics">
                <div class="metric">
                    <div class="metric-value">${site.uptime.toFixed(1)}%</div>
                    <div class="metric-label">Uptime</div>
                </div>
                <div class="metric">
                    <div class="metric-value">${site.response_time_ms}ms</div>
                    <div class="metric-label">Response Time</div>
                </div>
            </div>
            <div class="site-footer">
                <small class="last-check">Last check: ${lastCheckTime}</small>
            </div>
            <div class="site-chart"></div>
        ';
        
        return card;
    }
    
    initCharts() {
        this.initResponseTimeChart();
        this.initUptimeChart();
    }
    
    async initResponseTimeChart() {
        const ctx = document.getElementById('response-time-chart');
        if (!ctx) return;
        
        try {
            const history = await fetch('/api/history?since=24h&limit=100').then(r => r.json());
            const chartData = this.processResponseTimeData(history);
            
            this.charts.responseTime = new Chart(ctx, {
                type: 'line',
                data: {
                    labels: chartData.labels,
                    datasets: chartData.datasets
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    interaction: {
                        intersect: false,
                        mode: 'index'
                    },
                    plugins: {
                        legend: {
                            display: true,
                            position: 'bottom'
                        }
                    },
                    scales: {
                        x: {
                            type: 'time',
                            time: {
                                displayFormats: {
                                    hour: 'HH:mm',
                                    minute: 'HH:mm'
                                }
                            },
                            title: {
                                display: true,
                                text: 'Time'
                            }
                        },
                        y: {
                            beginAtZero: true,
                            title: {
                                display: true,
                                text: 'Response Time (ms)'
                            }
                        }
                    }
                }
            });
        } catch (error) {
            console.error('Failed to initialize response time chart:', error);
        }
    }
    
    async initUptimeChart() {
        const ctx = document.getElementById('uptime-chart');
        if (!ctx) return;
        
        try {
            const stats = await fetch('/api/stats').then(r => r.json());
            const chartData = this.processUptimeData(stats);
            
            this.charts.uptime = new Chart(ctx, {
                type: 'doughnut',
                data: {
                    labels: chartData.labels,
                    datasets: [{
                        data: chartData.data,
                        backgroundColor: [
                            '#059669', // Success
                            '#dc2626', // Error
                            '#d97706'  // Warning
                        ],
                        borderWidth: 0,
                        hoverOffset: 4
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: {
                            display: true,
                            position: 'bottom'
                        }
                    }
                }
            });
        } catch (error) {
            console.error('Failed to initialize uptime chart:', error);
        }
    }
    
    processResponseTimeData(history) {
        const siteData = {};
        const colors = ['#2563eb', '#059669', '#d97706', '#dc2626', '#7c3aed'];
        
        history.forEach(entry => {
            if (!siteData[entry.site_name]) {
                siteData[entry.site_name] = [];
            }
            siteData[entry.site_name].push({
                x: entry.timestamp,
                y: entry.duration / 1000000 // Convert from nanoseconds to milliseconds
            });
        });
        
        const datasets = Object.keys(siteData).map((siteName, index) => ({
            label: siteName,
            data: siteData[siteName].slice(-50), // Last 50 points
            borderColor: colors[index % colors.length],
            backgroundColor: colors[index % colors.length] + '20',
            tension: 0.4,
            fill: false
        }));
        
        return {
            labels: [], // Chart.js will handle time labels automatically
            datasets
        };
    }
    
    processUptimeData(stats) {
        let totalSuccess = 0;
        let totalFailure = 0;
        
        Object.values(stats).forEach(site => {
            totalSuccess += site.successful_checks;
            totalFailure += site.failed_checks;
        });
        
        const total = totalSuccess + totalFailure;
        
        return {
            labels: ['Successful', 'Failed'],
            data: [
                total > 0 ? (totalSuccess / total * 100).toFixed(1) : 0,
                total > 0 ? (totalFailure / total * 100).toFixed(1) : 0
            ]
        };
    }
    
    async updateCharts() {
        try {
            // Update response time chart
            if (this.charts.responseTime) {
                const history = await fetch('/api/history?since=24h&limit=100').then(r => r.json());
                const chartData = this.processResponseTimeData(history);
                this.charts.responseTime.data.datasets = chartData.datasets;
                this.charts.responseTime.update('none');
            }
            
            // Update uptime chart
            if (this.charts.uptime) {
                const stats = await fetch('/api/stats').then(r => r.json());
                const chartData = this.processUptimeData(stats);
                this.charts.uptime.data.datasets[0].data = chartData.data;
                this.charts.uptime.update('none');
            }
        } catch (error) {
            console.error('Failed to update charts:', error);
        }
    }
    
    async updateActivityFeed(history) {
        const activityList = document.getElementById('activity-list');
        activityList.innerHTML = '';
        
        // Show last 20 activities
        const recentHistory = history.slice(0, 20);
        
        if (recentHistory.length === 0) {
            activityList.innerHTML = '<div class="activity-loading">No recent activity found</div>';
            return;
        }
        
        recentHistory.forEach(entry => {
            const activityItem = this.createActivityItem(entry);
            activityList.appendChild(activityItem);
        });
    }
    
    createActivityItem(entry) {
        const item = document.createElement('div');
        item.className = 'activity-item';
        
        const iconClass = entry.success ? 'success' : 'error';
        const statusText = entry.success ? 'UP' : 'DOWN';
        const timeAgo = this.formatTimeAgo(new Date(entry.timestamp));
        
        let details = 'Response time: ${Math.round(entry.duration / 1000000)}ms';
        if (!entry.success && entry.error) {
            details = 'Error: ${entry.error}';
        }
        
        item.innerHTML = '
            <div class="activity-icon ${iconClass}"></div>
            <div class="activity-content">
                <div class="activity-message">${this.escapeHtml(entry.site_name)} is ${statusText}</div>
                <div class="activity-details">${this.escapeHtml(details)}</div>
            </div>
            <div class="activity-time">${timeAgo}</div>
        ';
        
        return item;
    }
    
    startPeriodicUpdates() {
        // Update data every 30 seconds
        setInterval(async () => {
            try {
                const overview = await fetch('/api/overview').then(r => r.json());
                this.updateOverview(overview);
                this.updateSitesGrid(overview.sites);
                
                // Update activity feed every minute
                if (!this.lastActivityUpdate || Date.now() - this.lastActivityUpdate > 60000) {
                    const history = await fetch('/api/history?limit=20').then(r => r.json());
                    this.updateActivityFeed(history);
                    this.lastActivityUpdate = Date.now();
                }
                
            } catch (error) {
                console.error('Failed to update dashboard:', error);
            }
        }, 30000);
    }
    
    hideLoadingOverlay() {
        const overlay = document.getElementById('loading-overlay');
        overlay.classList.add('hidden');
    }
    
    showToast(message, type = 'info') {
        const container = document.getElementById('toast-container');
        const toast = document.createElement('div');
        toast.className = 'toast ${type}';
        
        const icons = {
            success: '✅',
            error: '❌',
            warning: '⚠️',
            info: 'ℹ️'
        };
        
        toast.innerHTML = '
            <div class="toast-icon">${icons[type] || icons.info}</div>
            <div class="toast-message">${this.escapeHtml(message)}</div>
            <button class="toast-close" onclick="this.parentElement.remove()">×</button>
        ';
        
        container.appendChild(toast);
        
        // Auto-remove after 5 seconds
        setTimeout(() => {
            toast.remove();
        }, 5000);
    }
    
    formatNumber(num) {
        if (num >= 1000000) {
            return (num / 1000000).toFixed(1) + 'M';
        }
        if (num >= 1000) {
            return (num / 1000).toFixed(1) + 'K';
        }
        return num.toString();
    }
    
    formatTimeAgo(date) {
        const now = new Date();
        const diff = now - date;
        const minutes = Math.floor(diff / 60000);
        const hours = Math.floor(minutes / 60);
        const days = Math.floor(hours / 24);
        
        if (days > 0) return days + 'd ago';
        if (hours > 0) return hours + 'h ago';
        if (minutes > 0) return minutes + 'm ago';
        return 'Just now';
    }
    
    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// Global functions
function refreshData() {
    window.dashboard.loadInitialData();
    window.dashboard.updateCharts();
    
    // Visual feedback
    const refreshBtn = document.querySelector('.btn-refresh');
    const icon = refreshBtn.querySelector('.refresh-icon');
    icon.style.transform = 'rotate(360deg)';
    setTimeout(() => {
        icon.style.transform = 'rotate(0deg)';
    }, 500);
}

// Initialize dashboard when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.dashboard = new SiteMonitorDashboard();
});
`
