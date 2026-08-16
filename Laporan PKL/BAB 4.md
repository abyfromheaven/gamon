B. Implementasi Kegiatan PKL 

Berdasarkan hasil analisis permasalahan yang telah diuraikan pada subbab sebelumnya, penulis merancang dan mengimplementasikan sebuah prototype sistem monitoring infrastruktur teknologi informasi berbasis website sebagai alternatif solusi terhadap permasalahan yang ditemukan selama pelaksanaan Praktik Kerja Lapangan (PKL). Prototype tersebut dikembangkan dengan nama GAMON (Garda Monitoring) dan dirancang untuk membantu proses monitoring perangkat secara lebih terpusat, terstruktur, serta mudah diakses oleh administrator. Pada subbab ini akan dijelaskan proses implementasi yang meliputi gambaran umum aplikasi, teknologi yang digunakan, analisis kebutuhan sistem, perancangan sistem, implementasi fitur, hingga hasil implementasi prototype yang telah dikembangkan. 

1. Gambaran Umum Aplikasi  

GAMON (Garda Monitoring) merupakan prototype aplikasi berbasis website yang digunakan sebagai media pengelolaan dan pemantauan infrastruktur teknologi informasi. Aplikasi ini menyediakan lingkungan pemantauan yang menampilkan informasi perangkat secara terorganisasi, sehingga pengguna dapat memperoleh gambaran mengenai kondisi perangkat yang terdaftar dalam sistem melalui satu aplikasi. GAMON ditujukan bagi administrator atau pengelola infrastruktur teknologi informasi yang bertanggung jawab terhadap pemantauan perangkat. 

Secara umum, GAMON terdiri atas beberapa bagian utama yang mendukung pengelolaan dan pemantauan perangkat. Bagian dashboard digunakan untuk menampilkan gambaran umum kondisi perangkat, sedangkan device management digunakan untuk mengelola data perangkat yang terdapat dalam sistem. Selain itu, terdapat bagian monitoring untuk menampilkan kondisi perangkat, alert notification untuk menyampaikan informasi mengenai perubahan kondisi perangkat, serta monitoring history untuk menampilkan catatan hasil pemantauan. Aplikasi juga menyediakan bagian settings untuk mengatur kebutuhan tertentu dalam penggunaan sistem. 

Ruang lingkup GAMON berfokus pada pengelolaan dan pemantauan perangkat infrastruktur teknologi informasi yang telah terdaftar dalam sistem. Informasi yang disajikan mencakup data perangkat, kondisi perangkat, pemberitahuan perubahan kondisi, dan riwayat pemantauan. Dengan demikian, GAMON mencakup fungsi pengelolaan informasi perangkat dan penyajian informasi hasil pemantauan dalam satu aplikasi berbasis website.  

2. Teknologi yang Digunakan 

Pengembangan aplikasi GAMON memanfaatkan beberapa teknologi utama yang saling terintegrasi untuk mendukung proses monitoring infrastruktur teknologi informasi. Setiap teknologi memiliki peran yang berbeda, mulai dari pengembangan antarmuka pengguna, pengolahan logika aplikasi, penyimpanan data, hingga proses komunikasi dalam sistem monitoring. Pemilihan teknologi tersebut disesuaikan dengan kebutuhan fungsional aplikasi serta karakteristik sistem yang dikembangkan. Uraian mengenai teknologi yang digunakan beserta perannya dalam implementasi aplikasi dijelaskan pada bagian berikut. 

1. Go (Golang) 
    

Go atau Golang merupakan bahasa pemrograman yang dikembangkan oleh Google dan dirancang untuk mendukung pengembangan perangkat lunak yang sederhana, efisien, dan memiliki kinerja yang baik. Go banyak digunakan dalam pengembangan aplikasi backend, layanan berbasis jaringan, serta aplikasi yang membutuhkan pemrosesan secara efisien. Salah satu fitur yang menjadi karakteristik Go adalah dukungannya terhadap concurrency, yaitu kemampuan suatu program untuk menangani beberapa proses secara bersamaan atau saling tumpang tindih dalam pelaksanaannya (Hanif, 2024; Sekawan Media, 2025). 

