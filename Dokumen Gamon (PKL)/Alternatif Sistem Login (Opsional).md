# 1. Alternatif A: Sistem IP Whitelisting (Paling Rekomendasi & RPL Banget)

Daripada pakai _username/password_, lu kunci aplikasi GAMON agar **hanya bisa dibuka dari IP komputer ruang NOC saja**.

- **Cara Kerja di Golang:** Buat satu baris kode (_middleware_) untuk mengecek IP komputer yang mencoba membuka web GAMON. Jika IP-nya cocok dengan IP ruang NOC (misal `192.168.1.100`), web akan terbuka. Jika dibuka dari HP atau Wi-Fi luar, Golang akan langsung memblokir dan menampilkan tulisan _"Akses Ditolak: Hanya untuk Komputer Ruang NOC"_.
- **Nilai Ilmiah:** Ini poin plus yang sangat besar di mata penguji sidang. Lu bisa menjelaskan konsep **Network-Level Security** tanpa perlu pusing bikin sistem database _user_ dan fitur reset password.

# 2. Alternatif B: Sistem Single Global Password (Satu Password untuk Semua)

Jika penguji lu nanti tetap bersikeras minta pengaman, buat saja satu halaman _lockscreen_ sederhana sebelum masuk dashboard.

- **Cara Kerja:** Cukup buat satu kolom _input_ password di React (tanpa _username_). Password-nya di-_hardcode_ saja di backend (misal password-nya: `NOCGamon2026`).
- **Keuntungan:** Semua anak NOC tahu password yang sama ini. Tidak ada database _user_, tidak ada risiko lupa password perorangan karena password-nya bersifat global untuk satu ruangan.

---
# 3. Konsep "Reset via Terminal/File" untuk GAMON:

Karena GAMON adalah aplikasi _software_ (bukan router fisik yang punya tombol), jika anak NOC lupa dengan _Single Global Password_-nya, pemulihannya bisa dibuat seperti ini:

1. **Gunakan File Konfigurasi (`config.json`):** Simpan password global tersebut di dalam sebuah file teks biasa bernama `config.json` di folder backend Golang.
2. **Jika Lupa Password:** Anak NOC tidak bisa mereset lewat web panel. Mereka harus meminta akses langsung ke PC Server yang menjalankan aplikasi GAMON, lalu membuka file `config.json` tersebut lewat Notepad untuk melihat atau mengganti password-nya secara manual.
3. **Nilai Ilmiah:** Di laporan PKL, lu bisa menjelaskan ini sebagai **"Local Configuration Access Security"**. Konsep ini membuktikan bahwa hanya orang yang punya akses fisik ke server NOC saja yang bisa mengubah atau melihat password aplikasi GAMON. Ini sangat aman dan jauh lebih mudah di-coding!