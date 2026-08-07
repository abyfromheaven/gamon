# Teknologi: Telegram Bot API

## Rangkuman

| Teknologi | Digunakan pada |
|-----------|----------------|
| Telegram Bot API | Media penyampaian notifikasi perangkat: pemberitahuan gangguan dan pemberitahuan perangkat kembali normal |

---

## 1. Peran Teknologi dalam Project GAMON

Telegram Bot API merupakan antarmuka yang disediakan oleh layanan Telegram untuk menghubungkan sebuah bot dengan aplikasi. Dalam project GAMON, teknologi ini dimanfaatkan sebagai **media penyampaian notifikasi** kepada pengguna mengenai kondisi perangkat yang dipantau.

Keberadaan fitur ini didasari oleh kenyataan bahwa pengguna tidak selalu berada di depan layar aplikasi. Melalui Telegram, pengguna dapat menerima informasi kondisi perangkat secara langsung pada perangkat yang dibawanya, sehingga gangguan pada perangkat dapat diketahui lebih awal dan segera ditindaklanjuti.

Peran Telegram Bot API dalam GAMON antara lain:

| Peran | Penjelasan |
|-------|-----------|
| Memberitahu adanya gangguan | Menyampaikan kepada pengguna bahwa sebuah perangkat tidak merespons pemeriksaan |
| Memberitahu perangkat kembali normal | Menyampaikan kepada pengguna bahwa perangkat telah pulih |
| Menghubungkan akun pengguna | Menjembatani hubungan antara aplikasi dan akun Telegram pengguna yang ingin menerima notifikasi |

Hal penting yang perlu ditegaskan adalah bahwa Telegram pada GAMON **hanya bertugas sebagai media pengantar informasi**. Seluruh proses pemantauan dan penentuan kondisi perangkat tetap dilakukan oleh komponen monitoring pada aplikasi. Telegram hanya menerima hasil yang telah diproses kemudian menyampaikannya kepada pengguna. Dengan demikian, gangguan pada layanan Telegram tidak akan mengganggu berjalannya pemantauan utama.

**Bukti penggunaan:** aplikasi GAMON memiliki bagian yang bertanggung jawab untuk menyusun dan mengirim pesan notifikasi, menerima perintah dari pengguna pada bot, serta mengelola hubungan akun Telegram (terlihat pada struktur program dan tabel data koneksi Telegram).

---

## 2. Alasan Pemilihan Teknologi

Alasan pemilihan Telegram Bot API dapat disimpulkan dari dokumentasi perancangan project GAMON.

**a. Gratis dan tidak memerlukan biaya tambahan.** Layanan bot Telegram tersedia tanpa biaya dan tidak membutuhkan infrastruktur tambahan, sehingga sesuai untuk aplikasi berukuran menengah seperti GAMON.

**b. Mudah dihubungkan.** Bot Telegram dapat dikaitkan dengan aplikasi melalui antarmuka yang relatif sederhana, sehingga proses penyambungan ke dalam sistem tidak rumit.

**c. Mendukung penyampaian informasi secara cepat.** Pesan dapat sampai kepada pengguna segera setelah terjadi perubahan kondisi perangkat, sehingga gangguan dapat diketahui secara dini.

**d. Mendukung pemantauan dari mana saja.** Karena Telegram dapat diakses melalui perangkat yang dibawa pengguna, pengguna mendapat kabar kondisi perangkat tanpa harus berada di tempat pemantauan.

**e. Tidak mengganggu proses utama.** Karena peran Telegram hanya sebagai media penyampaian, kegagalan pada layanan ini tidak akan menghambat proses pemantauan utama aplikasi.

---

## 3. Cara Kerja Teknologi dalam Konteks Project (Penjelasan Lengkap)

Dalam konteks GAMON, cara kerja Telegram dalam menyampaikan notifikasi dapat diuraikan berikut ke beberapa tahap.

**a. Penerima data pemantauan yang akan diinformasikan.** Seluruh kegiatan pemantauan perangkat terus berjalan di dalam sistem. Ketika sistem menentukan bahwa kondisi sebuah perangkat berubah, sistem menyiapkan informasi yang menjelaskan kondisi tersebut sebagai bahan untuk disampaikan kepada pengguna.

**b. Penyiapan pesan notifikasi.** Informasi perubahan kondisi tersebut diolah menjadi sebuah pesan yang mudah dipahami, memuat keterangan tentang perangkat yang dimaksud serta status kejadiannya. Apabila sebuah perangkat bermasalah, pengguna akan memperoleh informasi bahwa perangkat tersebut tidak merespons; apabila perangkat sudah kembali normal, pengguna akan memperoleh pemberitahuan yang sesuai.

**c. Penyampaian melalui akun Telegram yang terhubung.** Pesan tersebut dikirim kepada akun Telegram pengguna yang sebelumnya telah dipasangkan dengan sistem. Dengan demikian, informasi sampai kepada pengguna secara langsung tanpa harus pengguna membuka aplikasi terlebih dahulu.

**d. Pencegahan pesan berulang.** Sistem hanya mengirim notifikasi pada saat terjadi perubahan kondisi perangkat. Jika kondisi tidak berubah, tidak ada pesan yang terkirim berulang-ulang, sehingga pengguna tidak terganggu oleh notifikasi yang sama secara terus-menerus.

**e. Penerimaan perintah dari pengguna.** Pengguna dapat berinteraksi dengan bot melalui perintah sederhana untuk menghubungkan akunnya atau memeriksa status hubungan. Sistem membaca perintah tersebut dan memberikan tanggapan berupa informasi yang dibutuhkan pengguna.

Dengan alur tersebut, Telegram Bot API menjadi media komunikasi dua arah yang menghubungkan sistem kepada pengguna, sehingga informasi kondisi perangkat dapat disampaikan secara langsung dan tepat waktu.

---

## 4. Bukti Penggunaan

1. **Modul pengiriman notifikasi** — bagian program yang berperan menyusun dan mengirim pesan pemberitahuan gangguan serta pemberitahuan perangkat kembali normal.

2. **Modul penerimaan perintah** — bagian yang membaca dan menanggapi perintah yang dikirim pengguna kepada bot.

3. **Modul pengelola koneksi** — bagian yang mengatur token hubungan akun Telegram dan mencatat status koneksi.

4. **Data penyimpanan pada database** — terdapat bagian data yang menyimpan informasi hubungan akun Telegram pada aplikasi.

5. **Catatan perancangan** — dokumentasi yang menguraikan alur dan aturan fungsi notifikasi Telegram pada project GAMON.

---

## Catatan

- Telegram pada GAMON hanya berfungsi sebagai media penyampaian informasi; proses pemantauan dan penyimpanan kondisi tetap dilakukan oleh komponen utama aplikasi.
- Fitur notifikasi melalui Telegram bersifat opsional. Apabila tidak diaktifkan, proses pemantauan perangkat tetap dapat berjalan seperti biasa.