Pada aplikasi GAMON, Go digunakan sebagai bahasa pemrograman backend untuk menangani logika aplikasi, menyediakan layanan REST API, serta menjalankan proses yang berkaitan dengan monitoring perangkat. Go dipilih karena memiliki dukungan terhadap concurrency melalui mekanisme goroutine, sehingga sesuai dengan kebutuhan GAMON yang melakukan pemantauan terhadap beberapa perangkat. Penggunaan goroutine memungkinkan pekerjaan tertentu dijalankan secara konkuren tanpa harus menunggu seluruh proses sebelumnya selesai secara berurutan (novalagung, 2023; Sekawan Media, 2025). 

2. WebSocket 
    

WebSocket merupakan protokol komunikasi yang memungkinkan pertukaran data antara client dan server secara real-time. Teknologi ini sesuai untuk aplikasi yang membutuhkan pembaruan informasi secara cepat dan berkelanjutan, karena data dapat disampaikan kepada pengguna tanpa harus melakukan permintaan secara berulang (AppMaster, 2022). 

Pada aplikasi GAMON, WebSocket digunakan untuk menyampaikan pembaruan kondisi perangkat dari backend kepada frontend secara real-time. GAMON dirancang sebagai aplikasi monitoring yang dapat digunakan untuk memantau kondisi perangkat secara berkelanjutan, sehingga informasi mengenai perubahan kondisi perangkat perlu ditampilkan pada aplikasi tanpa mengharuskan pengguna melakukan refresh halaman secara manual. Penggunaan WebSocket memungkinkan data hasil monitoring yang baru tersedia untuk langsung diperbarui pada tampilan aplikasi. 

Implementasi WebSocket pada sisi backend GAMON menggunakan Gorilla WebSocket, yaitu pustaka WebSocket untuk bahasa pemrograman Go. Penggunaan Gorilla WebSocket dipilih untuk mendukung kebutuhan komunikasi real-time antara backend dan frontend, khususnya dalam memperbarui informasi kondisi perangkat selama proses monitoring berlangsung. 

3. SQLite 
    

SQLite merupakan sistem manajemen basis data (Database Management System/DBMS) yang menggunakan pendekatan penyimpanan data dalam sebuah berkas (file) sehingga tidak memerlukan database server terpisah. Karakteristik tersebut membuat SQLite bersifat ringan, mudah digunakan, dan sesuai untuk aplikasi yang membutuhkan pengelolaan data secara lokal (IDCloudHost, 2022). 

Pada aplikasi GAMON, SQLite digunakan sebagai media penyimpanan data utama untuk menyimpan informasi yang dibutuhkan oleh sistem, seperti data perangkat, riwayat hasil monitoring, data alert, serta pengaturan aplikasi. SQLite dipilih karena tidak memerlukan proses instalasi maupun konfigurasi database server yang kompleks sehingga lebih mudah diimplementasikan pada tahap pengembangan prototype. Selain itu, seluruh data disimpan dalam satu berkas (file), memiliki penggunaan sumber daya yang ringan, serta mendukung transaksi ACID yang membantu menjaga konsistensi dan keandalan data. Karakteristik tersebut menjadikan SQLite sesuai dengan kebutuhan GAMON sebagai aplikasi monitoring yang dijalankan pada satu server (Koder.ai, n.d.). 

4. React 
    

React merupakan pustaka (library) JavaScript yang dikembangkan untuk membangun antarmuka pengguna (User Interface/UI) pada aplikasi web. React memungkinkan antarmuka aplikasi dikembangkan menggunakan komponen yang dapat digunakan kembali sehingga proses pengembangan menjadi lebih terstruktur dan efisien (Setiawan, 2022). 

