## 1. Latar Belakang

GAMON menggunakan ICMP sebagai salah satu metode utama untuk mengetahui host availability.

Secara konseptual:

```text
Monitoring Engine
       ↓
ICMP Checker
       ↓
ICMP Executor
       ↓
Target Device
       ↓
Monitoring Result
```

Sebelumnya pemeriksaan ICMP dapat dianggap hanya menggunakan satu mekanisme. Dengan adanya kebutuhan untuk menyediakan alternatif `ping` dan `fping`, ICMP Checker perlu dipisahkan dari mekanisme aktual yang digunakan untuk menjalankan pemeriksaan.

Dengan demikian:

```text
ICMP Checker
      │
      ├── Ping Executor
      │
      └── Fping Executor
```

ICMP Checker tetap menangani konsep pemeriksaan ICMP, sedangkan Executor menangani bagaimana pemeriksaan tersebut dijalankan.

---

# 2. Permasalahan

Penggunaan satu executor ICMP membuat GAMON bergantung pada satu metode eksekusi.

Hal tersebut dapat menjadi keterbatasan ketika:

- jumlah target monitoring bertambah;
    
- kebutuhan monitoring menjadi lebih besar;
    
- environment server berbeda;
    
- administrator memiliki kebutuhan eksekusi yang berbeda;
    
- executor tertentu tidak tersedia pada environment tertentu.
    

Selain itu, apabila logic `ping` ditempatkan langsung di Monitoring Engine, penambahan alternatif seperti `fping` akan membuat engine semakin penuh dengan logic spesifik.

Contoh desain yang kurang baik:

```text
Monitoring Engine
│
├── kalau pakai ping → logic ping
├── kalau pakai fping → logic fping
├── kalau Windows → logic Windows
├── kalau Linux → logic Linux
└── dst.
```

Desain tersebut membuat Monitoring Engine terlalu mengetahui detail implementasi.

Yang diinginkan:

```text
Monitoring Engine
        ↓
   ICMP Checker
        ↓
  Executor Selection
        ↓
 ┌──────┴──────┐
 ↓             ↓
Ping         Fping
Executor     Executor
```

---

# 3. Tujuan

Fitur Ping/Fping Executor bertujuan untuk:

1. Memisahkan mekanisme eksekusi ICMP dari ICMP Checker.
    
2. Menyediakan pilihan executor `Ping` atau `Fping`.
    
3. Memungkinkan administrator menentukan executor yang digunakan GAMON.
    
4. Mempertahankan hasil monitoring dalam bentuk yang konsisten meskipun executor berbeda.
    
5. Memudahkan pergantian executor tanpa mengubah konsep ICMP monitoring.
    
6. Menyiapkan arsitektur untuk kemungkinan penambahan executor lain di masa depan.
    
7. Menjaga Monitoring Engine tetap tidak bergantung pada detail command tertentu.
    

---

# 4. Konsep Utama

Perlu dibedakan tiga istilah:

```text
ICMP
↓
Jenis/protokol monitoring

ICMP Checker
↓
Komponen yang melakukan proses pemeriksaan ICMP

Executor
↓
Mekanisme yang digunakan untuk menjalankan pemeriksaan
```

Sehingga:

```text
Monitoring Engine
        ↓
ICMP Checker
        ↓
Executor
   ┌────┴────┐
   ↓         ↓
 Ping      Fping
```

Baik Ping maupun Fping tetap menghasilkan satu tujuan yang sama:

> Menentukan apakah target dapat dijangkau melalui ICMP dan memperoleh informasi hasil pemeriksaan yang dibutuhkan GAMON.

---

# 5. Posisi Executor dalam Arsitektur GAMON

Arsitektur konseptualnya:

