# Monitoring Page

## Overview

Halaman Monitoring merupakan halaman utama yang digunakan untuk melihat kondisi seluruh perangkat yang sedang dimonitor oleh Gamon secara terpusat. Halaman ini berfokus pada monitoring status perangkat, bukan pada pengelolaan perangkat ataupun analisis performa yang kompleks.

Tujuan utama dari halaman Monitoring adalah membantu pengguna mengetahui dengan cepat perangkat mana yang sedang Online, Offline, mengalami Warning, ataupun tidak sedang dimonitor (Unknown).

Halaman ini dirancang untuk menjawab beberapa pertanyaan berikut:

- Perangkat mana yang sedang mengalami gangguan?
- Berapa banyak perangkat yang sedang Online, Offline, atau Warning?
- Kapan terakhir kali perangkat dilakukan pengecekan?
- Bagaimana kondisi perangkat saat ini?
- Informasi dasar apa saja yang dimiliki oleh perangkat tersebut?

> Halaman Monitoring tidak ditujukan untuk melakukan konfigurasi perangkat ataupun menampilkan analisis performa yang kompleks seperti CPU Usage, RAM Usage, Network Traffic, maupun grafik monitoring jangka panjang. Fokus utamanya adalah monitoring status perangkat secara sederhana dan terpusat.


---

## Objectives

- Menampilkan kondisi seluruh perangkat yang sedang dimonitor oleh Gamon.
- Membantu pengguna mengetahui perangkat yang mengalami gangguan dengan cepat.
- Menyediakan informasi monitoring perangkat secara sederhana dan mudah dipahami.
- Menjadi pusat monitoring status perangkat yang telah didaftarkan pada Gamon.


---

## Monitoring Mechanism

Setiap perangkat yang berstatus **Active** akan dilakukan monitoring secara berkala sesuai dengan Check Interval yang telah ditentukan pada Device Management.

Gamon akan melakukan pengecekan terhadap perangkat menggunakan Monitoring Method yang dipilih (misalnya Ping atau HTTP Request), kemudian memperbarui status monitoring perangkat secara otomatis.

Apabila perangkat berhasil dilakukan pengecekan, maka status perangkat akan berubah menjadi:

- Online → Perangkat berhasil diakses dan berjalan dengan normal.
- Offline → Perangkat gagal diakses atau mengalami gangguan.
- Warning → Perangkat berhasil diakses namun ditemukan kondisi tertentu yang memerlukan perhatian.
- Unknown → Status perangkat tidak diketahui karena perangkat sedang tidak dimonitor (Inactive).

Status monitoring akan diperbarui secara berkala dan ditampilkan secara realtime pada halaman Monitoring.


---

## Information Displayed

Informasi yang akan ditampilkan pada halaman Monitoring meliputi:

- Device Name
- Device Type
- Monitoring Status
- Last Check Time
- Monitoring Method
- Check Interval
- Device Configuration Status (Active / Inactive)
- IP Address atau URL perangkat


---

## Sections

### 1. Monitoring Summary

Section ini digunakan untuk memberikan gambaran singkat mengenai kondisi perangkat yang sedang dimonitor.

Informasi yang ditampilkan:

- Total Device
- Online Device
- Offline Device
- Warning Device
- Unknown Device

Tujuan:
- Membantu pengguna mengetahui kondisi monitoring secara keseluruhan dengan cepat.


---

### 2. Filter & Search

Section ini digunakan untuk mempermudah pengguna dalam mencari perangkat tertentu maupun melakukan filtering berdasarkan status perangkat.

Fitur yang disediakan:

- Search Device
- Filter berdasarkan Monitoring Status:
    - All
    - Online
    - Offline
    - Warning
    - Unknown
- Filter berdasarkan Device Type:
    - Server
    - Router
    - Switch
    - Access Point
    - Website
    - Printer
    - Other

Tujuan:
- Membantu pengguna menemukan perangkat yang sedang dicari dengan lebih cepat.
- Mempermudah proses identifikasi perangkat yang mengalami gangguan.


---

### 3. Device List

Section ini merupakan bagian utama dari halaman Monitoring yang menampilkan seluruh perangkat yang sedang dimonitor oleh Gamon.

Informasi yang ditampilkan untuk setiap perangkat:

- Device Name
- Device Type
- Monitoring Status
- Last Check Time

