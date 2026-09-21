## 1. Latar Belakang

GAMON merupakan sistem monitoring infrastruktur jaringan yang dirancang untuk membantu NOC mengetahui kondisi perangkat dan mendeteksi adanya gangguan pada jaringan.

Pada implementasi saat ini, GAMON menggunakan ICMP sebagai metode utama untuk melakukan pengecekan ketersediaan perangkat. Sistem secara berkala melakukan pengecekan terhadap alamat IP perangkat dan menentukan apakah perangkat berada dalam kondisi online atau offline.

Pendekatan tersebut efektif untuk mengetahui host availability, tetapi memiliki keterbatasan. Perangkat yang masih dapat merespons ICMP belum tentu seluruh service yang berjalan pada perangkat tersebut dalam kondisi normal.

Contohnya:

```text
Server Web
IP: 192.168.1.10

ICMP Ping     → ONLINE
TCP Port 80   → CLOSED
```

Dalam kondisi tersebut, server masih dapat dijangkau secara jaringan, tetapi service web tidak dapat diakses.

Artinya, apabila GAMON hanya menggunakan ICMP, kondisi tersebut dapat dianggap sebagai perangkat yang normal meskipun dari perspektif layanan sebenarnya terdapat gangguan.

TCP Port Monitoring dirancang untuk melengkapi keterbatasan tersebut.

---

## 2. Permasalahan

Permasalahan utama yang ingin diselesaikan adalah keterbatasan monitoring berbasis ICMP dalam mendeteksi gangguan pada service tertentu.

Secara konseptual terdapat dua kondisi berbeda:

```text
Host Availability
Apakah perangkat dapat dijangkau?
                ↓
              ICMP
```

dan:

```text
Service Availability
Apakah service tertentu dapat diakses?
                ↓
           TCP Port
```

Tanpa TCP Port Monitoring, GAMON hanya dapat mengetahui kondisi host secara umum.

Akibatnya terdapat kemungkinan kondisi:

```text
ICMP = ONLINE
Service = DOWN
```

tidak terdeteksi sebagai anomali oleh sistem.

Hal tersebut berhubungan langsung dengan permasalahan utama GAMON, yaitu membantu NOC mengetahui adanya kondisi abnormal pada infrastruktur tanpa harus melakukan pengecekan secara manual terhadap setiap perangkat dan service.

---

## 3. Tujuan Fitur

TCP Port Monitoring bertujuan untuk menambahkan kemampuan GAMON dalam memonitor ketersediaan service berbasis TCP pada suatu perangkat.

Fitur ini diharapkan mampu:

1. Memeriksa apakah port TCP tertentu dapat diakses.
    
2. Menentukan kondisi service berdasarkan hasil pemeriksaan port.
    
3. Membedakan kondisi host yang aktif dengan service yang aktif.
    
4. Menyimpan hasil pemeriksaan sebagai bagian dari data monitoring.
    
5. Mendeteksi perubahan kondisi service.
    
6. Menghasilkan alert ketika service mengalami gangguan.
    
7. Mengintegrasikan hasil monitoring dengan sistem alert dan notification GAMON yang sudah ada.
    
8. Memberikan informasi yang lebih spesifik kepada NOC mengenai sumber anomali.
    

---

# 4. Konsep Sistem

TCP Port Monitoring tidak menggantikan ICMP Monitoring.

Keduanya merupakan metode monitoring yang memiliki fungsi berbeda.

```text
                    DEVICE
                       │
             ┌─────────┴─────────┐
             │                   │
            ICMP                TCP
             │                   │
      Host Availability     Service Availability
             │                   │
      "Host reachable?"    "Port accessible?"
```

Dengan demikian, sebuah perangkat dapat memiliki beberapa jenis pemeriksaan.

Contoh:

```text
Web Server 01
192.168.1.10

ICMP
ONLINE

TCP 22
OPEN

TCP 80
OPEN

TCP 443
OPEN
```