Pada aplikasi GAMON, React digunakan untuk membangun seluruh antarmuka pengguna, seperti dashboard, device management, monitoring, alert center, settings, serta berbagai komponen pendukung lainnya yang digunakan dalam proses monitoring. React dipilih karena memudahkan pengembangan antarmuka yang terstruktur, konsisten, dan mudah dipelihara. Selain itu, React mampu menampilkan perubahan data secara dinamis sehingga informasi kondisi perangkat pada GAMON dapat diperbarui dengan lebih responsif. 

5. TypeScript 
    

TypeScript merupakan bahasa pemrograman yang dikembangkan oleh Microsoft sebagai superset dari JavaScript dengan menambahkan sistem tipe data (static typing). Penambahan sistem tipe tersebut membantu pengembang menulis kode yang lebih terstruktur, mengurangi potensi kesalahan, serta mempermudah pengembangan aplikasi, terutama yang memiliki skala cukup besar (Hidayat, 2016). 

Pada aplikasi GAMON, TypeScript digunakan pada sisi frontend yang dibangun menggunakan React untuk mengembangkan halaman, komponen, serta proses komunikasi data antara frontend dan backend. TypeScript dipilih karena membantu menjaga konsistensi struktur data yang digunakan dalam aplikasi, mengurangi potensi kesalahan selama proses pengembangan, serta mempermudah pemeliharaan kode seiring dengan bertambahnya fitur pada GAMON. 

6. Tailwind CSS 
    

Tailwind CSS merupakan framework CSS yang menerapkan pendekatan utility-first, yaitu menyediakan berbagai kelas utilitas yang dapat digunakan secara langsung untuk membangun dan mengatur tampilan antarmuka pengguna. Pendekatan tersebut memudahkan proses pengembangan antarmuka menjadi lebih cepat, konsisten, dan mudah dipelihara (IDCloudHost, 2025). Pada aplikasi GAMON, Tailwind CSS digunakan untuk membangun tampilan antarmuka pada seluruh halaman aplikasi, seperti dashboard, device management, monitoring, alert center, dan settings. Tailwind CSS dipilih karena mempermudah pengembangan antarmuka yang konsisten serta responsif tanpa memerlukan penulisan kode CSS yang kompleks. 

7. REST API 
    

REST API (Representational State Transfer Application Programming Interface) merupakan antarmuka yang digunakan untuk memungkinkan komunikasi dan pertukaran data antara aplikasi yang berbeda melalui jaringan. REST API banyak digunakan dalam pengembangan aplikasi berbasis website untuk menghubungkan bagian antarmuka pengguna dengan layanan yang berada pada sisi server, sehingga aplikasi dapat memperoleh maupun mengirimkan data sesuai kebutuhan (Prasatya, 2026). 

Pada aplikasi GAMON, REST API digunakan sebagai penghubung komunikasi antara frontend yang dibangun menggunakan React dan backend yang menggunakan Go. Melalui REST API, frontend dapat meminta dan mengirimkan data kepada backend, sedangkan backend memproses permintaan tersebut dan memberikan data yang dibutuhkan oleh aplikasi. REST API digunakan dalam berbagai kebutuhan GAMON, seperti pengelolaan data perangkat, pengambilan informasi monitoring, pengelolaan alert, dan pengaturan sistem. 

REST API dipilih karena memungkinkan komunikasi antara frontend dan backend dilakukan melalui antarmuka yang terstruktur. Pemisahan tersebut membuat bagian antarmuka pengguna dan proses pengolahan data dapat dikelola secara terpisah, sehingga sesuai dengan kebutuhan pengembangan GAMON sebagai aplikasi berbasis website. 

8. ICMP Ping 
    

ICMP (Internet Control Message Protocol) merupakan salah satu protokol jaringan yang digunakan untuk membantu proses komunikasi dan diagnostik pada jaringan komputer. Salah satu implementasinya adalah ping, yang digunakan untuk memeriksa apakah suatu perangkat dapat dijangkau melalui jaringan sekaligus mengukur waktu respons (latency) dari perangkat tersebut (Codepolitan, 2023). 

