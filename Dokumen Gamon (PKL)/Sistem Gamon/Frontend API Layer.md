# Frontend API Layer — Gamon (Garda Monitoring)

## 1. Pendahuluan

### 1.1 Pengertian Frontend API Layer

Frontend API layer adalah lapisan perangkat lunak pada sisi client yang bertugas mengatur komunikasi antara antarmuka pengguna dan layanan backend. Pada aplikasi web, lapisan ini menerima kebutuhan data dari komponen tampilan, membentuk HTTP request, mengirimkannya ke server, membaca response JSON, lalu mengembalikan hasil yang sudah terstruktur kepada komponen yang memanggilnya.

Dalam Gamon, frontend dibangun menggunakan React dan TypeScript, sedangkan backend menggunakan Go. Oleh karena itu, frontend API layer berfungsi sebagai **adapter** yang menerjemahkan kebutuhan aplikasi React menjadi request ke REST API Go. Implementasi Fase 4 ditempatkan pada file `frontend/src/lib/api.ts`.

### 1.2 Kondisi Sebelum Fase 4

Sebelum Fase 4, tampilan Dashboard, Device Management, dan Alert Center masih mengambil informasi dari file dummy di frontend. Data pada layar dapat terlihat lengkap, tetapi belum berasal dari SQLite, tidak berubah ketika backend menerima data baru, dan operasi CRUD belum benar-benar diteruskan ke server.

| Aspek | Sebelum Fase 4 | Setelah Fase 4 |
|---|---|---|
| Sumber data frontend | File dummy/hardcoded | REST API backend Go siap diakses |
| Pemanggilan HTTP | Tidak terpusat | Terpusat pada `api.ts` |
| Tipe data API | Belum didefinisikan | Didefinisikan dengan TypeScript |
| Penanganan error | Belum konsisten | Menggunakan `APIError` |
| Konfigurasi alamat backend | Belum ada | Environment variable + fallback lokal |
| Integrasi halaman | Belum dilakukan | Dilanjutkan pada Fase 5 |

### 1.3 Tujuan Fase 4

Tujuan Fase 4 adalah membangun fondasi komunikasi REST API yang konsisten, dapat digunakan ulang, dan mudah diuji. Fase ini tidak mengganti isi seluruh halaman frontend. Penggantian data dummy menjadi data nyata merupakan tanggung jawab Fase 5.

Target utama Fase 4 adalah sebagai berikut:

1. Menyediakan fungsi TypeScript untuk seluruh endpoint REST utama Gamon.
2. Menyatukan pembentukan URL, header JSON, parsing response, dan penanganan error.
3. Menentukan kontrak tipe data antara React dan backend Go.
4. Menyiapkan frontend agar integrasi UI pada Fase 5 tidak perlu menulis ulang logika HTTP.

---

## 2. Posisi dalam Arsitektur Gamon

### 2.1 Arsitektur Full-Stack

Frontend API layer berada di antara komponen React dan REST API backend. Lapisan ini tidak menggantikan backend Go dan tidak berhubungan langsung dengan database. Semua akses database tetap dilakukan oleh handler Go melalui `database/sql` dan SQLite.

```text
┌────────────────────────┐
│  React Pages/Components │
│ Dashboard, Devices,     │
│ Alerts, Monitoring      │
└────────────┬───────────┘
             │ memanggil fungsi TypeScript
             ▼
┌────────────────────────┐
│ Frontend API Layer      │
│ frontend/src/lib/api.ts │
│ - request<T>()          │
│ - jsonRequest()         │
│ - APIError              │
└────────────┬───────────┘
             │ HTTP/JSON (REST API)
             ▼
┌────────────────────────┐
│ Go Backend              │
│ http.ServeMux + handler │
└────────────┬───────────┘
             │ SQL
             ▼
┌────────────────────────┐
│ SQLite                  │
│ devices, ping_history,  │
│ alerts                  │
└────────────────────────┘
```

### 2.2 Pembagian REST API dan WebSocket

Gamon menggunakan arsitektur hybrid. REST API digunakan untuk operasi yang bersifat request-response, seperti mengambil daftar device, membuat device, menghapus device, atau mengambil alert. WebSocket digunakan untuk data monitoring yang harus dikirim secara real-time, misalnya hasil ping terbaru atau perubahan status perangkat.

