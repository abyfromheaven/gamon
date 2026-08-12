# KEBUTUHAN NON-FUNGSIONAL SISTEM GAMON (GARDA MONITORING)
## Penjelasan Ilmiah sebagai Pelengkap Analisis Kebutuhan Sistem

---

## A. Konteks Permasalahan yang Menjadi Landasan Kebutuhan

Kebutuhan non-fungsional GAMON dirumuskan berdasarkan tuntutan penggunaan sistem di lingkungan *Network Operation Center* (NOC) Universitas Pertahanan RI. Sistem monitoring harus berjalan terus-menerus (*24/7*) di jaringan internal, menyajikan informasi kondisi perangkat secara langsung kepada administrator, serta tidak boleh menghambat proses monitoring utama ketika terjadi kegagalan pada komponen pendukung seperti media notifikasi. Tuntutan tersebut melahirkan kebutuhan-kebutuhan kualitas yang tidak terkait langsung dengan fungsi spesifik, melainkan dengan bagaimana sistem menjalankan fungsinya.

---

## B. Definisi Kebutuhan Non-Fungsional

Kebutuhan non-fungsional (*non-functional requirements*) adalah batasan-batasan pada layanan atau fungsi yang ditawarkan sistem, seperti batasan waktu, batasan proses pengembangan, dan standar yang berlaku. Kebutuhan non-fungsional sering kali lebih kritis dibandingkan kebutuhan fungsional karena pelanggaran terhadap kebutuhan non-fungsional dapat membuat sistem tidak dapat digunakan sama sekali (Sommerville, 2011). Dalam konteks pengukuran kualitas, kebutuhan non-fungsional merujuk pada karakteristik mutu perangkat lunak yang dirumuskan dalam standar ISO/IEC 25010, antara lain *performance efficiency*, *reliability*, *usability*, *security*, *maintainability*, dan *portability*.

Berdasarkan kajian tersebut, kebutuhan non-fungsional GAMON dirumuskan ke dalam tujuh kategori berikut.

---

## C. Perumusan Kebutuhan Non-Fungsional Sistem

### NFR-01: Performansi (*Performance Efficiency*)

Sistem harus menyajikan perubahan status perangkat kepada seluruh pengguna secara langsung tanpa permintaan berulang dari pengguna. Mekanisme penyampaian data dilakukan secara *push* melalui koneksi WebSocket, bukan *pull* melalui polling HTTP REST, sehingga penggunaan *bandwidth* dan beban kerja server (CPU dan RAM) lebih efisien dibandingkan pengambilan data secara berkala. Selain itu, proses pengecekan setiap perangkat dilaksanakan secara konkuren menggunakan *goroutine* sehingga banyak perangkat dapat dipantau secara paralel dalam satu proses server.

### NFR-02: Keandalan (*Reliability*)

Sistem harus tetap berjalan dan mempertahankan data meskipun terjadi kegagalan pada komponen tertentu. Seluruh proses monitoring dijalankan secara otomatis kembali ketika server dinyalakan tanpa intervensi manual. Akses basis data dikonfigurasi menggunakan mode WAL (*Write-Ahead Logging*) dan batas waktu tunggu untuk mengurangi konflik baca-tulis. Proses pengiriman notifikasi eksternal dilaksanakan secara *asynchronous* sehingga kegagalan media notifikasi tidak menghambat proses monitoring. Alert yang telah dibuat tetap tersimpan di basis data meskipun notifikasi eksternal gagal dikirim.

### NFR-03: Ketersediaan (*Availability*)

Sistem harus tersedia untuk digunakan oleh administrator NOC melalui jaringan internal instansi. Penerapan basis data berbasis file (SQLite) tanpa memerlukan server basis data terpisah dan tanpa infrastruktur tambahan menjamin sistem mudah dijalankan dan tersedia di lingkungan yang terbatas. Seluruh informasi kondisi perangkat dan riwayat pengecekan tersedia kapan pun diperlukan melalui halaman *dashboard* dan halaman monitoring.

### NFR-04: Keamanan (*Security*)

Sistem harus melindungi akses terhadap layanan monitoring dan mencegah penerima notifikasi yang tidak berwenang. Koneksi notifikasi eksternal diamankan melalui mekanisme pemasangan akun berbasis token yang memiliki masa berlaku (*expiry*), sehingga hanya akun yang memiliki token yang dapat menerima notifikasi. Untuk lingkungan jaringan internal NOC pada tahap *prototype*, sistem tidak menerapkan otentikasi pengguna secara penuh karena dinilai dapat diterima, dengan pertimbangan utama keamanan difokuskan pada saluran notifikasi eksternal yang berada di luar lingkup jaringan internal.

### NFR-05: Kemudahan Penggunaan (*Usability*)

Sistem harus mudah digunakan oleh administrator dengan latar belakang teknis yang beragam. Antarmuka disajikan melalui *website* dengan tampilan *dashboard* yang intuitif, menggunakan bahasa yang mudah dipahami, dan menampilkan status perangkat dalam bentuk indikator visual yang jelas. Fitur pencarian dan pemfilteran disediakan agar administrator dapat menemukan perangkat atau informasi yang dibutuhkan dengan cepat.

### NFR-06: Kemudahan Pemeliharaan dan Pengembangan (*Maintainability*)

Sistem harus mudah dipelihara dan dikembangkan. Arsitektur sistem dipisahkan ke dalam modul-modul yang memiliki tanggung jawab berbeda, yaitu mesin monitoring, penanganan permintaan HTTP, dan layanan notifikasi, sehingga setiap modul dapat dikembangkan tanpa mengubah modul lainnya. Struktur data hasil pengecekan dirancang generik agar metode monitoring baru dapat ditambahkan tanpa mengubah arsitektur inti.

### NFR-07: Portabilitas (*Portability*)

Sistem harus dapat dijalankan pada sistem operasi yang berbeda. Driver basis data yang digunakan bersifat murni bahasa Go tanpa ketergantungan pada pustaka C (*CGO-free*), sehingga proses kompilasi tidak memerlukan perangkat lunak eksternal dan hasil *build* bersifat portabel. Perintah pengecekan jaringan disesuaikan secara otomatis dengan sistem operasi tempat sistem dijalankan (Windows maupun Linux).

---

## D. Pemetaan Kebutuhan Non-Fungsional terhadap Karakteristik Kualitas ISO/IEC 25010

| Karakteristik Kualitas | Kebutuhan Non-Fungsional |
|---|---|
| *Performance efficiency* | NFR-01 |
| *Reliability* | NFR-02, NFR-03 |
| *Security* | NFR-04 |
| *Usability* | NFR-05 |
| *Maintainability* | NFR-06 |
| *Portability* | NFR-07 |

---

## E. Kesimpulan

Kebutuhan non-fungsional GAMON menetapkan standar kualitas sistem yang menjadi batasan dalam perancangan dan implementasi. Ketujuh kebutuhan yang dirumuskan menjamin bahwa sistem monitoring tidak hanya mampu melaksanakan fungsinya, tetapi juga menjalankannya secara efisien, andal, tersedia, aman, mudah digunakan, mudah dipelihara, dan portabel. Perumusan ini melengkapi kebutuhan fungsional sehingga keseluruhan analisis kebutuhan sistem menjadi utuh sebagai dasar perancangan GAMON.
