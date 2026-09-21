## 1. Latar Belakang

Monitoring Engine merupakan komponen inti GAMON yang bertanggung jawab menjalankan pemeriksaan terhadap perangkat secara berkala dan menghasilkan informasi mengenai kondisi perangkat maupun service yang dipantau.

Pada versi awal, mekanisme monitoring berfokus pada ICMP Ping. Dengan penambahan TCP Port Monitoring, Monitoring Engine akan memiliki lebih dari satu jenis metode pemeriksaan.

Apabila setiap metode dibuat dengan alur dan mekanisme sendiri-sendiri, sistem dapat berkembang menjadi sulit dipelihara:

```text
ICMP
→ logic sendiri
→ status sendiri
→ history sendiri
→ alert sendiri

TCP
→ logic sendiri
→ status sendiri
→ history sendiri
→ alert sendiri
```

Desain tersebut berpotensi menghasilkan duplikasi logic dan perilaku monitoring yang tidak konsisten.

Oleh karena itu, Monitoring Engine perlu dirapikan sehingga:

```text
Monitoring Engine
       │
       ├── ICMP Checker
       │
       └── TCP Checker
```

memiliki tanggung jawab masing-masing untuk melakukan pemeriksaan, tetapi hasilnya dapat diproses menggunakan alur monitoring yang sama.

---

# 2. Permasalahan

Penambahan metode monitoring baru menimbulkan beberapa permasalahan desain.

Pertama, setiap metode memiliki karakteristik pemeriksaan yang berbeda.

```text
ICMP
→ memeriksa reachability host

TCP
→ memeriksa accessibility service/port
```

Kedua, hasil pemeriksaan juga memiliki informasi berbeda.

ICMP misalnya dapat menghasilkan:

```text
Status
Latency
TTL
```

sedangkan TCP dapat menghasilkan:

```text
Status
Port
Response Time
```

Ketiga, meskipun cara pemeriksaannya berbeda, keduanya tetap membutuhkan proses yang serupa setelah pemeriksaan:

```text
Check
 ↓
Result
 ↓
Evaluation
 ↓
History
 ↓
Status
 ↓
Alert
 ↓
Notification
```

Maka dibutuhkan pemisahan antara **cara melakukan pemeriksaan** dan **cara memproses hasil pemeriksaan**.

---

# 3. Tujuan

Upgrade Monitoring Engine bertujuan untuk membuat sistem monitoring memiliki struktur yang lebih terorganisasi sehingga:

1. ICMP dan TCP memiliki checker/executor yang jelas.
    
2. Setiap checker bertanggung jawab hanya terhadap metode monitoringnya.
    
3. Monitoring Engine menjadi pengatur proses monitoring, bukan pelaksana seluruh detail pemeriksaan.
    
4. Hasil ICMP dan TCP memiliki format konseptual yang dapat diproses oleh alur sistem yang konsisten.
    
5. Status, history, alert, dan notification tidak dibuat ulang untuk setiap metode.
    
6. Penambahan metode monitoring berikutnya seperti HTTP/HTTPS dapat dilakukan tanpa merombak keseluruhan engine.
    
7. Monitoring tetap dapat berjalan secara berkala dan paralel.
    
8. Kegagalan salah satu checker tidak menghentikan monitoring perangkat lainnya.
    

---

# 4. Prinsip Desain

Desain Monitoring Engine menggunakan prinsip:

> **Checker menentukan bagaimana suatu target diperiksa. Monitoring Engine menentukan bagaimana proses monitoring dijalankan dan hasilnya diproses.**

Dengan demikian:

```text
ICMP Checker
→ tahu bagaimana melakukan ICMP check.

TCP Checker
→ tahu bagaimana melakukan TCP check.

Monitoring Engine
→ tahu kapan harus melakukan check dan apa yang dilakukan terhadap hasilnya.
```

Ini adalah pemisahan tanggung jawab utama yang ingin dicapai.

---

# 5. Arsitektur Konseptual

Struktur tingkat tinggi:

```text
                    MONITORING ENGINE
                           │
              ┌────────────┴────────────┐
              │                         │
        ICMP Monitoring            TCP Monitoring
              │                         │
         ICMP Checker               TCP Checker
              │                         │
              └────────────┬────────────┘
                           ↓
                    Monitoring Result
                           ↓
                    Result Processing
                           ↓
                    Status Evaluation
                           ↓
                  History / Alert Engine
                           ↓
                       Notification
```

