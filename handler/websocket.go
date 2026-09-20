package handler

// ============================================================
// MODULE WEBSOCKET - KOMUNIKASI REAL-TIME
// ============================================================
// Module ini mengelola koneksi real-time antara server (Go) dan browser (React).
//
// Apa itu WebSocket?
// WebSocket adalah teknologi yang memungkinkan komunikasi dua arah
// antara server dan browser secara real-time (langsung/tanpa refresh).
//
// Contoh penggunaan di GAMON:
// - Ketika status perangkat berubah (online -> offline), browser langsung
//   menampilkan perubahan tersebut tanpa perlu refresh halaman.
// - Admin bisa melihat pemantauan secara langsung (live) di browser.
// ============================================================

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Upgrader mengubah koneksi HTTP biasa menjadi WebSocket.
// Setiap kali browser membuka koneksi WebSocket, koneksi akan di-upgrade.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Izinkan semua browser terhubung
	},
}

// ============================================================
// STRUKTUR DATA CLIENT (BROWSER)
// ============================================================

// Client merepresentasikan satu browser/pengguna yang sedang terhubung.
// Setiap tab browser yang membuka GAMON akan membuat satu Client baru.
type Client struct {
	hub  *Hub                // Referensi ke Hub (pusat koordinasi)
	conn *websocket.Conn     // Koneksi WebSocket ke browser
	send chan []byte          // Channel untuk mengirim pesan ke browser
	done chan struct{}        // Channel untuk menandakan koneksi selesai
}

// ============================================================
// STRUKTUR HUB (PUSAT KOORDINASI)
// ============================================================

// Hub adalah pusat koordinasi untuk semua koneksi WebSocket.
// Hub bertugas:
// 1. Mencatat semua browser yang sedang terhubung
// 2. Mengirim pesan ke SEMUA browser yang terhubung (broadcast)
// 3. Mengelola koneksi baru dan yang terputus
type Hub struct {
	clients    map[*Client]bool // Daftar semua browser yang terhubung
	broadcast  chan []byte       // Channel untuk mengirim pesan ke semua client
	register   chan *Client      // Channel untuk mendaftarkan client baru
	unregister chan *Client      // Channel untuk menghapus client yang terputus
	mu         sync.RWMutex     // Pengaman agar data tidak rusak saat diakses bersamaan
	db         *sql.DB          // Koneksi database
}

// Message adalah format paket pesan yang dikirim melalui WebSocket.
// Setiap pesan memiliki type (jenis) dan data (isi).
type Message struct {
	Type string          `json:"type"` // Jenis pesan (contoh: "check_result", "status_change")
	Data json.RawMessage `json:"data"` // Isi pesan dalam format JSON
}

// DeviceStatus adalah data status perangkat yang dikirim ke browser.
// Data ini digunakan untuk update tampilan di halaman monitoring.
type DeviceStatus struct {
	DeviceID   int     `json:"device_id"`   // ID perangkat
	Name       string  `json:"name"`        // Nama perangkat
	Type       string  `json:"type"`        // Jenis perangkat
	IP         string  `json:"ip"`          // Alamat IP
	Method     string  `json:"method"`      // Metode pengecekan
	Status     string  `json:"status"`      // Status: online/offline
	LatencyMs  float64 `json:"latency_ms"`  // Waktu respons (milidetik)
	LastCheck  string  `json:"last_check"`  // Kapan terakhir kali dicek
	LastOnline string  `json:"last_online"` // Kapan terakhir kali online
}

// ============================================================
// FUNGSI-FUNGSI HUB
// ============================================================

// NewHub membuat Hub baru yang siap digunakan.
func NewHub(db *sql.DB) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		db:         db,
	}
}

// Run adalah perulangan utama Hub yang berjalan terus menerus.
// Hub mendengarkan 3 jenis event:
// 1. Client baru terhubung (register)
// 2. Client terputus (unregister)
// 3. Pesan baru perlu dikirim ke semua client (broadcast)
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Browser baru terhubung, catat di daftar client
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client WebSocket terhubung. Total terhubung: %d", len(h.clients))

			// Kirim status awal semua perangkat ke browser baru
			go h.sendInitialState(client)

		case client := <-h.unregister:
			// Browser terputus, hapus dari daftar client
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.done)
			}
			h.mu.Unlock()
			log.Printf("Client WebSocket terputus. Total terhubung: %d", len(h.clients))

		case message := <-h.broadcast:
			// Kirim pesan ke SEMUA browser yang sedang terhubung
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Jika client lambat, putuskan koneksi
					close(client.done)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// sendInitialState mengirim data status awal semua perangkat ke browser baru.
