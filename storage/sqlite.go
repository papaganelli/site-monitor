package storage

import (
	"database/sql"
	"fmt"
	"site-monitor/monitor"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// SQLiteStorage implements Storage interface using SQLite
type SQLiteStorage struct {
	db   *sql.DB
	mu   sync.RWMutex
	path string
}

// NewSQLiteStorage creates a new SQLite storage instance
func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	storage := &SQLiteStorage{
		db:   db,
		path: dbPath,
	}

	return storage, nil
}

// Init creates the necessary tables and indexes
func (s *SQLiteStorage) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create results table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		site_name TEXT NOT NULL,
		url TEXT NOT NULL,
		status_code INTEGER DEFAULT 0,
		response_time_ns INTEGER NOT NULL,
		success BOOLEAN NOT NULL,
		error_message TEXT DEFAULT '',
		timestamp DATETIME NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	if _, err := s.db.Exec(createTableSQL); err != nil {
		return fmt.Errorf("failed to create results table: %w", err)
	}

	// Create indexes for better query performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_site_timestamp ON results(site_name, timestamp DESC);",
		"CREATE INDEX IF NOT EXISTS idx_timestamp ON results(timestamp DESC);",
		"CREATE INDEX IF NOT EXISTS idx_site_success ON results(site_name, success);",
		"CREATE INDEX IF NOT EXISTS idx_success_timestamp ON results(success, timestamp DESC);",
	}

	for _, indexSQL := range indexes {
		if _, err := s.db.Exec(indexSQL); err != nil {
			return fmt.Errorf("failed to create index: %w", err)
		}
	}

	return nil
}

