package metrics

import (
	"fmt"
	"math"
	"site-monitor/storage"
	"sort"
	"strings"
	"time"
)

// AdvancedMetrics represents advanced statistical metrics
type AdvancedMetrics struct {
	SiteName string `json:"site_name"`
	Period   string `json:"period"`

	// Response Time Percentiles
	P50  time.Duration `json:"p50_response_time"` // Median
	P90  time.Duration `json:"p90_response_time"`
	P95  time.Duration `json:"p95_response_time"`
	P99  time.Duration `json:"p99_response_time"`
	P999 time.Duration `json:"p999_response_time"`

	// Reliability Metrics
	MTTR time.Duration `json:"mttr"` // Mean Time To Recovery
	MTBF time.Duration `json:"mtbf"` // Mean Time Between Failures

	// Availability Metrics
	UptimePercent     float64 `json:"uptime_percent"`
	DowntimePercent   float64 `json:"downtime_percent"`
	AvailabilityNines int     `json:"availability_nines"` // 99.9% = 3 nines

	// Error Analysis
	ErrorRate      float64                    `json:"error_rate_percent"`
	ErrorBreakdown map[string]ErrorStatistics `json:"error_breakdown"`

	// Performance Trends
	ResponseTimeStdDev time.Duration  `json:"response_time_std_dev"`
	ResponseTimeTrend  TrendDirection `json:"response_time_trend"`
	UptimeTrend        TrendDirection `json:"uptime_trend"`

	// SLA Compliance
	SLACompliance map[string]SLAResult `json:"sla_compliance"`

	// Time-based Analysis
	PeakHours     []HourlyStats `json:"peak_hours"`
	WeeklyPattern WeeklyPattern `json:"weekly_pattern"`

	// Calculated Fields
	TotalChecks       int64     `json:"total_checks"`
	SuccessfulChecks  int64     `json:"successful_checks"`
	FailedChecks      int64     `json:"failed_checks"`
	FirstCheck        time.Time `json:"first_check"`
	LastCheck         time.Time `json:"last_check"`
	AnalysisTimestamp time.Time `json:"analysis_timestamp"`
}

// ErrorStatistics represents error analysis for a specific error type
type ErrorStatistics struct {
	Count      int64     `json:"count"`
	Percentage float64   `json:"percentage"`
	FirstSeen  time.Time `json:"first_seen"`
	LastSeen   time.Time `json:"last_seen"`
	Pattern    string    `json:"pattern"`
}

// TrendDirection represents the direction of a trend
type TrendDirection string

const (
	TrendImproving TrendDirection = "improving"
	TrendStable    TrendDirection = "stable"
	TrendDegrading TrendDirection = "degrading"
	TrendUnknown   TrendDirection = "unknown"
)

// SLAResult represents SLA compliance result
type SLAResult struct {
	Target    float64       `json:"target_percent"`
	Actual    float64       `json:"actual_percent"`
	Compliant bool          `json:"compliant"`
	Violation time.Duration `json:"violation_duration"`
}

// HourlyStats represents statistics for a specific hour of the day
type HourlyStats struct {
	Hour            int           `json:"hour"`
	AvgResponseTime time.Duration `json:"avg_response_time"`
	SuccessRate     float64       `json:"success_rate"`
	CheckCount      int64         `json:"check_count"`
}

// WeeklyPattern represents weekly performance patterns
type WeeklyPattern struct {
	MondayStats    DayStats `json:"monday"`
	TuesdayStats   DayStats `json:"tuesday"`
	WednesdayStats DayStats `json:"wednesday"`
	ThursdayStats  DayStats `json:"thursday"`
	FridayStats    DayStats `json:"friday"`
	SaturdayStats  DayStats `json:"saturday"`
	SundayStats    DayStats `json:"sunday"`
	BestDay        string   `json:"best_day"`
	WorstDay       string   `json:"worst_day"`
}

// DayStats represents statistics for a specific day of the week
type DayStats struct {
	AvgResponseTime time.Duration `json:"avg_response_time"`
	SuccessRate     float64       `json:"success_rate"`
	CheckCount      int64         `json:"check_count"`
}

// AdvancedMetricsCalculator calculates advanced metrics
type AdvancedMetricsCalculator struct {
	storage storage.Storage
}

// NewAdvancedMetricsCalculator creates a new calculator
func NewAdvancedMetricsCalculator(storage storage.Storage) *AdvancedMetricsCalculator {
	return &AdvancedMetricsCalculator{
		storage: storage,
	}
}

