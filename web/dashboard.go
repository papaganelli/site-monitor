package web

// dashboardHTML contains the main dashboard HTML template
const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Site Monitor Dashboard</title>
    <link rel="stylesheet" href="/static/dashboard.css">
    <!-- Chart.js v3.9.1 - Version stable et compatible -->
    <script src="https://cdn.jsdelivr.net/npm/chart.js@3.9.1/dist/chart.min.js"></script>
</head>
<body>
    <div class="dashboard">
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

        <main class="main-content">
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

            <section class="sites-section">
                <h2 class="section-title">Sites Status</h2>
                <div class="sites-grid" id="sites-grid">
                </div>
            </section>

            <section class="charts-section">
                <div class="charts-grid">
                    <div class="chart-container">
                        <h3 class="chart-title">Response Time Trends (Last 24h)</h3>
                        <canvas id="response-time-chart" width="400" height="200"></canvas>
                    </div>
                    
                    <div class="chart-container">
                        <h3 class="chart-title">Uptime Distribution</h3>
                        <canvas id="uptime-chart" width="400" height="200"></canvas>
                    </div>
                </div>
            </section>

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

    <div class="loading-overlay" id="loading-overlay">
        <div class="loading-spinner"></div>
        <div class="loading-text">Loading Dashboard...</div>
    </div>

    <div class="toast-container" id="toast-container"></div>

    <script src="/static/dashboard.js"></script>
</body>
</html>`

// dashboardCSS contains the CSS styles
const dashboardCSS = `:root {
    --primary-color: #2563eb;
    --primary-dark: #1d4ed8;
    --secondary-color: #64748b;
    --success-color: #059669;
    --warning-color: #d97706;
    --error-color: #dc2626;
    --background-primary: #ffffff;
    --background-secondary: #f8fafc;
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
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    background: var(--background-secondary);
    color: var(--text-primary);
    line-height: 1.6;
    overflow-x: hidden;
}

.dashboard {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
}

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
}

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

.fade-in {
    animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
}`

// dashboardJS contains the JavaScript code for dashboard functionality
const dashboardJS = `
var SiteMonitorDashboard = function() {
    this.ws = null;
    this.charts = {};
    this.lastUpdate = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    
    this.init();
};

SiteMonitorDashboard.prototype.init = function() {
    var self = this;
    
    // Attendre que Chart.js soit chargé
    var checkChart = function() {
        if (typeof Chart !== 'undefined') {
            console.log('Chart.js loaded successfully, version:', Chart.version);
            self.loadInitialData().then(function() {
                self.initWebSocket();
                self.initCharts();
                self.startPeriodicUpdates();
                self.hideLoadingOverlay();
            });
        } else {
            console.log('Waiting for Chart.js to load...');
            setTimeout(checkChart, 100);
        }
    };
    checkChart();
};

SiteMonitorDashboard.prototype.loadInitialData = function() {
    var self = this;
    return Promise.all([
        fetch('/api/overview').then(function(r) { return r.json(); }),
        fetch('/api/history?limit=100').then(function(r) { return r.json(); })
    ]).then(function(results) {
        self.updateOverview(results[0]);
        self.updateSitesGrid(results[0].sites);
        self.updateActivityFeed(results[1]);
    }).catch(function(error) {
        console.error('Failed to load initial data:', error);
        self.showToast('Failed to load dashboard data', 'error');
    });
};

SiteMonitorDashboard.prototype.initCharts = function() {
    this.initResponseTimeChart();
    this.initUptimeChart();
};