Apabila kemudian port 443 tidak dapat diakses:

```text
ICMP
ONLINE

TCP 22
OPEN

TCP 80
OPEN

TCP 443
CLOSED
```

GAMON tidak menyimpulkan bahwa seluruh perangkat mati.

GAMON justru dapat mengidentifikasi bahwa **host masih tersedia tetapi service pada port 443 mengalami gangguan**.

Ini merupakan perbedaan utama antara ICMP Monitoring dan TCP Port Monitoring.

---

# 5. Scope Fitur

Untuk versi awal, TCP Port Monitoring difokuskan pada pemeriksaan konektivitas terhadap port TCP tertentu.

Setiap konfigurasi monitoring TCP minimal memiliki:

```text
Device
IP / Host
TCP Port
Service Name
Monitoring Interval
Status
```

Contoh:

```text
Device       : Web Server 01
Host         : 192.168.1.10
Port         : 443
Service      : HTTPS
Interval     : 30 seconds
Status       : OPEN
```

Service Name pada tahap awal dapat berfungsi sebagai identitas/deskripsi bagi NOC. Sistem tidak harus mengetahui secara otomatis bahwa port `443` pasti merupakan HTTPS.

---

# 6. Status Monitoring

Status TCP perlu dipisahkan dari status ICMP.

Contohnya:

```text
Device Status
     │
     └── ICMP → ONLINE


Service Status
     │
     ├── TCP 22  → OPEN
     ├── TCP 80  → OPEN
     └── TCP 443 → CLOSED
```

Dengan demikian, status perangkat dan status service tidak tercampur.

Model konseptualnya:

```text
DEVICE
│
├── Host Status
│     └── Online / Offline
│
└── Services
      ├── TCP 22
      │     └── Open / Closed / Timeout
      │
      ├── TCP 80
      │     └── Open / Closed / Timeout
      │
      └── TCP 443
            └── Open / Closed / Timeout
```

---

# 7. Logika Monitoring

Siklus monitoring TCP secara konseptual:

```text
Start Monitoring Cycle
        ↓
Ambil konfigurasi TCP
        ↓
Tentukan Device
        ↓
Tentukan Host/IP
        ↓
Tentukan TCP Port
        ↓
Lakukan TCP Connection Check
        ↓
Terima hasil pemeriksaan
        ↓
Klasifikasikan hasil
        ↓
Simpan hasil monitoring
        ↓
Bandingkan dengan status sebelumnya
        ↓
Ada perubahan status?
       / \
     Ya   Tidak
     ↓      ↓
 Generate   Lanjut
 Alert      Monitoring
     ↓
 Notification
     ↓
Monitoring berikutnya
```

Yang penting adalah TCP check bukan hanya menghasilkan “berhasil/gagal”, tetapi menghasilkan kondisi yang dapat diproses oleh monitoring engine.

Contohnya:

```text
OPEN
CLOSED / REFUSED
TIMEOUT
UNREACHABLE
```

Namun klasifikasi detail status tersebut bisa ditentukan kemudian pada tahap technical design.

---

# 8. Status Transition

Salah satu bagian terpenting adalah perubahan status.

Misalnya kondisi awal:

```text
TCP 443
OPEN
```

Kemudian pada pemeriksaan berikutnya:

```text
TCP 443
CLOSED
```

Maka sistem mengenali:

```text
OPEN
  ↓
CLOSED
  ↓
Status Change
  ↓
Potential Anomaly
```

Sebaliknya:

```text
CLOSED
  ↓
OPEN
  ↓
Recovery
```

Jadi sistem tidak hanya menyimpan kondisi saat ini, tetapi juga dapat mengenali perubahan kondisi.

---

# 9. Failure Threshold

Untuk menghindari alert akibat gangguan sesaat atau hasil pemeriksaan yang tidak konsisten, TCP Monitoring dapat menggunakan konsep failure threshold yang konsisten dengan mekanisme monitoring GAMON.

Contoh:

```text
Failure Threshold = 3
```

