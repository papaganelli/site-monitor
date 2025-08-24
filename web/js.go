package web

// dashboardJS contains the JavaScript code
const dashboardJS = `
var SiteMonitorDashboard = function() {
    this.ws = null;
    this.charts = {};
    this.lastUpdate = null;
    this.reconnectAttempts = 0;
    this.maxReconnectAttempts = 5;
    this.lastActivityUpdate = null;
    
    this.init();
};

SiteMonitorDashboard.prototype.init = function() {
    var self = this;
    
    var checkChart = function() {
        if (typeof Chart !== 'undefined') {
            self.loadInitialData().then(function() {
                self.initWebSocket();
                self.initCharts();
                self.startPeriodicUpdates();
                self.hideLoadingOverlay();
            });
        } else {
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

// ✅ FONCTION CORRIGÉE: Graphique des temps de réponse avec meilleur groupage
SiteMonitorDashboard.prototype.initResponseTimeChart = function() {
    var ctx = document.getElementById('response-time-chart');
    if (!ctx) return;
    
    var self = this;
    
    // CORRECTION: Demander plus de données sur une période plus courte pour plus de densité
    fetch('/api/history?since=6h&limit=100')
        .then(function(r) { return r.json(); })
        .then(function(history) {
            if (history.length === 0) {
                self.showChartError('response-time-chart', 'No monitoring data available yet. Start monitoring to collect data.');
                return;
            }

            var validEntries = history.filter(function(entry) {
                var responseTime = entry.response_time_ms || entry.duration || entry.response_time_ns;
                return entry.success && responseTime && responseTime > 0;
            });
            
            if (validEntries.length === 0) {
                self.showChartError('response-time-chart', 'No successful requests with valid response times found.');
                return;
            }

            // CORRECTION: Utiliser la logique de groupage améliorée
            var groupedData = validEntries.length > 15 ? 
                self.groupDataByTimeInterval(validEntries) : 
                validEntries.sort(function(a, b) { return new Date(a.timestamp) - new Date(b.timestamp); });
            
            if (groupedData.length === 0) {
                self.showChartError('response-time-chart', 'Unable to process data for display.');
                return;
            }

            var labels = [];
            var dataPoints = [];
            
            groupedData.forEach(function(point) {
                var date = new Date(point.timestamp);
                var timeLabel = date.getHours().toString().padStart(2, '0') + ':' + 
                               date.getMinutes().toString().padStart(2, '0');
                
                var rawDuration = point.response_time_ms || point.duration || point.response_time_ns || 0;
                var responseTimeMs;
                
                if (rawDuration > 1000000) {
                    responseTimeMs = rawDuration / 1000000;
                } else if (rawDuration > 1000) {
                    responseTimeMs = rawDuration / 1000;
                } else {
                    responseTimeMs = rawDuration;
                }
                
                if (responseTimeMs > 0) {
                    labels.push(timeLabel);
                    dataPoints.push(Math.round(responseTimeMs * 100) / 100);
                }
            });
            
            if (dataPoints.length === 0) {
                self.showChartError('response-time-chart', 'No valid response time data after conversion.');
                return;
            }
            
            var minValue = Math.min.apply(Math, dataPoints);
            var maxValue = Math.max.apply(Math, dataPoints);
            var range = maxValue - minValue;
            var margin = Math.max(range * 0.2, Math.max(minValue * 0.1, 1));
            var yMin = Math.max(0, minValue - margin);
            var yMax = maxValue + margin;
            
            self.charts.responseTime = new Chart(ctx, {
                type: 'line',
                data: {
                    labels: labels,
                    datasets: [{
                        label: 'Response Time (ms)',
                        data: dataPoints,
                        borderColor: '#059669',
                        backgroundColor: 'rgba(5, 150, 105, 0.1)',
                        tension: 0.3,
                        fill: true,
                        pointRadius: 4,
                        pointHoverRadius: 7,
                        pointBackgroundColor: '#059669',
                        pointBorderColor: '#ffffff',
                        pointBorderWidth: 2,
                        borderWidth: 2,
                    }]
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
                                text: 'Time (24h format)'
                            }
                        },
                        y: {
                            display: true,
                            beginAtZero: false,
                            min: yMin,
                            max: yMax,
                            title: {
                                display: true,
                                text: 'Response Time (milliseconds)'
                            },
                            ticks: {
                                callback: function(value) {
                                    return Math.round(value * 100) / 100 + ' ms';
                                }
                            }
                        }
                    },
                    animation: {
                        duration: 750
                    }
                }
            });
        })
        .catch(function(error) {
            console.error('Failed to initialize response time chart:', error);
            self.showChartError('response-time-chart', 'Network error: ' + error.message);
        });
};

// ✅ FONCTION CORRIGÉE: Groupage amélioré des données
SiteMonitorDashboard.prototype.groupDataByTimeInterval = function(history) {
    if (history.length === 0) return [];
    
    // Trier par timestamp (plus récent en premier)
    history.sort(function(a, b) {
        return new Date(b.timestamp) - new Date(a.timestamp);
    });
    
    // CORRECTION: Prendre plus de données et réduire l'intervalle
    var recentHistory = history.slice(0, 50); // Plus de données
    
    // CORRECTION: Si peu de données, ne pas grouper du tout
    if (recentHistory.length <= 10) {
        return recentHistory.reverse(); // Chronologique
    }
    
    // CORRECTION: Intervalles plus petits (15 minutes au lieu de 30)
    var grouped = {};
    var interval = 15 * 60 * 1000; // 15 minutes
    
    recentHistory.forEach(function(entry) {
        var timestamp = new Date(entry.timestamp).getTime();
        var slotTime = Math.floor(timestamp / interval) * interval;
        var slotKey = slotTime.toString();
        
        // CORRECTION: Faire une moyenne au lieu de garder le dernier
        if (!grouped[slotKey]) {
            grouped[slotKey] = {
                timestamp: entry.timestamp,
                response_time_ms: entry.response_time_ms || entry.duration || entry.response_time_ns,
                site_name: entry.site_name,
                success: entry.success,
                count: 1
            };
        } else {
            // Moyenne pondérée des temps de réponse
            var existingTime = grouped[slotKey].response_time_ms || 0;
            var newTime = entry.response_time_ms || entry.duration || entry.response_time_ns || 0;
            var count = grouped[slotKey].count;
            
            grouped[slotKey].response_time_ms = ((existingTime * count) + newTime) / (count + 1);
            grouped[slotKey].count = count + 1;
            
            // Garder le timestamp le plus récent
            if (new Date(entry.timestamp) > new Date(grouped[slotKey].timestamp)) {
                grouped[slotKey].timestamp = entry.timestamp;
            }
        }
    });
    
    // Convertir en array et trier chronologiquement
    return Object.values(grouped).sort(function(a, b) {
        return new Date(a.timestamp) - new Date(b.timestamp);
    });
};

SiteMonitorDashboard.prototype.initUptimeChart = function() {
    var ctx = document.getElementById('uptime-chart');
    if (!ctx) return;
    
    var self = this;
    
    fetch('/api/stats')
        .then(function(r) { return r.json(); })
        .then(function(stats) {
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
    
    return {
        labels: ['Successful', 'Failed'],
        data: [successPercent.toFixed(1), failurePercent.toFixed(1)]
    };
};

// ✅ CORRECTION: Fonction updateCharts avec le nouveau groupage
SiteMonitorDashboard.prototype.updateCharts = function() {
    var self = this;
    
    if (this.charts.responseTime) {
        fetch('/api/history?since=6h&limit=100')
            .then(function(r) { return r.json(); })
            .then(function(history) {
                if (history.length === 0) return;

                var validEntries = history.filter(function(entry) {
                    var responseTime = entry.response_time_ms || entry.duration || entry.response_time_ns;
                    return entry.success && responseTime && responseTime > 0;
                });

                var groupedData = validEntries.length > 15 ? 
                    self.groupDataByTimeInterval(validEntries) : 
                    validEntries.sort(function(a, b) { return new Date(a.timestamp) - new Date(b.timestamp); });
                
                if (groupedData.length === 0) return;

                var labels = [];
                var dataPoints = [];
                
                groupedData.forEach(function(entry) {
                    var date = new Date(entry.timestamp);
                    var timeLabel = date.getHours().toString().padStart(2, '0') + ':' + 
                                   date.getMinutes().toString().padStart(2, '0');
                    labels.push(timeLabel);
                    
                    var rawDuration = entry.response_time_ms || entry.duration || entry.response_time_ns || 0;
                    var responseTimeMs;
                    
                    if (rawDuration > 1000000) {
                        responseTimeMs = rawDuration / 1000000;
                    } else if (rawDuration > 1000) {
                        responseTimeMs = rawDuration / 1000;
                    } else {
                        responseTimeMs = rawDuration;
                    }
                    
                    dataPoints.push(Math.round(responseTimeMs * 100) / 100);
                });
                
                if (dataPoints.length === 0) return;
                
                var minValue = Math.min.apply(Math, dataPoints);
                var maxValue = Math.max.apply(Math, dataPoints);
                var range = maxValue - minValue;
                var margin = Math.max(range * 0.2, Math.max(minValue * 0.1, 1));
                var yMin = Math.max(0, minValue - margin);
                var yMax = maxValue + margin;
                
                self.charts.responseTime.data.labels = labels;
                self.charts.responseTime.data.datasets[0].data = dataPoints;
                self.charts.responseTime.options.scales.y.min = yMin;
                self.charts.responseTime.options.scales.y.max = yMax;
                self.charts.responseTime.update('none');
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
                self.charts.uptime.update('none');
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
            '<div style="display: flex; align-items: center; justify-content: center; height: 250px; color: #dc2626; text-align: center; flex-direction: column; background: linear-gradient(145deg, #fef2f2 0%, #fee2e2 100%); border-radius: 12px; border: 2px dashed #dc2626; margin: 10px 0;">' +
                '<div style="font-size: 3rem; margin-bottom: 15px; opacity: 0.7;">📊</div>' +
                '<div style="font-weight: 600; font-size: 1.1em; margin-bottom: 8px; color: #991b1b;">Chart Data Unavailable</div>' +
                '<div style="font-size: 0.9em; color: #dc2626; max-width: 300px; line-height: 1.4;">' + message + '</div>' +
                '<div style="margin-top: 15px; padding: 8px 16px; background: #dc2626; color: white; border-radius: 6px; font-size: 0.8em; font-weight: 500;">💡 Tip: Let monitoring run for at least 5-10 minutes</div>' +
            '</div>';
    }
};

SiteMonitorDashboard.prototype.initWebSocket = function() {
    var protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    var wsUrl = protocol + '//' + window.location.host + '/ws';
    var self = this;
    
    this.ws = new WebSocket(wsUrl);
    
    this.ws.onopen = function() {
        self.updateConnectionStatus(true);
        self.reconnectAttempts = 0;
    };
    
    this.ws.onmessage = function(event) {
        var message = JSON.parse(event.data);
        self.handleWebSocketMessage(message);
    };
    
    this.ws.onclose = function() {
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
    
    var responseTime = entry.response_time_ms || entry.duration || entry.response_time_ns || 0;
    var responseTimeMs = responseTime > 1000000 ? responseTime / 1000000 : 
                        responseTime > 1000 ? responseTime / 1000 : responseTime;
    
    var details = 'Response time: ' + Math.round(responseTimeMs) + 'ms';
    if (!entry.success && (entry.error || entry.error_message)) {
        details = 'Error: ' + (entry.error || entry.error_message);
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