SiteMonitorDashboard.prototype.initResponseTimeChart = function() {
    var ctx = document.getElementById('response-time-chart');
    if (!ctx) {
        console.error('Response time chart canvas not found');
        return;
    }
    
    var self = this;
    
    fetch('/api/history?since=1h&limit=20')
        .then(function(r) { return r.json(); })
        .then(function(history) {
            console.log('Raw history data:', history);
            
            if (history.length === 0) {
                self.showChartError('response-time-chart', 'No data available. Let monitoring run for a few minutes.');
                return;
            }

            var labels = [];
            var dataPoints = [];
            
            var recentHistory = history.slice(0, 10).reverse();
            
            recentHistory.forEach(function(entry, index) {
                console.log('Processing entry:', entry);
                
                var date = new Date(entry.timestamp);
                var timeLabel = date.getHours().toString().padStart(2, '0') + ':' + 
                               date.getMinutes().toString().padStart(2, '0');
                labels.push(timeLabel);
                
                var responseTimeMs = entry.duration / 1000000;
                console.log('Response time for', timeLabel, ':', responseTimeMs, 'ms');
                dataPoints.push(Math.round(responseTimeMs * 100) / 100);
            });
            
            console.log('Chart labels:', labels);
            console.log('Chart data points:', dataPoints);
            
            // Calculer l'échelle dynamique
            var minValue = Math.min.apply(Math, dataPoints);
            var maxValue = Math.max.apply(Math, dataPoints);
            var range = maxValue - minValue;
            
            // Ajouter une marge de 20%
            var margin = range * 0.2;
            var yMin = Math.max(0, minValue - margin);
            var yMax = maxValue + margin;
            
            console.log('Y scale:', yMin, 'to', yMax);
            
            self.charts.responseTime = new Chart(ctx, {
                type: 'line',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Response Time (ms)',
                        data: dataPoints,
                        borderColor: '#059669',
                        backgroundColor: 'rgba(5, 150, 105, 0.1)',
                        tension: 0.4,
                        fill: true,
                        pointRadius: 4,
                        pointHoverRadius: 7,
                        pointBackgroundColor: '#059669',
                        pointBorderColor: '#ffffff',
                        pointBorderWidth: 2,
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: {
                            display: true,
                            position: 'bottom'
                        },
                        tooltip: {
                            callbacks: {
                                label: function(context) {
                                    return 'Response Time: ' + context.raw + ' ms';
                                }
                            }
                        }
                    },
                    scales: {
                        x: {
                            display: true,
                            title: {
                                display: true,
                                text: 'Time'
                            }
                        },
                        y: {
                            display: true,
                            beginAtZero: false,
                            min: yMin,
                            max: yMax,
                            title: {
                                display: true,
                                text: 'Response Time (ms)'
                            },
                            ticks: {
                                callback: function(value) {
                                    return value.toFixed(1) + ' ms';
                                }
                            }
                        }
                    }
                }
            });
            
            console.log('Response time chart created successfully!');
        })
        .catch(function(error) {
            console.error('Failed to initialize response time chart:', error);
            self.showChartError('response-time-chart', 'Error: ' + error.message);
        });
};

SiteMonitorDashboard.prototype.initUptimeChart = function() {
    var ctx = document.getElementById('uptime-chart');
    if (!ctx) {
        console.error('Uptime chart canvas not found');
        return;
    }
    
    var self = this;
    
    fetch('/api/stats')
        .then(function(r) { return r.json(); })
        .then(function(stats) {
            console.log('Loaded stats data for uptime chart:', Object.keys(stats).length, 'sites');
            var chartData = self.processUptimeData(stats);
            
            self.charts.uptime = new Chart(ctx, {
                type: 'doughnut',
                data: {
                    labels: chartData.labels,
                    datasets: [{
                        data: chartData.data,
                        backgroundColor: [
                            '#059669',
                            '#dc2626'
                        ],
                        borderWidth: 2,
                        borderColor: '#ffffff',
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
                        },
                        tooltip: {
                            callbacks: {
                                label: function(context) {
                                    var label = context.label || '';
                                    if (label) {
                                        label += ': ';
                                    }
                                    label += context.parsed + '%';
                                    return label;
                                }
                            }
                        }
                    }
                }
            });
            
            console.log('Uptime chart created successfully');
        })
        .catch(function(error) {
            console.error('Failed to initialize uptime chart:', error);
            self.showChartError('uptime-chart', 'Failed to load uptime data');
        });
};

