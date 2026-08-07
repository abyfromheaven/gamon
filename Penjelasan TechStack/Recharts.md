# Teknologi: Recharts (Recharts 3.10.1)

## Rangkuman

| Teknologi | Digunakan pada |
|-----------|----------------|
| Recharts | Grafik riwayat latency (waktu tanggap) perangkat pada halaman Monitoring |

---

## 1. Peran Teknologi dalam Project GAMON

Recharts merupakan pustaka (library) pembuatan grafik untuk aplikasi berbasis React. Dalam project GAMON, Recharts digunakan untuk menampilkan data pemantauan dalam bentuk visual, khususnya **grafik riwayat latency** suatu perangkat. Grafik ini membantu pengguna memahami perkembangan waktu tanggap sebuah perangkat dari waktu ke waktu secara lebih mudah dibandingkan hanya membaca angka.

Peran Recharts dalam GAMON antara lain:

| Peran | Penjelasan |
|-------|-----------|
| Menampilkan riwayat latency secara visual | Menggambarkan perubahan nilai latency (dalam milidetik) suatu perangkat sepanjang waktu |
| Mendukung analisis kondisi perangkat | Membantu pengguna melihat pola gangguan atau penurunan kinerja perangkat dari data pemantauan |
| Menyajikan informasi yang ringkas dan informatif | Data yang banyak disajikan dalam satu grafik sehingga mudah dipahami sekilas |

Dalam implementasinya, Recharts dipakai pada komponen `LatencyChart` yang menampilkan data dalam bentuk grafik area (area chart). Grafik ini menampilkan sumbu horizontal untuk waktu dan sumbu vertikal untuk nilai latency, dengan keterangan (tooltip) yang muncul saat pengguna mengarahkan penunjuk pada titik tertentu dalam grafik.

**Bukti penggunaan:**
- `frontend/package.json` — dependensi `recharts ^3.10.1`
- `frontend/src/components/LatencyChart.tsx` — komponen yang menggunakan Recharts untuk menampilkan grafik latency
- Halaman `MonitoringPage.tsx` — memanfaatkan komponen `LatencyChart` untuk menampilkan riwayat perangkat yang dipilih

---

## 2. Alasan Pemilihan Teknologi

Alasan pemilihan Recharts dapat disimpulkan dari kebutuhan dan implementasi project GAMON.

**a. Dibangun khusus untuk React.** Recharts dirancang khusus untuk aplikasi React, sehingga mudah diintegrasikan dengan komponen-komponen React yang sudah ada pada GAMON tanpa memerlukan penyesuaian yang rumit.

**b. Mudah digunakan dan disesuaikan.** Recharts menyediakan komponen grafik yang mudah disusun, seperti sumbu, kisi, keterangan, dan area, yang dapat diatur sesuai kebutuhan. Hal ini memungkinkan GAMON membuat grafik latency sesuai dengan tampilan antarmukanya.

**c. Mendukung tampilan responsif.** Recharts dapat menyesuaikan ukuran grafik terhadap ukuran layar melalui komponen *ResponsiveContainer*, sehingga grafik tetap terlihat baik pada berbagai perangkat.

**d. Cukup untuk kebutuhan visualisasi data monitoring.** Kebutuhan GAMON untuk menampilkan riwayat latency dapat dipenuhi dengan baik oleh Recharts tanpa menambah beban yang berlebihan pada aplikasi.

---

## 3. Cara Kerja Teknologi dalam Konteks Project (Penjelasan Lengkap)

Dalam konteks GAMON, cara kerja Recharts dapat diuraikan sebagai berikut.

**a. Penerimaan data riwayat latency.** Komponen `LatencyChart` menerima data berupa kumpulan nilai yang memuat informasi waktu dan besar latency suatu perangkat. Data ini berasal dari riwayat hasil pemantauan yang tersimpan pada database dan diambil melalui lapisan API aplikasi.

**b. Penyusunan grafik area.** Data tersebut disusun menjadi grafik jenis area. Sumbu horizontal (X) menyatakan waktu, sedangkan sumbu vertikal (Y) menyatakan nilai latency dalam milidetik. Bagian dalam grafik diberi warna dengan tingkat transparansi tertentu agar nilai yang ditampilkan mudah dibaca.

**c. Keterangan dan kisi.** Grafik dilengkapi kisi (grid) sebagai pembantu pembacaan nilai serta keterangan (tooltip) yang menampilkan nilai latency pada saat pengguna mengarahkan penunjuk ke suatu titik. Penyajian ini memudahkan pengguna membaca nilai spesifik dari grafik.

**d. Penyesuaian ukuran.** Grafik dibungkus dalam *ResponsiveContainer* sehingga ukurannya menyesuaikan diri dengan ruang tampilan yang tersedia pada halaman. Ketika pengguna memilih perangkat lain pada halaman Monitoring, data riwayat yang ditampilkan diganti dengan data perangkat yang baru dipilih, dan grafik diperbarui menampilkan kondisi terbaru.

Dengan cara ini, Recharts membantu mengubah data pemantauan yang berupa angka menjadi informasi visual yang mudah dipahami pengguna.

---

## 4. Bukti Penggunaan

1. **File `frontend/package.json`** — Mencantumkan dependensi `recharts ^3.10.1` pada kebutuhan utama aplikasi *frontend*.

2. **File `frontend/src/components/LatencyChart.tsx`** — Komponen pembuat grafik latency yang memanfaatkan komponen Recharts, seperti `AreaChart`, `XAxis`, `YAxis`, `CartesianGrid`, `Tooltip`, dan `ResponsiveContainer`.

3. **File `frontend/src/pages/MonitoringPage.tsx`** — Halaman Monitoring yang menggunakan komponen `LatencyChart` untuk menampilkan riwayat latency perangkat yang dipilih pengguna.

---

## Catatan

- Recharts merupakan pustaka grafik yang khusus digunakan pada bagian *frontend* dan hanya berperan untuk visualisasi data; proses pengambilan dan penyimpanan data tetap dilakukan oleh komponen lain.
- Grafik yang dihasilkan adalah grafik area (area chart) untuk menampilkan perubahan latency perangkat sepanjang waktu.