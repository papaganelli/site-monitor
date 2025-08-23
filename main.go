package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"site-monitor/config"
	"site-monitor/monitor"
	"site-monitor/storage"
	"sync"
	"syscall"
)

func main() {
	// Load configuration from file
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize storage
	db, err := storage.NewSQLiteStorage("site-monitor.db")
	if err != nil {
		log.Fatal("Failed to initialize storage:", err)
	}
	defer db.Close()

	// Initialize database tables
	if err := db.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	fmt.Printf("🚀 Starting monitoring for %d sites\n", len(cfg.Sites))
	fmt.Printf("💾 Database initialized: site-monitor.db\n")

	var wg sync.WaitGroup

	// Channel to handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start monitoring each site in a separate goroutine
	for _, site := range cfg.Sites {
		wg.Add(1)

		go func(s config.Site) {
			defer wg.Done()

			interval, err := s.GetInterval()
			if err != nil {
				log.Printf("Invalid interval for %s: %v", s.Name, err)
				return
			}

			timeout, err := s.GetTimeout()
			if err != nil {
				log.Printf("Invalid timeout for %s: %v", s.Name, err)
				return
			}

			m := monitor.New(s.URL, interval)
			m.SetName(s.Name)
			m.SetTimeout(timeout)
			m.SetStorage(db) // Attach storage to monitor

			fmt.Printf("📍 Starting %s (%s) - checking every %s\n",
				s.Name, s.URL, s.Interval)

			if err := m.Start(); err != nil {
				log.Printf("Error monitoring %s: %v", s.Name, err)
			}
		}(site)
	}

	// Handle graceful shutdown
	go func() {
		<-sigChan
		fmt.Println("\n🛑 Received shutdown signal, stopping monitors...")
		// Note: In a real implementation, we would send a stop signal to monitors
		// For now, the program will exit and goroutines will be terminated
		os.Exit(0)
	}()

	wg.Wait()
}