```text
                         GAMON
                           │
                           ↓
                  MONITORING ENGINE
                           │
                           ↓
                     ICMP CHECKER
                           │
                  Executor Selection
                           │
              ┌────────────┴────────────┐
              ↓                         ↓
        PING EXECUTOR             FPING EXECUTOR
              │                         │
              └────────────┬────────────┘
                           ↓
                     TARGET HOST
                           ↓
                    RAW CHECK RESULT
                           ↓
                 NORMALIZED RESULT
                           ↓
                  RESULT PROCESSING
                           ↓
              STATUS / HISTORY / ALERT
```

Hal pentingnya adalah **hasil dari Ping dan Fping harus masuk ke processing pipeline yang sama**.

---

# 6. Ping Executor

Ping Executor merupakan executor ICMP yang menggunakan mekanisme `ping` sebagai pelaksana pemeriksaan.

Konsepnya:

```text
ICMP Checker
      ↓
Ping Executor
      ↓
Target
      ↓
Ping Result
```

Contoh hasil konseptual:

```text
Target : 192.168.1.10
Status : ONLINE
Latency: 2 ms
TTL    : 64
```

Ping dapat menjadi executor default karena merupakan metode yang umum digunakan dan sederhana untuk pemeriksaan satu target.

Namun PRD ini tidak menetapkan command atau implementasi spesifik yang digunakan untuk menjalankannya.

---

# 7. Fping Executor

Fping Executor merupakan alternatif executor ICMP yang ditujukan terutama untuk kebutuhan monitoring terhadap banyak target.

Konsep:

```text
ICMP Checker
      ↓
Fping Executor
      ↓
Multiple Targets
      ↓
Fping Results
```

Keberadaan Fping memberikan alternatif mekanisme eksekusi ketika GAMON digunakan untuk memonitor jumlah target yang lebih banyak.

Namun sistem tidak boleh mengasumsikan bahwa Fping selalu lebih cepat atau selalu lebih efisien daripada Ping dalam semua kondisi.

Efisiensi bergantung pada:

- jumlah target;
    
- interval monitoring;
    
- concurrency;
    
- konfigurasi executor;
    
- environment;
    
- karakteristik jaringan.
    

Karena itu, pilihan executor sebaiknya dipandang sebagai **opsi konfigurasi**, bukan klaim bahwa salah satu metode selalu superior.

---

# 8. Pemilihan Executor

Pemilihan executor dilakukan melalui konfigurasi GAMON.

Contoh desain Settings:

```text
ICMP Executor

● Ping
○ Fping
```

atau:

```text
ICMP Executor: [ Ping ▼ ]
```

Administrator dapat mengubah executor sesuai kebutuhan environment.

Konsepnya:

```text
Settings
   ↓
ICMP Executor
   ↓
Ping / Fping
   ↓
ICMP Checker
```

Perubahan konfigurasi tersebut hanya mengubah mekanisme eksekusi ICMP.

Tidak mengubah:

```text
Device
Monitoring Method
Alert Engine
Telegram
Dashboard
History
```

---

# 9. Scope Konfigurasi

Untuk versi awal, konfigurasi executor sebaiknya sederhana.

Minimal:

```text
ICMP Executor
├── Ping
└── Fping
```

Tidak perlu langsung menyediakan banyak parameter executor.

Contohnya jangan langsung membuat:

```text
Packet Count
Packet Size
TTL
Interval
Timeout
Retry
Interface
Source Address
Fping Mode
Ping Mode
...
```

karena sebagian parameter tersebut lebih cocok menjadi konfigurasi monitoring umum atau advanced configuration.

Tujuan fitur ini hanya memberikan **pemilihan mekanisme eksekusi**.

---

# 10. Global atau Per-Device?

Untuk desain awal GAMON, pilihan executor lebih baik dianggap sebagai **konfigurasi ICMP secara global**.

Contoh:

```text
Settings

ICMP Executor:
Fping
```

Maka:

```text
Device A → ICMP → Fping
Device B → ICMP → Fping
Device C → ICMP → Fping
```

Hal ini jauh lebih sederhana daripada:

