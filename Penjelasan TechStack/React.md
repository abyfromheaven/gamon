# Teknologi: React (React 19.2.7 + ReactDOM 19.2.7)

## Rangkuman

| Teknologi | Digunakan pada |
|-----------|----------------|
| React | Dashboard, Device Management, Monitoring, Alert Center, Settings, dan seluruh komponen antarmuka pengguna |

---

## 1. Peran Teknologi dalam Project GAMON

React merupakan pustaka (library) JavaScript/TypeScript yang membangun seluruh sisi *frontend* atau antarmuka pengguna GAMON. Berbeda dengan backend Go yang berjalan di server, React berjalan di peramban pengguna dan bertanggung jawab menampilkan data, menerima masukan, serta berkomunikasi dengan backend.

Dalam konteks GAMON, React tidak hanya sekadar menampilkan halaman statis, melainkan membangun sebuah *Single Page Application* (Aplikasi Halaman Tunggal) tempat pengguna dapat berpindah antarhalaman tanpa di-reload oleh peramban. Seluruh halaman inti aplikasi dan komponen-komponen pendukungnya dibangun dengan model komponen (component) khas React:

| Halaman/Komponen | Peran |
|------------------|-------|
| Dashboard | Ringkasan kondisi keseluruhan (total, online, offline, alert terbaru) |
| Device Management | Pengelolaan data perangkat (CRUD, aktif/nonaktif) |
| Monitoring | Tampilan status real-time perangkat + grafik riwayat latency |
| Alert Center | Daftar, filter, detail, dan tindakan terhadap alert |
| Settings | Konfigurasi threshold, interval, dan integrasi Telegram |
| Komponen pendukung | Sidebar, TopBar, form modal, chart, badge status, notifikasi, dsb. |

React berkomunikasi dengan backend Go melalui dua jalur: **REST API** untuk operasi data (ambil, kirim, ubah, hapus) dan **WebSocket** untuk menerima pembaruan *real-time* agar tampilan selalu mengikuti kondisi perangkat yang terkini tanpa memerlukan *refresh* manual.

**Bukti penggunaan:**
- `frontend/package.json` — react `^19.2.7`, react-dom `^19.2.7`
- `frontend/src/main.tsx` — render `<App />` ke elemen DOM dengan ReactDOM
- `frontend/src/App.tsx` — komponen akar yang merangkai seluruh halaman
- Folder `frontend/src/components/` berisi puluhan komponen React (`Dashboard.tsx`, `DeviceTable.tsx`, `AlertList.tsx`, `LatencyChart.tsx`, dan lainnya)
- Folder `frontend/src/pages/` berisi halaman-halaman aplikasi

---

## 2. Alasan Pemilihan Teknologi

Alasan pemilihan React dapat disimpulkan dari kebutuhan project dan bukti pada implementasi *frontend*.

### a. Mendukung tampilan Real-time

GAMON merupakan aplikasi yang menekankan penyajian data secara *real-time*. React mampu menampilkan pembaruan data WebSocket secara dinamis dan efisien. Dengan React, ketika server menerima hasil ping baru atau perubahan status perangkat, komponen yang menampilkannya langsung dapat diperbarui secara *reactive* tanpa harus memuat ulang seluruh halaman. Bukti pada `frontend/src/hooks/useWebSocket.ts` dan pemanfaatan state di berbagai komponen.

### b. Pembangunan UI berbasis komponen yang dapat digunakan ulang

React mendorong pembangunan antarmuka sebagai kumpulan komponen kecil yang dapat digunakan kembali. Dalam GAMON, komponen seperti `StatusIndicator`, `MetricCard`, `StatusBadge`, `Toast`, dan lainnya dibangun sekali lalu dipakai ulang di beberapa halaman. Pendekatan ini membuat kode lebih rapi, mudah dirawat, dan konsisten tampilannya.

### c. Ekosistem dan kompatibilitas dengan TypeScript

React dipakai bersama TypeScript sehingga pengembangan *frontend* memperoleh *type safety* (keamanan tipe data). Ini menurunkan kemungkinan *bug* yang berkaitan dengan bentuk data yang salah, terutama karena banyaknya struktur data (device, alert, monitoring) yang dipertukarkan dengan backend. Bukti pada `frontend/src/types/index.ts` dan `frontend/tsconfig.json`.

### d. Mendukung tooling modern (Vite + Tailwind CSS)

React terintegrasi dengan **Vite** (pembangun/build dan *development server*) serta **Tailwind CSS** (kerangka CSS). Kombinasi ini mempercepat proses pengembangan dan menghasilkan aplikasi yang ringan serta mudah di-*styling*. Bukti pada `frontend/vite.config.ts` (plugin-react dan plugin-tailwind) dan `frontend/package.json`.

---

