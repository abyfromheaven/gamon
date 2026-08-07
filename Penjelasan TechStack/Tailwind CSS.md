# Teknologi: Tailwind CSS (Tailwind CSS 4.3.3)

## Rangkuman

| Teknologi | Digunakan pada |
|-----------|----------------|
| Tailwind CSS | Seluruh tampilan antarmuka pengguna (styling) aplikasi GAMON |

---

## 1. Peran Teknologi dalam Project GAMON

Tailwind CSS merupakan kerangka kerja (*framework*) CSS yang digunakan untuk mengatur seluruh gaya tampilan (styling) pada antarmuka pengguna GAMON. Tailwind bekerja dengan pendekatan *utility-first*, yaitu menyediakan sekumpulan kelas siap pakai yang langsung diterapkan pada elemen untuk mengatur jarak, warna, ukuran, dan aspek penampilan lainnya. Dengan cara ini, tampilan aplikasi dapat diatur langsung di dalam kode antarmuka tanpa harus menulis banyak aturan CSS terpisah.

Dalam GAMON, Tailwind CSS berperan penting dalam beberapa hal.

| Peran | Penjelasan |
|-------|-----------|
| Menata tampilan seluruh halaman dan komponen | Dipakai untuk mengatur tata letak, warna, ukuran teks, jarak, dan responsivitas pada setiap halaman dan komponen |
| Menerapkan tema visual yang konsisten | Memakai turunan dari Tailwind, yaitu design token pada `index.css`, agar seluruh antarmuka memiliki warna yang seragam dan sejalan dengan identitas aplikasi |
| Mendukung tampilan responsif | Utility kelas praktis untuk menyesuaikan tampilan pada berbagai ukuran layar |

Tailwind CSS dalam project ini tidak hanya digunakan secara bawaan, tetapi dikustomisasi untuk menyesuaikan kebutuhan aplikasi. GAMON menerapkan tema bernama "Warm Dark Dashboard" dengan palet warna gelap hangat dan aksen oranye, yang didefinisikan melalui sistem tema Tailwind pada `frontend/src/index.css`. Tampilan GAMON juga memanfaatkan keluarga huruf (font) bernama **Geist** yang dipasang sebagai aset lokal.

**Bukti penggunaan:**
- `frontend/package.json` — dependensi `tailwindcss` dan `@tailwindcss/vite`
- `frontend/vite.config.ts` — plugin Tailwind untuk Vite
- `frontend/src/index.css` — impor `@import "tailwindcss"` dan definisi tema
- Kelas-kelas utilitas Tailwind yang tersebar pada seluruh halaman dan komponen

---

## 2. Alasan Pemilihan Teknologi

Alasan pemilihan Tailwind CSS dapat disimpulkan dari kebutuhan dan implementasi project GAMON.

**a. Kecepatan dalam membangun antarmuka.** Pendekatan *utility-first* memungkinkan tampilan dibuat langsung melalui kelas-kelas yang siap pakai, sehingga mempercepat proses pembangunan *frontend* dan mengurangi kebutuhan menulis file CSS khusus untuk tiap elemen.

**b. Konsistensi antarmuka.** Dengan tema yang terdefinisi pada satu tempat, seluruh komponen memakai warna, jarak, dan gaya yang sama. Hal Ini menjaga tampilan GAMON tetap seragam antar halaman.

**c. Tampilan yang ringan dan dapat dikustomisasi.** Tailwind hanya menghasilkan *style* yang benar-benar digunakan, sehingga hasil akhir aplikasi ringan. Tema juga mudah dikustomisasi menyesuaikan identitas visual aplikasi.

**d. Integrasi sederhana dengan Vite.** Tailwind dihubungkan langsung dengan Vite melalui plugin resmi, sehingga proses pembangunan aplikasi *frontend* menjadi lebih mudah dan terintegrasi.

---

## 3. Cara Kerja Teknologi dalam Konteks Project (Penjelasan Lengkap)

Dalam konteks GAMON, cara kerja Tailwind CSS dapat diurai sebagai berikut.

**a. Pemasangan langsung pada tampilan dengan kelas utilitarian.** Setiap elemen pada halaman dan komponen diberi kelas-kelas Tailwind tertentu untuk mengatur tampilannya. Misalnya, kelas untuk mengatur warna latar, ukuran teks, jarak antar elemen, dan tata letak menggunakan susunan kolom dan baris. Proses ini membuat tampilan dapat dirancang secara langsung dalam kode antarmuka.

**b. Pendefinisian tema**

Tema khusus aplikasi didefinisikan pada `index.css` dengan dukungan sintaks tema pada Tailwind. Tema ini berisi sejumlah warna dan huruf yang menjadi identitas visual, seperti warna untuk latar, permukaan, aksen, keadaan sukses, peringatan, dan bahaya. Dependansi warna ini kemudian direferensikan dalam kelas-kelas Tailwind pada komponen, sehingga warna yang dipakai di seluruh aplikasi mengikuti tema yang sama.

**c. Sampling cepat**

 Karena Tailwind melibatkan pengolah khusus saat membangun aplikasi, hanya kelas-kelas yang dipergunakan yang benar-benar dihasilkan ke dalam file CSS akhir. Hasilnya, ukuran file gaya menjadi lebih efisien di homepage, dan tampilan tetap lengkap.

**d. Tampilan responsif dan indah.**

Tailwind menyediakan kelas untuk menyesuaikan tampilan pada berbagai ukuran layar. Pada GAMON, hal ini menghasilkan tampilan yang dapat menyesuaian diri, misal saat perangkat dibuka pada layar komputer versus layar ponsel.

---

## 4. Kunjungan Pemakaian

1. **File `frontend/package.json`** — Mencantumkan `tailwindcss@^4.3.3` dan plugin `@tailwindcss/vite@^4.3.3` sebagai ketergantungan pengembangan.

2. **File `frontend/vite.config.ts`** — Menambahkan plugin Tailwind CSS ke dalam konfigurasi Vite.

3. **File `frontend/src/index.css`** — Diawali dengan `@import "tailwindcss"`, memuat tema dengan warna dan huruf *Geist*.

4. **Folder `frontend/src/pages/` dan `frontend/src/components/`** — Seluruh halaman dan komponen memakai kelas-kelas utilitas Tailwind untuk mengatur tampilan.

---

## Catatan

- Tailwind CSS digunakan bersama Vite sebagai kerangka *styling*, dengan pengaturan utama pada file `index.css` dan `vite.config.ts`.
- Desain tampilan GAMON dirancang dengan tema gelap hangat (Warm Dark Dashboard) yang dibuat melalui kustomisasi Tailwind.