Monitoring Engine menjadi pusat orkestrasi.

Checker menjadi komponen yang menangani pemeriksaan spesifik.

---

# 6. Tanggung Jawab Monitoring Engine

Monitoring Engine bertanggung jawab terhadap proses tingkat tinggi.

Secara konseptual:

```text
Monitoring Engine
│
├── membaca konfigurasi monitoring
├── menentukan target yang aktif
├── menentukan kapan target diperiksa
├── menentukan checker yang digunakan
├── menjalankan pemeriksaan
├── menerima hasil pemeriksaan
├── meneruskan hasil ke processing
└── menjalankan siklus monitoring berikutnya
```

Monitoring Engine **tidak perlu mengetahui detail bagaimana ICMP atau TCP bekerja**.

Misalnya engine tidak perlu memiliki logic:

```text
"kalau ICMP, lakukan X"
"kalau TCP, lakukan Y"
"kalau error tertentu, lakukan Z"
```

secara tersebar di seluruh bagian engine.

Detail tersebut berada pada checker masing-masing.

---

# 7. ICMP Checker

ICMP Checker bertanggung jawab terhadap pemeriksaan availability host menggunakan ICMP.

Konsep:

```text
ICMP Checker
     ↓
Target Host/IP
     ↓
ICMP Check
     ↓
ICMP Result
```

Hasil konseptual dapat berupa:

```text
Host
Status
Latency
TTL
Timestamp
```

Contoh:

```text
Device: Router Core
Host: 192.168.1.1

Status: ONLINE
Latency: 2 ms
TTL: 64
```

ICMP Checker tidak bertanggung jawab untuk:

- membuat alert;
    
- mengirim Telegram;
    
- memainkan sound;
    
- memperbarui dashboard;
    
- menentukan lifecycle alert.
    

Ia hanya menghasilkan hasil pemeriksaan.

---

# 8. TCP Checker

TCP Checker memiliki tanggung jawab yang sama pada level konsep, tetapi menggunakan metode TCP.

```text
TCP Checker
     ↓
Target Host/IP + Port
     ↓
TCP Check
     ↓
TCP Result
```

Hasil konseptual:

```text
Host
Port
Status
Response Time
Timestamp
```

Contoh:

```text
Device: Web Server
Host: 192.168.1.10
Port: 443

Status: OPEN
Response Time: 8 ms
```

TCP Checker juga tidak bertanggung jawab terhadap alert dan notification.

---

# 9. Persamaan ICMP dan TCP

Meskipun pemeriksaannya berbeda, keduanya mempunyai konsep dasar yang sama:

```text
Target
   ↓
Check
   ↓
Result
   ↓
Timestamp
```

Kemudian hasil masuk ke proses yang lebih umum:

```text
Result
 ↓
Evaluation
 ↓
Status
 ↓
History
 ↓
Alert
```

Dengan demikian, sistem tidak membutuhkan dua pipeline monitoring yang sepenuhnya terpisah.

---

# 10. Perbedaan ICMP dan TCP

Checker tetap harus mempertahankan karakteristik masing-masing.

|Aspek|ICMP Checker|TCP Checker|
|---|---|---|
|Tujuan|Host availability|Service availability|
|Target|Host/IP|Host/IP + port|
|Protokol|ICMP|TCP|
|Informasi utama|Reachability, latency|Port accessibility, response time|
|Status|Online/Offline|Open/Closed/Timeout, sesuai desain|
|Fungsi|Mengetahui host reachable|Mengetahui service/port reachable|

Jadi “checker yang jelas” bukan berarti ICMP dan TCP dipaksa memiliki logic yang identik.

Yang diseragamkan adalah **alur hasilnya**, bukan mekanisme pemeriksaannya.

---

# 11. Monitoring Configuration

Setiap target monitoring perlu menentukan metode yang digunakan.

Konsep:

```text
Device
│
├── Monitoring Method: ICMP
│
atau
│
└── Monitoring Method: TCP
       └── Port: 443
```

Contoh:

```text
Web Server 01

Monitoring:
├── ICMP
└── TCP 443
```

Ini memungkinkan satu perangkat memiliki lebih dari satu monitoring check.

Hal ini penting karena:

```text
ICMP = Host Availability
TCP 443 = HTTPS Service Availability
```

