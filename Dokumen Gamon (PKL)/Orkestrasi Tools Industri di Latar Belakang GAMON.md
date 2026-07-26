1. Mesin Nmap (Untuk Fitur "Auto Discovery" & "TCP Port Scan")

Daripada lu pusing bikin perulangan ping untuk mendeteksi IP baru atau membuka koneksi satu per satu ke port TCP, serahkan semuanya ke Nmap.

- **Fitur Subnet Scan / Cari Alat Baru:** Golang menjalankan perintah `nmap -sn 192.168.1.0/24`. Hasilnya berupa daftar host yang hidup dalam satu ruangan.
- **Fitur Deteksi Layanan (Port Scanning):** Jika user ingin memeriksa layanan apa saja yang hidup di Server B, Golang lu bisa memanggil `nmap -p 80,443,3306 --open <IP_Target>`. Nmap akan mengembalikan data apakah port web atau database-nya sedang terbuka atau tertutup.

2. Mesin Fping (Untuk Fitur "Advanced ICMP Analytics")

Di dunia jaringan Linux, ada _tool_ modifikasi dari ping biasa bernama **fping**. Ini adalah _tool_ dewa untuk urusan monitoring massal karena didesain khusus agar tidak memblokir (_non-blocking_).

- **Kenapa pakai fping?:** Ping biasa hanya bisa mengecek satu IP per perintah. Sedangkan fping bisa mengecek puluhan IP sekaligus dalam satu baris perintah.
- **Cara Kerja di Go:** Golang lu cukup memanggil `fping -c 4 -q 192.168.1.1 192.168.1.2 192.168.1.3`. Perintah ini menyuruh fping mengirimkan 4 paket sekaligus ke tiga IP berbeda secara serentak. Output-nya langsung memberikan statistik _Packet Loss_ dan rata-rata _Latency_. Lu terbebas dari masalah nomor 2 dan 11 di analisis lu kemarin!

3. Mesin Snmpwalk / Snmpget (Untuk Fitur "Fisik Perangkat")

Untuk mengambil data RAM, CPU, dan suhu dari _router_ Mikrotik/Cisco, di Linux sudah ada program bawaan bernama `snmpwalk`.

- **Cara Kerja di Go:** Golang lu tinggal menjalankan perintah `snmpwalk -v 2c -c public <IP_Target> <Kode_OID_CPU_atau_RAM>`. Nilai angka yang keluar langsung ditangkap oleh Go dan dikirim ke React untuk dijadikan grafik persentase yang keren.

---

Alasan Ilmiah Mengapa Konsep Ini Sangat Kuat untuk Laporan PKL

Jika penguji sidang bertanya, _"Kenapa kamu pakai Nmap dan Fping di dalam program kamu?"_, lu bisa mematahkan pertanyaan mereka dengan argumen ilmiah kelas atas ini:

1. **Menerapkan Konsep Reusability & Standarisasi Industri:** _"Aplikasi GAMON dirancang dengan prinsip efisiensi rekayasa perangkat lunak. Daripada membuat algoritma pemindaian soket mentah (raw socket) dari nol yang belum teruji kestabilannya, GAMON memanfaatkan biner mesin Nmap dan Fping yang sudah menjadi standar emas (Gold Standard) keamanan industri untuk akurasi data pemindaian."_
2. **Peran RPL sebagai Data Parser & Interface:** _"Tugas utama pengembang RPL di sini adalah sebagai Orchestrator. Golang bertindak sebagai interpreter yang merubah teks mentah (Unstructured Text Output) dari terminal Nmap/Fping menjadi format terstruktur JSON. Data JSON tersebut kemudian dialirkan secara realtime melalui WebSockets agar ramah bagi antarmuka web modern (React)."_