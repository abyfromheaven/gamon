# KEBUTUHAN FUNGSIONAL SISTEM GAMON (GARDA MONITORING)
## Penjelasan Ilmiah sebagai Dasar Analisis Kebutuhan Sistem

---

## A. Konteks Permasalahan yang Menjadi Landasan Kebutuhan

Berdasarkan hasil observasi selama Praktik Kerja Lapangan (PKL) di *Network Operation Center* (NOC) Universitas Pertahanan RI, ditemukan bahwa proses monitoring infrastruktur jaringan masih dilakukan secara manual. Administrator jaringan harus memantau sejumlah besar perangkat yang berasal dari berbagai vendor melalui antarmuka yang berbeda-beda, sehingga perangkat yang mengalami gangguan tidak selalu terdeteksi secara cepat. Kondisi ini berpotensi memperlambat waktu respons penanganan gangguan dan berdampak pada ketersediaan (*availability*) layanan jaringan yang menjadi tanggung jawab instansi.

Dari permasalahan tersebut, GAMON dirancang sebagai sistem monitoring jaringan berbasis *website* yang terpusat. Kebutuhan fungsional sistem ini diturunkan secara langsung dari masalah yang diangkat, bukan dari kebutuhan pengembangan tambahan. Artinya, setiap kebutuhan fungsional yang dirumuskan merupakan jawaban teknis terhadap satu atau lebih permasalahan yang teridentifikasi di lapangan.

---

## B. Definisi Kebutuhan Fungsional

Secara teoretis, kebutuhan fungsional (*functional requirements*) adalah pernyataan tentang layanan yang harus disediakan sistem, bagaimana sistem bereaksi terhadap *input* tertentu, dan bagaimana sistem berperilaku pada situasi tertentu (Sommerville, 2011). Kebutuhan fungsional menjawab pertanyaan **"apa yang harus dilakukan sistem"** dan menjadi dasar dalam merancang *use case* serta diagram aktivitas sistem.

Berdasarkan kajian tersebut, kebutuhan fungsional GAMON dikelompokkan ke dalam empat ranah layanan yang merepresentasikan alur kerja sistem, yaitu **manajemen perangkat**, **monitoring dan deteksi**, **alert serta notifikasi**, dan **penyajian informasi terpusat**.

---

## C. Perumusan Kebutuhan Fungsional Sistem

### 1. Ranah Manajemen Perangkat

Kebutuhan fungsional pada ranah ini menjawab permasalahan banyaknya perangkat yang harus dipantau, yaitu dengan menyediakan media untuk mendaftarkan perangkat ke dalam satu sistem terpusat.

**FR-01: Registrasi perangkat.**

Sistem harus mampu menerima pendaftaran perangkat baru beserta atributnya, meliputi nama, alamat IP, tipe perangkat (server, router, switch, *access point*, atau *website*), metode monitoring, interval pengecekan, lokasi (opsional), dan deskripsi. Kebutuhan ini merupakan prasyarat mutlak karena sistem tidak dapat memantau perangkat yang belum terdaftar.

**FR-02: Penyajian daftar perangkat.**

Sistem harus mampu menampilkan seluruh perangkat yang terdaftar beserta informasi konfigurasinya.

**FR-03: Pengubahan data perangkat.**

Sistem harus mampu memperbarui konfigurasi perangkat yang telah terdaftar, misalnya perubahan alamat IP, interval, atau metode monitoring, tanpa menghapus perangkat tersebut.

**FR-04: Penghapusan perangkat.**

Sistem harus mampu menghapus perangkat beserta seluruh data terkaitnya (riwayat pengecekan dan alert) secara konsisten.

**FR-05: Pengendalian status monitoring.**

Sistem harus mampu mengaktifkan dan menonaktifkan proses monitoring pada perangkat tertentu tanpa menghapus perangkat dari daftar, misalnya saat perangkat sedang dalam masa perawatan (*maintenance*).

### 2. Ranah Monitoring dan Deteksi

Ranah ini merupakan inti layanan yang menjawab permasalahan keterlambatan deteksi gangguan.

**FR-06: Pengecekan status secara berkala.**

Sistem harus mampu melakukan pengecekan ketersediaan perangkat secara otomatis dan periodik sesuai interval yang dikonfigurasi menggunakan protokol ICMP (*Internet Control Message Protocol*), yaitu dengan mengirimkan paket *Echo Request* dan menunggu *Echo Reply*.

**FR-07: Penentuan status perangkat.**

Sistem harus mampu menetapkan status perangkat menjadi *online* apabila perangkat merespons pengecekan dan *offline* apabila tidak merespons dalam batas waktu tertentu (*timeout*), serta menghitung latensi respons sebagai parameter kualitas koneksi.

**FR-08: Deteksi dini melalui ambang kegagalan.**

