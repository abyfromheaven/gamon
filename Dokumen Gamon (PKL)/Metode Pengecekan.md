# 1. Metode ICMP Ping
Metode ini adalah cara paling dasar dan wajib ada di setiap aplikasi monitoring jaringan.
- **Cara Kerja:** Aplikasi GAMON secara berkala mengirimkan paket **ICMP Echo Request** (Ping) ke alamat IP atau _domain_ perangkat target.
- **Logika Sistem:** Jika perangkat membalas (`Reply`), artinya hidup. Jika tidak membalas dalam waktu tertentu (_Request Timed Out_ atau RTO), GAMON langsung menandai perangkat tersebut **Down** dan mengirimkan alarm.
- **Kelebihan:** Bisa digunakan di perangkat apa saja (Router, Switch, Server Linux/Windows, PC, bahkan Access Point) asalkan punya IP Address.
# 2. Metode Port Scanning / TCP Checking
Ada kalanya perangkatnya hidup (bisa di-ping), tetapi layanan atau aplikasi di dalamnya mati (misal web server-nya _crash_).
- **Cara Kerja:** GAMON mencoba membuka koneksi ke _port_ spesifik pada perangkat tersebut.
- **Contoh:**
    - Memeriksa Port `80` / `443` (untuk memastikan Web Server aktif).
    - Memeriksa Port `22` (SSH) atau `3389` (RDP) untuk memastikan _remote management_ server bisa diakses.
    - Memeriksa Port `3306` (MySQL) untuk memastikan _database_ hidup.
- **Logika Sistem:** Jika GAMON gagal melakukan _handshake_ ke _port_ tersebut, sistem langsung mendeteksi adanya kegagalan layanan.
# 3. Metode SNMP (Simple Network Management Protocol)
Ini adalah protokol standar industri khusus untuk memantau perangkat jaringan tanpa API vendor. Hampir semua _router_ (seperti Mikrotik, Cisco) dan _switch_ _manageable_ mendukung SNMP.

- **Cara Kerja:** GAMON bertindak sebagai **SNMP Manager** yang meminta data (_polling_) ke perangkat (**SNMP Agent**).
- **Logika Sistem:** Selain bisa tahu perangkat itu hidup atau mati, SNMP juga bisa mengambil data performa fisik perangkat secara langsung, seperti:
    - Persentase penggunaan CPU dan RAM server.
    - Suhu perangkat.
    - Status _port_ ethernet (apakah kabelnya tercolok atau lepas).