| Saluran | Karakteristik | Penggunaan pada Gamon |
|---|---|---|
| REST API | Request-response, satu kali komunikasi per aksi | CRUD device, dashboard awal, alert, riwayat monitoring |
| WebSocket | Koneksi persisten dan event-driven | Hasil ping, status perangkat, notifikasi perubahan status |

Dengan pembagian ini, REST API menjadi saluran kontrol dan pengambilan data awal, sedangkan WebSocket menjadi saluran pembaruan cepat. Fase 4 menangani REST API; pemakaian hasil WebSocket sebagai state multi-device merupakan bagian Fase 5.

### 2.3 Alur Request Umum

```text
Pengguna melakukan aksi pada halaman React
              │
              ▼
Komponen memanggil fungsi, misalnya createDevice(data)
              │
              ▼
api.ts membuat HTTP POST + Content-Type: application/json
              │
              ▼
Backend Go memvalidasi request dan memproses database/engine
              │
              ▼
Backend mengirim response JSON
              │
              ▼
api.ts memeriksa HTTP status, success, data/message
              │
              ├── sukses: mengembalikan data atau menyelesaikan Promise
              └── gagal : melempar APIError
```

---

## 3. Desain Frontend API Layer

### 3.1 Prinsip Satu Pintu Akses API

Seluruh endpoint REST pada frontend dipusatkan pada `frontend/src/lib/api.ts`. Komponen React tidak perlu mengetahui detail URL endpoint, cara membuat JSON body, ataupun format error dari backend.

Pendekatan ini dipilih untuk menghindari pengulangan kode. Tanpa API layer terpusat, setiap halaman berpotensi menulis `fetch()` sendiri dengan header, URL, dan penanganan error yang berbeda. Perbedaan kecil tersebut dapat menyebabkan bug sulit dilacak, terutama ketika jumlah halaman bertambah.

| Pendekatan | Kekurangan/Kelebihan |
|---|---|
| `fetch()` langsung pada setiap komponen | Cepat dibuat pada awal proyek, tetapi URL, header, dan error handling mudah terduplikasi |
| API client terpusat | Perlu satu file tambahan, tetapi kontrak, error handling, dan pemeliharaan lebih konsisten |

Gamon memilih API client terpusat karena aplikasi memiliki beberapa halaman yang akan memakai endpoint yang sama dan perlu berkembang dari prototype dummy menjadi sistem terintegrasi.

### 3.2 Konfigurasi Base URL

Alamat backend disimpan melalui `VITE_API_BASE_URL`. Jika variabel tersebut belum tersedia, aplikasi memakai fallback berikut:

```ts
const DEFAULT_API_BASE_URL = 'http://localhost:8080';
```

Pada pengembangan lokal, frontend Vite biasanya berjalan pada port yang berbeda dari backend Go. Fallback tersebut membuat proyek dapat langsung dijalankan pada komputer developer tanpa konfigurasi tambahan. Jika alamat backend berubah, misalnya menjadi server lokal NOC, nilai environment variable dapat diganti tanpa mengubah kode TypeScript.

Contoh konfigurasi lokal:

```env
VITE_API_BASE_URL=http://localhost:8080
```

Vite hanya mengekspos variabel dengan prefix `VITE_` ke source code browser. Karena nilainya akan menjadi bagian dari hasil build, variabel ini hanya boleh digunakan untuk konfigurasi publik seperti alamat API. Password database, token bot, atau secret lain tidak boleh disimpan dengan prefix tersebut.[^1]

### 3.3 Request Helper Generik

Fungsi inti pada API layer adalah `request<T>()`. Huruf `T` adalah generic TypeScript yang menyatakan tipe data hasil request. Contohnya, `request<Device[]>()` menyatakan bahwa hasil yang diharapkan berupa array device.

```ts
async function request<T>(path: string, init: RequestInit = {}, expectsData = true): Promise<T>
```

Tanggung jawab fungsi ini adalah:

1. Menggabungkan base URL dan path endpoint.
2. Menambahkan header `Accept: application/json`.
3. Menjalankan browser Fetch API.
4. Mengubah body response menjadi JSON.
5. Memeriksa status HTTP dan properti `success` dari backend.
6. Mengembalikan `data` jika tersedia atau melempar error yang mudah dipahami.