// CalculateAdvancedMetrics calculates comprehensive metrics for a site
func (calc *AdvancedMetricsCalculator) CalculateAdvancedMetrics(siteName string, since time.Time, period string) (*AdvancedMetrics, error) {
	// Get raw historical data
	history, err := calc.storage.GetHistory(siteName, since)
	if err != nil {
		return nil, err
	}

	if len(history) == 0 {
		return nil, fmt.Errorf("no data available for site %s since %s", siteName, since.Format("2006-01-02"))
	}

	metrics := &AdvancedMetrics{
		SiteName:          siteName,
		Period:            period,
		AnalysisTimestamp: time.Now(),
		FirstCheck:        history[len(history)-1].Timestamp, // Last item is earliest
		LastCheck:         history[0].Timestamp,              // First item is most recent
		TotalChecks:       int64(len(history)),
	}

	// Calculate basic counts
	var successfulChecks, failedChecks int64
	var responseTimes []time.Duration
	var errorCounts map[string]int64 = make(map[string]int64)
	var downtimeEvents []DowntimeEvent

	for _, entry := range history {
		if entry.Success {
			successfulChecks++
			responseTimes = append(responseTimes, entry.Duration)
		} else {
			failedChecks++
			// Count errors by message
			errorMsg := entry.Error
			if errorMsg == "" {
				errorMsg = "Unknown Error"
			}
			errorCounts[errorMsg]++
		}
	}

	metrics.SuccessfulChecks = successfulChecks
	metrics.FailedChecks = failedChecks

	// Calculate uptime/downtime percentages
	if metrics.TotalChecks > 0 {
		metrics.UptimePercent = float64(successfulChecks) / float64(metrics.TotalChecks) * 100
		metrics.DowntimePercent = float64(failedChecks) / float64(metrics.TotalChecks) * 100
		metrics.ErrorRate = metrics.DowntimePercent
	}

	// Calculate availability nines
	metrics.AvailabilityNines = calc.calculateNines(metrics.UptimePercent)

	// Calculate response time percentiles
	if len(responseTimes) > 0 {
		sort.Slice(responseTimes, func(i, j int) bool {
			return responseTimes[i] < responseTimes[j]
		})

		metrics.P50 = calc.percentile(responseTimes, 50)
		metrics.P90 = calc.percentile(responseTimes, 90)
		metrics.P95 = calc.percentile(responseTimes, 95)
		metrics.P99 = calc.percentile(responseTimes, 99)
		metrics.P999 = calc.percentile(responseTimes, 99.9)

		// Calculate standard deviation
		metrics.ResponseTimeStdDev = calc.standardDeviation(responseTimes)
	}

	// Calculate MTTR and MTBF
	downtimeEvents = calc.identifyDowntimeEvents(history)
	if len(downtimeEvents) > 0 {
		metrics.MTTR = calc.calculateMTTR(downtimeEvents)
		metrics.MTBF = calc.calculateMTBF(downtimeEvents, time.Since(metrics.FirstCheck))
	}

	// Calculate trends
	metrics.ResponseTimeTrend = calc.calculateResponseTimeTrend(history)
	metrics.UptimeTrend = calc.calculateUptimeTrend(history)

	// Calculate error breakdown
	metrics.ErrorBreakdown = calc.calculateErrorBreakdown(errorCounts, history, failedChecks)

	// Calculate SLA compliance
	metrics.SLACompliance = calc.calculateSLACompliance(metrics.UptimePercent)

	// Calculate hourly patterns
	metrics.PeakHours = calc.calculateHourlyPattern(history)

	// Calculate weekly patterns
	metrics.WeeklyPattern = calc.calculateWeeklyPattern(history)

	return metrics, nil
}

// percentile calculates the nth percentile from sorted duration slice
func (calc *AdvancedMetricsCalculator) percentile(sortedDurations []time.Duration, percentile float64) time.Duration {
	if len(sortedDurations) == 0 {
		return 0
	}

	index := percentile / 100.0 * float64(len(sortedDurations)-1)
	lower := int(math.Floor(index))
	upper := int(math.Ceil(index))

	if lower == upper {
		return sortedDurations[lower]
	}

	// Linear interpolation
	weight := index - float64(lower)
	lowerVal := float64(sortedDurations[lower])
	upperVal := float64(sortedDurations[upper])

	return time.Duration(lowerVal + weight*(upperVal-lowerVal))
}

