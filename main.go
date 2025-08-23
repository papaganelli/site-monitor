package main

import (
	"fmt"
	"log"
	"site-monitor/config"
	"site-monitor/monitor"
	"sync"
)

func main() {
	// Load configuration from file
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	fmt.Printf("🚀 Starting monitoring for %d sites\n", len(cfg.Sites))

	var wg sync.WaitGroup

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

			fmt.Printf("📍 Starting %s (%s) - checking every %s\n",
				s.Name, s.URL, s.Interval)

			if err := m.Start(); err != nil {
				log.Printf("Error monitoring %s: %v", s.Name, err)
			}
		}(site)
	}

	wg.Wait()
}