Fetch API mengembalikan objek `Response` ketika server sudah memberikan header, termasuk saat server membalas status error HTTP. Oleh sebab itu, aplikasi tetap perlu memeriksa `response.ok`; kegagalan HTTP tidak cukup ditangani hanya dengan `catch`.[^2]

### 3.4 Helper JSON Request

Fungsi `jsonRequest()` dipakai oleh operasi POST dan PUT. Fungsi tersebut membuat `RequestInit` dengan method HTTP yang benar, header JSON, serta body hasil `JSON.stringify()`.

```ts
function jsonRequest(method: 'POST' | 'PUT' | 'DELETE', body?: unknown): RequestInit
```

Dengan cara ini, fungsi seperti `createDevice()` dan `updateDevice()` tidak perlu mengulang kode untuk menyetel `Content-Type: application/json`.

### 3.5 APIError

`APIError` adalah class error khusus yang memiliki pesan dan HTTP status. Error ini memisahkan masalah API dari error JavaScript umum sehingga halaman Fase 5 dapat menampilkan pesan yang sesuai kepada pengguna.

| Kondisi | Status pada `APIError` | Contoh pesan |
|---|---:|---|
| Backend tidak bisa dihubungi | `0` | Tidak dapat terhubung ke backend Gamon |
| Response bukan JSON | Status HTTP asli | Backend Gamon mengirim respons yang bukan JSON |
| HTTP/error envelope | Status HTTP asli | Pesan dari backend atau request gagal |
| Response sukses tanpa `data` saat data diperlukan | Status HTTP asli | Respons backend tidak memiliki data yang diharapkan |

---

## 4. Kontrak Data TypeScript

### 4.1 Pengertian Kontrak Data

Kontrak data adalah kesepakatan bentuk JSON yang dipertukarkan oleh frontend dan backend. Pada Gamon, backend Go menghasilkan properti dengan format `snake_case`, misalnya `device_id`, `check_interval`, dan `latency_ms`. API layer mempertahankan format tersebut agar tidak ada perubahan data tersembunyi di tengah perjalanan.

TypeScript interface digunakan untuk mendeskripsikan bentuk objek tanpa menambah kode pada runtime. TypeScript melakukan static type checking sebelum aplikasi dijalankan sehingga kesalahan seperti salah nama field atau tipe nilai yang tidak sesuai dapat diketahui pada saat build.[^3]

### 4.2 Kelompok Tipe yang Didefinisikan

| Kelompok | Tipe | Kegunaan |
|---|---|---|
| Device | `Device`, `DeviceInput`, `DeviceUpdate` | Membaca, membuat, dan memperbarui perangkat |
| Alert | `Alert`, `AlertFilters` | Membaca alert dan membentuk filter query |
| Dashboard | `Dashboard`, `DashboardSummary`, `DashboardAlert` | Membaca ringkasan sistem dan alert terbaru |
| Monitoring | `MonitoringRecord`, `PingHistoryRecord` | Membaca status terakhir dan riwayat pengecekan |
| Status | `DeviceStatus`, `MonitorStatus`, `AlertStatus`, `AlertSeverity` | Membatasi nilai status yang sah |

### 4.3 Nilai Nullable dan Optional

Tidak semua data selalu tersedia. Misalnya port dapat bernilai `null` untuk ICMP Ping dan `resolved_at` bernilai `null` ketika alert masih berlangsung. Interface API layer menyatakan kondisi tersebut secara eksplisit.

```ts
port: number | null;
resolved_at: string | null;
last_check: string | null;
```

Beberapa field seperti `created_at` dan `updated_at` bersifat optional pada tipe `Device`, karena endpoint daftar device mengirim field tersebut, sedangkan response create/update backend saat ini berfokus pada field konfigurasi device. Penulisan tipe ini mencerminkan kontrak implementasi backend yang sebenarnya.

### 4.4 Kesiapan Multi-Method Monitoring

Tipe `DeviceMethod` sudah mengizinkan `ICMP Ping`, `HTTP Check`, dan `TCP Port`. Pada MVP saat ini engine yang aktif masih ICMP Ping, tetapi bentuk data tersebut sudah menyiapkan konfigurasi URL dan port bagi pengembangan metode berikutnya. Penambahan opsi tipe tidak berarti HTTP dan TCP checker sudah aktif; implementasi engine tetap diperlukan pada fase lanjutan.