Pada aplikasi GAMON, ICMP Ping digunakan untuk melakukan pemeriksaan kondisi perangkat secara berkala dengan mengecek keterjangkauan perangkat serta memperoleh nilai latency. Hasil pemeriksaan tersebut digunakan untuk menentukan status perangkat, menyimpan riwayat monitoring, dan mendukung proses pemberian alert apabila terdeteksi adanya perubahan kondisi perangkat. 

ICMP Ping dipilih karena merupakan metode yang sederhana, ringan, dan dapat digunakan pada berbagai perangkat yang terhubung ke jaringan, sehingga sesuai dengan kebutuhan GAMON dalam melakukan monitoring ketersediaan perangkat. 

9. Telegram Bot API 
    

Telegram Bot API merupakan antarmuka pemrograman aplikasi (Application Programming Interface/API) yang disediakan oleh Telegram untuk memungkinkan aplikasi berkomunikasi dengan bot Telegram. Melalui Telegram Bot API, aplikasi dapat mengirimkan informasi secara otomatis kepada pengguna sehingga banyak dimanfaatkan sebagai media penyampaian notifikasi pada berbagai sistem (grammY, 2022). 

Pada aplikasi GAMON, Telegram Bot API digunakan sebagai media pengiriman notifikasi kepada administrator mengenai perubahan kondisi perangkat yang sedang dipantau, seperti ketika perangkat mengalami gangguan maupun kembali ke kondisi normal. Dengan adanya notifikasi tersebut, administrator dapat memperoleh informasi mengenai kondisi perangkat tanpa harus terus-menerus memantau halaman aplikasi. Telegram Bot API dipilih karena menyediakan mekanisme pengiriman pesan secara otomatis, mudah diintegrasikan dengan aplikasi, serta mendukung penyampaian informasi secara cepat sehingga sesuai dengan kebutuhan GAMON sebagai sistem monitoring infrastruktur teknologi informasi. 

10. Vite 
    

Vite merupakan build tool modern yang digunakan untuk mendukung proses pengembangan aplikasi web, baik sebagai development server maupun sebagai alat untuk membangun (build) aplikasi agar siap digunakan. Vite dirancang untuk memberikan proses pengembangan yang cepat serta mendukung berbagai teknologi frontend modern seperti React dan TypeScript (BuildWithAngga, 2024). Pada aplikasi GAMON, Vite digunakan untuk mendukung proses pengembangan frontend yang dibangun menggunakan React, TypeScript, dan Tailwind CSS, serta menghasilkan aplikasi yang siap dijalankan setelah proses build. Vite dipilih karena memiliki proses pengembangan yang cepat dengan memanfaatkan toolchain modern berbasis Rust yang dirancang untuk memberikan performa tinggi dan efisiensi. Selain itu, Vite juga mudah diintegrasikan dengan React, TypeScript, dan Tailwind CSS serta memiliki konfigurasi yang sederhana sehingga memudahkan proses pengembangan dan pemeliharaan aplikasi. 

11. Node.js dan NPM 
    

Node.js merupakan runtime JavaScript yang memungkinkan kode JavaScript dijalankan di luar browser, sedangkan NPM (Node Package Manager) merupakan pengelola paket bawaan Node.js yang digunakan untuk menginstal, mengelola dependensi, dan menjalankan skrip pengembangan aplikasi (Faradilla A., 2025; Dicoding Indonesia, 2021). Pada aplikasi GAMON, Node.js digunakan sebagai lingkungan untuk menjalankan berbagai tool pengembangan frontend, seperti Vite dan TypeScript, sedangkan NPM digunakan untuk mengelola pustaka yang dibutuhkan serta menjalankan proses development dan build. Node.js dan NPM dipilih karena menyediakan ekosistem pengembangan yang lengkap, memudahkan pengelolaan dependensi, serta mendukung proses pengembangan frontend secara lebih efisien dan terstruktur. 

12. Git dan Github 
    