// standardDeviation calculates standard deviation of response times
func (calc *AdvancedMetricsCalculator) standardDeviation(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	// Calculate mean
	var sum float64
	for _, d := range durations {
		sum += float64(d)
	}
	mean := sum / float64(len(durations))

	// Calculate variance
	var variance float64
	for _, d := range durations {
		variance += math.Pow(float64(d)-mean, 2)
	}
	variance /= float64(len(durations))

	return time.Duration(math.Sqrt(variance))
}

// calculateNines determines availability nines (99%, 99.9%, 99.99%, etc.)
func (calc *AdvancedMetricsCalculator) calculateNines(uptimePercent float64) int {
	if uptimePercent >= 99.999 {
		return 5
	} else if uptimePercent >= 99.99 {
		return 4
	} else if uptimePercent >= 99.9 {
		return 3
	} else if uptimePercent >= 99.0 {
		return 2
	} else if uptimePercent >= 90.0 {
		return 1
	}
	return 0
}

// DowntimeEvent represents a period of downtime
type DowntimeEvent struct {
	StartTime time.Time
	EndTime   time.Time
	Duration  time.Duration
}

// identifyDowntimeEvents identifies periods of continuous downtime
func (calc *AdvancedMetricsCalculator) identifyDowntimeEvents(history []storage.HistoryEntry) []DowntimeEvent {
	var events []DowntimeEvent
	var currentEvent *DowntimeEvent

	// Process in reverse order (oldest first)
	for i := len(history) - 1; i >= 0; i-- {
		entry := history[i]

		if !entry.Success {
			// Start or continue downtime event
			if currentEvent == nil {
				currentEvent = &DowntimeEvent{
					StartTime: entry.Timestamp,
					EndTime:   entry.Timestamp,
				}
			}
			currentEvent.EndTime = entry.Timestamp
		} else {
			// End downtime event if one was active
			if currentEvent != nil {
				currentEvent.Duration = currentEvent.EndTime.Sub(currentEvent.StartTime)
				events = append(events, *currentEvent)
				currentEvent = nil
			}
		}
	}

	// Handle ongoing downtime
	if currentEvent != nil {
		currentEvent.Duration = time.Since(currentEvent.StartTime)
		events = append(events, *currentEvent)
	}

	return events
}

// calculateMTTR calculates Mean Time To Recovery
func (calc *AdvancedMetricsCalculator) calculateMTTR(events []DowntimeEvent) time.Duration {
	if len(events) == 0 {
		return 0
	}

	var totalDuration time.Duration
	for _, event := range events {
		totalDuration += event.Duration
	}

	return totalDuration / time.Duration(len(events))
}

// calculateMTBF calculates Mean Time Between Failures
func (calc *AdvancedMetricsCalculator) calculateMTBF(events []DowntimeEvent, totalPeriod time.Duration) time.Duration {
	if len(events) <= 1 {
		return totalPeriod
	}

	return totalPeriod / time.Duration(len(events))
}

// calculateResponseTimeTrend analyzes response time trends
func (calc *AdvancedMetricsCalculator) calculateResponseTimeTrend(history []storage.HistoryEntry) TrendDirection {
	if len(history) < 10 {
		return TrendUnknown
	}

	// Split history into two halves and compare averages
	mid := len(history) / 2
	recent := history[:mid] // More recent data
	older := history[mid:]  // Older data

	recentAvg := calc.averageResponseTime(recent)
	olderAvg := calc.averageResponseTime(older)

	// Calculate percentage change
	if olderAvg == 0 {
		return TrendUnknown
	}

	change := (float64(recentAvg) - float64(olderAvg)) / float64(olderAvg) * 100

	if change > 10 {
		return TrendDegrading
	} else if change < -10 {
		return TrendImproving
	}
	return TrendStable
}

// calculateUptimeTrend analyzes uptime trends
func (calc *AdvancedMetricsCalculator) calculateUptimeTrend(history []storage.HistoryEntry) TrendDirection {
	if len(history) < 20 {
		return TrendUnknown
	}

	// Split into quarters for trend analysis
	quarter := len(history) / 4
	q1 := history[:quarter]              // Most recent
	q4 := history[3*quarter : 4*quarter] // Oldest

	q1Success := calc.successRate(q1)
	q4Success := calc.successRate(q4)

	diff := q1Success - q4Success

	if diff > 5 {
		return TrendImproving
	} else if diff < -5 {
		return TrendDegrading
	}
	return TrendStable
}