// Ini dilakukan agar browser langsung menampilkan data saat pertama kali dibuka.
func (h *Hub) sendInitialState(client *Client) {
	// Tunggu sebentar agar koneksi benar-benar siap
	time.Sleep(100 * time.Millisecond)

	// Ambil data status semua perangkat dari database
	statuses := h.getAllDeviceStatuses()

	// Format data dalam JSON
	initialState := map[string]interface{}{
		"type": "initial_state",
		"data": statuses,
	}

	dataBytes, err := json.Marshal(initialState)
	if err != nil {
		log.Printf("Gagal memformat data awal WebSocket: %v", err)
		return
	}

	// Kirim data ke browser
	select {
	case <-client.done:
		log.Printf("Client terputus sebelum data awal terkirim")
	case client.send <- dataBytes:
		log.Printf("Berhasil mengirim data status awal ke client (%d perangkat)", len(statuses))
	}
}

// getAllDeviceStatuses mengambil data status semua perangkat aktif dari database.
// Data ini mencakup status terkini, waktu pengecekan terakhir, dan last online.
func (h *Hub) getAllDeviceStatuses() []DeviceStatus {
	var statuses []DeviceStatus

	// Query: ambil data perangkat + hasil ping terakhir
	rows, err := h.db.Query(`
		SELECT d.id, d.name, d.type, d.ip, d.method,
			COALESCE(ph.status, 'unknown') as last_status,
			COALESCE(ph.latency_ms, 0) as last_latency,
			COALESCE(ph.timestamp, d.created_at) as last_check,
			COALESCE(d.last_online, '') as last_online
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

	// Baca hasil query
	for rows.Next() {
		var ds DeviceStatus
		var lastCheck, lastOnline string
		if err := rows.Scan(&ds.DeviceID, &ds.Name, &ds.Type, &ds.IP, &ds.Method,
			&ds.Status, &ds.LatencyMs, &lastCheck, &lastOnline); err != nil {
			log.Printf("Gagal membaca baris status perangkat: %v", err)
			continue
		}
		ds.LastCheck = lastCheck
		ds.LastOnline = lastOnline
		statuses = append(statuses, ds)
	}

	return statuses
}

// Broadcast mengirim pesan ke SEMUA browser yang sedang terhubung.
// Fungsi ini dipanggil oleh Engine setiap kali ada perubahan status.
func (h *Hub) Broadcast(msgType string, data interface{}) {
	// Ubah data menjadi JSON
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Gagal merubah data ke JSON: %v", err)
		return
	}

	// Format pesan dengan type dan data
	msg := Message{
		Type: msgType,
		Data: dataBytes,
	}

	// Ubah pesan menjadi JSON dan kirim ke channel broadcast
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Gagal merubah pesan WebSocket ke JSON: %v", err)
		return
	}

	h.broadcast <- msgBytes
}

// ============================================================
// PENGELOLAAN KONEKSI WEBSOCKET
// ============================================================

// HandleWebSocket menangani permintaan koneksi WebSocket dari browser.
// Fungsi ini mengubah koneksi HTTP biasa menjadi WebSocket.
func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection ke WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Gagal upgrade HTTP ke WebSocket: %v", err)
		return
	}

	// Buat client baru untuk browser ini
	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
		done: make(chan struct{}),
	}

	// Daftarkan client ke Hub
	hub.register <- client

	// Jalankan perulangan kirim dan terima pesan di latar belakang
	go client.writePump()
	go client.readPump()
}

// ============================================================
// PERULANGAN KIRIM DAN TERIMA PESAN
// ============================================================

// readPump membaca pesan yang dikirim oleh browser (jika ada).
// Saat ini, pesan dari browser belum diproses secara khusus.
func (c *Client) readPump() {
	defer func() {
		// Bersihkan saat koneksi ditutup
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Baca pesan terus menerus dari browser
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break // Koneksi terputus, hentikan perulangan
		}

		// Coba parse pesan dari browser (untuk logging)
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

// writePump mengirim pesan dari server ke browser.
// Perulangan ini berjalan terus menerus sampai koneksi ditutup.
func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case <-c.done:
			// Koneksi selesai, hentikan perulangan
			return
		case message, ok := <-c.send:
			if !ok {
				// Channel ditutup, hentikan perulangan
				return
			}
			// Kirim pesan ke browser
			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return // Error, hentikan perulangan
			}
		}
	}
}