```text
Device A → Ping
Device B → Fping
Device C → Ping
Device D → Fping
```

yang akan membuat konfigurasi setiap device menjadi lebih kompleks.

Per-device executor dapat dipertimbangkan di masa depan apabila memang ada kebutuhan operasional yang nyata.

---

# 11. Hubungan dengan Monitoring Interval

Executor tidak menentukan kapan monitoring dilakukan.

Interval tetap menjadi tanggung jawab Monitoring Engine.

```text
Monitoring Engine
       ↓
Interval
       ↓
ICMP Checker
       ↓
Selected Executor
       ↓
Check
```

Misalnya:

```text
Interval = 30 detik
Executor = Fping
```

Maka setiap siklus monitoring:

```text
00:00 → Fping
00:30 → Fping
01:00 → Fping
01:30 → Fping
```

Jika executor diganti:

```text
Executor = Ping
```

maka siklusnya menjadi:

```text
00:00 → Ping
00:30 → Ping
01:00 → Ping
```

Interval tidak berubah.

---

# 12. Normalisasi Hasil

Ini merupakan bagian desain yang sangat penting.

Ping dan Fping dapat memberikan output yang berbeda.

Contohnya secara konseptual:

```text
Ping
→ output A

Fping
→ output B
```

Tetapi GAMON tidak boleh membuat processing layer harus memahami kedua output tersebut.

Maka:

```text
Ping Executor
       ↓
Raw Ping Result
       ↓
Normalize
       ↓
Standard ICMP Result
```

dan:

```text
Fping Executor
       ↓
Raw Fping Result
       ↓
Normalize
       ↓
Standard ICMP Result
```

Kemudian keduanya menjadi:

```text
Standard ICMP Result
├── Target
├── Status
├── Latency
├── TTL (jika tersedia)
└── Timestamp
```

Informasi yang tidak tersedia dari executor tertentu tidak boleh dipaksakan.

---

# 13. Contoh Normalisasi

Misalnya Ping menghasilkan:

```text
Target: 192.168.1.10
Reachable: YES
Latency: 2 ms
TTL: 64
```

Fping menghasilkan informasi:

```text
Target: 192.168.1.10
Reachable: YES
Latency: 2.1 ms
```

Keduanya dapat dinormalisasi menjadi:

```text
Target      : 192.168.1.10
Status      : ONLINE
Latency     : 2.1 ms
TTL         : available / unavailable
Timestamp   : ...
```

Result Processor tidak perlu tahu apakah hasil tersebut berasal dari Ping atau Fping.

---

# 14. Failure Handling

Executor harus memiliki kondisi ketika pemeriksaan tidak menghasilkan response.

Misalnya:

```text
Target
 ↓
Executor
 ↓
No Response
```

Kemudian hasil dinormalisasi menjadi kondisi yang dapat diproses engine.

Contoh konseptual:

```text
ONLINE
OFFLINE
TIMEOUT
ERROR
```

Namun definisi final status dapat ditentukan dalam Technical Design.

Yang penting:

> Executor tidak langsung membuat Alert.

Executor hanya melaporkan hasil pemeriksaan.

---

# 15. Error Executor vs Target Down

Ini perlu dibedakan.

Misalnya:

```text
Target tidak merespons
```

berbeda dengan:

```text
Executor tidak dapat dijalankan
```

Contoh:

```text
Ping Executor
     ↓
Target tidak merespons
     ↓
Monitoring Result
     ↓
Potential Host Down
```

berbeda dengan:

```text
Ping Executor
     ↓
Executor tidak tersedia / gagal dijalankan
     ↓
Execution Error
```

Kondisi kedua tidak boleh langsung dianggap sebagai:

```text
DEVICE OFFLINE
```

karena masalahnya mungkin berada pada executor/environment GAMON, bukan pada target.

Ini sangat penting agar sistem tidak menghasilkan false alert.

---

# 16. Executor Availability

Ketika administrator memilih:

```text
Fping
```

GAMON harus dapat mengetahui apakah executor tersebut tersedia pada environment tempat aplikasi berjalan.

Konsep:

```text
Settings
   ↓
Select Fping
   ↓
Check Executor Availability
   ↓
Available?
 ┌──────┴──────┐
YES           NO
 ↓             ↓
Use Fping   Configuration Error
```

Jika Fping tidak tersedia, GAMON sebaiknya tidak diam-diam menganggap device offline.

Sistem harus membedakan:

```text
Executor unavailable
```

dengan:

```text
Target unavailable
```

---

# 17. Runtime Environment

Executor dapat bergantung pada environment tempat GAMON berjalan.

Karena itu konsep runtime environment dapat digambarkan:

```text
GAMON Server
      ↓
Runtime Environment
      ↓
Available Executors
      ↓
ICMP Checker
```

Contohnya:

```text
Runtime
OS: Linux

Available:
✓ Ping
✓ Fping
```

atau:

```text
Runtime
OS: Environment lain

Available:
✓ Ping
✗ Fping
```

Detail OS detection sendiri merupakan bagian terpisah dari fitur Runtime Environment.

Executor hanya perlu mendapatkan informasi apakah executor yang dipilih dapat digunakan.

---

# 18. Fallback

Untuk versi awal, gua **tidak menyarankan automatic fallback**.

Contohnya jangan langsung:

```text
Selected: Fping

Fping unavailable
      ↓
Automatically use Ping
```

Karena administrator memilih Fping dengan alasan tertentu.

Jika sistem diam-diam menggantinya dengan Ping, behavior monitoring menjadi tidak sesuai konfigurasi.

Lebih jelas:

```text
Selected: Fping
Fping unavailable
      ↓
Configuration / Executor Error
      ↓
Administrator memperbaiki environment
```

Automatic fallback bisa menjadi fitur lanjutan jika nanti memang dibutuhkan.

---

# 19. Pengaruh terhadap Alert

Executor tidak boleh menentukan alert secara langsung.

Misalnya:

```text
Ping
 ↓
ONLINE
```

dan:

```text
Fping
 ↓
ONLINE
```

keduanya menghasilkan:

```text
Host Status = ONLINE
```

Kemudian ketika:

```text
Fping
 ↓
OFFLINE
```

hasil masuk ke:

```text
Status Evaluation
 ↓
Failure Threshold
 ↓
Host Down
 ↓
Alert Engine
```

Jadi pergantian executor tidak mengubah konsep alert.

---

# 20. Pengaruh terhadap History

History juga harus tetap konsisten.

Misalnya:

```text
Ping
→ ONLINE
→ 2 ms
```

kemudian executor diganti:

```text
Fping
→ ONLINE
→ 2.1 ms
```

Data history tetap merepresentasikan monitoring ICMP.

Jika diperlukan untuk audit/debugging, executor yang digunakan dapat dicatat sebagai metadata:

```text
Method : ICMP
Executor: Fping
```

Namun apakah executor perlu menjadi field history permanen dapat ditentukan pada tahap desain database.

---

# 21. Pergantian Executor

Perubahan executor tidak boleh dianggap sebagai perubahan device.

Contoh:

```text
Sebelum:
ICMP → Ping

Settings:
Ping → Fping

Sesudah:
ICMP → Fping
```

Device tetap sama.

Monitoring method tetap sama:

```text
ICMP
```

Yang berubah hanya:

```text
Executor
```

Ini mempertegas pemisahan:

```text
Monitoring Method ≠ Executor
```

---

# 22. Use Case

### Use Case 1 — Default Ping

```text
Administrator
      ↓
Settings
      ↓
ICMP Executor = Ping
      ↓
Monitoring Engine
      ↓
ICMP Checker
      ↓
Ping Executor
      ↓
Target
```

### Use Case 2 — Administrator memilih Fping

