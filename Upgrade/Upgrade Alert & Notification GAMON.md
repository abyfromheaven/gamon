## 1. Latar Belakang

GAMON dirancang untuk membantu NOC mengetahui kondisi infrastruktur jaringan secara realtime dan mendeteksi perubahan kondisi perangkat maupun service.

Sistem monitoring yang hanya mampu mendeteksi kondisi abnormal belum cukup apabila informasi tersebut tidak sampai kepada operator.

Contohnya:

```text
Monitoring
    ↓
Device / Service DOWN
    ↓
Alert dibuat
    ↓
???
```

Apabila operator tidak sedang melihat dashboard ketika kejadian berlangsung, anomaly dapat terlewat.

Oleh karena itu, sistem Alert & Notification perlu berfungsi sebagai lapisan yang mengubah hasil monitoring menjadi informasi yang secara aktif disampaikan kepada NOC.

Konsep upgrade:

```text
Monitoring
    ↓
Anomaly Detection
    ↓
Alert Management
    ↓
Notification
    ├── Dashboard / Realtime
    ├── Telegram
    └── Sound Alert
```

---

## 2. Permasalahan

Permasalahan utama yang ingin diselesaikan adalah kemungkinan NOC tidak menyadari adanya anomaly yang telah terdeteksi oleh sistem.

Kondisi tersebut dapat terjadi apabila:

- NOC tidak sedang membuka halaman monitoring.
    
- Perubahan status terjadi tanpa perhatian operator.
    
- Informasi anomaly kurang spesifik.
    
- Alert hanya muncul secara visual.
    
- Tidak terdapat indikator yang cukup kuat untuk membedakan kondisi normal dan abnormal.
    
- Operator tidak mengetahui apakah anomaly masih berlangsung atau sudah pulih.
    

Maka sistem perlu memberikan mekanisme notification yang bersifat lebih aktif dan informatif.

---

## 3. Tujuan

Upgrade Alert & Notification bertujuan untuk:

1. Memastikan anomaly yang terdeteksi monitoring engine menghasilkan alert.
    
2. Memberikan informasi anomaly secara realtime kepada NOC.
    
3. Memberikan informasi yang cukup untuk mengidentifikasi sumber masalah.
    
4. Menyediakan notification melalui beberapa channel.
    
5. Memberikan sound alert sebagai indikator lokal ketika terjadi anomaly.
    
6. Memungkinkan sound alert diaktifkan atau dinonaktifkan melalui Settings.
    
7. Membedakan notification untuk anomaly dan recovery.
    
8. Menyediakan lifecycle alert yang jelas dari muncul sampai terselesaikan.
    
9. Mengurangi ketergantungan NOC terhadap pemantauan dashboard secara terus-menerus.
    

---

# 4. Konsep Alert vs Notification

Ini sebaiknya dibedakan sejak awal.

**Alert** adalah representasi kejadian/anomaly di dalam sistem.

**Notification** adalah cara sistem menyampaikan alert tersebut kepada NOC.

Jadi:

```text
                 MONITORING
                     ↓
               Anomaly Event
                     ↓
                ALERT ENGINE
                     ↓
                   ALERT
                     │
          ┌──────────┼──────────┐
          ↓          ↓          ↓
      Dashboard   Telegram    Sound
```

Contoh:

```text
TCP 443 Web Server DOWN
```

adalah event/alert.

Sedangkan:

```text
Telegram message
Dashboard notification
Sound
```

adalah notification channel.

Pemisahan ini penting karena satu alert dapat dikirim melalui beberapa channel sekaligus.

---

# 5. Jenis Event

Sistem perlu memiliki konsep event agar Alert Engine tidak hanya mengetahui “ada error”, tetapi mengetahui jenis perubahan yang terjadi.

Untuk tahap awal:

```text
ALERT
│
├── Host Down
├── Service Down
└── Monitoring Failure

RECOVERY
│
├── Host Recovered
└── Service Recovered
```

Contoh dengan TCP Monitoring:

```text
TCP 443
OPEN
 ↓
CLOSED
 ↓
SERVICE DOWN
```

menghasilkan alert.

Kemudian:

```text
TCP 443
CLOSED
 ↓
OPEN
 ↓
SERVICE RECOVERED
```

menghasilkan recovery event.

---

# 6. Alert Lifecycle

Alert perlu memiliki lifecycle yang jelas.

Konsep:

```text
                    Detection
                       ↓
                    CREATED
                       ↓
                    ONGOING
                       ↓
                ┌──────┴──────┐
                │             │
           Acknowledged    Still Down
                │             │
                │             │
                └──────┬──────┘
                       ↓
                    RESOLVED
```