keduanya dapat memberikan informasi berbeda terhadap perangkat yang sama.

---

# 12. Satu Device, Banyak Check

Desain engine sebaiknya tidak menganggap:

> satu device = satu monitoring.

Lebih fleksibel apabila:

> satu device dapat memiliki beberapa monitoring check.

Contoh:

```text
Web Server 01
│
├── ICMP
│
├── TCP 22
│
├── TCP 80
│
└── TCP 443
```

Hasilnya:

```text
Web Server 01
│
├── Host      → ONLINE
├── SSH 22    → OPEN
├── HTTP 80   → OPEN
└── HTTPS 443 → DOWN
```

Ini membuat model GAMON lebih siap dikembangkan ke HTTP/HTTPS maupun metode lainnya.

---

# 13. Siklus Monitoring

Monitoring Engine menjalankan siklus secara berulang.

Konsep:

```text
Start
  ↓
Load Active Checks
  ↓
Determine Schedule
  ↓
Run Checkers
  ↓
Collect Results
  ↓
Process Results
  ↓
Store Results
  ↓
Evaluate Status
  ↓
Trigger Alert if Required
  ↓
Next Monitoring Cycle
```

Contohnya:

```text
Cycle 01
├── ICMP Router
├── ICMP Server
├── TCP Server:443
└── TCP Server:22

Cycle 02
├── ICMP Router
├── ICMP Server
├── TCP Server:443
└── TCP Server:22
```

Interval masing-masing check dapat berbeda apabila konfigurasi GAMON memang mengizinkannya.

---

# 14. Parallel Monitoring

Karena GAMON dapat memonitor banyak perangkat, checker tidak seharusnya menyebabkan satu target menghambat seluruh monitoring.

Konsep:

```text
Monitoring Engine
        │
        ├── ICMP Device A
        ├── ICMP Device B
        ├── TCP Device C:443
        ├── TCP Device D:22
        └── TCP Device E:80
```

Jika:

```text
TCP Device C
```

mengalami timeout, pemeriksaan:

```text
ICMP Device A
ICMP Device B
TCP Device D
```

tetap dapat berjalan.

Detail concurrency belum dibahas di PRD ini.

Yang ditetapkan di level desain hanyalah:

> Kegagalan atau keterlambatan satu monitoring check tidak boleh menghentikan keseluruhan Monitoring Engine.

---

# 15. Result Processing

Setelah checker selesai, hasil pemeriksaan diteruskan ke proses berikutnya.

```text
Checker
   ↓
Monitoring Result
   ↓
Result Processing
```

Result Processor bertanggung jawab untuk menerjemahkan hasil checker menjadi kondisi monitoring yang dapat digunakan oleh sistem.

Misalnya:

```text
ICMP Result
ONLINE
 ↓
Host Status = ONLINE
```

atau:

```text
TCP Result
OPEN
 ↓
Service Status = AVAILABLE
```

Kemudian:

```text
Previous Status
        +
Current Status
        ↓
Status Evaluation
```

---

# 16. Status Change Detection

Salah satu fungsi penting engine adalah mendeteksi perubahan kondisi.

Contoh ICMP:

```text
ONLINE
 ↓
OFFLINE
```

menghasilkan:

```text
Host Down Event
```

Contoh TCP:

```text
OPEN
 ↓
CLOSED
```

menghasilkan:

```text
Service Down Event
```

Recovery:

```text
OFFLINE
 ↓
ONLINE
```

atau:

```text
CLOSED
 ↓
OPEN
```

menghasilkan recovery event.

Dengan demikian, Alert Engine tidak perlu mengetahui cara pemeriksaan dilakukan.

Alert Engine hanya menerima event:

```text
Host Down
Service Down
Host Recovery
Service Recovery
```

---

# 17. Failure Threshold

Failure threshold berada pada level **monitoring state**, bukan menjadi tanggung jawab khusus ICMP atau TCP.

Contoh:

```text
Threshold = 3
```

TCP:

```text
Check 1 → CLOSED
Check 2 → CLOSED
Check 3 → CLOSED
```

kemudian:

```text
Service Down Event
```

ICMP juga menggunakan prinsip yang sama:

```text
Check 1 → Failed
Check 2 → Failed
Check 3 → Failed
```

kemudian:

```text
Host Down Event
```

Dengan desain ini, mekanisme threshold tidak perlu dibuat ulang dalam ICMP Checker dan TCP Checker.