Sistem harus mampu mencegah *false alarm* dengan menetapkan perangkat berstatus *offline* hanya apabila terjadi kegagalan pengecekan beruntun sebanyak ambang batas yang dikonfigurasi (*failure threshold*). Kebutuhan ini menjamin bahwa keputusan *offline* didasarkan pada bukti berulang, bukan pada satu kegagalan yang mungkin bersifat sementara.

**FR-09: Pencatatan riwayat pengecekan.**

Sistem harus mampu menyimpan setiap hasil pengecekan (status, latensi, *time-to-live*, dan waktu pengecekan) ke dalam basis data sebagai riwayat yang dapat ditelusuri.

### 3. Ranah Alert dan Notifikasi

Ranah ini menjawab permasalahan belum adanya mekanisme notifikasi terpusat dan lambatnya respons administrator.

**FR-10: Pembuatan alert otomatis.**

Sistem harus mampu membuat alert secara otomatis (severity *critical*) ketika suatu perangkat dinyatakan *offline* berdasarkan ambang kegagalan, sehingga gangguan terangkat tanpa harus dipantau manual oleh administrator.

**FR-11: Resolusi alert otomatis.**

Sistem harus mampu menyelesaikan (meresolusi) alert yang masih berstatus *ongoing* secara otomatis ketika perangkat kembali *online*.

**FR-12: Pengelolaan alert.**

Sistem harus mampu menampilkan daftar alert, memfilter berdasarkan status, tingkat keparahan, dan tipe perangkat, serta mendukung tindakan *acknowledge* (pengakuan) dan resolusi manual oleh administrator.

**FR-13: Notifikasi real-time pada aplikasi.**

Sistem harus mampu menyampaikan perubahan status perangkat kepada seluruh pengguna aplikasi secara langsung melalui koneksi WebSocket, ditampilkan dalam bentuk *banner* notifikasi dan indikator jumlah alert berjalan.

**FR-14: Notifikasi eksternal.**

Sistem harus mampu meneruskan alert melalui media di luar aplikasi (Telegram) sehingga administrator tetap dapat mengetahui gangguan meskipun tidak sedang membuka *dashboard*. Kebutuhan ini dilengkapi mekanisme pemasangan akun (*pairing*) berbasis token dan hanya mengirimkan notifikasi pada saat terjadi perubahan status guna mencegah notifikasi berulang (*spam prevention*).

### 4. Ranah Penyajian Informasi Terpusat

Ranah ini menjawab permasalahan pemantauan melalui banyak antarmuka yang berbeda.

**FR-15: Penyajian ringkasan sistem.**

Sistem harus mampu menampilkan ringkasan kondisi keseluruhan infrastruktur, meliputi jumlah total perangkat, perangkat *online*, perangkat *offline*, serta alert terbaru, pada halaman *dashboard*.

**FR-16: Penyajian status real-time.**

Sistem harus mampu menampilkan status terkini seluruh perangkat secara *real-time* pada satu halaman monitoring, dilengkapi dengan fitur pencarian, filter berdasarkan status dan tipe perangkat, serta panel detail dan riwayat pengecekan per perangkat.

---

## D. Pemetaan Kebutuhan Fungsional terhadap Permasalahan

| Permasalahan di Lapangan | Kebutuhan Fungsional yang Menjawab |
|---|---|
| Banyak perangkat yang harus dipantau | FR-01 s.d. FR-05 |
| Proses monitoring masih manual | FR-06, FR-07, FR-08 |
| Gangguan tidak selalu terdeteksi cepat | FR-06, FR-08, FR-10 |
| Monitoring melalui banyak antarmuka yang berbeda | FR-15, FR-16 |
| Belum ada mekanisme notifikasi terpusat | FR-10, FR-11, FR-12, FR-13, FR-14 |

Pemetaan tersebut menunjukkan bahwa seluruh kebutuhan fungsional GAMON bersifat responsif terhadap permasalahan yang diangkat, sehingga sistem layak dijadikan solusi terhadap kondisi nyata di lingkungan instansi.

---

## E. Kesimpulan

Kebutuhan fungsional GAMON dirumuskan berdasarkan analisis permasalahan yang ditemukan selama observasi di NOC, bukan merupakan kebutuhan pengembangan lanjutan. Keenam belas kebutuhan fungsional yang telah diuraikan merepresentasikan layanan inti sistem: kemampuan mengelola perangkat yang dipantau, melaksanakan pengecekan berkala secara otomatis, mendeteksi gangguan secara dini, menghasilkan dan mendistribusikan alert, serta menyajikan informasi kondisi infrastruktur secara terpusat dan *real-time*. Perumusan ini menjadi landasan untuk tahap perancangan sistem selanjutnya, yaitu pembuatan *use case*, *activity diagram*, dan perancangan basis data.