Pemeriksaan:

```text
Check 1 → CLOSED
Check 2 → CLOSED
Check 3 → CLOSED
```

Baru setelah threshold tercapai:

```text
Service Down
     ↓
Create Alert
     ↓
Send Notification
```

Dengan demikian satu kali kegagalan belum langsung dianggap sebagai gangguan permanen.

Konsep ini juga membuat perilaku ICMP dan TCP lebih konsisten di dalam Monitoring Engine.

---

# 10. Integrasi dengan Alert Engine

Ini bagian yang menurut gua paling penting untuk dikaitkan dengan tujuan awal GAMON.

TCP Monitoring tidak boleh berhenti pada:

```text
TCP 443 = CLOSED
```

Data tersebut harus diteruskan ke mekanisme deteksi anomali.

Konsep keseluruhannya:

```text
TCP Monitor
     ↓
TCP Result
     ↓
Result Processor
     ↓
Status Evaluation
     ↓
Failure Threshold
     ↓
Anomaly Detected
     ↓
Alert Engine
```

Kemudian:

```text
Alert Engine
     ├── Dashboard
     ├── WebSocket
     └── Telegram
```

Dengan begitu, NOC tidak perlu membuka halaman monitoring secara terus-menerus untuk mengetahui adanya gangguan.

---

# 11. Desain Alert

Alert TCP harus memberikan informasi yang cukup agar NOC langsung memahami apa yang bermasalah.

Contoh konseptual:

```text
SERVICE DOWN

Device  : Web Server 01
Host    : 192.168.1.10
Service : HTTPS
Port    : TCP 443

Status  : DOWN
Detected: 14:32:17
```

Bukan hanya:

```text
Device Offline
```

Karena informasi tersebut salah apabila ICMP masih menunjukkan perangkat online.

Alert yang benar:

```text
Host      → ONLINE
Service   → DOWN
Port      → TCP 443
```

Ini membuat alert lebih actionable.

---

# 12. Recovery Notification

Monitoring juga harus memperhatikan kondisi ketika service kembali normal.

Contoh:

```text
14:32
TCP 443 → DOWN
        ↓
Alert Created
        ↓
Telegram Notification
```

Kemudian:

```text
14:35
TCP 443 → OPEN
        ↓
Recovery Detected
        ↓
Alert Resolved
        ↓
Recovery Notification
```

Contoh informasi:

```text
SERVICE RECOVERED

Device  : Web Server 01
Service : HTTPS
Port    : TCP 443

Status  : UP
Recovered: 14:35:41
```

Dengan begitu NOC mengetahui bukan hanya kapan gangguan terjadi, tetapi juga kapan service kembali normal.

---

# 13. Integrasi dengan Dashboard

Dashboard tidak perlu membuat halaman terpisah khusus untuk TCP jika tidak diperlukan.

Informasi dapat diintegrasikan ke monitoring device.

Contoh:

```text
Web Server 01
192.168.1.10

Host
● ONLINE
Latency: 2 ms

Services
● TCP 22     OPEN
● TCP 80     OPEN
● TCP 443    DOWN
```

Dengan desain tersebut, NOC dapat melihat hubungan antara:

```text
Device
  ↓
Host Status
  ↓
Service Status
```

secara langsung.

---

# 14. Konsep Data Monitoring

Secara konseptual, setiap hasil TCP monitoring memiliki informasi seperti:

```text
Device
Host
Port
Protocol
Status
Timestamp
Response Time
```

Contoh:

```text
Device       : Web Server 01
Host         : 192.168.1.10
Protocol     : TCP
Port         : 443
Status       : OPEN
Response Time: 12 ms
Timestamp    : 14:30:00
```

Data tersebut nantinya dapat digunakan untuk:

```text
Current Status
History
Alert
Recovery
Dashboard
```

Detail struktur database belum perlu ditentukan dalam PRD ini.

---

# 15. Hubungan dengan ICMP

