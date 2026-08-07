# Teknologi: ICMP Ping (Internet Control Message Protocol)

## Rangkuman

| Teknologi | Digunakan pada |
|-----------|----------------|
| ICMP Ping | Modul monitoring perangkat jaringan (pemeriksaan status online/offline dan pengukuran latency) |

---

## 1. Peran Teknologi dalam Project GAMON

ICMP (Internet Control Message Protocol) merupakan sebuah protokol jaringan yang digunakan GAMON sebagai **metode utama pemeriksaan kondisi perangkat**. Secara khusus, GAMON memanfaatkan perintah *ping* yang bekerja atas dasar ICMP untuk mengirimkan sebuah paket permintaan (echo request) ke alamat IP perangkat target. Apabila perangkat membalas (echo reply), perangkat dianggap hidup (online); sebaliknya, apabila tidak membalas dalam batas waktu tertentu, perangkat dianggap tidak aktif (offline).

Peran ICMP dalam GAMON antara lain:

| Peran | Penjelasan |
|-------|-----------|
| Penentuan status perangkat | Menentukan apakah sebuah perangkat berada dalam keadaan online atau offline |
| Pengukuran latency | Menghitung waktu tanggap perangkat (dalam milidetik) yang digunakan sebagai salah satu informasi monitoring |
| Dasar pembuatan alert | Hasil pemeriksaan yang berulang kali gagal menjadi dasar dalam pembangkitan alert otomatis |

ICMP berperan sebagai "sensor" yang dikirim GAMON untuk mengetahui kondisi tiap perangkat. Setiap hasil pemeriksaan kemudian disimpan, ditampilkan pada halaman monitoring, dan dimanfaatkan untuk menentukan pembuatan alert apabila perangkat dianggap bermasalah.

**Bukti penggunaan:**
- `monitor/ping.go` — implementasi pemeriksaan ICMP dengan menjalankan perintah `ping` sistem
- `monitor/engine.go` — proses yang menjalankan pemeriksaan (ping) dan memproses hasilnya
- `database` — hasil ping disimpan pada tabel `ping_history`
- Dokumen `Dokumen Gamon (PKL)/Metode Pengecekan.md` dan `Evaluasi ICMP Ping Method.md`

---

## 2. Alasan Pemilihan Teknologi

Alasan pemilihan ICMP Ping dapat disimpulkan dari dokumentasi dan implementasi project GAMON.

**a. Metode dasar dan paling umum untuk monitoring jaringan.** ICMP merupakan cara yang paling mendasar dan lazim dipakai untuk memeriksa ketersediaan perangkat dalam aplikasi monitoring jaringan. Hampir semua perangkat yang memiliki alamat IP dapat diperiksa dengan metode ini.

**b. Kesederhanaan dan kemudahan.** Melalui perintah `ping`, sistem dapat langsung mengetahui apakah perangkat merespons tanpa perlu berurusan dengan detail pembuatan paket jaringan secara manual. Ini membuat obyek pemeriksaannya ringan dan mudah diterapkan.

**c. Keamanan pada level prototype.** GAMON menjalankan perintah `ping` bawaan sistem operasi lewat modul `os/exec` alih-alih membuat paket ICMP sendiri memakai *raw socket*. Pendekatan ini tidak memerlukan hak akses *root* atau *administrator* dan lebih aman bagi server (dijelaskan pada `Library Ping vs Command Ping OS.md`).

**d. Keberagaman perangkat yang dapat dipantau.** Karena hanya memerlukan alamat IP, ICMP dapat diterapkan ke berbagai perangkat seperti router, server, switch, dan access point. Ini sesuai dengan kebutuhan GAMON yang memantau banyak jenis perangkat.

---

## 3. Cara Kerja Teknologi dalam Konteks Project (Penjelasan Lengkap)

Dalam konteks GAMON, cara kerja ICMP dapat diuraikan sebagai berikut.

**a. Pemeriksaan berkala oleh mesin monitoring.** Mesin monitoring (pada `monitor/engine.go`) menjalankan satu proses pemeriksaan untuk setiap perangkat aktif secara berkala sesuai interval yang diatur. Dalam satu siklus, sistem melakukan pemeriksaan pada perangkat tersebut.

**b. Eksekusi perintah ping.** Pemeriksaan dilakukan dengan menjalankan perintah `ping` bawaan sistem operasi melalui modul `os/exec`, seperti yang terlihat pada file `monitor/ping.go`. Perintah dijalankan sesuai platform: pada Linux memakai `ping -c 1 -W 3 <ip>`, sedangkan pada Windows memakai `ping -n 1 -w 3000 <ip>`. Pendekatan ini memanfaatkan program ping yang sudah dimiliki dan dipercaya oleh sistem, sehingga tidak memerlukan hak akses khusus (diuraikan dalam `Evaluasi ICMP Ping Method.md` dan `Library Ping vs Command Ping OS.md`).

**c. Penguraian dan pemrosesan hasil.** Keluaran (output) perintah ping diuraikan untuk memperoleh dua informasi utama, yaitu status perangkat (dari ada tidaknya balasan) dan besar latency (dalam milidetik). Hasil ini diolah menjadi sebuah struktur data yang siap disimpan dan ditampilkan.

**d. Penyimpanan dan penentuan status.** Hasil pemeriksaan disimpan pada tabel `ping_history`. Sebuah perangkat baru dianggap tidak aktif setelah mengalami beberapa kegagalan berturut-turut sesuai ambang yang diatur (failure threshold), untuk mengurangi salah deteksi akibat satu kali kegagalan yang bersifat sementara.

**e. Dampak pada tampilan dan alert.** Setelah status berubah, hasilnya siarkan ke antarmuka melalui WebSocket dan dapat memicu pembuatan atau pembatalan alert secara otomatis. Dengan demikian, ICMP tidak hanya menentukan status, tetapi juga menjadi dasar dari proses monitoring dan peringatan secara menyeluruh.

Dengan alur tersebut, ICMP menjadi salah satu teknologi yang menopang seluruh proses pemantauan kesehatan jaringan pada GAMON.

---

## 4. Bukti Penggunaan

1. **File `monitor/ping.go`** — Prasarana perintah ICMP Ping. Menggunakan `os/exec` untuk menjalankan `ping` sistem dan memproses keluaran menjadi status beserta latency.

2. **File `monitor/engine.go`** — Mesin yang menjalankan siklus pemeriksaan perangkat secara berkala dan menerapkan ambang kegagalan.

3. **Tabel `ping_history`** — Menyimpan setiap hasil pemeriksaan (status, latency, ttl, waktu) pada database, diinisialisasi pada `database/db.go`.

4. **Dokumentasi perancangan** — `Metode Pengecekan.md`, `Evaluasi ICMP Ping Method.md`, dan `Library Ping vs Command Ping OS.md` yang menjelaskan rancangan dan alasan penggunaan ICMP.

5. **Model device** — pada tabel `devices`, field `method` memuat nilai "ICMP Ping" sebagai metode pemeriksaan perangkat (lihat `database/db.go`).

---

## Catatan

- GAMON menjalankan perintah `ping` bawaan sistem operasi dan bukan menyusun paket ICMP sendiri, yang menjaga keamanan karena tidak memerlukan hak akses root.
- Waktu jeda kecil akibat membuka proses terminal dianggap tidak mengganggu karena perangkat hanya dipemeriksa sekali tiap beberapa detik (interval ringan).