SiteMonitorDashboard.prototype.processUptimeData = function(stats) {
    var totalSuccess = 0;
    var totalFailure = 0;
    
    for (var key in stats) {
        if (stats.hasOwnProperty(key)) {
            var site = stats[key];
            totalSuccess += site.successful_checks;
            totalFailure += site.failed_checks;
        }
    }
    
    var total = totalSuccess + totalFailure;
    
    if (total === 0) {
        return {
            labels: ['No Data'],
            data: [100]
        };
    }
    
    var successPercent = (totalSuccess / total * 100);
    var failurePercent = (totalFailure / total * 100);
    
    console.log('Uptime data:', successPercent.toFixed(1) + '% success,', failurePercent.toFixed(1) + '% failure');
    
    return {
        labels: ['Successful', 'Failed'],
        data: [successPercent.toFixed(1), failurePercent.toFixed(1)]
    };
};

SiteMonitorDashboard.prototype.updateCharts = function() {
    var self = this;
    
    console.log('Updating charts...');
    
    if (this.charts.responseTime) {
        fetch('/api/history?since=1h&limit=20')
            .then(function(r) { return r.json(); })
            .then(function(history) {
                if (history.length === 0) return;

                var labels = [];
                var dataPoints = [];
                
                var recentHistory = history.slice(0, 10).reverse();
                
                recentHistory.forEach(function(entry) {
                    var date = new Date(entry.timestamp);
                    var timeLabel = date.getHours().toString().padStart(2, '0') + ':' + 
                                   date.getMinutes().toString().padStart(2, '0');
                    labels.push(timeLabel);
                    
                    var responseTimeMs = entry.duration / 1000000;
                    dataPoints.push(Math.round(responseTimeMs * 100) / 100);
                });
                
                // Recalculer l'échelle
                var minValue = Math.min.apply(Math, dataPoints);
                var maxValue = Math.max.apply(Math, dataPoints);
                var range = maxValue - minValue;
                var margin = range * 0.2;
                var yMin = Math.max(0, minValue - margin);
                var yMax = maxValue + margin;
                
                self.charts.responseTime.data.labels = labels;
                self.charts.responseTime.data.datasets[0].data = dataPoints;
                self.charts.responseTime.options.scales.y.min = yMin;
                self.charts.responseTime.options.scales.y.max = yMax;
                self.charts.responseTime.update();
                
                console.log('Response time chart updated with', dataPoints.length, 'points');
            })
            .catch(function(error) {
                console.error('Failed to update response time chart:', error);
            });
    }
    
    if (this.charts.uptime) {
        fetch('/api/stats')
            .then(function(r) { return r.json(); })
            .then(function(stats) {
                var chartData = self.processUptimeData(stats);
                self.charts.uptime.data.labels = chartData.labels;
                self.charts.uptime.data.datasets[0].data = chartData.data;
                self.charts.uptime.update();
                console.log('Uptime chart updated');
            })
            .catch(function(error) {
                console.error('Failed to update uptime chart:', error);
            });
    }
};

SiteMonitorDashboard.prototype.showChartError = function(canvasId, message) {
    var canvas = document.getElementById(canvasId);
    if (canvas) {
        var container = canvas.parentNode;
        var title = container.querySelector('.chart-title');
        var titleText = title ? title.textContent : 'Chart';
        
        container.innerHTML = 
            '<h3 class="chart-title">' + titleText + '</h3>' +
            '<div style="display: flex; align-items: center; justify-content: center; height: 200px; color: #dc2626; text-align: center; flex-direction: column; background: #fef2f2; border-radius: 8px; border: 1px dashed #dc2626;">' +
                '<div style="font-size: 2rem; margin-bottom: 10px;">📊</div>' +
                '<div style="font-weight: 500; margin-bottom: 5px;">Chart Error</div>' +
                '<div style="font-size: 0.9em; color: #991b1b;">' + message + '</div>' +
            '</div>';
    }
};

