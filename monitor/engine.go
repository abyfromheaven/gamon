package monitor

import (
	"context"
	"log"
	"sync"
	"time"
)

type HubInterface interface {
	Broadcast(msgType string, data interface{})
}

type Engine struct {
	hub      HubInterface
	targets  map[string]context.CancelFunc
	mu       sync.Mutex
}

func NewEngine(hub HubInterface) *Engine {
	return &Engine{
		hub:     hub,
		targets: make(map[string]context.CancelFunc),
	}
}

func (e *Engine) Start(ip string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.targets[ip]; exists {
		log.Printf("Already monitoring %s", ip)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.targets[ip] = cancel

	go e.pingLoop(ctx, ip)

	log.Printf("Started monitoring %s", ip)
}

func (e *Engine) Stop(ip string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if cancel, exists := e.targets[ip]; exists {
		cancel()
		delete(e.targets, ip)
		log.Printf("Stopped monitoring %s", ip)
	}
}

func (e *Engine) pingLoop(ctx context.Context, ip string) {
	seq := 0
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	// Ping pertama langsung
	seq++
	result := PingOnce(ip, seq)
	log.Printf("[%s] seq=%d status=%s latency=%.2fms", ip, seq, result.Status, result.Latency)
	e.hub.Broadcast("ping_result", result)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Ping loop stopped for %s", ip)
			return
		case <-ticker.C:
			seq++
			result := PingOnce(ip, seq)
			log.Printf("[%s] seq=%d status=%s latency=%.2fms", ip, seq, result.Status, result.Latency)
			e.hub.Broadcast("ping_result", result)
		}
	}
}
