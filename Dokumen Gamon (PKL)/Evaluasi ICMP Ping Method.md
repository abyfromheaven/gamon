# 1. Masalah "Goroutine Leak" (Kebocoran Memori)
**Kelemahan Sistem Lu Sekarang:**  
Setiap kali frontend menembak `POST /api/monitor`, backend Go lu akan membuat satu _Goroutine_ baru yang berjalan dalam _looping_ abadi (`for { time.Sleep(3 * time.Second) }`).

- **Masalahnya:** Jika user me-_refresh_ browser React, atau menekan tombol monitor untuk IP yang sama berkali-kali, backend Go akan membuat _Goroutine_ baru lagi untuk IP yang sama, sementara _Goroutine_ yang lama **tetap hidup** di latar belakang dan tidak bisa mati.
- **Akibatnya:** RAM server NOC akan terkuras habis secara perlahan (_Memory Leak_) karena ratusan _Goroutine_ hantu berjalan sendirian di latar belakang melakukan ping ke IP yang sama.

**Solusi Pengembangan:**  
Lu harus menerapkan sistem **Registry/Map Tracker** di Golang. Sebelum membuat _Goroutine_ baru, Go harus mengecek dulu apakah IP tersebut sudah dimonitor atau belum.

```
// Gunakan Map untuk mencatat IP yang sedang aktif di-ping
var ActiveMonitors = make(map[string]context.CancelFunc)
```

Jika IP sudah ada di dalam Map, jangan bikin _Goroutine_ baru lagi. Jika koneksi WebSocket terputus, gunakan `context.Context` untuk mematikan (_kill_) _Goroutine_ tersebut dengan bersih.

---
# 2. Masalah "Bloking & Antrean" Jaringan (Concurrent Ping)

**Kelemahan Sistem Lu Sekarang:**  
Jika nanti IP perangkat yang lu monitoring bertambah menjadi 20 atau 50 IP, dan semuanya dijalankan dalam satu alur yang kaku, sistem lu akan mengalami masalah waktu tunggu.

- **Masalahnya:** Perintah `ping` dengan batas tunggu `-W 3` (3 detik) bersifat memblokir (_blocking execution_). Jika ada 5 perangkat yang _Down_ bersamaan, Go harus menunggu total 15 detik hanya untuk menyelesaikan satu putaran _looping_ pengecekan. Akibatnya, dashboard React lu akan terlihat macet dan tidak _realtime_ lagi.

**Solusi Pengembangan:**  
Pastikan setiap 1 IP perangkat benar-benar memiliki 1 _Goroutine_ mandiri yang berjalan secara paralel (_Worker Pool_ sederhana), sehingga jika IP A sedang RTO (menunggu 3 detik), pengecekan IP B dan IP C tidak akan ikut melambat.

---
# 3. Struktur Pengiriman Data WebSocket

**Kelemahan Sistem Lu Sekarang:**  
Frontend lu saat ini mengirim data lewat `POST /api/monitor`, lalu backend mengirim hasilnya lewat WebSocket.

- **Masalahnya:** Alur ini kurang efisien. Ketika anak NOC pertama kali membuka web GAMON, layar dashboard-nya akan kosong melompong sampai ada perubahan data dari WebSocket.

**Solusi Pengembangan (Konsep State Inisialisasi):**  
Ubah alurnya menjadi seperti ini:

1. Saat pertama kali koneksi WebSocket terhubung (`onopen` di React), Backend Go langsung mengirimkan status terakhir dari semua perangkat yang ada saat itu juga (_Current State_).
2. Setelah itu, baru jalankan _looping_ ping tiap 3 detik untuk mengirim data pembaruan saja (_Delta Update_). Ini bikin UI React lu langsung responsif begitu web dibuka.