SiteMonitorDashboard.prototype.initWebSocket = function() {
    var protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    var wsUrl = protocol + '//' + window.location.host + '/ws';
    var self = this;
    
    this.ws = new WebSocket(wsUrl);
    
    this.ws.onopen = function() {
        console.log('WebSocket connected');
        self.updateConnectionStatus(true);
        self.reconnectAttempts = 0;
    };
    
    this.ws.onmessage = function(event) {
        var message = JSON.parse(event.data);
        self.handleWebSocketMessage(message);
    };
    
    this.ws.onclose = function() {
        console.log('WebSocket disconnected');
        self.updateConnectionStatus(false);
        self.attemptReconnect();
    };
    
    this.ws.onerror = function(error) {
        console.error('WebSocket error:', error);
        self.updateConnectionStatus(false);
    };
};

SiteMonitorDashboard.prototype.handleWebSocketMessage = function(message) {
    switch (message.type) {
        case 'overview_update':
            this.updateOverview(message.data);
            this.updateSitesGrid(message.data.sites);
            this.updateCharts();
            break;
        default:
            console.log('Unknown message type:', message.type);
    }
};

SiteMonitorDashboard.prototype.updateConnectionStatus = function(connected) {
    var statusIndicator = document.getElementById('connection-status');
    var statusDot = statusIndicator.querySelector('.status-dot');
    var statusText = statusIndicator.querySelector('span');
    
    if (connected) {
        statusDot.style.background = '#059669';
        statusText.textContent = 'Connected';
    } else {
        statusDot.style.background = '#dc2626';
        statusText.textContent = 'Disconnected';
    }
};

SiteMonitorDashboard.prototype.updateOverview = function(data) {
    document.getElementById('total-sites').textContent = data.total_sites;
    document.getElementById('healthy-sites').textContent = data.healthy_sites;
    document.getElementById('overall-uptime').textContent = data.overall_uptime.toFixed(1) + '%';
    document.getElementById('total-checks').textContent = this.formatNumber(data.total_checks);
    
    this.lastUpdate = new Date(data.last_update);
};

SiteMonitorDashboard.prototype.updateSitesGrid = function(sites) {
    var grid = document.getElementById('sites-grid');
    grid.innerHTML = '';
    var self = this;
    
    sites.forEach(function(site) {
        var siteCard = self.createSiteCard(site);
        grid.appendChild(siteCard);
    });
};

SiteMonitorDashboard.prototype.createSiteCard = function(site) {
    var card = document.createElement('div');
    card.className = 'site-card fade-in';
    
    var lastCheckTime = this.formatTimeAgo(new Date(site.last_check));
    
    card.innerHTML = 
        '<div class="site-header">' +
            '<div class="site-name">' + this.escapeHtml(site.name) + '</div>' +
            '<div class="site-status ' + site.status + '">' + site.status + '</div>' +
        '</div>' +
        '<div class="site-metrics">' +
            '<div class="metric">' +
                '<div class="metric-value">' + site.uptime.toFixed(1) + '%</div>' +
                '<div class="metric-label">Uptime</div>' +
            '</div>' +
            '<div class="metric">' +
                '<div class="metric-value">' + site.response_time_ms + 'ms</div>' +
                '<div class="metric-label">Response Time</div>' +
            '</div>' +
        '</div>';
    
    return card;
};

SiteMonitorDashboard.prototype.updateActivityFeed = function(history) {
    var activityList = document.getElementById('activity-list');
    activityList.innerHTML = '';
    var self = this;
    
    var recentHistory = history.slice(0, 20);
    
    if (recentHistory.length === 0) {
        activityList.innerHTML = '<div class="activity-loading">No recent activity found</div>';
        return;
    }
    
    recentHistory.forEach(function(entry) {
        var activityItem = self.createActivityItem(entry);
        activityList.appendChild(activityItem);
    });
};

SiteMonitorDashboard.prototype.createActivityItem = function(entry) {
    var item = document.createElement('div');
    item.className = 'activity-item';
    
    var iconClass = entry.success ? 'success' : 'error';
    var statusText = entry.success ? 'UP' : 'DOWN';
    var timeAgo = this.formatTimeAgo(new Date(entry.timestamp));
    
    var details = 'Response time: ' + Math.round(entry.duration / 1000000) + 'ms';
    if (!entry.success && entry.error) {
        details = 'Error: ' + entry.error;
    }
    
    item.innerHTML = 
        '<div class="activity-icon ' + iconClass + '"></div>' +
        '<div class="activity-content">' +
            '<div class="activity-message">' + this.escapeHtml(entry.site_name) + ' is ' + statusText + '</div>' +
            '<div class="activity-details">' + this.escapeHtml(details) + '</div>' +
        '</div>' +
        '<div class="activity-time">' + timeAgo + '</div>';
    
    return item;
};