---

## 5. Endpoint dan Fungsi Client

### 5.1 Device Management

| Fungsi | HTTP Method | Endpoint | Hasil |
|---|---|---|---|
| `fetchDevices()` | GET | `/api/devices` | `Device[]` |
| `createDevice(data)` | POST | `/api/devices` | `Device` |
| `updateDevice(id, data)` | PUT | `/api/devices/{id}` | `Device` |
| `deleteDevice(id)` | DELETE | `/api/devices/{id}` | Konfirmasi sukses |
| `startMonitor(id)` | POST | `/api/devices/{id}/start` | Konfirmasi sukses |
| `stopMonitor(id)` | POST | `/api/devices/{id}/stop` | Konfirmasi sukses |

Contoh penambahan device:

```text
React Device Form
       │
       ▼
createDevice(deviceInput)
       │
       ▼
POST /api/devices
Content-Type: application/json
       │
       ▼
Go DeviceHandler → SQLite INSERT
       │
       ▼
{ "success": true, "data": { ...device } }
```

Device yang berstatus `active` akan diikutkan ke monitoring ketika backend menjalankan proses auto-start saat boot. API client juga menyediakan `startMonitor(id)` dan `stopMonitor(id)` untuk mengirim perintah monitoring secara eksplisit.

### 5.2 Alert Center

| Fungsi | HTTP Method | Endpoint | Hasil |
|---|---|---|---|
| `fetchAlerts(filters)` | GET | `/api/alerts` | `Alert[]` |
| `resolveAlert(id)` | PUT | `/api/alerts/{id}/resolve` | Konfirmasi sukses |

`fetchAlerts()` membentuk query parameter secara aman menggunakan `URLSearchParams`. Filter hanya dikirim jika nilainya tersedia.

```text
fetchAlerts({ status: 'ongoing', severity: 'critical' })
                         │
                         ▼
GET /api/alerts?status=ongoing&severity=critical
```

### 5.3 Dashboard dan Monitoring

| Fungsi | HTTP Method | Endpoint | Hasil |
|---|---|---|---|
| `fetchDashboard()` | GET | `/api/dashboard` | Ringkasan device dan lima alert terbaru |
| `fetchMonitoring()` | GET | `/api/monitoring` | Status monitoring terakhir seluruh device |
| `fetchDeviceHistory(id)` | GET | `/api/monitoring/{id}/history` | Maksimum 50 riwayat pengecekan |

Fungsi-fungsi tersebut memberikan initial state dari database ketika halaman pertama kali dimuat. Sesudah itu, Fase 5 dapat memakai WebSocket untuk memperbarui informasi secara real-time tanpa menunggu refresh halaman.

---

## 6. Format Response dan Penanganan Kesalahan

### 6.1 Response Data

Sebagian besar endpoint baca dan create/update mengembalikan envelope berikut:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Core Router-01",
      "ip": "192.168.1.1"
    }
  ]
}
```

API client hanya mengembalikan nilai di dalam properti `data` kepada pemanggil. Dengan demikian komponen React dapat langsung menggunakan array atau objek hasil tanpa perlu berulang kali mengakses `response.data`.

### 6.2 Response Konfirmasi

Endpoint aksi seperti delete device, start/stop monitoring, dan resolve alert tidak mengembalikan data resource. Backend mengirim konfirmasi berikut:

```json
{
  "success": true,
  "message": "Device deleted successfully"
}
```

API layer membedakan response tersebut melalui parameter `expectsData`. Jika suatu fungsi hanya membutuhkan konfirmasi keberhasilan, Promise diselesaikan tanpa memaksa field `data` tersedia. Pemisahan ini penting agar kontrak frontend cocok dengan handler Go.

### 6.3 Response Gagal

Contoh response gagal dari backend:

```json
{
  "success": false,
  "message": "Device not found"
}
```

Ketika response gagal, `request<T>()` membentuk `APIError`. Pada Fase 5, halaman dapat menangkap error tersebut untuk menampilkan notifikasi kepada user atau mengembalikan state loading ke kondisi normal.

```text
Request gagal
     │
     ├── jaringan putus → APIError status 0
     ├── server mengirim bukan JSON → APIError
     ├── HTTP 4xx/5xx → APIError dengan pesan backend
     └── success: false → APIError dengan pesan backend
