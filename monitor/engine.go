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

type DeviceConfig struct {
	DeviceID int
	IP       string
	URL      string
	Port     int
	Method   string
	Interval int
}

type Engine struct {
	hub      HubInterface
	targets  map[int]context.CancelFunc
	mu       sync.Mutex
}

func NewEngine(hub HubInterface) *Engine {
	return &Engine{
		hub:     hub,
		targets: make(map[int]context.CancelFunc),
	}
}

func (e *Engine) Start(config DeviceConfig) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, exists := e.targets[config.DeviceID]; exists {
		log.Printf("Already monitoring device %d", config.DeviceID)
		return
	}

	interval := config.Interval
	if interval <= 0 {
		interval = 3
	}

	ctx, cancel := context.WithCancel(context.Background())
	e.targets[config.DeviceID] = cancel

	go e.checkLoop(ctx, config, time.Duration(interval)*time.Second)

	log.Printf("Started monitoring device %d (%s)", config.DeviceID, config.IP)
}

func (e *Engine) Stop(deviceID int) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if cancel, exists := e.targets[deviceID]; exists {
		cancel()
		delete(e.targets, deviceID)
		log.Printf("Stopped monitoring device %d", deviceID)
	}
}

func (e *Engine) checkLoop(ctx context.Context, config DeviceConfig, interval time.Duration) {
	seq := 0
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	seq++
	result := PingOnce(config.IP, seq)
	result.DeviceID = config.DeviceID
	result.Method = config.Method
	log.Printf("[device=%d] seq=%d status=%s latency=%.2fms", config.DeviceID, seq, result.Status, result.Latency)
	e.hub.Broadcast("ping_result", result)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Check loop stopped for device %d", config.DeviceID)
			return
		case <-ticker.C:
			seq++
			result := PingOnce(config.IP, seq)
			result.DeviceID = config.DeviceID
			result.Method = config.Method
			log.Printf("[device=%d] seq=%d status=%s latency=%.2fms", config.DeviceID, seq, result.Status, result.Latency)
			e.hub.Broadcast("ping_result", result)
		}
	}
}