---

# 18. History

History juga sebaiknya diproses setelah checker menghasilkan result.

```text
ICMP Checker
     ↓
ICMP Result
     ↓
History

TCP Checker
     ↓
TCP Result
     ↓
History
```

Namun data yang disimpan dapat memiliki informasi spesifik sesuai metode.

Misalnya ICMP:

```text
Status
Latency
TTL
Timestamp
```

TCP:

```text
Status
Port
Response Time
Timestamp
```

Detail model database tidak menjadi bagian dari PRD ini.

---

# 19. Hubungan dengan Alert Engine

Monitoring Engine tidak mengirim Telegram atau memainkan sound secara langsung.

Alurnya:

```text
Checker
  ↓
Result
  ↓
Monitoring Engine
  ↓
Status Evaluation
  ↓
Event
  ↓
Alert Engine
  ↓
Notification
```

Misalnya:

```text
TCP 443
OPEN → CLOSED
        ↓
Service Down Event
        ↓
Alert Engine
        ↓
Telegram + Dashboard + Sound
```

Dengan pemisahan tersebut, Monitoring Engine tetap fokus pada monitoring.

---

# 20. Executor vs Checker

Dalam desain ini ada baiknya dibedakan dua istilah.

**Checker** menjelaskan jenis pemeriksaan:

```text
ICMP Checker
TCP Checker
```

Sedangkan **Executor** menjelaskan mekanisme yang digunakan untuk menjalankan pemeriksaan tersebut.

Contoh untuk ICMP:

```text
ICMP Checker
      │
      ├── Ping Executor
      └── Fping Executor
```

Dengan demikian, ketika fitur Ping/Fping ditambahkan, arsitekturnya tidak menjadi:

```text
Monitoring Engine
├── Ping
├── Fping
├── TCP
```

tetapi:

```text
Monitoring Engine
│
├── ICMP Checker
│      ├── Ping Executor
│      └── Fping Executor
│
└── TCP Checker
```

Ini jauh lebih bersih secara konsep.

TCP juga nantinya dapat memiliki executor tertentu apabila memang diperlukan, tetapi detail tersebut belum ditentukan dalam PRD ini.

---

# 21. Runtime Environment

Jika GAMON nantinya menggunakan executor berbasis external tool, environment tempat GAMON berjalan dapat memengaruhi executor yang tersedia.

Konsepnya:

```text
GAMON Runtime
      ↓
Runtime Environment
      ↓
Available Executor
```

Misalnya:

```text
ICMP Checker
      ↓
Executor Selection
      ↓
Ping / Fping
```

Namun **OS Detection bukan tanggung jawab Monitoring Engine**.

Engine hanya perlu mengetahui executor mana yang tersedia dan dapat digunakan.

Dengan demikian konsepnya tetap modular:

```text
Monitoring Engine
       ↓
ICMP Checker
       ↓
ICMP Executor
       ↓
Runtime Environment
```

Hal ini mencegah Monitoring Engine dipenuhi logic khusus Windows/Linux/macOS.

---

# 22. Penambahan HTTP/HTTPS di Masa Depan

Salah satu tujuan desain ini adalah agar metode monitoring berikutnya tidak memerlukan perubahan besar.

Setelah ICMP dan TCP:

```text
Monitoring Engine
│
├── ICMP Checker
│
├── TCP Checker
│
└── HTTP Checker
```

HTTP Checker kemudian memiliki tanggung jawab:

```text
Target URL
 ↓
HTTP Request
 ↓
HTTP Result
```

Kemudian hasilnya tetap masuk:

```text
Monitoring Result
 ↓
Result Processing
 ↓
Status Evaluation
 ↓
Alert
 ↓
Notification
```

Dengan demikian penambahan HTTP tidak memerlukan pembuatan Alert Engine baru.

---

# 23. Error Isolation

Salah satu requirement penting:

> Error pada satu checker tidak boleh menjatuhkan Monitoring Engine.

Contoh:

```text
TCP Checker
Device A
     ↓
Unexpected Error
```

tidak boleh menyebabkan:

```text
Monitoring Engine
     ↓
STOP
```

Melainkan:

```text
TCP Device A
→ Check Failed

Monitoring Engine
→ Continue

ICMP Device B
→ Checked

TCP Device C
→ Checked
```