```

HTTP mendefinisikan semantics method dan status code untuk komunikasi client-server. Oleh karena itu, API layer memeriksa baik status HTTP maupun field `success` dari envelope Gamon.[^4]

---

## 7. Alasan Pemilihan Teknologi dan Pendekatan

### 7.1 Mengapa Menggunakan Fetch API

Fetch API merupakan API browser standar untuk membuat request jaringan. Gamon tidak menggunakan library tambahan seperti Axios karena kebutuhan Fase 4 masih mencakup method HTTP dasar, header JSON, parsing response, dan error handling sederhana. Menambahkan library eksternal untuk kebutuhan tersebut akan meningkatkan dependency tanpa manfaat yang sebanding pada prototype.

| Opsi | Kelebihan | Kekurangan | Keputusan untuk Gamon |
|---|---|---|---|
| Fetch API | Tersedia di browser, tanpa dependency tambahan, cukup untuk REST sederhana | Error handling perlu ditulis eksplisit | Dipilih |
| Axios | Fitur interceptor dan konfigurasi bawaan | Menambah dependency | Belum diperlukan |
| Request per komponen | Cepat untuk satu halaman | Duplikasi dan tidak konsisten | Tidak dipilih |

### 7.2 Mengapa Menggunakan TypeScript

TypeScript dipakai untuk memberikan pemeriksaan tipe statis pada pertukaran data. Jika backend menggunakan `latency_ms` tetapi komponen salah menulis `latency`, compiler dapat membantu menemukan ketidaksesuaian sebelum aplikasi dijalankan.

Manfaat utama bagi Gamon:

1. Mengurangi kesalahan nama field JSON.
2. Membuat status yang valid lebih jelas melalui union type.
3. Menjelaskan kontrak data bagi pengembang berikutnya.
4. Memudahkan refactor saat Fase 5 mengintegrasikan halaman.

### 7.3 Mengapa Menggunakan Environment Variable

Alamat server adalah konfigurasi yang dapat berubah antara komputer developer dan jaringan NOC. Menulis alamat langsung di banyak file akan meningkatkan risiko salah konfigurasi. Dengan satu variabel `VITE_API_BASE_URL`, perubahan cukup dilakukan pada environment atau file `.env`.

Pendekatan tersebut membuat source code tetap sama antara development dan deployment. Namun, karena konfigurasi Vite dibundel ke browser, variabel ini bukan tempat penyimpanan secret.[^1]

---

## 8. Pengujian dan Hasil Implementasi

### 8.1 Validasi Frontend

Implementasi Fase 4 divalidasi dengan build dan lint frontend:

```bash
cd frontend
npm run build
npm run lint
```

Perintah build menjalankan pemeriksaan TypeScript serta proses bundling Vite. Lint memeriksa aturan kualitas kode yang telah dikonfigurasi. Kedua perintah dijalankan dari folder `frontend` karena `package.json` berada pada folder tersebut.

Alternatif menjalankan dari root repository:

```bash
npm --prefix frontend run build
npm --prefix frontend run lint
```

### 8.2 Validasi Backend Terkait

Meskipun Fase 4 tidak mengubah kode Go, package backend tetap diuji untuk memastikan kontrak endpoint yang dikonsumsi API client tidak merusak build utama:

```bash
go test . ./database ./handler ./monitor
```

Hasil validasi menunjukkan package utama, database, handler, dan monitor berhasil diuji tanpa error kompilasi. Saat ini package tersebut belum memiliki file unit test, sehingga output `no test files` menunjukkan kompilasi berhasil, bukan berarti seluruh perilaku bisnis sudah diuji otomatis.

### 8.3 Perbaikan Error Build Sebelumnya

Saat validasi Fase 4 ditemukan tiga error frontend yang sudah ada sebelum API client dibuat. Ketiganya diperbaiki agar build keseluruhan dapat berjalan.

| Error awal | Perbaikan | Alasan |
|---|---|---|
| `StatusCard` tidak ditemukan | Menambahkan komponen `StatusCard.tsx` | Dashboard memiliki komponen tampilan hasil ping yang direferensikan |
| `JSX.Element` tidak dikenali | Menggunakan `ReactNode` pada Sidebar | Sesuai penggunaan tipe React modern |
| Status dummy `warning` tidak sesuai tipe device | Mengubah menjadi `active` | Status konfigurasi device hanya `active`/`inactive`; warning adalah hasil monitoring |

---

## 9. Batasan Prototype dan Pengembangan Lanjutan

Fase 4 adalah fondasi API client, bukan implementasi lengkap pengalaman pengguna. Batasan berikut perlu dicatat secara ilmiah agar hasil proyek tidak diklaim melebihi implementasinya.

| Fitur | Status pada Fase 4 | Rencana/Fase terkait |
|---|---|---|
| API client REST typed | Selesai | Dipakai oleh Fase 5 |
| Halaman memakai data database | Belum | Fase 5 |
| CRUD UI nyata | Belum | Fase 5 |
| State WebSocket multi-device | Belum | Fase 5 |
| Retry otomatis request | Belum | Pengembangan lanjutan |
| Timeout request khusus | Belum | Pengembangan lanjutan |
| Autentikasi dan otorisasi | Belum | Di luar scope MVP |
| Unit test frontend | Belum | Pengembangan lanjutan |

### 9.1 Hubungan dengan Fase 5

Fase 5 menggunakan fungsi yang telah dibuat pada Fase 4 sebagai berikut:

```text
Fase 4: api.ts menyediakan komunikasi REST yang typed
                         │
                         ▼
