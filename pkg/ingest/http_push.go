package ingest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	
	"gateway/pkg/config"
	"gateway/pkg/models"
	"time"
)

type SpcHttpPushPlugin struct {
	config *config.SpcHttpPushConfig
	server *http.Server
	ch     chan<- models.NormalizedData
}

func NewSpcHttpPushPlugin(cfg *config.SpcHttpPushConfig) *SpcHttpPushPlugin {
	return &SpcHttpPushPlugin{config: cfg}
}

func (p *SpcHttpPushPlugin) Start(ctx context.Context, ch chan<- models.NormalizedData) error {
	p.ch = ch
	mux := http.NewServeServeMux()
	mux.HandleFunc(p.config.Path, p.handlePush)

	addr := fmt.Sprintf("%s:%d", p.config.ListenHost, p.config.ListenPort)
	p.server = &http.Server{Addr: addr, Handler: mux}

	go func() {
		log.Printf("Starting HTTP Push Server on %s", addr)
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	return nil
}

func (p *SpcHttpPushPlugin) Stop() error {
	if p.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return p.server.Shutdown(ctx)
	}
	return nil
}

func (p *SpcHttpPushPlugin) handlePush(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data models.NormalizedData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Assuming data is somewhat raw, normalize it here or directly use if it matches
	p.ch <- data
	w.WriteHeader(http.StatusAccepted)
}