Git merupakan Version Control System (VCS), yaitu sistem yang digunakan untuk mencatat, mengelola, dan melacak setiap perubahan pada kode sumber selama proses pengembangan perangkat lunak sehingga riwayat perubahan dapat terdokumentasi dengan baik dan versi sebelumnya dapat dipulihkan apabila diperlukan (Prasatya, 2025). GitHub merupakan platform berbasis cloud yang digunakan untuk menyimpan repositori Git serta memudahkan sinkronisasi dan pengelolaan proyek pengembangan (Faradilla A., 2023). Pada pengembangan aplikasi GAMON, Git digunakan untuk mengelola riwayat perubahan kode selama proses pengembangan, sedangkan GitHub digunakan sebagai repositori daring (remote repository) untuk menyimpan cadangan proyek dan melakukan sinkronisasi kode. Git dan GitHub dipilih karena memudahkan pengelolaan versi aplikasi, menjaga keamanan data melalui penyimpanan terpusat, serta mendukung proses pengembangan yang lebih terstruktur sesuai dengan praktik pengembangan perangkat lunak. 

13. Zed Code Editor 
    

Zed merupakan code editor modern yang dirancang untuk mendukung proses pengembangan perangkat lunak dengan performa yang cepat dan responsif. Zed dikembangkan menggunakan bahasa pemrograman Rust serta memanfaatkan GPU rendering untuk proses tampilan, sehingga mampu memberikan pengalaman penyuntingan kode yang ringan dan efisien (Innocent, 2025). Pada pengembangan aplikasi GAMON, Zed digunakan sebagai code editor untuk menulis, mengelola, dan memodifikasi kode pada sisi frontend maupun backend. Zed dipilih karena memiliki performa yang lebih cepat, penggunaan memori yang lebih ringan dibandingkan code editor berbasis Electron, serta mampu memberikan respons yang lebih baik ketika menangani proyek yang terdiri atas banyak berkas, sehingga mendukung proses pengembangan aplikasi secara lebih nyaman dan produktif. 

3. Analisis Kebutuhan Sistem 

Sebelum proses perancangan dan implementasi sistem dilakukan, diperlukan analisis kebutuhan sistem untuk mengidentifikasi fungsi serta karakteristik yang harus dipenuhi oleh aplikasi GAMON. Analisis ini disusun berdasarkan hasil identifikasi permasalahan, tujuan pengembangan aplikasi, serta kebutuhan pengguna terhadap proses monitoring infrastruktur teknologi informasi. Hasil analisis kebutuhan tersebut menjadi dasar dalam penyusunan rancangan sistem pada bagian selanjutnya sehingga implementasi yang dilakukan sesuai dengan tujuan pengembangan aplikasi. Analisis kebutuhan sistem pada aplikasi GAMON terdiri atas kebutuhan fungsional dan kebutuhan non-fungsional yang dijelaskan pada subbab berikut. 

a. Kebutuhan Fungsional 

Kebutuhan fungsional merupakan kebutuhan yang berkaitan dengan fungsi atau layanan yang harus disediakan oleh aplikasi agar dapat memenuhi tujuan pengembangannya. Kebutuhan ini disusun berdasarkan hasil analisis terhadap proses monitoring serta fitur-fitur yang diperlukan untuk mendukung pengelolaan perangkat dan penyajian informasi secara efektif. Adapun kebutuhan fungsional aplikasi GAMON disajikan pada Tabel 4.1. 

Tabel 4.1 Kebutuhan Fungsional Aplikasi GAMON 