Contoh informasi yang ingin disampaikan kepada pengguna:

- Perangkat sedang Online.
- Perangkat mengalami Offline.
- Perangkat mengalami Warning.
- Perangkat sedang tidak dimonitor.

Tujuan:
- Menampilkan kondisi seluruh perangkat secara sederhana dan mudah dipahami.
- Membantu pengguna mengetahui perangkat yang mengalami gangguan dengan cepat.


---

### 4. Device Detail

Section ini akan ditampilkan ketika pengguna memilih salah satu perangkat pada Device List.

Informasi yang ditampilkan meliputi:

- Device Name
- Device Type
- IP Address / URL
- Monitoring Status
- Monitoring Method
- Last Check Time
- Check Interval
- Device Configuration Status
- Monitoring Description (informasi sederhana mengenai hasil monitoring)

Tujuan:
- Memberikan informasi monitoring yang lebih lengkap terhadap perangkat yang dipilih.
- Membantu pengguna mengetahui kondisi perangkat secara lebih detail.


---

## User Flow

1. Pengguna membuka halaman Monitoring.
2. Pengguna melihat Monitoring Summary untuk mengetahui kondisi monitoring secara keseluruhan.
3. Pengguna dapat melakukan pencarian ataupun filtering perangkat yang ingin dilihat.
4. Pengguna memilih salah satu perangkat pada Device List.
5. Pengguna melihat informasi monitoring perangkat secara lebih lengkap melalui Device Detail.
6. Apabila perangkat mengalami gangguan, pengguna dapat menuju halaman Alert Center untuk melihat informasi alert yang dihasilkan oleh Gamon.


---

## Scope Limitation (MVP)

Halaman Monitoring pada prototype Gamon hanya berfokus pada monitoring status perangkat dan tidak mencakup fitur monitoring tingkat enterprise.

Fitur yang tidak termasuk dalam MVP:

- CPU Usage Monitoring
- RAM Usage Monitoring
- Network Traffic Monitoring
- Network Topology Visualization
- Device Mapping
- Performance Analytics
- Predictive Monitoring
- AI Recommendation System
- Advanced Vendor Integration

Prototype ini berfokus pada pembuktian konsep (Proof of Concept) bahwa Gamon mampu melakukan monitoring status perangkat secara terpusat serta membantu pengguna mengetahui adanya gangguan dengan lebih cepat.


HALAMAN DASHBOARD VS HALAMAN MONITORING
Menurut gua yang paling penting adalah kita harus menentukan terlebih dahulu "tanggung jawab" dari masing-masing halaman. Karena kalau tidak, nantinya Dashboard dan Monitoring akan saling mengambil peran satu sama lain.

Awalnya kita sempat berpikir bahwa Dashboard juga akan menampilkan status perangkat, sedangkan Monitoring juga menampilkan status perangkat. Akibatnya muncul pertanyaan seperti:

> "Kalau status perangkat sudah ada di Dashboard, buat apa ada halaman Monitoring?"

atau sebaliknya,

> "Kalau Monitoring sudah bisa melihat jumlah perangkat yang Offline dan Online, mengapa Dashboard juga menampilkannya?"

Untuk menghindari hal tersebut, kita perlu memposisikan kedua halaman tersebut sebagai berikut.

### Dashboard

Dashboard merupakan halaman ringkasan (summary page) yang berfungsi untuk memberikan gambaran umum mengenai kondisi monitoring yang sedang berlangsung.

Dashboard hanya menjawab pertanyaan:

* Apakah terdapat perangkat yang mengalami gangguan?
* Berapa banyak perangkat yang sedang dimonitor?
* Berapa banyak perangkat yang Online, Offline, Warning, dan Unknown?
* Apakah terdapat alert terbaru?
* Apakah sistem monitoring Gamon sedang berjalan dengan baik?

Dashboard tidak bertujuan untuk memberitahu:

* Perangkat mana yang mengalami gangguan.
* Informasi detail setiap perangkat.
* Kondisi monitoring seluruh perangkat secara lengkap.

Dengan kata lain, Dashboard hanya berfungsi sebagai "pintu masuk" untuk mengetahui kondisi sistem monitoring secara keseluruhan.

> Dashboard = "Apa yang sedang terjadi?"

---

### Monitoring

Monitoring merupakan halaman operasional yang digunakan untuk melihat kondisi seluruh perangkat yang sedang dimonitor oleh Gamon.