Error harus dianggap sebagai hasil/kondisi dari monitoring tertentu dan ditangani sesuai mekanisme yang dirancang.

---

# 24. Observability Monitoring Engine

Karena Monitoring Engine merupakan komponen inti, sistem sebaiknya dapat membedakan:

```text
Monitoring Target
        ↓
Checker
        ↓
Execution
        ↓
Result
```

Contoh informasi internal:

```text
Device: Web Server 01
Check: TCP 443
Started: 14:30:00
Finished: 14:30:01
Result: OPEN
```

Hal ini akan membantu debugging apabila terdapat masalah pada engine.

Tidak berarti semua informasi tersebut harus ditampilkan ke pengguna.

---

# 25. Non-Functional Requirement

Monitoring Engine yang dirancang ulang harus memenuhi prinsip:

**Modular**  
Penambahan checker baru tidak membutuhkan perubahan besar pada komponen lain.

**Consistent**  
ICMP dan TCP mengikuti alur processing yang konsisten.

**Isolated**  
Kegagalan satu check tidak menghentikan monitoring lainnya.

**Extensible**  
HTTP/HTTPS dan metode lainnya dapat ditambahkan kemudian.

**Maintainable**  
Logic monitoring tidak tersebar di Alert Engine, Notification, Dashboard, atau Database.

**Configurable**  
Metode dan interval monitoring ditentukan berdasarkan konfigurasi.

---

# 26. Batasan Scope

PRD ini **tidak menentukan**:

- bahasa pemrograman;
    
- library ICMP;
    
- library TCP;
    
- penggunaan `ping`;
    
- penggunaan `fping`;
    
- struktur folder;
    
- function;
    
- class;
    
- interface;
    
- API endpoint;
    
- schema database;
    
- mekanisme concurrency;
    
- implementasi goroutine;
    
- timeout value;
    
- retry implementation.
    

Semua itu masuk ke tahap Technical Design/Implementation Design.

PRD hanya menetapkan bagaimana komponen seharusnya berperilaku dan berhubungan.

---

# 27. Desain Akhir

Secara keseluruhan, desain Monitoring Engine GAMON menjadi:

```text
                         GAMON
                           │
                           ↓
                  MONITORING ENGINE
                           │
             ┌─────────────┴─────────────┐
             │                           │
       ICMP CHECKER                 TCP CHECKER
             │                           │
       ┌─────┴─────┐                 TCP Check
       │           │                     │
      Ping       Fping                   │
       │           │                     │
       └─────┬─────┘                     │
             │                           │
             └─────────────┬─────────────┘
                           ↓
                  MONITORING RESULT
                           ↓
                  RESULT PROCESSING
                           ↓
                  STATUS EVALUATION
                           ↓
                   FAILURE THRESHOLD
                           ↓
                    EVENT DETECTION
                           │
              ┌────────────┴────────────┐
              ↓                         ↓
           HISTORY                  ALERT ENGINE
                                        │
                              ┌─────────┼─────────┐
                              ↓         ↓         ↓
                          Dashboard  Telegram   Sound
```

Jadi inti upgrade ini bukan sekadar “pisahkan file ICMP dan TCP”.

Konsep besarnya adalah **memisahkan tiga lapisan tanggung jawab**:

```text
1. CHECKING
   Bagaimana target diperiksa?
   → ICMP Checker / TCP Checker

2. ORCHESTRATION & PROCESSING
   Kapan diperiksa dan bagaimana hasilnya diproses?
   → Monitoring Engine

3. RESPONSE
   Apa yang dilakukan ketika kondisi berubah?
   → History / Alert / Notification
```

Dengan desain ini, GAMON punya fondasi yang jauh lebih masuk akal untuk berkembang:

```text
              Monitoring Engine
                     │
       ┌─────────────┼─────────────┐
       ↓             ↓             ↓
      ICMP           TCP        HTTP/HTTPS
       │             │             │
       └─────────────┼─────────────┘
                     ↓
             Unified Processing
                     ↓
             Unified Alerting
                     ↓
              Unified Notification
```

Dan menurut gua, **ini justru upgrade arsitektur yang paling penting** sebelum lu mulai menambah Nmap, HTTP, SNMP, dan lain-lain. Kalau fondasi engine-nya sudah dipisahkan seperti ini, fitur-fitur tersebut nantinya menjadi penambahan checker baru, bukan bongkar ulang seluruh Monitoring Engine.