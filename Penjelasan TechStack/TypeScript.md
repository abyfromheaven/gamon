# Teknologi: TypeScript (TypeScript ~6.0.2)

## Rangkuman

| Teknologi | Digunakan pada |
|-----------|----------------|
| TypeScript | Seluruh pengembangan antarmuka pengguna (frontend) GAMON |

---

## 1. Peran Teknologi dalam Project GAMON

TypeScript merupakan bahasa pemrograman yang menjadi dasar penulisan seluruh kode pada sisi *frontend* aplikasi GAMON. TypeScript dapat dipahami sebagai pengembangan dari bahasa JavaScript yang dilengkapi dengan sistem tipe data, sehingga penulisan kode menjadi lebih terstruktur dan terkendali. Dalam project ini, TypeScript digunakan bersama dengan React untuk membangun seluruh halaman, komponen, serta logika antarmuka yang ditampilkan kepada pengguna melalui peramban.

Peran utama TypeScript pada GAMON adalah menjaga konsistensi bentuk data yang mengalir di dalam aplikasi. Aplikasi ini banyak menerima data dari server, seperti data perangkat, hasil pemantauan, dan daftar peringatan, yang umumnya berbentuk data terstruktur. Dengan adanya sistem tipe, GAMON dapat memastikan bahwa setiap data yang masuk dan keluar memiliki bentuk yang sesuai, sehingga potensi kesalahan yang disebabkan oleh data yang tidak cocok dapat dikurangi sejak awal pengembangan.

Selain itu, TypeScript juga berperan dalam meningkatkan keterbacaan dan kemudahan pemeliharaan kode. Ketika sebuah fungsi atau komponen memiliki tipe yang dinyatakan secara jelas, siapa pun yang membaca kode dapat memahami bentuk data yang diterima dan dihasilkan tanpa harus menelusuri seluruh alur implementasinya. Hal ini sangat membantu pada aplikasi yang memiliki banyak halaman dan komponen seperti GAMON.

**Bukti penggunaan:** seluruh kode *frontend* ditulis dalam file berformat TypeScript (`.ts` dan `.tsx`), mencakup halaman, komponen, *hook*, lapisan pemanggilan API, dan definisi tipe data pada folder `frontend/`.

---

## 2. Alasan Pemilihan Teknologi

Pemilihan TypeScript didasarkan pada pertimbangan yang berkaitan dengan kualitas dan keandalan perangkat lunak, sebagai berikut.

**a. Kepastian bentuk data (type safety).** GAMON bertukar banyak data terstruktur dengan server, seperti data perangkat, peringatan, dan hasil pemantauan. TypeScript memungkinkan bentuk data tersebut dinyatakan secara tegas, sehingga kesalahan akibat data yang bentuknya tidak sesuai dapat dideteksi lebih awal. Hal ini sejalan dengan prinsip rekayasa perangkat lunak yang menekankan pencegahan kesalahan daripada perbaikannya.

**b. Kemudahan pemeliharaan.** Dengan tipe data yang dinyatakan secara jelas pada setiap fungsi dan komponen, kode antarmuka yang jumlahnya banyak menjadi lebih mudah dipahami, dirawat, dan dikembangkan oleh pengembang.

**c. Integrasi yang baik dengan React.** TypeScript terintegrasi dengan baik dengan React, sehingga seluruh komponen antarmuka dapat ditulis dengan dukungan tipe yang kuat tanpa mengurangi kemampuan React dalam membangun tampilan yang dinamis.

**d. Pencegahan kesalahan pada tahap pembangunan.** Melalui pemeriksaan tipe yang dilakukan pada proses pembangunan aplikasi, kesalahan akibat ketidakcocokan data dapat teridentifikasi sebelum aplikasi dijalankan. Pendekatan ini mengurangi kemungkinan gangguan yang muncul ketika aplikasi sudah digunakan.

---

## 3. Cara Kerja Teknologi dalam Konteks Project (Penjelasan Lengkap)

Dalam konteks GAMON, TypeScript bekerja dengan cara menyediakan seperangkat aturan yang memastikan seluruh data dan fungsi pada aplikasi memiliki bentuk yang konsisten.

Pertama, aplikasi mendefinisikan seluruh bentuk data yang digunakan pada satu tempat, misalnya data perangkat, peringatan, dan hasil pemantauan. Definisi ini menjadi acuan bersama bagi seluruh halaman dan komponen. Ketika server mengirimkan data, aplikasi mengetahui bentuk data yang seharusnya diterima; ketika aplikasi mengirimkan data, aplikasi juga mengetahui bentuk data yang seharusnya dikirim. Dengan demikian, komunikasi antara aplikasi dan server menjadi lebih terkendali.

Kedua, pada saat menulis tampilan antarmuka, setiap komponen dinyatakan dengan jelas jenis data yang diterimanya dan dikembalikannya. Hal ini membuat komponen-komponen tersebut dapat disusun dan digunakan kembali dengan aman, karena hubungan antar komponen sudah terjamin oleh sistem tipe.

Ketiga, karena peramban tidak dapat menjalankan TypeScript secara langsung, seluruh kode tersebut terlebih dahulu diubah menjadi JavaScript pada saat aplikasi dibangun. Dengan demikian, hasil akhir yang berjalan pada peramban adalah JavaScript, sementara kode yang dikembangkan dan dipelihara adalah TypeScript yang lebih aman dan terstruktur.

Melalui ketiga tahapan tersebut, TypeScript berperan sebagai fondasi yang menjaga kualitas dan keteraturan seluruh kode *frontend* GAMON, sekaligus memastikan bahwa aplikasi dapat dikembangkan dan dipelihara secara berkelanjutan.

---

## 4. Bukti Penggunaan

1. **File `frontend/package.json`** — mencantumkan TypeScript sebagai salah satu kebutuhan pengembangan *frontend*.

2. **File konfigurasi TypeScript** — terdapat `tsconfig.json`, `tsconfig.app.json`, dan `tsconfig.node.json` yang mengatur cara kerja TypeScript pada project.

3. **Folder `frontend/src/`** — seluruh isi kode *frontend*, termasuk halaman, komponen, *hook*, dan lapisan API, ditulis dalam format TypeScript.

4. **File `frontend/src/types/index.ts`** — tempat pendefinisian bentuk data aplikasi (data perangkat, peringatan, dan lainnya).

5. **Folder `frontend/src/pages/` dan `frontend/src/components/`** — seluruh halaman dan komponen antarmuka ditulis dalam format yang mengandung tampilan React dan didefinisikan dengan tipe data.

---

## Catatan

- TypeScript merupakan pengembangan dari JavaScript. Kode yang ditulis dalam TypeScript diubah menjadi JavaScript pada proses pembangunan agar dapat dijalankan oleh peramban.
- Penggunaan TypeScript pada project ini bertujuan untuk meningkatkan kualitas, keteraturan, dan kemudahan pemeliharaan kode antarmuka.