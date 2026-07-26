# 1. Dashboard Utama
### Tujuan Halaman Dashboard
Halaman ini murni berfungsi sebagai **Monitor Dinding (Wallboard Monitor)** di ruang NOC. Tujuannya adalah memberikan informasi visual instan kepada anak NOC mengenai kondisi jaringan detik ini juga, sehingga masalah _dashboard fatigue_ terselesaikan. Artinya Dashboard hanya menjawab 4 pertanyaan.
> - Berapa perangkat yang sedang dimonitor?
> - Berapa perangkat yang online?
> - Apakah ada perangkat yang bermasalah?
> - Apakah terdapat alert terbaru?

- **Informasi yang ingin ditampilkan:** Total Device, Online Device, Offline Device, Warning, Device Summary (Server, Router, Switch, Access Point, Website)
- **Latest Alert:** ketika di klik maka akan masuk ke alert center
- **System Information (Opsional):** Menampilkan seperti apakah monitoringnya berjalan, Check Intervalnya berjalan, Last Scan kapan, Notification nya berjalan? cuma untuk memastika "ohh monitoringnya berjalan"
- **Quick Action (Opsional):** Add Device, View Monitoring, View Alert, Refresh Status.
- **Auto-Trigger Audio Alarm (Fitur Solusi Utama):** Kode JavaScript sederhana di halaman ini yang mendengarkan WebSocket. Jika ada salah satu data perangkat berubah menjadi `status: "Down"`, browser di PC NOC langsung otomatis membunyikan suara alarm (beep/sirine) untuk memutus keheningan ruangan.
- **Status Koneksi Backend:** Lampu indikator kecil di pojok kanan atas web bertuliskan _"Golang Connected"_ atau _"Golang Disconnected"_ untuk membuktikan ke penguji bahwa pipa WebSocket lu sedang bekerja secara _realtime_.


---
# 2. Device Management
### Tujuan Halaman Device Management
Halaman ini digunakan oleh staf NOC (sebagai Administrator sistem) untuk **mengelola inventaris data perangkat jaringan yang ingin dimonitor**. Di sini tempat anak NOC mendaftarkan IP baru, mengubah informasi, atau menghapus perangkat yang sudah pensiun tanpa perlu menyentuh kode program backend. Nah berarti halaman Device Management mempunyai 4 fungsi utama:
1. Melihat seluruh perangkat yang terdaftar.
2. Menambahkan perangkat baru.
3. Mengubah konfigurasi perangkat.
4. Menghapus perangkat.

- **Informasi yang ingin ditampilkan:**
	1. Perangkat yang terdaftar (Device Registered)
	2. Nama Perangkat
	3. Jenis/Tipe Perangkat (Server, Acces Point, Router, dll)
	4. IP Address
	5. Method
	6. Status
	7. Port (jika ada)
	8. Action (Edit & Delete)
- **Tombol Tambah Perangkat**
- **Modal Form untuk Tambah Perangkat**
- **Tabel List Perangkat**
- **Search & Filter**
- **Active & Inactive Status**
---
# 3. Halaman Monitoring
Pertanyaan yang ingin dijawab oleh halaman Monitoring adalah:
> - Perangkat mana yang sedang bermasalah?
> - Bagaimana status perangkat tersebut?
> - Kapan terakhir dilakukan pengecekan?
> - Apakah perangkat tersebut online atau offline?
> - Bagaimana kondisi perangkat tersebut saat ini?


---

1. Alert Center
2. Log Peringatan & Riwayat
3. Settings
4. Login