package services

import (
	"context"
	"incident-platform/backend/internal/models"
	"log"
	"net/http"
	"time"
)

// MonitorService periodically checks every registered service and
// automatically opens/resolves incidents when status changes.
type MonitorService struct {
	serviceService  *ServiceService
	incidentService *IncidentService
	client          *http.Client
	interval        time.Duration
}

func NewMonitorService(ss *ServiceService, is *IncidentService, interval time.Duration) *MonitorService {
	return &MonitorService{
		serviceService:  ss,
		incidentService: is,
		// The single most important line in this file: a hard timeout.
		// Without it, one dead server hangs the whole monitor.
		client:   &http.Client{Timeout: 5 * time.Second},
		interval: interval,
	}
}

// Start launches the monitor loop in the background. It returns immediately.
func (m *MonitorService) Start(ctx context.Context) {
	go func() {
		// Run once immediately, then on every tick. (Otherwise you'd wait
		// a full interval before the first check.)
		m.tick(ctx)
		ticker := time.NewTicker(m.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				m.tick(ctx)
			case <-ctx.Done():
				log.Println("monitor stopped")
				return
			}
		}
	}()
}

// tick runs ONE full monitoring cycle: check every service once.
func (m *MonitorService) tick(ctx context.Context) {
	services, err := m.serviceService.GetAll(ctx)
	if err != nil {
		log.Printf("monitor: failed to fetch services: %v", err)
		return
	}

	for _, svc := range services {
		m.checkService(ctx, svc)
	}
}

// checkService pings ONE service and reacts to a status change.
func (m *MonitorService) checkService(ctx context.Context, svc models.Service) {
	newStatus := m.ping(svc.URL)
	oldStatus := svc.Status

	// No change? Nothing to do. (This is why we removed the UP->HEALTHY
	// rewrite earlier — it would have wrecked this comparison.)
	if newStatus == oldStatus {
		return
	}

	log.Printf("monitor: %s changed %s -> %s", svc.Name, oldStatus, newStatus)

	// 1. Persist the new status on the service.
	statusArg := newStatus
	if _, err := m.serviceService.Update(ctx, svc.ID, models.UpdateServiceInput{Status: &statusArg}); err != nil {
		log.Printf("monitor: failed to update service %d: %v", svc.ID, err)
		return
	}

	// 2. React to the transition.
	switch {
	case oldStatus != "DOWN" && newStatus == "DOWN":
		// UP -> DOWN: open a fresh incident.
		_, err := m.incidentService.Create(ctx, models.CreateIncidentInput{
			ServiceID: svc.ID,
			Status:    "OPEN",
			Message:   "Service became unreachable",
		})
		if err != nil {
			log.Printf("monitor: failed to open incident for %d: %v", svc.ID, err)
		}

	case oldStatus == "DOWN" && newStatus != "DOWN":
		// DOWN -> UP: resolve the active incident (stamps resolved_at).
		inc, err := m.incidentService.GetOpenByServiceID(ctx, svc.ID)
		if err != nil {
			// ErrNotFound just means there was no open incident — fine.
			if err != ErrNotFound {
				log.Printf("monitor: failed to find open incident for %d: %v", svc.ID, err)
			}
			return
		}
		if err := m.incidentService.Resolve(ctx, inc.ID); err != nil {
			log.Printf("monitor: failed to resolve incident %d: %v", inc.ID, err)
		}
	}
}

// ping returns "UP" if the URL responds with a 2xx/3xx, else "DOWN".
// A network error, a timeout, or a 5xx all count as DOWN.
func (m *MonitorService) ping(url string) string {
	resp, err := m.client.Get(url)
	if err != nil {
		return "DOWN"
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "DOWN"
	}
	return "UP"
}
