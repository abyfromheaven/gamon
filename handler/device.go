package handler

// ============================================================================
// MODUL HANDLER PERANGKAT (handler/device.go)
// ============================================================================
// Modul ini mengelola seluruh API Endpoint yang berkaitan dengan data Perangkat
// Jaringan (Devices), seperti:
// 1. Menampilkan daftar perangkat & detail 1 perangkat.
// 2. Menambah, mengedit, dan menghapus perangkat.
// 3. Mengaktifkan / menonaktifkan status pemantauan (start/stop monitoring).
// ============================================================================

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gamon/monitor"
)

// DeviceHandler menyimpan dependency koneksi database, engine pemantau, dan websocket hub.
type DeviceHandler struct {
	db     *sql.DB
	engine *monitor.Engine
	hub    *Hub
}

// NewDeviceHandler membuat instansi baru DeviceHandler.
func NewDeviceHandler(db *sql.DB, engine *monitor.Engine, hub *Hub) *DeviceHandler {
	return &DeviceHandler{db: db, engine: engine, hub: hub}
}

// Format standar data perangkat yang dikembalikan sebagai balasan JSON ke Frontend React.
type DeviceResponse struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Type          string `json:"type"`
	IP            string `json:"ip"`
	URL           string `json:"url"`
	Port          *int   `json:"port"`
	Method        string `json:"method"`
	Location      string `json:"location"`
	CheckInterval int    `json:"check_interval"`
	Status        string `json:"status"`
	Description   string `json:"description"`
	CreatedAt     string `json:"created_at,omitempty"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

// Struct untuk menerima data penambahan perangkat baru dari request body JSON.
type CreateDeviceRequest struct {
	Name          string `json:"name"`
	Type          string `json:"type"`
	IP            string `json:"ip"`
	URL           string `json:"url"`
	Port          *int   `json:"port"`
	Method        string `json:"method"`
	Location      string `json:"location"`
	CheckInterval int    `json:"check_interval"`
	Status        string `json:"status"`
	Description   string `json:"description"`
}

// Struct untuk menerima data pembaruan perangkat.
type UpdateDeviceRequest struct {
	Name          *string `json:"name"`
	Type          *string `json:"type"`
	IP            *string `json:"ip"`
	URL           *string `json:"url"`
	Port          *int    `json:"port"`
	Method        *string `json:"method"`
	Location      *string `json:"location"`
	CheckInterval *int    `json:"check_interval"`
	Status        *string `json:"status"`
	Description   *string `json:"description"`
}

// Struct untuk menerima ubah status aktif/nonaktif.
type ToggleStatusRequest struct {
	Status string `json:"status"`
}

// HandleDevices memproses request rute /api/devices (GET untuk list, POST untuk tambah).
func (h *DeviceHandler) HandleDevices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listDevices(w, r)
	case http.MethodPost:
		h.createDevice(w, r)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
	}
}

// HandleDevice memproses request rute /api/devices/{id} (GET detail, PUT edit, DELETE hapus, & aksi start/stop/status).
func (h *DeviceHandler) HandleDevice(w http.ResponseWriter, r *http.Request) {
	// Ekstrak ID Perangkat dari URL Path
	idTeks := strings.TrimPrefix(r.URL.Path, "/api/devices/")
	idTeks = strings.Split(idTeks, "/")[0]
	idPerangkat, err := strconv.Atoi(idTeks)
	if err != nil {
		respondError(w, http.StatusBadRequest, "ID Perangkat tidak valid")
		return
	}

	bagianSubURL := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/devices/"+idTeks), "/")

	// Cek rute khusus seperti /start, /stop, /status
	if len(bagianSubURL) > 1 && bagianSubURL[1] == "start" {
		h.startMonitoring(w, r, idPerangkat)
		return
	}
	if len(bagianSubURL) > 1 && bagianSubURL[1] == "stop" {
		h.stopMonitoring(w, r, idPerangkat)
		return
	}
	if len(bagianSubURL) > 1 && bagianSubURL[1] == "status" {
		h.toggleStatus(w, r, idPerangkat)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getDevice(w, r, idPerangkat)
	case http.MethodPut:
		h.updateDevice(w, r, idPerangkat)
	case http.MethodDelete:
		h.deleteDevice(w, r, idPerangkat)
	default:
		respondError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
	}
}

// listDevices mengambil seluruh daftar perangkat dari database SQLite.
func (h *DeviceHandler) listDevices(w http.ResponseWriter, _ *http.Request) {
	barisData, err := h.db.Query("SELECT id, name, type, ip, url, port, method, location, check_interval, status, description, created_at, updated_at FROM devices ORDER BY created_at DESC")
	if err != nil {
		log.Printf("Gagal membaca daftar perangkat: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal mengambil daftar perangkat")
		return
	}
	defer barisData.Close()

	var daftarPerangkat []DeviceResponse
	for barisData.Next() {
		var d DeviceResponse
		var waktuDibuat, waktuDiubah string
		if err := barisData.Scan(&d.ID, &d.Name, &d.Type, &d.IP, &d.URL, &d.Port, &d.Method, &d.Location, &d.CheckInterval, &d.Status, &d.Description, &waktuDibuat, &waktuDiubah); err != nil {
			log.Printf("Gagal membaca baris perangkat: %v", err)
			continue
		}
		d.CreatedAt = waktuDibuat
		d.UpdatedAt = waktuDiubah
		daftarPerangkat = append(daftarPerangkat, d)
	}

	if daftarPerangkat == nil {
		daftarPerangkat = []DeviceResponse{}
	}
	respondData(w, daftarPerangkat)
}

// getDevice mengambil rincian 1 perangkat berdasarkan ID.
func (h *DeviceHandler) getDevice(w http.ResponseWriter, _ *http.Request, idPerangkat int) {
	var d DeviceResponse
	var waktuDibuat, waktuDiubah string
	err := h.db.QueryRow("SELECT id, name, type, ip, url, port, method, location, check_interval, status, description, created_at, updated_at FROM devices WHERE id = ?", idPerangkat).
		Scan(&d.ID, &d.Name, &d.Type, &d.IP, &d.URL, &d.Port, &d.Method, &d.Location, &d.CheckInterval, &d.Status, &d.Description, &waktuDibuat, &waktuDiubah)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Perangkat tidak ditemukan")
		return
	}
	if err != nil {
		log.Printf("Gagal mengambil data perangkat ID %d: %v", idPerangkat, err)
		respondError(w, http.StatusInternalServerError, "Gagal mengambil detail perangkat")
		return
	}
	d.CreatedAt = waktuDibuat
	d.UpdatedAt = waktuDiubah
	respondData(w, d)
}

// createDevice menambahkan perangkat baru ke database dan langsung memulai pemantauan jika status 'active'.
func (h *DeviceHandler) createDevice(w http.ResponseWriter, r *http.Request) {
	var req CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	// Validasi input wajib
	if req.Name == "" {
		respondError(w, http.StatusBadRequest, "Nama perangkat wajib diisi")
		return
	}
	if req.Type == "" {
		respondError(w, http.StatusBadRequest, "Tipe perangkat wajib diisi")
		return
	}
	if req.IP == "" {
		respondError(w, http.StatusBadRequest, "Alamat IP wajib diisi")
		return
	}
	if req.Method == "" {
		req.Method = "ICMP Ping"
	}
	if req.CheckInterval <= 0 {
		req.CheckInterval = 3
	}
	if req.Status == "" {
		req.Status = "active"
	}
	if req.Status != "active" && req.Status != "inactive" {
		respondError(w, http.StatusBadRequest, "Status harus 'active' atau 'inactive'")
		return
	}

	hasilSimpan, err := h.db.Exec(
		"INSERT INTO devices (name, type, ip, url, port, method, location, check_interval, status, description) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		req.Name, req.Type, req.IP, req.URL, req.Port, req.Method, req.Location, req.CheckInterval, req.Status, req.Description,
	)
	if err != nil {
		log.Printf("Gagal menyimpan perangkat baru: %v", err)
		respondError(w, http.StatusInternalServerError, "Gagal menambahkan perangkat")
		return
	}

	idBaru, _ := hasilSimpan.LastInsertId()
	if req.Status == "active" {
		h.engine.Start(bikinKonfigurasiPemantau(int(idBaru), req.IP, req.URL, req.Port, req.Method, req.CheckInterval))
	}

	perangkatBaru := DeviceResponse{
		ID:            int(idBaru),
		Name:          req.Name,
		Type:          req.Type,
		IP:            req.IP,
		URL:           req.URL,
		Port:          req.Port,
		Method:        req.Method,
		Location:      req.Location,
		CheckInterval: req.CheckInterval,
		Status:        req.Status,
		Description:   req.Description,
	}

	log.Printf("Perangkat berhasil ditambahkan: %s (%s)", req.Name, req.IP)
	respondData(w, perangkatBaru)
}

// updateDevice memperbarui data perangkat dan menyesuaikan pemantauan engine secara otomatis.
func (h *DeviceHandler) updateDevice(w http.ResponseWriter, r *http.Request, idPerangkat int) {
	var req UpdateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	var namaLama, tipeLama, ipLama, urlLama, metodeLama, lokasiLama, statusLama, deskripsiLama string
	var portLama *int
	var intervalLama int
	err := h.db.QueryRow("SELECT name, type, ip, url, port, method, location, check_interval, status, description FROM devices WHERE id = ?", idPerangkat).
		Scan(&namaLama, &tipeLama, &ipLama, &urlLama, &portLama, &metodeLama, &lokasiLama, &intervalLama, &statusLama, &deskripsiLama)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Perangkat tidak ditemukan")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal mengambil data perangkat")
		return
	}

	nama := namaLama
	tipe := tipeLama
	ip := ipLama
	url := urlLama
	port := portLama
	metode := metodeLama
	lokasi := lokasiLama
	interval := intervalLama
	status := statusLama
	deskripsi := deskripsiLama

	if req.Name != nil {
		nama = *req.Name
	}
	if req.Type != nil {
		tipe = *req.Type
	}
	if req.IP != nil {
		ip = *req.IP
	}
	if req.URL != nil {
		url = *req.URL
	}
	if req.Port != nil {
		port = req.Port
	}
	if req.Method != nil {
		metode = *req.Method
	}
	if req.Location != nil {
		lokasi = *req.Location
	}
	if req.CheckInterval != nil {
		interval = *req.CheckInterval
	}
	if req.Status != nil {
		if *req.Status != "active" && *req.Status != "inactive" {
			respondError(w, http.StatusBadRequest, "Status harus 'active' atau 'inactive'")
			return
		}
		status = *req.Status
	}
	if req.Description != nil {
		deskripsi = *req.Description
	}

	_, err = h.db.Exec(
		"UPDATE devices SET name = ?, type = ?, ip = ?, url = ?, port = ?, method = ?, location = ?, check_interval = ?, status = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		nama, tipe, ip, url, port, metode, lokasi, interval, status, deskripsi, idPerangkat,
	)
	if err != nil {
		log.Printf("Gagal memperbarui perangkat ID %d: %v", idPerangkat, err)
		respondError(w, http.StatusInternalServerError, "Gagal memperbarui data perangkat")
		return
	}

	// Hentikan dulu pemantauan lama, lalu nyalakan kembali jika status 'active'
	h.engine.Stop(idPerangkat)
	if status == "active" {
		h.engine.Start(bikinKonfigurasiPemantau(idPerangkat, ip, url, port, metode, interval))
	}

	perangkatDiubah := DeviceResponse{
		ID:            idPerangkat,
		Name:          nama,
		Type:          tipe,
		IP:            ip,
		URL:           url,
		Port:          port,
		Method:        metode,
		Location:      lokasi,
		CheckInterval: interval,
		Status:        status,
		Description:   deskripsi,
	}

	log.Printf("Perangkat ID %d berhasil diperbarui", idPerangkat)
	respondData(w, perangkatDiubah)
}

// deleteDevice menghapus perangkat dari database SQLite dan menghentikan goroutine pemantauannya.
func (h *DeviceHandler) deleteDevice(w http.ResponseWriter, _ *http.Request, idPerangkat int) {
	hasilExec, err := h.db.Exec("DELETE FROM devices WHERE id = ?", idPerangkat)
	if err != nil {
		log.Printf("Gagal menghapus perangkat ID %d: %v", idPerangkat, err)
		respondError(w, http.StatusInternalServerError, "Gagal menghapus perangkat")
		return
	}

	jumlahTerhapus, _ := hasilExec.RowsAffected()
	if jumlahTerhapus == 0 {
		respondError(w, http.StatusNotFound, "Perangkat tidak ditemukan")
		return
	}

	h.engine.Stop(idPerangkat)
	log.Printf("Perangkat ID %d berhasil dihapus", idPerangkat)
	respondSuccess(w, "Perangkat berhasil dihapus")
}

// startMonitoring mulai memantau IP perangkat tertentu di latar belakang.
func (h *DeviceHandler) startMonitoring(w http.ResponseWriter, r *http.Request, idPerangkat int) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}
	var ip, metode, url, status string
	var port *int
	var interval int
	err := h.db.QueryRow("SELECT ip, method, url, port, check_interval, status FROM devices WHERE id = ?", idPerangkat).
		Scan(&ip, &metode, &url, &port, &interval, &status)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Perangkat tidak ditemukan")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal mengambil data perangkat")
		return
	}

	if status != "active" {
		respondError(w, http.StatusBadRequest, "Perangkat dalam status nonaktif. Aktifkan terlebih dahulu.")
		return
	}

	h.engine.Start(bikinKonfigurasiPemantau(idPerangkat, ip, url, port, metode, interval))
	log.Printf("Pemantauan dimulai untuk perangkat ID %d (%s)", idPerangkat, ip)
	respondSuccess(w, "Pemantauan berhasil dimulai")
}

// stopMonitoring menghentikan proses pemantauan perangkat.
func (h *DeviceHandler) stopMonitoring(w http.ResponseWriter, r *http.Request, idPerangkat int) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}
	h.engine.Stop(idPerangkat)
	log.Printf("Pemantauan dihentikan untuk perangkat ID %d", idPerangkat)
	respondSuccess(w, "Pemantauan berhasil dihentikan")
}

// bikinKonfigurasiPemantau membuat struct monitor.DeviceConfig untuk dikirim ke Engine Pemantau.
func bikinKonfigurasiPemantau(id int, ip, url string, port *int, metode string, interval int) monitor.DeviceConfig {
	konfig := monitor.DeviceConfig{DeviceID: id, IP: ip, URL: url, Method: metode, Interval: interval}
	if port != nil {
		konfig.Port = *port
	}
	return konfig
}

// toggleStatus mengubah status aktif/nonaktif perangkat (active <-> inactive).
func (h *DeviceHandler) toggleStatus(w http.ResponseWriter, r *http.Request, idPerangkat int) {
	if r.Method != http.MethodPut {
		respondError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	var req ToggleStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	if req.Status != "active" && req.Status != "inactive" {
		respondError(w, http.StatusBadRequest, "Status harus 'active' atau 'inactive'")
		return
	}

	var statusSaatIni string
	err := h.db.QueryRow("SELECT status FROM devices WHERE id = ?", idPerangkat).Scan(&statusSaatIni)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Perangkat tidak ditemukan")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal membaca status perangkat")
		return
	}

	if statusSaatIni == req.Status {
		respondError(w, http.StatusBadRequest, "Perangkat sudah dalam status "+req.Status)
		return
	}

	_, err = h.db.Exec("UPDATE devices SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", req.Status, idPerangkat)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Gagal mengubah status perangkat")
		return
	}

	if req.Status == "active" {
		var ip, metode, url string
		var port *int
		var interval int
		h.db.QueryRow("SELECT ip, method, url, port, check_interval FROM devices WHERE id = ?", idPerangkat).
			Scan(&ip, &metode, &url, &port, &interval)

		h.engine.Start(bikinKonfigurasiPemantau(idPerangkat, ip, url, port, metode, interval))
		log.Printf("Pemantauan diaktifkan untuk perangkat ID %d (%s)", idPerangkat, ip)
	} else {
		h.engine.Stop(idPerangkat)
		log.Printf("Pemantauan dinonaktifkan untuk perangkat ID %d", idPerangkat)
	}

	type StatusResponse struct {
		ID      int    `json:"id"`
		Status  string `json:"status"`
		Message string `json:"message"`
	}

	pesan := "Perangkat berhasil diaktifkan"
	if req.Status == "inactive" {
		pesan = "Perangkat berhasil dinonaktifkan"
	}

	respondData(w, StatusResponse{
		ID:      idPerangkat,
		Status:  req.Status,
		Message: pesan,
	})
}