```text
Administrator
      ↓
Settings
      ↓
ICMP Executor = Fping
      ↓
Monitoring Engine
      ↓
ICMP Checker
      ↓
Fping Executor
      ↓
Target
```

### Use Case 3 — Fping tidak tersedia

```text
Administrator
      ↓
Select Fping
      ↓
System checks availability
      ↓
Fping unavailable
      ↓
Configuration Error
```

Target tidak langsung dianggap offline.

### Use Case 4 — Target memang down

```text
Fping Executor
      ↓
Target tidak merespons
      ↓
Normalized ICMP Result
      ↓
Status Evaluation
      ↓
Failure Threshold
      ↓
Host Down
      ↓
Alert
```

---

# 23. Non-Functional Requirement

**Consistency**  
Ping dan Fping harus menghasilkan informasi monitoring yang dapat diproses dengan alur yang sama.

**Modularity**  
Executor dapat diganti tanpa mengubah Monitoring Engine.

**Extensibility**  
Executor tambahan dapat ditambahkan tanpa mendesain ulang ICMP Checker.

**Environment Awareness**  
Sistem mengetahui apakah executor yang dipilih tersedia.

**Reliability**  
Executor failure tidak boleh dianggap sebagai target failure.

**Configurability**  
Administrator dapat memilih executor melalui Settings.

**Transparency**  
Sistem tidak diam-diam mengganti executor yang dipilih tanpa konfigurasi eksplisit.

---

# 24. Batasan Scope

PRD ini belum menentukan:

- command `ping`;
    
- command `fping`;
    
- parameter command;
    
- jumlah packet;
    
- timeout;
    
- retry;
    
- packet size;
    
- parsing output;
    
- library;
    
- function;
    
- class;
    
- interface;
    
- subprocess implementation;
    
- OS-specific command;
    
- concurrency implementation.
    

Semua itu masuk ke Technical Design.

PRD hanya menentukan **peran dan hubungan antar-komponen**.

---

# 25. Desain Akhir

Desain akhirnya:

```text
                         GAMON
                           │
                           ↓
                  MONITORING ENGINE
                           │
                           ↓
                     ICMP CHECKER
                           │
                           ↓
                  EXECUTOR SELECTOR
                           │
                ┌──────────┴──────────┐
                ↓                     ↓
          PING EXECUTOR          FPING EXECUTOR
                │                     │
                └──────────┬──────────┘
                           ↓
                      TARGET HOST
                           ↓
                      RAW RESULT
                           ↓
                  RESULT NORMALIZER
                           ↓
                 STANDARD ICMP RESULT
                           ↓
                  STATUS EVALUATION
                           ↓
                FAILURE THRESHOLD
                           ↓
                     EVENT/STATUS
                           │
              ┌────────────┴────────────┐
              ↓                         ↓
           HISTORY                  ALERT ENGINE
                                        │
                                ┌───────┼───────┐
                                ↓       ↓       ↓
                            Dashboard Telegram Sound
```

Jadi konsep utamanya bisa diringkas menjadi:

```text
ICMP
 ↓
ICMP Checker
 ↓
"Bagaimana saya melakukan pengecekan?"
 ↓
Executor
 ├── Ping
 └── Fping
 ↓
"Bagaimana hasilnya?"
 ↓
Normalized ICMP Result
 ↓
Monitoring Engine
 ↓
Status / History / Alert
```

Dengan desain ini, kalau suatu saat lu ingin menambahkan executor lain, misalnya implementasi ICMP native atau mekanisme lain, konsepnya cukup:

```text
ICMP Checker
│
├── Ping Executor
├── Fping Executor
└── Executor Baru
```

tanpa harus mengubah cara GAMON melakukan status evaluation, history, alert, WebSocket, Telegram, atau sound notification. Itu yang sebenarnya ingin dicapai dari upgrade ini: **Ping/Fping menjadi detail eksekusi, bukan menjadi dua sistem monitoring yang berbeda.**