SiteMonitorDashboard.prototype.startPeriodicUpdates = function() {
    var self = this;
    setInterval(function() {
        fetch('/api/overview')
            .then(function(r) { return r.json(); })
            .then(function(overview) {
                self.updateOverview(overview);
                self.updateSitesGrid(overview.sites);
                
                if (!self.lastActivityUpdate || Date.now() - self.lastActivityUpdate > 60000) {
                    fetch('/api/history?limit=20')
                        .then(function(r) { return r.json(); })
                        .then(function(history) {
                            self.updateActivityFeed(history);
                            self.lastActivityUpdate = Date.now();
                        });
                }
            })
            .catch(function(error) {
                console.error('Failed to update dashboard:', error);
            });
    }, 30000);
};

SiteMonitorDashboard.prototype.hideLoadingOverlay = function() {
    var overlay = document.getElementById('loading-overlay');
    overlay.classList.add('hidden');
};

SiteMonitorDashboard.prototype.showToast = function(message, type) {
    type = type || 'info';
    var container = document.getElementById('toast-container');
    var toast = document.createElement('div');
    toast.className = 'toast ' + type;
    
    var icons = {
        success: '✅',
        error: '❌',
        warning: '⚠️',
        info: 'ℹ️'
    };
    
    toast.innerHTML = 
        '<div class="toast-icon">' + (icons[type] || icons.info) + '</div>' +
        '<div class="toast-message">' + this.escapeHtml(message) + '</div>' +
        '<button class="toast-close" onclick="this.parentElement.remove()">×</button>';
    
    container.appendChild(toast);
    
    setTimeout(function() {
        if (toast.parentNode) {
            toast.remove();
        }
    }, 5000);
};

SiteMonitorDashboard.prototype.formatNumber = function(num) {
    if (num >= 1000000) {
        return (num / 1000000).toFixed(1) + 'M';
    }
    if (num >= 1000) {
        return (num / 1000).toFixed(1) + 'K';
    }
    return num.toString();
};

SiteMonitorDashboard.prototype.formatTimeAgo = function(date) {
    var now = new Date();
    var diff = now - date;
    var minutes = Math.floor(diff / 60000);
    var hours = Math.floor(minutes / 60);
    var days = Math.floor(hours / 24);
    
    if (days > 0) return days + 'd ago';
    if (hours > 0) return hours + 'h ago';
    if (minutes > 0) return minutes + 'm ago';
    return 'Just now';
};

SiteMonitorDashboard.prototype.escapeHtml = function(text) {
    var div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
};

SiteMonitorDashboard.prototype.attemptReconnect = function() {
    var self = this;
    if (this.reconnectAttempts < this.maxReconnectAttempts) {
        this.reconnectAttempts++;
        var delay = Math.pow(2, this.reconnectAttempts) * 1000;
        
        console.log('Attempting to reconnect in ' + delay + 'ms (attempt ' + this.reconnectAttempts + ')');
        
        setTimeout(function() {
            self.initWebSocket();
        }, delay);
    } else {
        this.showToast('Connection lost. Please refresh the page.', 'error');
    }
};

function refreshData() {
    if (window.dashboard) {
        window.dashboard.loadInitialData();
        window.dashboard.updateCharts();
        
        var refreshBtn = document.querySelector('.btn-refresh');
        var icon = refreshBtn.querySelector('.refresh-icon');
        icon.style.transform = 'rotate(360deg)';
        setTimeout(function() {
            icon.style.transform = 'rotate(0deg)';
        }, 500);
    }
}

document.addEventListener('DOMContentLoaded', function() {
    window.dashboard = new SiteMonitorDashboard();
});`
