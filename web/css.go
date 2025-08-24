package web

// dashboardCSS contains the CSS styles for the dashboard
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
