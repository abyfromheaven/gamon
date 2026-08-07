# Teknologi: Vite (Vite 8.1.1)

## Rangkuman

| Teknologi | Digunakan pada |
|-----------|----------------|
| Vite | Pembangunan dan pengelolaan sisi frontend GAMON (menjalankan server pengembangan, mengubah kode menjadi produk akhir, dan menghubungkan berbagai pustaka depan) |

---

## 1. Peran Teknologi dalam Project GAMON

Vite merupakan perangkat pembangunan (build tool) yang digunakan pada sisi *frontend* aplikasi GAMON. Perangkat ini bertanggung jawab atas proses pengelolaan dan penyusunan seluruh kode *frontend* agar siap dijalankan maupun digunakan sebagai produk jadi.

Dalam proses pengembangan, Vite berperan penting dalam beberapa hal.

| Peran | Penjelasan |
|-------|-----------|
| Menjalankan server pengembangan | Menyediakan lingkungan untuk menjalankan aplikasi selama proses pembuatan, sehingga perubahan kode dapat langsung terlihat |
| Menggabungkan dan mengolah kode | Mengubah kode-kode yang dikembangkan menjadi bentuk yang siap dijalankan oleh peramban |
| Menghubungkan berbagai teknologi frontend | Mengaitkan pustaka dan kerangka kerja yang dipakai, seperti React dan Tailwind CSS, ke dalam satu proses yang utuh |
| Menyiapkan produk akhir | Menghasilkan aplikasi jadi yang siap digunakan setelah proses pembangunan selesai |

Dengan demikian, Vite menjadi penopang seluruh alur pengembangan *frontend*, mulai dari tahap mengembangkan sampai tahap menghasilkan aplikasi yang siap dipakai.

**Bukti penggunaan:**
- Aplikasi *frontend* memiliki pengaturan untuk menjalankan proses pengembangan, pembangunan, dan pengecekan
- Alat yang menghubungkan React dan kerangka styling ke dalam aplikasi bergantung pada Vite
- Beberapa perintah untuk menjalankan aplikasi diatur melalui skrip yang memakai Vite

---

## 2. Alasan Pemilihan Teknologi

Alasan pemilihan Vite dapat disimpulkan dari kebutuhan dan implementasi project GAMON.

**a. Memberikan proses pengembangan yang cepat.** Vite dirancang untuk memberikan pengalaman pengembangan yang responsif, sehingga perubahan pada kode dapat langsung dinikmati oleh pengembang tanpa proses yang berbelit.

**b. Terintegrasi dengan alat dan pustaka yang dipakai.** Vite berhubungan baik dengan React dan kerangka styling yang digunakan pada *frontend*, sehingga seluruh kebutuhan pembangunan aplikasi dapat diurus dalam satu alat.

**c. Menghasilkan aplikasi yang ringan dan sesuai.** Vite hanya menggabungkan bagian yang benar-benar dipakai, sehingga hasil produk *frontend* menjadi efisien dan siap digunakan.

**d. Proses pengaturan yang sederhana.** Pengaturan Vite pada project tergolong sederhana, sehingga mudah dipahami dan dirawat selama pengembangan aplikasi.

---

## 3. Cara Kerja Teknologi dalam Konteks Project (Penjelasan Lengkap)

Dalam konteks GAMON, cara kerja Vite dapat diuraikan sebagai berikut.

**a. Menjalankan aplikasi pada tahap pengembangan.** Ketika pengembang menjalankan mode pengembangan, Vite membuka sebuah layanan yang menampilkan aplikasi di peramban. Selama tahap ini, perubahan yang dilakukan pada kode *frontend* dapat segera tampil tanpa proses yang rumit, sehingga pengelola perkembangan berjalan lebih lancar.

**b. Menghubungkan bahasa dan pustaka yang dipakai.** Kode *frontend* pada GAMON ditulis menggunakan aspek-aspek React, TypeScript, dan kerangka styling. Vite berperan menghubungkan seluruh aspek ini sehingga kode yang dikembangkan dapat diproses dengan benar menjadi satu aplikasi yang utuh.

**c. Menyiapkan produk akhir.** Pada tahap pembangunan untuk penggunaan, Vite mengolah seluruh kode menjadi bentuk yang siap dijalankan oleh peramban. Proses ini mencakup pemeriksaan tipe data dan penyusunan kode sehingga aplikasi dapat dijalankan dengan baik oleh pengguna akhir.

**d. Perintah menjalankan aplikasi.** Proses pengembangan dan pembangunan aplikasi dijalankan melalui perintah yang ditentukan pada pengaturan aplikasi. Perintah-perintah ini memanggil kemampuan Vite untuk menjalankan atau menyiapkan aplikasi sesuai kebutuhan.

Dengan alur tersebut, Vite menyatukan seluruh proses penyusunan *frontend* GAMON secara efisien dan teratur.

---

## 4. Bukti Penggunaan

1. **Perintah skrip aplikasi.** Pada pengaturan aplikasi *frontend*, terdapat skrip yang memanggil Vite untuk menjalankan proses pengembangan (`dev`), proses pembangunan akhir (`build`), dan proses menampilkan hasil (preview).

2. **Pengaturan alat Vite.** Terdapat file pengaturan khusus Vite yang mengonfigurasi cara Vite menghubungkan React dan kerangka styling ke dalam aplikasi.

3. **Ketergantungan pada perpaket aplikasi.** Pada daftar kebutuhan aplikasi *frontend*, Vite tercantum sebagai salah satu alat yang digunakan selama pengembangan.

4. **Folder kode *frontend*.** Seluruh folder kode *frontend* disusun dan diproses dalam kerangka kerja yang diatur oleh Vite.

---

## Catatan

- Vite berperan pada sisi *frontend* saja dan mengelola proses pengembangan serta pembangunan aplikasi.
- Vite dihubungkan dengan React dan kerangka styling sehingga seluruh proses pembuatan antarmuka berjalan dalam satu alat yang utuh.