// SaveResult stores a monitoring result in the database
func (s *SQLiteStorage) SaveResult(result monitor.Result) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	insertSQL := `
	INSERT INTO results (site_name, url, status_code, response_time_ns, success, error_message, timestamp)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := s.db.Exec(
		insertSQL,
		result.Name,
		result.URL,
		result.Status,
		result.Duration.Nanoseconds(),
		result.Success,
		result.Error,
		result.Timestamp,
	)

	if err != nil {
		return fmt.Errorf("failed to save result: %w", err)
	}

	return nil
}

// GetHistory retrieves monitoring history for a specific site
func (s *SQLiteStorage) GetHistory(siteName string, since time.Time) ([]HistoryEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	querySQL := `
	SELECT id, site_name, url, status_code, response_time_ns, success, error_message, timestamp, created_at
	FROM results
	WHERE site_name = ? AND timestamp >= ?
	ORDER BY timestamp DESC`

	rows, err := s.db.Query(querySQL, siteName, since)
	if err != nil {
		return nil, fmt.Errorf("failed to query history: %w", err)
	}
	defer rows.Close()

	return s.scanHistoryEntries(rows)
}

// GetAllHistory retrieves monitoring history for all sites
func (s *SQLiteStorage) GetAllHistory(since time.Time) ([]HistoryEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	querySQL := `
	SELECT id, site_name, url, status_code, response_time_ns, success, error_message, timestamp, created_at
	FROM results
	WHERE timestamp >= ?
	ORDER BY timestamp DESC`

	rows, err := s.db.Query(querySQL, since)
	if err != nil {
		return nil, fmt.Errorf("failed to query all history: %w", err)
	}
	defer rows.Close()

	return s.scanHistoryEntries(rows)
}

// GetStats calculates statistics for a specific site
func (s *SQLiteStorage) GetStats(siteName string, since time.Time) (Stats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	statsSQL := `
	SELECT 
		COUNT(*) as total_checks,
		COUNT(CASE WHEN success = 1 THEN 1 END) as successful_checks,
		COUNT(CASE WHEN success = 0 THEN 1 END) as failed_checks,
		COALESCE(AVG(CASE WHEN success = 1 THEN response_time_ns END), 0) as avg_response_time_ns,
		COALESCE(MIN(CASE WHEN success = 1 THEN response_time_ns END), 0) as min_response_time_ns,
		COALESCE(MAX(CASE WHEN success = 1 THEN response_time_ns END), 0) as max_response_time_ns,
		MAX(timestamp) as last_check,
		MIN(timestamp) as first_check
	FROM results
	WHERE site_name = ? AND timestamp >= ?`

	var stats Stats
	var avgNs, minNs, maxNs float64
	var lastCheckStr, firstCheckStr sql.NullString // ← Utiliser sql.NullString pour gérer les timestamps SQLite

	row := s.db.QueryRow(statsSQL, siteName, since)
	err := row.Scan(
		&stats.TotalChecks,
		&stats.SuccessfulChecks,
		&stats.FailedChecks,
		&avgNs,
		&minNs,
		&maxNs,
		&lastCheckStr,  // ← Scanner en tant que string
		&firstCheckStr, // ← Scanner en tant que string
	)

	if err != nil {
		return Stats{}, fmt.Errorf("failed to get stats: %w", err)
	}

	stats.SiteName = siteName
	if stats.TotalChecks > 0 {
		stats.SuccessRate = float64(stats.SuccessfulChecks) / float64(stats.TotalChecks) * 100
	}

	// Convertir les float64 en time.Duration
	stats.AvgResponseTime = time.Duration(int64(avgNs))
	stats.MinResponseTime = time.Duration(int64(minNs))
	stats.MaxResponseTime = time.Duration(int64(maxNs))

	// Convertir les strings en time.Time
	if lastCheckStr.Valid {
		if parsedTime, err := time.Parse(time.RFC3339, lastCheckStr.String); err == nil {
			stats.LastCheck = parsedTime
		} else {
			// Essayer d'autres formats de date si RFC3339 échoue
			if parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", lastCheckStr.String); err == nil {
				stats.LastCheck = parsedTime
			} else if parsedTime, err := time.Parse("2006-01-02T15:04:05Z", lastCheckStr.String); err == nil {
				stats.LastCheck = parsedTime
			}
		}
	}

	if firstCheckStr.Valid {
		if parsedTime, err := time.Parse(time.RFC3339, firstCheckStr.String); err == nil {
			stats.FirstCheck = parsedTime
		} else {
			// Essayer d'autres formats de date si RFC3339 échoue
			if parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", firstCheckStr.String); err == nil {
				stats.FirstCheck = parsedTime
			} else if parsedTime, err := time.Parse("2006-01-02T15:04:05Z", firstCheckStr.String); err == nil {
				stats.FirstCheck = parsedTime
			}
		}
	}

	// Calculate uptime/downtime (simplified calculation)
	if !stats.FirstCheck.IsZero() && !stats.LastCheck.IsZero() {
		totalDuration := stats.LastCheck.Sub(stats.FirstCheck)
		stats.Uptime = time.Duration(float64(totalDuration) * stats.SuccessRate / 100)
		stats.Downtime = totalDuration - stats.Uptime
	}

	return stats, nil
}

// GetAllStats calculates statistics for all sites
func (s *SQLiteStorage) GetAllStats(since time.Time) (map[string]Stats, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// First, get all unique site names
	sitesSQL := "SELECT DISTINCT site_name FROM results WHERE timestamp >= ?"
	rows, err := s.db.Query(sitesSQL, since)
	if err != nil {
		return nil, fmt.Errorf("failed to get site names: %w", err)
	}
	defer rows.Close()

	var siteNames []string
	for rows.Next() {
		var siteName string
		if err := rows.Scan(&siteName); err != nil {
			return nil, fmt.Errorf("failed to scan site name: %w", err)
		}
		siteNames = append(siteNames, siteName)
	}

	// Get stats for each site
	allStats := make(map[string]Stats)
	for _, siteName := range siteNames {
		stats, err := s.GetStats(siteName, since)
		if err != nil {
			return nil, fmt.Errorf("failed to get stats for %s: %w", siteName, err)
		}
		allStats[siteName] = stats
	}

	return allStats, nil
}

// Close closes the database connection
func (s *SQLiteStorage) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// scanHistoryEntries is a helper function to scan database rows into HistoryEntry structs
func (s *SQLiteStorage) scanHistoryEntries(rows *sql.Rows) ([]HistoryEntry, error) {
	var entries []HistoryEntry

	for rows.Next() {
		var entry HistoryEntry
		var responseTimeNs int64
		var timestampStr, createdAtStr string

		err := rows.Scan(
			&entry.ID,
			&entry.SiteName,
			&entry.URL,
			&entry.Status,
			&responseTimeNs,
			&entry.Success,
			&entry.Error,
			&timestampStr,
			&createdAtStr,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		entry.Duration = time.Duration(responseTimeNs)

		// Convertir les timestamps strings en time.Time
		if parsedTime, err := time.Parse(time.RFC3339, timestampStr); err == nil {
			entry.Timestamp = parsedTime
		} else if parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", timestampStr); err == nil {
			entry.Timestamp = parsedTime
		} else if parsedTime, err := time.Parse("2006-01-02T15:04:05Z", timestampStr); err == nil {
			entry.Timestamp = parsedTime
		}

		if parsedTime, err := time.Parse(time.RFC3339, createdAtStr); err == nil {
			entry.CreatedAt = parsedTime
		} else if parsedTime, err := time.Parse("2006-01-02 15:04:05.999999999-07:00", createdAtStr); err == nil {
			entry.CreatedAt = parsedTime
		} else if parsedTime, err := time.Parse("2006-01-02T15:04:05Z", createdAtStr); err == nil {
			entry.CreatedAt = parsedTime
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return entries, nil
}
