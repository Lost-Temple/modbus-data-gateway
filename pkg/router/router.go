package router

import (
	"context"
	"fmt"
	"log"
	
	"gateway/pkg/config"
	"gateway/pkg/core"
	"gateway/pkg/db"
	"gateway/pkg/models"
)

type Router struct {
	config *config.Config
	store  *db.Store
	egress map[string]core.EgressPlugin
}

func NewRouter(cfg *config.Config, store *db.Store) *Router {
	return &Router{
		config: cfg,
		store:  store,
		egress: make(map[string]core.EgressPlugin),
	}
}

func (r *Router) RegisterEgress(name string, plugin core.EgressPlugin) {
	r.egress[name] = plugin
}

func (r *Router) Dispatch(ctx context.Context, ch <-chan models.NormalizedData) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Router dispatching stopped")
			return
		case data := <-ch:
			// Save to local cache first
			if err := r.store.Save(data); err != nil {
				log.Printf("Failed to save data locally: %v", err)
			}
			
			// Simple routing logic based on device tags and configuration
			r.route(data)
		}
	}
}

func (r *Router) route(data models.NormalizedData) {
	for _, routeCfg := range r.config.Routes {
		if r.match(data, routeCfg.Match.DeviceTags) {
			// Apply transform if necessary (skipped for brevity in this example)
			for _, out := range routeCfg.Outputs {
				if plugin, exists := r.egress[out.Type]; exists {
					if err := plugin.Send(data); err != nil {
						log.Printf("Failed to send data to %s via %s: %v", out.Target, out.Type, err)
					}
				} else {
					log.Printf("Egress plugin not found for type: %s", out.Type)
				}
			}
		}
	}
}

func (r *Router) match(data models.NormalizedData, tags []string) bool {
	if len(tags) == 0 {
		return true // Match all if no tags specified
	}
	
	for _, tag := range tags {
		if _, ok := data.Tags[tag]; ok {
			return true
		}
	}
	return false
}