// averageResponseTime calculates average response time for successful checks
func (calc *AdvancedMetricsCalculator) averageResponseTime(entries []storage.HistoryEntry) time.Duration {
	var total time.Duration
	var count int

	for _, entry := range entries {
		if entry.Success {
			total += entry.Duration
			count++
		}
	}

	if count == 0 {
		return 0
	}

	return total / time.Duration(count)
}

// successRate calculates success rate for a slice of entries
func (calc *AdvancedMetricsCalculator) successRate(entries []storage.HistoryEntry) float64 {
	if len(entries) == 0 {
		return 0
	}

	var successCount int
	for _, entry := range entries {
		if entry.Success {
			successCount++
		}
	}

	return float64(successCount) / float64(len(entries)) * 100
}

// calculateErrorBreakdown analyzes error patterns
func (calc *AdvancedMetricsCalculator) calculateErrorBreakdown(errorCounts map[string]int64, history []storage.HistoryEntry, totalFailed int64) map[string]ErrorStatistics {
	breakdown := make(map[string]ErrorStatistics)

	// Track first and last occurrence of each error
	errorFirstSeen := make(map[string]time.Time)
	errorLastSeen := make(map[string]time.Time)

	for _, entry := range history {
		if !entry.Success {
			errorMsg := entry.Error
			if errorMsg == "" {
				errorMsg = "Unknown Error"
			}

			// Update first seen (processing in reverse chronological order)
			errorFirstSeen[errorMsg] = entry.Timestamp

			// Update last seen if not set
			if _, exists := errorLastSeen[errorMsg]; !exists {
				errorLastSeen[errorMsg] = entry.Timestamp
			}
		}
	}

	for errorMsg, count := range errorCounts {
		percentage := float64(count) / float64(totalFailed) * 100

		breakdown[errorMsg] = ErrorStatistics{
			Count:      count,
			Percentage: percentage,
			FirstSeen:  errorFirstSeen[errorMsg],
			LastSeen:   errorLastSeen[errorMsg],
			Pattern:    calc.classifyErrorPattern(errorMsg),
		}
	}

	return breakdown
}

// classifyErrorPattern attempts to classify error types
func (calc *AdvancedMetricsCalculator) classifyErrorPattern(errorMsg string) string {
	errorMsg = strings.ToLower(errorMsg)

	if strings.Contains(errorMsg, "timeout") || strings.Contains(errorMsg, "deadline") {
		return "Timeout"
	} else if strings.Contains(errorMsg, "connection") || strings.Contains(errorMsg, "network") {
		return "Network"
	} else if strings.Contains(errorMsg, "dns") || strings.Contains(errorMsg, "resolve") {
		return "DNS"
	} else if strings.Contains(errorMsg, "ssl") || strings.Contains(errorMsg, "tls") || strings.Contains(errorMsg, "certificate") {
		return "SSL/TLS"
	} else if strings.Contains(errorMsg, "refused") {
		return "Connection Refused"
	} else if strings.Contains(errorMsg, "5") && strings.Contains(errorMsg, "0") {
		return "Server Error (5xx)"
	} else if strings.Contains(errorMsg, "4") && strings.Contains(errorMsg, "0") {
		return "Client Error (4xx)"
	}

	return "Other"
}

// calculateSLACompliance calculates compliance against standard SLA targets
func (calc *AdvancedMetricsCalculator) calculateSLACompliance(actualUptime float64) map[string]SLAResult {
	slaTargets := map[string]float64{
		"99.9% (8.77h downtime/month)":    99.9,
		"99.95% (4.38h downtime/month)":   99.95,
		"99.99% (52.6min downtime/month)": 99.99,
		"99.5% (3.65d downtime/month)":    99.5,
		"95% (36.5h downtime/month)":      95.0,
	}

	compliance := make(map[string]SLAResult)

	for name, target := range slaTargets {
		violation := time.Duration(0)
		if actualUptime < target {
			// Rough estimation of violation time based on a 30-day month
			violationPercent := target - actualUptime
			violation = time.Duration(violationPercent / 100 * float64(30*24*time.Hour))
		}

		compliance[name] = SLAResult{
			Target:    target,
			Actual:    actualUptime,
			Compliant: actualUptime >= target,
			Violation: violation,
		}
	}

	return compliance
}