|   |   |   |
|---|---|---|
|No|Kebutuhan Fungsional|Deskripsi|
|FR-1|Registrasi Perangkat|Sistem dapat menambahkan perangkat baru beserta informasi seperti nama, alamat IP, tipe perangkat, metode monitoring, interval pengecekan, lokasi (opsional), dan deskripsi.|
|FR-2|Menampilkan Daftar Perangkat|Sistem dapat menampilkan seluruh perangkat yang telah terdaftar beserta informasi konfigurasinya.|
|FR-3|Mengubah Data Perangkat|Sistem dapat memperbarui konfigurasi perangkat yang telah terdaftar tanpa menghapus data perangkat tersebut.|
|FR-4|Menghapus Perangkat|Sistem dapat menghapus perangkat beserta data terkait, seperti riwayat monitoring dan alert.|
|FR-5|Mengaktifkan atau Menonaktifkan Monitoring|Sistem dapat mengaktifkan maupun menonaktifkan proses monitoring pada perangkat tertentu tanpa menghapus perangkat dari sistem.|
|FR-6|Melakukan Monitoring Berkala|Sistem dapat melakukan pengecekan status perangkat secara otomatis dan berkala menggunakan metode ICMP Ping sesuai interval yang ditentukan.|
|FR-7|Menentukan Status Perangkat|Sistem dapat menentukan status perangkat (online atau offline) serta mengukur nilai latency berdasarkan hasil monitoring.|
|FR-8|Mendeteksi Gangguan Berdasarkan Failure Threshold|Sistem dapat menetapkan perangkat berstatus offline setelah mencapai batas kegagalan pengecekan yang telah dikonfigurasi.|
|FR-9|Menyimpan Riwayat Monitoring|Sistem dapat menyimpan hasil setiap proses monitoring sebagai riwayat untuk kebutuhan penelusuran dan analisis.|
|FR-10|Membuat Alert Otomatis|Sistem dapat membuat alert secara otomatis ketika perangkat terdeteksi mengalami gangguan.|
|FR-11|Menyelesaikan Alert Otomatis|Sistem dapat mengubah status alert menjadi selesai ketika perangkat kembali dalam kondisi online.|
|FR-12|Mengelola Data Alert|Sistem dapat menampilkan, memfilter, melakukan acknowledge, dan menyelesaikan (resolve) alert secara manual.|
|FR-13|Menampilkan Notifikasi Realtime|Sistem dapat menampilkan perubahan status perangkat secara langsung melalui koneksi WebSocket.|
|FR-14|Mengirim Notifikasi Telegram|Sistem dapat mengirimkan notifikasi gangguan melalui Telegram kepada pengguna yang telah melakukan pairing akun.|
|FR-15|Menampilkan Ringkasan Kondisi Sistem|Sistem dapat menampilkan informasi ringkasan pada dashboard, seperti jumlah perangkat, status perangkat, dan alert terbaru.|
|FR-16|Menampilkan Informasi Monitoring|Sistem dapat menampilkan status perangkat secara real-time beserta fitur pencarian, penyaringan, detail perangkat, dan riwayat monitoring.|

b. Kebutuhan Non-Fungsional 

Kebutuhan non-fungsional merupakan kebutuhan yang berkaitan dengan karakteristik dan kualitas sistem dalam mendukung proses operasional aplikasi. Kebutuhan ini tidak menjelaskan fungsi yang disediakan oleh sistem, melainkan aspek-aspek yang memengaruhi kinerja, keandalan, kemudahan penggunaan, serta kemudahan pengembangan aplikasi. Adapun kebutuhan non-fungsional aplikasi GAMON disajikan pada Tabel 4.2. 

Tabel 4.2 Kebutuhan Fungsional Aplikasi GAMON 

