package web

// dashboardHTML contains the main dashboard HTML template
const dashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Site Monitor Dashboard</title>
    <link rel="stylesheet" href="/static/dashboard.css">
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