// calculateHourlyPattern analyzes performance patterns by hour of day
func (calc *AdvancedMetricsCalculator) calculateHourlyPattern(history []storage.HistoryEntry) []HourlyStats {
	hourlyData := make(map[int][]storage.HistoryEntry)

	// Group entries by hour
	for _, entry := range history {
		hour := entry.Timestamp.Hour()
		hourlyData[hour] = append(hourlyData[hour], entry)
	}

	var hourlyStats []HourlyStats

	for hour := 0; hour < 24; hour++ {
		entries := hourlyData[hour]
		if len(entries) == 0 {
			continue
		}

		stats := HourlyStats{
			Hour:       hour,
			CheckCount: int64(len(entries)),
		}

		var totalResponseTime time.Duration
		var successCount int

		for _, entry := range entries {
			if entry.Success {
				successCount++
				totalResponseTime += entry.Duration
			}
		}

		if successCount > 0 {
			stats.AvgResponseTime = totalResponseTime / time.Duration(successCount)
		}
		stats.SuccessRate = float64(successCount) / float64(len(entries)) * 100

		hourlyStats = append(hourlyStats, stats)
	}

	// Sort by hour
	sort.Slice(hourlyStats, func(i, j int) bool {
		return hourlyStats[i].Hour < hourlyStats[j].Hour
	})

	return hourlyStats
}

// calculateWeeklyPattern analyzes performance patterns by day of week
func (calc *AdvancedMetricsCalculator) calculateWeeklyPattern(history []storage.HistoryEntry) WeeklyPattern {
	weeklyData := make(map[time.Weekday][]storage.HistoryEntry)

	// Group entries by day of week
	for _, entry := range history {
		weekday := entry.Timestamp.Weekday()
		weeklyData[weekday] = append(weeklyData[weekday], entry)
	}

	// Calculate stats for each day
	dayStats := make(map[time.Weekday]DayStats)

	for weekday, entries := range weeklyData {
		if len(entries) == 0 {
			continue
		}

		var totalResponseTime time.Duration
		var successCount int

		for _, entry := range entries {
			if entry.Success {
				successCount++
				totalResponseTime += entry.Duration
			}
		}

		stats := DayStats{
			CheckCount:  int64(len(entries)),
			SuccessRate: float64(successCount) / float64(len(entries)) * 100,
		}

		if successCount > 0 {
			stats.AvgResponseTime = totalResponseTime / time.Duration(successCount)
		}

		dayStats[weekday] = stats
	}

	// Find best and worst days
	var bestDay, worstDay string
	var bestRate, worstRate float64 = -1, 101

	dayNames := map[time.Weekday]string{
		time.Monday:    "Monday",
		time.Tuesday:   "Tuesday",
		time.Wednesday: "Wednesday",
		time.Thursday:  "Thursday",
		time.Friday:    "Friday",
		time.Saturday:  "Saturday",
		time.Sunday:    "Sunday",
	}

	for weekday, stats := range dayStats {
		dayName := dayNames[weekday]
		if stats.SuccessRate > bestRate {
			bestRate = stats.SuccessRate
			bestDay = dayName
		}
		if stats.SuccessRate < worstRate {
			worstRate = stats.SuccessRate
			worstDay = dayName
		}
	}

	return WeeklyPattern{
		MondayStats:    dayStats[time.Monday],
		TuesdayStats:   dayStats[time.Tuesday],
		WednesdayStats: dayStats[time.Wednesday],
		ThursdayStats:  dayStats[time.Thursday],
		FridayStats:    dayStats[time.Friday],
		SaturdayStats:  dayStats[time.Saturday],
		SundayStats:    dayStats[time.Sunday],
		BestDay:        bestDay,
		WorstDay:       worstDay,
	}
}

// String returns a formatted summary of advanced metrics
func (m *AdvancedMetrics) String() string {
	return fmt.Sprintf(
		"📊 Advanced Metrics for %s (%s)\n"+
			"🎯 Uptime: %.2f%% (%d nines)\n"+
			"⚡ Response Times: P95=%v, P99=%v\n"+
			"🔧 MTTR: %v, MTBF: %v\n"+
			"📈 Trends: Response=%s, Uptime=%s",
		m.SiteName, m.Period,
		m.UptimePercent, m.AvailabilityNines,
		m.P95, m.P99,
		m.MTTR, m.MTBF,
		m.ResponseTimeTrend, m.UptimeTrend,
	)
}