Dalam implementasi GAMON yang sudah memiliki status `ongoing`, `resolved`, dan `acknowledged`, upgrade ini dapat memperjelas fungsi masing-masing status daripada membuat sistem alert baru yang sepenuhnya berbeda.

### Ongoing

Anomaly masih berlangsung.

```text
Server Web
TCP 443
DOWN

Alert: ONGOING
```

### Acknowledged

NOC sudah mengetahui/mengakui anomaly tersebut.

```text
NOC
 ↓
Open Alert
 ↓
Acknowledge
```

Acknowledgement tidak berarti masalah sudah selesai.

```text
ACKNOWLEDGED ≠ RESOLVED
```

### Resolved

Monitoring mendeteksi bahwa kondisi telah kembali normal.

```text
TCP 443
DOWN
 ↓
OPEN
 ↓
RESOLVED
```

---

# 7. Desain Informasi Alert

Alert harus memberikan informasi yang cukup untuk menjawab tiga pertanyaan NOC:

> Apa yang bermasalah?

> Perangkat mana yang bermasalah?

> Kapan masalah terjadi?

Untuk TCP:

```text
SERVICE DOWN

Device  : Web Server 01
IP      : 192.168.1.10
Service : HTTPS
Port    : TCP 443

Status  : DOWN
Detected: 14:32:17
```

Untuk ICMP:

```text
HOST DOWN

Device  : Router Core
IP      : 192.168.1.1

Status  : OFFLINE
Detected: 14:32:17
```

Jadi notification tidak hanya berbunyi:

```text
"Anomaly detected"
```

karena informasi tersebut terlalu umum untuk operator.

---

# 8. Sound Alert

Sound Alert menjadi notification channel tambahan untuk NOC yang sedang berada di depan komputer monitoring.

Konsepnya:

```text
Anomaly
   ↓
Alert Created
   ↓
Sound Alert
   +
Dashboard Alert
   +
Telegram
```

Sound berfungsi sebagai **attention mechanism**, bukan sebagai pengganti Telegram atau dashboard.

Tujuannya sederhana:

> Ketika anomaly terjadi saat NOC tidak sedang memperhatikan dashboard, suara dapat menarik perhatian operator.

---

# 9. Sound Alert Setting

Sound Alert harus dapat dikontrol melalui Settings.

Contoh:

```text
Notification Settings

Sound Alert
[ ON ]

Telegram Notification
[ ON ]
```

atau:

```text
Sound Alert
● Enabled
○ Disabled
```

Ketika enabled:

```text
Alert Created
     ↓
Sound dimainkan
```

Ketika disabled:

```text
Alert Created
     ↓
Sound tidak dimainkan
```

Namun alert tetap dibuat.

Ini penting:

```text
Sound OFF
≠
Alert OFF
```

Sound hanya salah satu notification channel.

---

# 10. Perilaku Sound Alert

Gua sarankan sound **tidak dimainkan pada setiap monitoring cycle yang gagal**.

Contoh yang buruk:

```text
10:00:00 → DOWN → beep
10:00:30 → DOWN → beep
10:01:00 → DOWN → beep
10:01:30 → DOWN → beep
10:02:00 → DOWN → beep
...
```

Ini bisa menjadi sangat mengganggu.

Lebih baik sound dikaitkan dengan **event**, bukan setiap hasil monitoring.

Contoh:

```text
OPEN
 ↓
OPEN
 ↓
OPEN
 ↓
CLOSED
 ↓
ANOMALY DETECTED
 ↓
🔊 Sound
```

Kemudian:

```text
CLOSED
 ↓
CLOSED
 ↓
CLOSED
```

tidak menghasilkan sound baru karena anomaly yang sama masih berlangsung.

Ketika recovery:

```text
CLOSED
 ↓
OPEN
 ↓
RECOVERY
```

sistem dapat memberikan notification recovery. Sound recovery dapat dibuat opsional pada pengembangan berikutnya; untuk versi awal, lebih sederhana apabila sound hanya digunakan untuk anomaly baru.

---

# 11. Deduplication Alert

Ini penting supaya satu masalah tidak menghasilkan puluhan alert.

Misalnya:

```text
TCP 443 DOWN
```

dan monitoring berjalan setiap 30 detik.

Sistem tidak boleh membuat:

```text
Alert #1
Alert #2
Alert #3
Alert #4
...
```

selama anomaly yang sama masih berlangsung.

Konsep:

```text
OPEN
 ↓
CLOSED
 ↓
Create Alert
 ↓
ONGOING
 ↓
CLOSED
 ↓
Update existing Alert
 ↓
CLOSED
 ↓
Update existing Alert
```

Baru ketika:

```text
CLOSED
 ↓
OPEN
```

alert dianggap resolved.

Jika kemudian:

```text
OPEN
 ↓
CLOSED
```

lagi, maka itu merupakan incident baru.

---

# 12. Notification Routing

Alert Engine sebaiknya tidak langsung bergantung kepada satu jenis notification.

Konsep:

```text
                   ALERT ENGINE
                        │
                        ↓
                 Notification Manager
                        │
             ┌──────────┼──────────┐
             ↓          ↓          ↓
         Dashboard   Telegram    Sound
```

Dengan demikian setiap channel dapat dikontrol secara independen.

Contoh Settings:

```text
Notification

Dashboard Notification    [ON]
Telegram Notification     [ON]
Sound Alert               [ON]
```

Kemudian:

```text
Telegram OFF
```

tidak berarti alert tidak dibuat.

Yang berubah hanya:

```text
Telegram notification → tidak dikirim
```

Dashboard dan sound tetap dapat bekerja.

---

# 13. Realtime Dashboard Notification

Ketika alert dibuat, dashboard harus menerima perubahan secara realtime melalui mekanisme WebSocket yang sudah digunakan GAMON.

Konsep:

```text
Monitoring Engine
       ↓
Anomaly
       ↓
Alert Engine
       ↓
WebSocket
       ↓
Dashboard
       ↓
New Alert
```

Contohnya jumlah alert dapat berubah:

```text
Alerts
[ 3 ]
```

menjadi:

```text
Alerts
[ 4 ]
```

tanpa operator melakukan refresh halaman.

Ini mempertahankan konsep realtime yang sudah menjadi bagian dari GAMON.

---

# 14. Telegram Notification

Telegram digunakan sebagai notification channel eksternal.

Konsep:

```text
Anomaly
   ↓
Alert Engine
   ↓
Telegram Notification
   ↓
NOC / Telegram
```

Pesan harus memiliki konteks yang cukup.

Contoh:

```text
🚨 SERVICE DOWN

Device: Web Server 01
IP: 192.168.1.10
Service: HTTPS
Port: TCP 443

Detected: 14:32:17
Status: DOWN
```

Recovery:

```text
SERVICE RECOVERED

Device: Web Server 01
Service: HTTPS
Port: TCP 443

Recovered: 14:35:41
Status: UP
```

Dengan begitu Telegram bukan sekadar “alarm”, tetapi membawa informasi yang dapat digunakan untuk menentukan tindakan selanjutnya.

---

# 15. Notification Priority

Untuk pengembangan awal, tidak perlu membuat sistem severity yang terlalu kompleks.

Cukup membedakan event berdasarkan dampaknya.

Contoh:

```text
HOST DOWN
     ↓
Critical

SERVICE DOWN
     ↓
Warning / High
```

Tetapi penentuan severity sebaiknya mengikuti kebutuhan operasional GAMON dan konfigurasi yang sudah ada, bukan membuat klasifikasi baru yang terlalu rumit.

Yang penting NOC dapat langsung membedakan:

```text
Normal
Warning
Critical
Recovery
```

jika memang diperlukan pada UI.

---

# 16. Alert History

Alert yang telah selesai tidak boleh hilang begitu saja.

Konsep:

```text
Current Alerts
       +
Alert History
```

Contoh:

```text
Alert History

14:32  Web Server 01 — TCP 443 DOWN
14:35  Web Server 01 — TCP 443 RECOVERED
15:12  Router Core — HOST DOWN
15:14  Router Core — HOST RECOVERED
```

History memungkinkan NOC atau administrator mengetahui kejadian yang pernah terjadi.

---

# 17. Hubungan dengan Monitoring Engine

Alert Engine tidak melakukan monitoring secara langsung.

Pembagian tanggung jawab:

```text
Monitoring Engine
→ melakukan pemeriksaan

Result Processor
→ memahami hasil

Alert Engine
→ menentukan apakah terjadi anomaly/event

Notification Manager
→ menentukan bagaimana event diberitahukan
```

Secara konseptual:

```text
             MONITORING ENGINE
                    ↓
              Monitoring Result
                    ↓
             Result Evaluation
                    ↓
              Alert Detection
                    ↓
                Alert Engine
                    ↓
          ┌─────────┼──────────┐
          ↓         ↓          ↓
      WebSocket  Telegram    Sound
```

Ini membuat desain lebih bersih dan nantinya memudahkan penambahan notification channel baru.

---

# 18. Skenario Utama

### Skenario 1 — Host Down

```text
ICMP
ONLINE
 ↓
OFFLINE
 ↓
Failure Threshold
 ↓
Host Down Alert
 ↓
Dashboard
 ↓
Telegram
 ↓
Sound
```