|   |   |   |
|---|---|---|
|No|Kebutuhan Fungsional|Deskripsi|
|NFR-1|Performance Efficiency (Efisiensi Performa)|Sistem mampu melakukan proses monitoring perangkat secara berkala dan menyajikan perubahan status secara real-time kepada pengguna dengan penggunaan sumber daya yang efisien.|
|NFR-2|Reliability (Keandalan)|Sistem mampu menjalankan proses monitoring secara berkelanjutan, menyimpan hasil monitoring secara konsisten, serta tetap beroperasi meskipun terjadi kegagalan pada komponen pendukung seperti layanan notifikasi.|
|NFR-3|Availability (Ketersediaan)|Sistem dapat diakses melalui jaringan internal dan menyediakan informasi kondisi perangkat beserta riwayat monitoring kapan pun diperlukan selama aplikasi berjalan.|
|NFR-4|Security (Keamanan)|Sistem mampu menjaga keamanan proses pengiriman notifikasi melalui mekanisme pairing berbasis token sehingga hanya pengguna yang terdaftar dapat menerima notifikasi Telegram.|
|NFR-5|Usability (Kemudahan Penggunaan)|Sistem menyediakan antarmuka berbasis website yang mudah dipahami, dilengkapi indikator visual, fitur pencarian, dan penyaringan data untuk memudahkan proses monitoring.|
|NFR-6|Maintainability (Kemudahan Pemeliharaan)|Sistem dirancang menggunakan struktur modul yang terpisah sehingga memudahkan proses pemeliharaan, pengembangan, dan penambahan fitur baru.|
|NFR-7|Portability (Portabilitas)|Sistem dapat dijalankan pada berbagai sistem operasi yang didukung tanpa memerlukan perubahan pada kode program utama.|

4. Perancangan Sistem 

Perancangan sistem merupakan tahapan yang dilakukan untuk menerjemahkan hasil analisis kebutuhan sistem ke dalam rancangan teknis sebagai acuan pada proses implementasi aplikasi GAMON. Pada tahap ini, setiap kebutuhan fungsional dan nonfungsional yang telah diidentifikasi pada subbab sebelumnya diwujudkan ke dalam model perancangan yang menggambarkan struktur sistem, interaksi pengguna, alur proses, logika kerja, serta rancangan penyimpanan data. Dengan adanya perancangan sistem, proses implementasi dapat dilakukan secara lebih terarah sehingga aplikasi yang dibangun sesuai dengan kebutuhan dan tujuan yang telah ditetapkan. 

a. Gambaran Umum Arsitektur Client & Server 

GAMON dirancang menggunakan arsitektur client dan server, yaitu sebuah pola arsitektur perangkat lunak yang membagi aplikasi menjadi dua pihak utama, yakni client sebagai pihak yang menampilkan data kepada pengguna, dan server sebagai pihak yang menyediakan serta memproses data. Kedua pihak tersebut berjalan secara terpisah dan saling berkomunikasi melalui jaringan komputer menggunakan protokol komunikasi tertentu. 

Pemisahan sisi client dan server memberikan pembagian tanggung jawab yang jelas. Client berfokus pada penyajian antarmuka (user interface) dan penyampaian permintaan pengguna, sedangkan server berfokus pada pemrosesan data, penyimpanan data, serta pengiriman hasil kepada client. Dengan pembagian ini, masing-masing pihak dapat dikembangkan, dijalankan, dan dikelola secara terpisah tanpa bergantung pada sisi yang lain. 

Pada GAMON, pola ini diwujudkan oleh frontend sebagai client dan backend sebagai server. Frontend berjalan di dalam peramban (browser) dan bertanggung jawab menampilkan antarmuka monitoring kepada administrator. Sementara itu, backend bertanggung jawab memproses permintaan, menjalankan proses monitoring, dan mengelola data pada database. Keduanya terhubung melalui jaringan sehingga data hasil monitoring dapat sampai dan ditampilkan pada antarmuka pengguna. 

Seluruh komponen yang membentuk arsitektur GAMON dapat dikelompokkan menjadi beberapa lapisan (tier) berdasarkan perannya masing-masing. Lapisan client terdiri atas peramban dan frontend yang bertugas menampilkan antarmuka kepada pengguna. Lapisan server ditempati oleh backend yang mengandung layanan REST API, WebSocket, dan monitoring engine sebagai pusat pemrosesan sistem. Lapisan data berisi database SQLite sebagai tempat penyimpanan data yang berjalan secara lokal tanpa memerlukan server basis data terpisah. Adapun komponen eksternal meliputi mekanisme pemeriksaan ICMP, perangkat target yang dipantau, serta Telegram Bot API beserta aplikasi Telegram sebagai kanal notifikasi. Hubungan dan aliran komunikasi antarkomponen tersebut digambarkan secara lebih jelas pada diagram arsitektur berikut.