Desain GAMON sebaiknya tidak menganggap ICMP dan TCP sebagai dua sistem yang terpisah.

Keduanya berada di dalam satu Monitoring Engine.

```text
              Monitoring Engine
                     │
          ┌──────────┴──────────┐
          │                     │
      ICMP Monitor          TCP Monitor
          │                     │
     Host Status            Service Status
          │                     │
          └──────────┬──────────┘
                     ↓
               Result Processing
                     ↓
                Alert Engine
```

Hal ini membuat GAMON memiliki konsep monitoring berlapis.

```text
Layer 1
Host Availability
ICMP

Layer 2
Service Availability
TCP

Layer 3
Application Availability
HTTP/HTTPS
```

HTTP/HTTPS dapat ditambahkan kemudian tanpa harus mengubah konsep dasar Monitoring Engine.

---

# 16. Contoh Skenario

Skenario normal:

```text
Web Server
        │
        ├── ICMP → ONLINE
        ├── TCP 22 → OPEN
        └── TCP 443 → OPEN
```

Tidak ada anomaly.

Skenario service mengalami gangguan:

```text
Web Server
        │
        ├── ICMP → ONLINE
        ├── TCP 22 → OPEN
        └── TCP 443 → CLOSED
                         ↓
                    Threshold
                         ↓
                  Service Anomaly
                         ↓
                  Alert + Telegram
```

Skenario server mati:

```text
Web Server
        │
        ├── ICMP → OFFLINE
        ├── TCP 22 → gagal
        └── TCP 443 → gagal
```

Dalam kondisi ini, GAMON dapat mengetahui bahwa permasalahan berada pada level host, bukan hanya service.

Skenario recovery:

```text
Sebelumnya:

ICMP → ONLINE
TCP 443 → CLOSED

Kemudian:

ICMP → ONLINE
TCP 443 → OPEN
```

Maka:

```text
Service Recovery
        ↓
Alert Resolved
        ↓
Recovery Notification
```

---

# 17. Batasan Fitur

TCP Port Monitoring tidak dimaksudkan untuk mengetahui kondisi internal service secara penuh.

Jika:

```text
TCP 443 → OPEN
```

hal tersebut hanya menunjukkan bahwa koneksi TCP ke port tersebut dapat dilakukan.

Belum berarti:

```text
Website berjalan dengan benar
Database normal
API tidak error
Aplikasi tidak mengalami bug
```

Untuk kebutuhan tersebut diperlukan monitoring pada layer yang lebih tinggi seperti HTTP/HTTPS atau mekanisme lain.

Jadi batasan sistem:

```text
ICMP
→ Host Availability

TCP
→ Service/Port Availability

HTTP/HTTPS
→ Application-Level Availability
```

Pembagian ini penting supaya GAMON tidak mengklaim kemampuan yang sebenarnya belum dimiliki.

---

# 18. Hasil yang Diharapkan

Setelah TCP Port Monitoring diterapkan, GAMON tidak lagi hanya memberikan informasi:

```text
"Apakah perangkat hidup?"
```

tetapi juga:

```text
"Apakah service yang dipantau pada perangkat tersebut tersedia?"
```

Sehingga alur deteksi anomali menjadi:

```text
Infrastructure
      ↓
Monitoring
      ↓
Host Availability
      +
Service Availability
      ↓
Anomaly Detection
      ↓
Alert
      ↓
Notification
      ↓
NOC Awareness
      ↓
Troubleshooting
```

Ini menurut gua adalah positioning yang paling kuat untuk fitur TCP tersebut.

Intinya, **TCP Port Monitoring bukan sekadar penambahan metode check baru**. Secara desain, dia memperluas model GAMON dari monitoring perangkat menjadi monitoring kondisi layanan pada perangkat. Dan karena hasilnya masuk ke Alert Engine + Telegram/WebSocket yang sudah ada, fitur ini juga memperkuat tujuan utama GAMON: mengurangi kemungkinan NOC tidak mengetahui adanya anomali.