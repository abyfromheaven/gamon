# BAB 1 - PENDAHULUAN
## A. Latar Belakang Pemilihan Judul
1. [x] Perkembangan Teknologi Informasi #REFERENSI 
	   Perlu sumber referensi untuk klaim mengenai perkembangan TI, digitalisasi, ketergantungan organisasi terhadap TI, dan pentingnya availability.
2. [x] Pentingnya Monitoring Infrastruktur TI #REFERENSI 
	Perlu sumber untuk konsep monitoring, fungsi monitoring, deteksi gangguan, dan pentingnya monitoring infrastruktur.
## B. Metode Penyusunan dan Pengumpulan Data
1. [x] Pendahuluan Metode, Observasi, Wawancara, Studi Literatur #REFERENSI 
# BAB 2 - TRANSFORMASI PKL
1. [x] Paragraf Pembuka BAB 2 #REFERENSI 
       Karena ada klaim mengenai PKL sebagai bagian dari pembelajaran SMK/Kurikulum Merdeka, gunakan sumber resmi pemerintah.
2. [x] Bagian Tujuan/Manfaat PKL #REFERENSI 
       Karena ada klaim mengenai PKL sebagai bagian dari pembelajaran SMK/Kurikulum Merdeka, gunakan sumber resmi pemerintah. 
3. [x] Bagian Tujuan/Manfaat PKL #REVISI-FOKUS 
       Jangan terlalu banyak teori pendidikan. BAB II tetap laporan PKL, bukan makalah tentang pendidikan vokasi.
# BAB 3 - TINJAUAN UMUM INSTANSI
1. [ ] Prinsip BAB 3 #REVISI-FOKUS 
       Fokus pada:
		- identitas instansi;
		- sejarah;
		- visi dan misi;
		- struktur organisasi;
		- bidang/kegiatan instansi;
		- unit yang relevan dengan PKL jika memang diperlukan.
2. [ ] Sumber Referensi BAB 3 #REFERENSI 
       gunakan sumber resmi instansi apabila tersedia.
# BAB 4 - PEMBAHASAN MASALAH
## A. Pembahasan Masalah yang Mengacu pada Judul
1. [ ] Kondisi Permasalahan #REVISI-FOKUS 
       Fokus: "Apa kondisi yang ditemukan?" (Berisikan fakta hasil observasi)
       Jangan menjelaskan penyebab secara terlalu jauh, jangan menjelaskan solusinya juga.
2. [ ] Analisis Permasalahaan #REVISI-FOKUS 
       Fokus: "Mengapa kondisi tersebut menjadi masalah?" Di sini boleh mulai melakukan analisis terhadap proses manual, banyaknya perangkat, perbedaan antarmuka, dan sebagainya. beberapa klaim pernyataan bisa ditambahkan referensi (jika ada), tetapi analisis utamanya berdasarkan hasil observasi selama pkl.
3. [ ] Dampak Permasalahan #REVISI-FOKUS 
       Fokus: "Apa akibatnya?" jangan mengulang poin 4.2 (Analisis Permasalahan). Contoh: masalah → dampak terhadap pemantauan → dampak terhadap respons gangguan. Untuk Referensi tidak wajib jika dampaknya merupakan hasil analisis berdasarkan observasi, jika menggunakan klaim ilmiah/kuantitatif baru wajib sumber.
4. [ ] Solusi yang Diusulkan #REVISI-FOKUS 
       Ini cukup menjelaskan: "GAMON diusulkan sebagai prototype sistem monitoring berbasis website untuk mengatasi permasalahan tersebut." Kemudian sebutkan secara singkat terkait montoring terpusat, device management, alert notification, dan history. jangan membahas tentang Golang, React, ICMP, WebSocket, Database ataupun algoritma alert, itu nanti ada bagiannya tersendiri.
## B. Implementasi Kegiatan PKL
1. [ ] Gambaran Umum Aplikasi #REVISI-FOKUS 
       Ini harus dipangkas, cukup menjawab:
       Apa itu Gamon? Apa tujuan utamanya? Apa ruang lingkupnya? Siapa pengguna sistemnya? Jangan memasukkan mekanisme teknis, icmp secara panjang, websocket, failure threshold, telegram pairing ataupun arsitektur lainnya. Definisi umum seperti "sistem monitoring infrastruktur TI" bisa diberi sumber jika mengambil definisi dari literatur. Tetapi untuk Gamon tidak perlu sumber karena memang berasal dari project sendiri.
2. [ ] Teknologi yang Digunakan #REVISI-FOKUS 
       Standarnya cukup: "Apa Teknologi ini?", "Digunakan untuk apa di GAMON?", dan "Mengapa dipilih?". Jangan membahas cara kerja internal secara panjang.
3. [ ] Golang #REVISI #REFERENSI 
       Kurangi pembahasan mengenai goroutine/concurrency, cukup satu kalimat mengenai concurrency sebagai alasan pemilihan.
       Untuk referensinya perlu terkait "Definisi Go" dan "Concurrency/Goroutine".
4. [ ] REST API #REVISI #REFERENSI 
       Kurangi penjelasan HTTP Method dan mekanisme request-response. Cukup definisi + peran dalam gamon + alasannya.
       Untuk Referensi: "Definisi Rest API"
5. [ ] SQLite #REVISI #REFERENSI 
       Jangan masuk WAL, Connection pool, migration dan detai internal.
       Fokus ke Definisi dan alasan penggunaan
6. [ ] TypeScript #REVISI #REFERENSI 
       Perlu di pangkas, tidak perlu menjelaskan .tsx, compiler, konfigurasi dan sebagainya, fokus menjelaskan definisi, peran dalam gamon dan alasannya memilihnya.
7. [ ] React #REFERENSI #REVISI 
       Pendekkan konsep component-based, tidak perlu menjelaskan SPA secara panjang. fokus menjelaskan definisi, peran dalam gamon dan alasannya memilihnya.
8. [ ] Tailwind CSS #REFERENSI #REVISI 
       Satu paragraf sebenarnya cukup, fokus pada utility-first, digunakan untuk UI Gamon dan alasan pemilihannya.
9. [ ] ICMP Ping #REVISI #REFERENSI 
       Jangan menjelaskan Echo Request/Echo Reply terlalu panjang. Cukup fokus ke kalau digunakan untuk mengecek keterjangkauan perangkat dan latency.
10. [ ] WebSocket #REFERENSI #REVISI 
        Jangan masuk ke hub, readpump ataupun writepump. cukup fokus di > WebSocket → komunikasi real-time → peran di dalam GAMON → Gorilla WebSocket → Alasan Pemilihannya.
11. [ ] Telegram Bot API #REFERENSI #REVISI 
        Jangan membahas pairing token disini, cukup membahas > Telegram Bot API → media notifikasi → digunakan GAMON → alasan pemilihan.