Fase 5: React pages memanggil fungsi API tersebut
                         │
                         ▼
Data dummy diganti dengan data SQLite melalui backend Go
                         │
                         ▼
WebSocket memperbarui status monitoring secara real-time
```

Dengan urutan ini, perubahan halaman dapat difokuskan pada state, loading, error, dan tampilan. Detail HTTP tetap berada di satu lokasi.

---

## 10. Kesimpulan

Fase 4 berhasil membangun frontend API layer Gamon pada `frontend/src/lib/api.ts`. Lapisan ini menyediakan fungsi REST untuk device, alert, dashboard, monitoring, dan riwayat pengecekan. Selain itu, implementasi memiliki type contract TypeScript, konfigurasi alamat backend, serta mekanisme penanganan error yang konsisten.

Prinsip desain utama pada Fase 4 adalah **pemisahan tanggung jawab**. Komponen React berfokus menampilkan informasi dan merespons aksi pengguna, API layer berfokus pada komunikasi HTTP, backend Go menangani logika bisnis, sedangkan SQLite menyimpan data permanen. Pemisahan ini memudahkan pengembangan Fase 5 dan mengurangi duplikasi kode.

### 10.1 Poin Jawaban untuk Sidang

| Pertanyaan penguji | Jawaban ringkas |
|---|---|
| Mengapa perlu API layer? | Agar seluruh request HTTP, parsing JSON, dan error handling konsisten serta tidak ditulis ulang di setiap halaman. |
| Mengapa memakai TypeScript? | Untuk memeriksa kesesuaian bentuk data API sebelum aplikasi dijalankan dan mengurangi kesalahan field JSON. |
| Mengapa memakai Fetch API, bukan Axios? | Fetch sudah tersedia di browser dan kebutuhan prototype Gamon belum memerlukan dependency tambahan. |
| Mengapa URL backend memakai environment variable? | Agar alamat backend dapat diubah antar lingkungan tanpa mengedit kode aplikasi. |
| Apa perbedaan REST API dan WebSocket di Gamon? | REST API digunakan untuk CRUD dan data awal; WebSocket digunakan untuk pembaruan monitoring secara real-time. |
| Apakah Fase 4 sudah membuat halaman memakai database? | Belum. Fase 4 menyiapkan jalur komunikasi; pemakaian pada halaman dilakukan pada Fase 5. |

---

## Daftar Pustaka

[^1]: Vite. *Env Variables and Modes*. https://vite.dev/guide/env-and-mode (diakses 30 Juli 2026).
[^2]: MDN Web Docs. *Fetch API*. https://developer.mozilla.org/en-US/docs/Web/API/Fetch_API (diakses 30 Juli 2026).
[^3]: TypeScript. *Handbook: Interfaces*. https://www.typescriptlang.org/docs/handbook/interfaces.html (diakses 30 Juli 2026).
[^4]: Fielding, R., et al. *RFC 9110: HTTP Semantics*. Internet Engineering Task Force, 2022. https://www.rfc-editor.org/rfc/rfc9110.pdf (diakses 30 Juli 2026).
