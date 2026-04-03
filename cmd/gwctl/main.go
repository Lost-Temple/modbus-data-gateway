package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gateway/pkg/config"
	"gateway/pkg/core"
	"gateway/pkg/db"
	"gateway/pkg/egress"
	"gateway/pkg/ingest"
	"gateway/pkg/models"
	"gateway/pkg/router"
	"gopkg.in/yaml.v3"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "c", "config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	// 1. Load Configuration
	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize Database
	store, err := db.NewStore(cfg.Ingest.SpcSqlite.DbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer store.Close()

	// 3. Initialize Router
	r := router.NewRouter(cfg, store)

	// 4. Initialize Plugins based on Config
	var ingester core.IngestPlugin
	
	switch cfg.Ingest.Mode {
	case "spc_http_push":
		ingester = ingest.NewSpcHttpPushPlugin(cfg.Ingest.SpcHttpPush)
	default:
		log.Fatalf("Unsupported ingest mode: %s", cfg.Ingest.Mode)
	}

	// Register egress adapters based on configured routes
	hzwsCfg := &cfg.Hzws
	r.RegisterEgress("hzws", egress.NewHzwsAdapter(hzwsCfg))

	// Add more plugins here based on configuration (e.g., MQTT)

	// 5. Start the Engine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dataCh := make(chan models.NormalizedData, 100)

	// Start Ingester
	if err := ingester.Start(ctx, dataCh); err != nil {
		log.Fatalf("Failed to start ingester: %v", err)
	}

	// Start Dispatcher
	go r.Dispatch(ctx, dataCh)

	// Start Egress plugins (if needed)
	for _, p := range cfg.Routes {
		for _, out := range p.Outputs {
			if egressPlugin, exists := r.egress[out.Type]; exists {
				if err := egressPlugin.Start(ctx); err != nil {
					log.Printf("Failed to start egress plugin %s: %v", out.Type, err)
				}
			}
		}
	}

	// 6. Handle OS Signals for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	
	<-sigCh
	log.Println("Shutting down gateway...")
	
	cancel()
	if err := ingester.Stop(); err != nil {
		log.Printf("Error stopping ingester: %v", err)
	}
	
	// Stop egress plugins
	for _, p := range r.egress {
		if err := p.Stop(); err != nil {
			log.Printf("Error stopping egress plugin: %v", err)
		}
	}

	log.Println("Gateway stopped gracefully.")
}

func loadConfig(path string) (*config.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file error: %w", err)
	}

	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unmarshal yaml error: %w", err)
	}

	// Apply defaults if needed
	if cfg.Ingest.SpcSqlite == nil {
		cfg.Ingest.SpcSqlite = &config.SpcSqliteConfig{DbPath: "history.sqlite3"}
	}

	return &cfg, nil
}
