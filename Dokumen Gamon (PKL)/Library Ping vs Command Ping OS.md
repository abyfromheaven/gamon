Secara ilmiah, perbedaannya terletak pada tingkat interaksi aplikasi dengan kernel Sistem Operasi. Berikut adalah penjelasan mendalam mengenai perbedaan, alasan _root permission_, dan perbandingannya:

---

## 1. Mengapa Pakai Library Butuh _Root Permission_ (Hak Akses Administrator)?

Saat _library_ (baik di Go, Python, atau Node.js) melakukan ICMP Ping secara mandiri, _library_ tersebut harus membuat sebuah koneksi jaringan khusus bernama Raw Socket.

- Secara Ilmiah: _Raw Socket_ memungkinkan aplikasi untuk menyusun paket data jaringannya sendiri (membuat _header_ ICMP sendiri) dari nol, tanpa melewati lapisan standar Sistem Operasi.
- Masalah Keamanan (Security Risk): Karena _Raw Socket_ bisa memanipulasi paket data sesuka hati, Sistem Operasi (seperti Linux/Windows) menganggap tindakan ini berbahaya (bisa digunakan untuk serangan siber). Oleh karena itu, OS mengunci fitur ini. Aplikasi hanya boleh membuka _Raw Socket_ jika dijalankan dengan hak akses tertinggi yaitu Root (Linux) atau Run as Administrator (Windows).

Jika aplikasi GAMON lu nanti harus dijalankan sebagai `root` terus-menerus di server NOC, itu sangat berbahaya dari sisi keamanan server.

---

## 2. Mengapa Menggunakan Command OS (`os/exec`) Tidak Butuh Root?

Rancangan dari AI lu menyarankan untuk menggunakan perintah bawaan OS seperti `ping 192.168.1.1`. Aplikasi Golang lu hanya bertugas memanggil perintah tersebut lewat terminal di latar belakang menggunakan modul `os/exec`.

- Kenapa tidak butuh root?: Karena program `ping` bawaan Windows atau Linux sudah memiliki izin khusus dari Sistem Operasi sejak awal. OS sudah percaya dengan program `ping` miliknya sendiri.
- Cara Kerjanya: Golang lu bertindak seperti anak magang yang mengetik di cmd/terminal, lalu membaca hasilnya (_parse stdout_), dan mengubah teks terminal tersebut menjadi data JSON.

---

## 3. Tabel Perbandingan (Bahan Bab 3 Laporan Lu)

Agar penguji sidang lu terpukau, lu bisa membandingkan kedua metode ini di dalam dokumen laporan lu seperti ini:

|Aspek Perbandingan|Menggunakan Library (Raw Socket)|Menggunakan Perintah OS (`os/exec`)|
|---|---|---|
|Kebutuhan Hak Akses|Wajib Root / Administrator. Jika tidak, aplikasi akan _error_ dan _crash_.|User Biasa (Aman). Tidak perlu hak akses khusus karena memanfaatkan program bawaan OS.|
|Keamanan Sistem|Risiko Tinggi. Berbahaya jika kode backend memiliki celah keamanan karena berjalan sebagai `root`.|Risiko Rendah. Aplikasi berjalan dengan hak akses terbatas, lebih aman untuk server NOC.|
|Ketergantungan OS|Cross-Platform Mudah. Kode program yang sama bisa langsung jalan di Windows/Linux tanpa ubah kode.|Terikat Sintaks OS. Harus membedakan perintah `ping -c 1` (Linux) dan `ping -n 1` (Windows).|
|Kecepatan Proses|Sedikit lebih cepat karena data diproses langsung di memori program.|Ada sedikit jeda (_overhead_) karena Go harus membuka proses terminal baru setiap 3 detik.|

---

## Kesimpulan untuk Pilihan Proyek GAMON

Untuk level Prototype Tugas Akhir SMK RPL, pilihan menggunakan Command OS (`os/exec`) jauh lebih bijak dan aman. Kelemahan jeda waktunya tidak akan terasa karena perangkat hanya di-ping setiap 3 detik sekali (sangat ringan). Lu juga terhindar dari pusingnya mengonfigurasi hak akses _root_ di komputer jaringan NOC tempat lu PKL.