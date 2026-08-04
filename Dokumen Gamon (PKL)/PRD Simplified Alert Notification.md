# PRD — Simplified Alert Notification System

---

## 1. Overview

Sistem notifikasi alert yang **sederhana, fungsional, dan mudah dijelaskan** untuk keperluan demo PKL.

Prinsip: **KISS** — Keep It Simple, Stupid.

---

## 2. Alert Flow (Simplified)

```
Device Online
    ↓ (3x gagal ping berturut-turut)
Critical Alert dibuat → Notifikasi muncul
    ↓
Device Offline (tunggu)
    ↓ (1x berhasil ping)
Alert otomatis resolved → Tidak ada alert baru
```

**Hanya ada 1 arah alert:**
- **Buat alert** → saat device offline (setelah 3x gagal)
- **Resolve alert** → saat device online kembali (1x berhasil)

**Tidak ada:**
- ❌ Recovery alert (alert baru saat online)
- ❌ Success threshold (cukup 1x online langsung resolve)
- ❌ Cooldown机制
- ❌ Multi-severity (hanya critical)

---

## 3. Alert Generation Rules

### Critical Alert

| Kondisi | Aksi |
|---------|------|
| Status sebelumnya: Online/Unknown | |
| Status sekarang: Offline | |
| Gagal berturut-turut: >= 3 kali | **Buat Critical Alert** |

### Auto-Resolve

| Kondisi | Aksi |
|---------|------|
| Status sebelumnya: Offline | |
| Status sekarang: Online | |
| | **Resolve semua ongoing alert device ini** |

### No Alert

| Kondisi | Aksi |
|---------|------|
| Online → Online | Tidak ada alert |
| Offline → Offline | Tidak ada alert (sudah ada alert) |
| Warning → Online | Tidak ada alert |
| Online → Warning | Tidak ada alert |

---

## 4. Alert Data

```sql
CREATE TABLE alerts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    device_id INTEGER NOT NULL,
    title TEXT NOT NULL,           -- "Device Offline"
    status TEXT DEFAULT 'ongoing', -- 'ongoing' | 'resolved'
    severity TEXT DEFAULT 'critical', -- selalu 'critical'
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    resolved_at DATETIME,
    description TEXT DEFAULT '',
    acknowledged BOOLEAN DEFAULT FALSE,
    acknowledged_at DATETIME
);
```

**Field yang DIHAPUS:**
- `alert_type` — tidak perlu (hanya ada critical)
- `high/medium/low/info severity` — tidak perlu

---

## 5. Notification Mechanism

### In-App Banner
- Muncul saat **status berubah** (online↔offline)
- Warna: **Merah** (offline), **Hijau** (online/recovery)
- Auto-dismiss: 10 detik
- Tombol: "Lihat Detail" + Close

### Alert Center
- Daftar semua alert (ongoing + resolved)
- Filter: status, device type
- Search: nama device
- Actions: Acknowledge, Mark as Resolved

### Sidebar Badge
- Badge angka = jumlah ongoing alerts
- Update real-time via WebSocket

---

## 6. WebSocket Events

Hanya 2 event yang di-broadcast:

| Event | Kapan | Data |
|-------|-------|------|
| `status_change` | Status device berubah | device_id, old_status, new_status |
| `check_result` | Setiap check selesai | device_id, status, latency_ms |

**Tidak ada:**
- ❌ `alert_created`
- ❌ `alert_resolved`

Alert di-refresh dari REST API setiap kali `status_change` diterima.

---

## 7. API Endpoints

| Endpoint | Method | Fungsi |
|----------|--------|--------|
| `GET /api/alerts` | GET | List alerts (filter: status, severity, device_type) |
| `GET /api/alerts/{id}` | GET | Detail alert |
| `PUT /api/alerts/{id}/resolve` | PUT | Resolve alert |
| `PUT /api/alerts/{id}/acknowledge` | PUT | Acknowledge alert |
| `GET /api/alerts/count` | GET | Jumlah ongoing alerts (untuk badge) |

---

## 8. Frontend Components

| Komponen | Fungsi |
|----------|--------|
| `AlertBanner` | Banner notifikasi saat status berubah |
| `AlertBannerContainer` | Queue manager untuk banner |
| `AlertCenterPage` | Halaman utama alert |
| `AlertList` | Daftar alert |
| `AlertDetailPanel` | Detail + aksi alert |
| `AlertFilters` | Filter search, status, severity |
| `AlertSummaryCards` | Ringkasan jumlah alert |
| `SeverityBadge` | Badge severity (hanya critical) |

---

## 9. State Transition Diagram

```
              failure >= 3
    ONLINE ──────────────> OFFLINE
       ▲                      │
       │                      │
       │   1x success         │
       └──────────────────────┘
              (auto-resolve)
```

---

## 10. Testing Checklist

### Critical Alert Flow
- [ ] Tambah device dengan IP yang bisa di-ping
- [ ] Matikan/hide device tersebut
- [ ] Tunggu 3x gagal check (≈9 detik)
- [ ] **Hasil:** Critical alert muncul di Alert Center
- [ ] **Hasil:** Banner merah muncul di atas layar
- [ ] **Hasil:** Badge count sidebar bertambah

### Auto-Resolve Flow
- [ ] Dari kondisi offline di atas
- [ ] Hidupkan kembali device
- [ ] Tunggu 1x check berhasil
- [ ] **Hasil:** Alert otomatis resolved
- [ ] **Hasil:** Banner hijau muncul (recovery)
- [ ] **Hasil:** Badge count sidebar berkurang

### Acknowledge Flow
- [ ] Buka alert ongoing
- [ ] Klik "Acknowledge"
- [ ] **Hasil:** Badge "ACK" muncul di list
- [ ] **Hasil:** Tombol Acknowledge hilang di detail

### No Spam Test
- [ ] Biarkan device online stabil selama 2 menit
- [ ] **Hasil:** Tidak ada alert baru muncul
- [ ] **Hasil:** Tidak ada banner muncul

---

## 11. File Changes Required

### Backend
| File | Perubahan |
|------|-----------|
| `database/db.go` | Hapus `alert_type` column, pertahankan `acknowledged` |
| `database/models.go` | Hapus `AlertType` field |
| `monitor/engine.go` | Hapus recovery alert logic, success threshold, cooldown |
| `monitor/engine_test.go` | Update test |
| `handler/alert.go` | Simplifikasi query |

### Frontend
| File | Perubahan |
|------|-----------|
| `hooks/useWebSocket.ts` | Hapus `alert_created`/`alert_resolved` handling |
| `App.tsx` | Simplifikasi props |
| `components/AlertBannerContainer.tsx` | Hapus alertCreated handling |
| `components/AlertBanner.tsx` | Hapus alert_type logic |
| `types/index.ts` | Hapus `alertType` field |
| `lib/api.ts` | Hapus `alert_type` dari Alert interface |
| `lib/presenters.ts` | Hapus `alertType` mapping |

---

**Versi:** 1.0 (PKL Demo)
**Tanggal:** 4 Agustus 2026