Monitoring menjawab pertanyaan:

* Perangkat mana yang sedang Online?
* Perangkat mana yang sedang Offline?
* Perangkat mana yang mengalami Warning?
* Kapan terakhir kali perangkat dilakukan pengecekan?
* Bagaimana kondisi perangkat tertentu saat ini?

Monitoring tidak bertujuan untuk:

* Menampilkan ringkasan sistem secara keseluruhan.
* Menampilkan alert terbaru.
* Mengelola konfigurasi perangkat.

Dengan kata lain, Monitoring berfungsi sebagai tempat pengguna melakukan aktivitas monitoring terhadap perangkat yang terdaftar pada Gamon.

> Monitoring = "Perangkat mana yang sedang mengalami kondisi tersebut?"

---

### Perbedaan Keduanya

| Dashboard                                            | Monitoring                                                       |
| ---------------------------------------------------- | ---------------------------------------------------------------- |
| Menampilkan ringkasan monitoring                     | Menampilkan kondisi seluruh perangkat                            |
| Berfokus pada kondisi sistem secara keseluruhan      | Berfokus pada kondisi masing-masing perangkat                    |
| Menampilkan jumlah perangkat berdasarkan status      | Menampilkan daftar perangkat beserta statusnya                   |
| Menampilkan latest alert                             | Tidak menampilkan latest alert                                   |
| Tidak menampilkan informasi monitoring secara detail | Menampilkan informasi monitoring perangkat secara lebih detail   |
| Digunakan untuk mengetahui apakah terdapat masalah   | Digunakan untuk mengetahui perangkat mana yang mengalami masalah |

---

### User Flow

Berikut merupakan alur penggunaan yang diharapkan pada Gamon:

> Dashboard → Monitoring → Device Detail → Alert Center (jika diperlukan)

Sebagai contoh:

1. Pengguna membuka Dashboard.
2. Dashboard menunjukkan bahwa terdapat 3 perangkat yang berstatus Offline.
3. Pengguna kemudian membuka halaman Monitoring.
4. Pengguna melakukan filter berdasarkan status Offline.
5. Monitoring menampilkan bahwa perangkat yang mengalami gangguan adalah DNS Server, Access Point Lt.3, dan Website Server.
6. Pengguna memilih salah satu perangkat untuk melihat informasi monitoring yang lebih lengkap.
7. Apabila diperlukan, pengguna dapat membuka Alert Center untuk melihat informasi alert yang dihasilkan oleh Gamon.

Pada alur tersebut dapat dilihat bahwa Dashboard tidak pernah memberitahu perangkat mana yang mengalami gangguan. Dashboard hanya memberitahu bahwa "terdapat gangguan". Untuk mengetahui "gangguan terjadi pada perangkat yang mana", pengguna harus membuka halaman Monitoring.

### Analogi Sederhana

Apabila dianalogikan dengan rumah sakit:

> Dashboard adalah papan informasi di lobby rumah sakit yang memberitahu bahwa terdapat 10 pasien yang sedang dirawat, 2 pasien dalam kondisi kritis, dan 5 pasien baru saja masuk.

Sedangkan:

> Monitoring adalah ruang perawatan yang digunakan untuk melihat siapa pasiennya, bagaimana kondisinya, serta informasi yang lebih detail mengenai pasien tersebut.

Dashboard tidak perlu mengetahui siapa pasiennya, sedangkan Monitoring memang dirancang untuk mengetahui informasi tersebut.

### Design Principle

Untuk menghindari redundansi informasi pada Gamon, kedua halaman dirancang dengan prinsip berikut:

* Dashboard berfokus pada **Summary Monitoring**.
* Monitoring berfokus pada **Device Monitoring**.
* Dashboard menjawab **"Apa yang sedang terjadi?"**.
* Monitoring menjawab **"Perangkat mana yang sedang mengalami kondisi tersebut?"**.
* Dashboard tidak menampilkan detail monitoring perangkat.
* Monitoring tidak mengambil peran Dashboard sebagai halaman ringkasan sistem.

Dengan pemisahan tanggung jawab tersebut, Dashboard dan Monitoring saling melengkapi tanpa menampilkan informasi yang sama secara berulang, sehingga antarmuka menjadi lebih sederhana, konsisten, dan sesuai dengan konsep MVP yang ingin diterapkan pada prototype Gamon.
