package handler

// Module WebSocket mengelola koneksi real-time dua arah antara server backend (Go) dan browser (React).
// Dengan WebSocket, browser dapat menerima perubahan status ping perangkat secara instan TANPA perlu me-refresh halaman web (live update).

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Upgrader mengubah (upgrade) HTTP Connection standar menjadi WebSocket Connection.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Mengizinkan semua origin terhubung
	},
}

// Client merepresentasikan 1 tab browser/pengguna yang sedang terhubung via WebSocket.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	done chan struct{}
}

// Hub mengelola seluruh client WebSocket yang aktif dan bertugas menyebarkan (broadcast) pesan.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	db         *sql.DB
}

// Message adalah struktur format paket pesan JSON WebSocket.
type Message struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// DeviceStatus format data status perangkat untuk dikirim ke WebSocket client.
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

// NewHub membuat objek Hub baru.
func NewHub(db *sql.DB) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		db:         db,
	}
}

// Run adalah perulangan utama Hub yang mendengarkan event terhubung, terputus, atau pesan broadcast di latar belakang.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Pendaftaran tab browser baru yang terhubung
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client WebSocket terhubung. Total terhubung: %d", len(h.clients))

			// Kirim status awal seluruh perangkat ke client baru
			go h.sendInitialState(client)

		case client := <-h.unregister:
			// Menghapus tab browser yang ditutup / terputus
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.done)
			}
			h.mu.Unlock()
			log.Printf("Client WebSocket terputus. Total terhubung: %d", len(h.clients))

		case message := <-h.broadcast:
			// Menyebarkan pesan baru ke SELURUH client WebSocket yang sedang aktif
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.done)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// sendInitialState mengirimkan snapshot data status awal saat browser pertama kali membuka web.
func (h *Hub) sendInitialState(client *Client) {
	time.Sleep(100 * time.Millisecond)

	statuses := h.getAllDeviceStatuses()

	initialState := map[string]interface{}{
		"type": "initial_state",
		"data": statuses,
	}

	dataBytes, err := json.Marshal(initialState)
	if err != nil {
		log.Printf("Gagal memformat data awal WebSocket: %v", err)
		return
	}

	select {
	case <-client.done:
		log.Printf("Client terputus sebelum data awal terkirim")
	case client.send <- dataBytes:
		log.Printf("Berhasil mengirim data status awal ke client (%d perangkat)", len(statuses))
	}
}

// getAllDeviceStatuses mengambil seluruh status perangkat aktif dari database beserta hasil ping terbarunya.
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
		log.Printf("Gagal mengambil status perangkat dari DB: %v", err)
		return statuses
	}
	defer rows.Close()

	for rows.Next() {
		var ds DeviceStatus
		var lastCheck string
		if err := rows.Scan(&ds.DeviceID, &ds.Name, &ds.Type, &ds.IP, &ds.Method, &ds.Status, &ds.LatencyMs, &lastCheck); err != nil {
			log.Printf("Gagal membaca baris status perangkat: %v", err)
			continue
		}
		ds.LastCheck = lastCheck
		statuses = append(statuses, ds)
	}

	return statuses
}

// Broadcast memformat pesan dan memasukkannya ke channel broadcast untuk dikirim ke seluruh client web.
func (h *Hub) Broadcast(msgType string, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Gagal merubah data ke JSON: %v", err)
		return
	}

	msg := Message{
		Type: msgType,
		Data: dataBytes,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Gagal merubah pesan WebSocket ke JSON: %v", err)
		return
	}

	h.broadcast <- msgBytes
}

// HandleWebSocket menerima permintaan HTTP dan mengubahnya menjadi koneksi WebSocket.
func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Gagal upgrade HTTP ke WebSocket: %v", err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
		done: make(chan struct{}),
	}

	hub.register <- client

	// Jalankan perulangan kirim dan terima di goroutine terpisah
	go client.writePump()
	go client.readPump()
}

// readPump membaca pesan yang dikirimkan oleh browser (jika ada).
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

		log.Printf("Pesan diterima dari browser: action=%s ip=%s", msg.Action, msg.IP)
	}
}

// writePump mengirimkan pesan dari channel 'send' ke browser.
func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case <-c.done:
			return
		case message, ok := <-c.send:
			if !ok {
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}
