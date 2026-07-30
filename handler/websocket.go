package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients        map[*Client]bool
	broadcast      chan []byte
	register       chan *Client
	unregister     chan *Client
	mu             sync.RWMutex
	db             *sql.DB
}

type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type DeviceStatus struct {
	DeviceID  int     `json:"device_id"`
	Name      string  `json:"name"`
	Type      string  `json:"type"`
	IP        string  `json:"ip"`
	Method    string  `json:"method"`
	Status    string  `json:"status"`
	LatencyMs float64 `json:"latency_ms"`
	LastCheck string  `json:"last_check"`
}

func NewHub(db *sql.DB) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		db:         db,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client connected. Total: %d", len(h.clients))

			go h.sendInitialState(client)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client disconnected. Total: %d", len(h.clients))

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) sendInitialState(client *Client) {
	time.Sleep(100 * time.Millisecond)

	statuses := h.getAllDeviceStatuses()

	initialState := map[string]interface{}{
		"type": "initial_state",
		"data": statuses,
	}

	dataBytes, err := json.Marshal(initialState)
	if err != nil {
		log.Printf("Error marshaling initial state: %v", err)
		return
	}

	select {
	case client.send <- dataBytes:
		log.Printf("Sent initial state to client (%d devices)", len(statuses))
	default:
		log.Printf("Failed to send initial state: client send buffer full")
	}
}

func (h *Hub) getAllDeviceStatuses() []DeviceStatus {
	var statuses []DeviceStatus

	rows, err := h.db.Query(`
		SELECT d.id, d.name, d.type, d.ip, d.method,
			COALESCE(ph.status, 'unknown') as last_status,
			COALESCE(ph.latency_ms, 0) as last_latency,
			COALESCE(ph.timestamp, d.created_at) as last_check
		FROM devices d
		LEFT JOIN ping_history ph ON ph.id = (
			SELECT id FROM ping_history WHERE device_id = d.id ORDER BY id DESC LIMIT 1
		)
		WHERE d.status = 'active'
		ORDER BY d.name
	`)
	if err != nil {
		log.Printf("Error querying device statuses: %v", err)
		return statuses
	}
	defer rows.Close()

	for rows.Next() {
		var ds DeviceStatus
		var lastCheck string
		if err := rows.Scan(&ds.DeviceID, &ds.Name, &ds.Type, &ds.IP, &ds.Method, &ds.Status, &ds.LatencyMs, &lastCheck); err != nil {
			log.Printf("Error scanning device status: %v", err)
			continue
		}
		ds.LastCheck = lastCheck
		statuses = append(statuses, ds)
	}

	return statuses
}

func (h *Hub) Broadcast(msgType string, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling data: %v", err)
		return
	}

	msg := Message{
		Type: msgType,
		Data: dataBytes,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	h.broadcast <- msgBytes
}

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}

	hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		var msg struct {
			Action string `json:"action"`
			IP     string `json:"ip"`
		}
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		log.Printf("Received from client: action=%s ip=%s", msg.Action, msg.IP)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			break
		}
	}
}