Jika sound enabled, suara dimainkan.

Jika sound disabled:

```text
Dashboard ✓
Telegram  ✓
Sound     ✗
```

Alert tetap ada.

### Skenario 2 — Service Down

```text
ICMP
ONLINE

TCP 443
OPEN
 ↓
CLOSED
 ↓
Failure Threshold
 ↓
Service Down Alert
 ↓
Dashboard
Telegram
Sound
```

Hal ini menunjukkan bahwa GAMON dapat membedakan:

```text
Host masih hidup
Service bermasalah
```

### Skenario 3 — Recovery

```text
Service DOWN
     ↓
Alert ONGOING
     ↓
Service OPEN
     ↓
Recovery Detected
     ↓
Alert RESOLVED
     ↓
Recovery Notification
```

### Skenario 4 — Sound Disabled

```text
Anomaly
   ↓
Alert Created
   ↓
Dashboard ✓
Telegram  ✓
Sound     ✗
```

Tidak ada dampak terhadap proses monitoring.

---

# 19. Pengaturan yang Dibutuhkan

Untuk versi awal, Settings jangan dibuat terlalu banyak.

Minimal:

```text
Notification Settings

Telegram Notification
[ ON / OFF ]

Sound Alert
[ ON / OFF ]
```

Jika sistem nantinya membutuhkan:

```text
Sound Volume
Sound Type
Repeat Sound
```

itu dapat menjadi pengembangan berikutnya.

Untuk H-1, cukup:

```text
Sound Alert
ON / OFF
```

karena itu sudah menjawab kebutuhan utama.

---

# 20. Non-Functional Requirement

Alert dan notification harus memenuhi beberapa prinsip:

**Realtime**  
Alert muncul sesegera mungkin setelah anomaly memenuhi kondisi deteksi.

**Reliable**  
Satu anomaly tidak menghasilkan notification berulang tanpa alasan.

**Informative**  
Notification menyebutkan perangkat, IP, service/port atau jenis anomaly, waktu, dan status.

**Configurable**  
NOC dapat mengaktifkan atau menonaktifkan channel tertentu.

**Non-blocking**  
Kegagalan satu notification channel seharusnya tidak menghentikan monitoring.

Contohnya:

```text
Telegram gagal
     ↓
Monitoring tetap berjalan
     ↓
Dashboard tetap menerima alert
     ↓
Sound tetap dapat bekerja
```

Ini penting secara arsitektur.

---

# 21. Batasan Fitur

Upgrade ini tidak bertujuan membuat sistem incident management penuh seperti platform enterprise.

Tidak perlu langsung memasukkan:

```text
Escalation Policy
On-call Scheduling
Multi-level Escalation
Email
SMS
WhatsApp
Pager
```

Fokus GAMON tetap pada:

```text
Detect
 ↓
Alert
 ↓
Notify
 ↓
Acknowledge
 ↓
Recover
```

Dengan scope tersebut, fitur tetap relevan dengan permasalahan yang ingin diselesaikan.

---

# 22. Konsep Akhir

Setelah upgrade, alur GAMON menjadi:

```text
                    INFRASTRUCTURE
                          │
                          ↓
                  MONITORING ENGINE
                          │
             ┌────────────┴────────────┐
             ↓                         ↓
            ICMP                       TCP
       Host Availability         Service Availability
             │                         │
             └────────────┬────────────┘
                          ↓
                   RESULT EVALUATION
                          ↓
                    ANOMALY DETECTED
                          ↓
                     ALERT ENGINE
                          │
                 ┌────────┼─────────┐
                 ↓        ↓         ↓
             Dashboard Telegram   Sound
             WebSocket
                 │
                 └──────────────┐
                                ↓
                              NOC
                                ↓
                          Acknowledgement
                                ↓
                            Monitoring
                                ↓
                            Recovery
                                ↓
                       Alert → RESOLVED
                                ↓
                       Recovery Notification
```

Dengan desain ini, upgrade Alert & Notification bukan cuma “nambah bunyi ketika error”. Konsepnya adalah membuat GAMON memiliki **alur awareness yang lengkap**: sistem mendeteksi anomaly, membuat incident/alert, menyampaikan informasi melalui beberapa channel, memungkinkan NOC mengetahui dan mengakui kejadian, lalu mendeteksi recovery.

Dan sound alert punya posisi yang jelas: **local attention mechanism untuk menarik perhatian NOC**, dengan `ON/OFF` sebagai konfigurasi independen. Ini membuat fitur tersebut relevan langsung dengan masalah awal GAMON tanpa membuat scope-nya melebar terlalu jauh.