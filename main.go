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
        log.Fatal("Impossible de charger la configuration:", err)
    }

    fmt.Printf("🚀 Démarrage du monitoring pour %d sites\n", len(cfg.Sites))
    
    var wg sync.WaitGroup
    
    // Start monitoring each site in a separate goroutine
    for _, site := range cfg.Sites {
        wg.Add(1)
        
        // Launch goroutine for each site
        go func(s config.Site) {
            defer wg.Done()
            
            // Convert string durations to time.Duration
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
            
            // Create and configure monitor
            m := monitor.New(s.URL, interval)
            m.SetName(s.Name)
            m.SetTimeout(timeout)
            
            fmt.Printf("📍 Démarrage de %s (%s) - vérification toutes les %s\n", 
                s.Name, s.URL, s.Interval)
            
            // Start monitoring (this blocks)
            if err := m.Start(); err != nil {
                log.Printf("Erreur lors du monitoring de %s: %v", s.Name, err)
            }
        }(site) // Important: pass site as parameter to avoid closure issues
    }
    
    // Wait for all goroutines to finish (never happens in this case)
    wg.Wait()
}