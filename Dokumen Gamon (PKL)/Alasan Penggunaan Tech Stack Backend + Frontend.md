# 1. Keunggulan Golang untuk Aplikasi Monitoring

- **Concurrency Ringan (Goroutines):** Ini senjata utama Go. Aplikasi monitoring harus melakukan ping atau pengecekan port ke puluhan atau ratusan IP perangkat secara **bersamaan (realtime)**. Di bahasa lain, ini memakan memori besar. Di Go, lu bisa memakai _Goroutines_ yang sangat hemat RAM untuk mengecek ratusan IP sekaligus tanpa membuat aplikasi _lag_ atau _freeze_.
- **Kecepatan Eksekusi Tinggi:** Go adalah bahasa pemrograman yang dicompile langsung ke bahasa mesin (_compiled language_), membuatnya jauh lebih cepat daripada Python atau Node.js dalam memproses data mentah jaringan. [[1](https://dibimbing.id/blog/detail/nodejs-vs-golang)]
- **Aplikasi Berupa Single Binary:** Hasil akhir _backend_ Go berupa satu file executable (`.exe` atau `.bin`). Di dunia NOC, ini sangat disukai karena aplikasinya sangat mudah dideploy ke server tanpa perlu menginstal _runtime_ rumit seperti Node.js atau Python di server tersebut.
# 2. Cara Kerja Realtime Antara Go dan React
Karena lu ingin sistemnya _realtime_, jangan pakai metode _HTTP Polling_ biasa (dimana React harus _refresh_ atau minta data tiap beberapa detik karena itu boros _resource_).

Gunakan **WebSockets**. Alurnya seperti ini:

1. **Backend Go** menjalankan _Goroutines_ untuk ping semua IP perangkat di latar belakang secara terus-menerus.
2. Begitu ada perubahan status (misal dari Online ke Down), **Go langsung mendorong (push) data terbaru** ke frontend melalui jalur _WebSocket_.
3. **Frontend React** menerima data tersebut dan memperbarui tampilan dashboard secara instan tanpa perlu _refresh_ halaman. 