## 3. Cara Kerja Teknologi dalam Konteks Project (Penjelasan Lengkap)

Secara umum, cara kerja React dalam GAMON dapat diuraikan sebagai berikut.

1. **Entry point dan rendering.** Aplikasi dimulai dari `frontend/src/main.tsx`. File ini memasang React ke dalam elemen root pada halaman HTML dan me-render komponen `App` sebagai komponen akar. Story dewasa itu, seluruh antarmuka dibangun dari komponen-komponen React.

2. **Pengaturan tata letak (layout).** Komponen akar `App.tsx` merangkai struktur utama: `Sidebar` (navigasi samping), `TopBar` (judul/jam/koneksi), dan area konten. Navigasi antar halaman dikelola dengan state React (`useState`) dalam bentuk *halaman* seperti dasboard, device, monitoring, alert, tanpa menggunakan pustaka router eksternal.

3. **State dan pengaturan data.** Setiap komponen mengelola *state* lokalnya sendiri dengan *hook* React (`useState`, `useCallback`, `useMemo`, `useRef`). Data yang didapat dari API atau WebSocket disimpan dalam *state* dan digunakan untuk me-render ulang tampilan saat data berubah.

4. **Komunikasi REST API.** Komponen memanggil fungsi-fungsi pada `frontend/src/lib/api.ts` untuk melakukan operasi terhadap backend. Misalnya mengambil daftar device, menambah perangkat baru, mengubah data, menghapus, atau mengubah pengaturan. Hasil operasi ini me-*update* *state* sehingga tampilan menyesuaikan.

5. **Komunikasi real-time (WebSocket).** `frontend/src/hooks/useWebSocket.ts` menghubungkan aplikasi ke server WebSocket Go. Hook ini menentukan koneksi, file statistik perubahan data monitoring secara otomatis (misalnya status online/offline, informasi latency), lalu memutakhirkan state yang mendorong komponen dashboard/monitoring untuk langsung menampilkan hasil terbaru. Hook ini juga menangani koneksi otomatis apabila koneksi terputus.

6. **Penggambaran grafik.** Untuk menampilkan riwayata latency dalam bentuk grafik (area chart), React menggunakan library **Recharts**. Komponen `LatencyChart.tsx` menghasilkan *chart* dari data riwayat ping suatu perangkat.

7. **Responsif dan interaktif.** Antarmuka dibuat selai responsif agar nyaman dipakai di berbagai ukuran layar, dan interaktif dengan modal form, dialog konfirmasi, dan notifikasi toast yang seluruhnya dibangun sebagai komponen React.

Dengan pola ini, seluruh alur pengguna di GAMON dapat berjalan secara dinamis: pengguna memilih device di halaman monitoring, sistem menampilkan status real-time dan grafik; perubahan status menampilkan alert; pengelolaan device bisa dilakukan langsung dari Device Management tanpa reload.

---

## 4. Bukti Penggunaan

1. **File `frontend/package.json`** — depende react `^19.2.7` dan react-dom `^19.2.7`, script `dev`, `build`, `lint`, `preview`.

2. **File `frontend/src/main.tsx`** — titik masuk yang me-render `<App />` menggunakan ReactDOM dalam mode '*StrictMode*'.

3. **File `frontend/src/App.tsx`** — komponen akar yang merangkai Sidebar, TopBar, halaman-halaman, modal Settings, dan container alert banner; mengelola navigasi dengan *state* React.

4. **Folder `frontend/src/pages/`** — `DashboardPage.tsx`, `DeviceManagementPage.tsx`, `MonitoringPage.tsx`, `AlertCenterPage.tsx` — halaman-halaman utama GAMON.

5. **Folder `frontend/src/components/`** — pul doz komponen seperti `Sidebar.tsx`, `TopBar.tsx`, `DeviceTable.tsx`, `DeviceFormModal.tsx`, `LatencyChart.tsx`, `AlertList.tsx`, `SettingsModal.tsx`, dan lainnya.

6. **File `frontend/src/hooks/useWebSocket.ts`** — integrasi WebSocket real-tim.

7. **File `frontend/src/lib/api.ts`** dan `frontend/src/lib/presenters.ts` — lapisan pemanggilan REST dan transformasi data.

8. **File `frontend/src/types/index.ts`** — definisi tipe data (TypeScript) untuk data yang dipertukarkan.

9. **File `frontend/vite.config.ts`** — plugin-react untuk mengintegrasikan React ke Vite.

---

## Catatan

- • Dependensi React langsung yang dipakai adalah `react` dan `react-dom`. Untuk grafik digunakan library **Recharts**, dan untuk styling digunakan Tailwind CSS (dibahas terpisah).
- Navigasi antarhalaman tidak memakai pustaka router eksternal seperti `react-router-dom`; perpindahan antar halaman diurus dengan *state* state di `App.tsx`.