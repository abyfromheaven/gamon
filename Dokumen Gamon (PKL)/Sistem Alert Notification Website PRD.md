# PRD
# Feature: Alert Notification

---

# 1. Overview

## Ringkasan

Alert Notification merupakan fitur pada aplikasi GAMON yang berfungsi untuk memberikan informasi secara real-time kepada administrator ketika terjadi perubahan kondisi perangkat jaringan yang sedang dimonitor.

Fitur ini dirancang agar administrator dapat mengetahui adanya gangguan tanpa harus melakukan pengecekan perangkat secara manual pada halaman Monitoring.

Alert Notification hanya menampilkan notifikasi ketika terjadi perubahan status perangkat (State Transition), sehingga tidak menghasilkan notifikasi yang berulang selama kondisi perangkat tidak berubah.

---

# 2. Background

Dalam proses monitoring infrastruktur TI, administrator sering kali harus membuka halaman monitoring secara terus-menerus untuk memastikan seluruh perangkat berada dalam kondisi normal.

Pendekatan tersebut memiliki beberapa kelemahan, antara lain:

- administrator dapat terlambat mengetahui gangguan;
- administrator harus melakukan monitoring secara aktif;
- potensi terlewatnya perangkat yang mengalami gangguan cukup tinggi.

Oleh karena itu dibutuhkan sebuah mekanisme Alert Notification yang mampu memberikan pemberitahuan secara otomatis ketika terjadi perubahan kondisi perangkat.

---

# 3. Objective

Fitur Alert Notification bertujuan untuk:

- Memberikan pemberitahuan secara otomatis.
- Mengurangi monitoring manual.
- Mempercepat respon administrator.
- Mengurangi kemungkinan perangkat down tidak diketahui.

---

# 4. Scope

Versi pertama Alert Notification mencakup:

✅ Critical Alert

✅ Recovery Alert

✅ Alert History

✅ Acknowledge Alert

Tidak termasuk:

❌ Telegram

❌ WhatsApp

❌ Email

❌ Push Notification Browser

---

# 5. User

Primary User

Administrator

---

# 6. Trigger

Alert hanya dibuat apabila terjadi perubahan status perangkat.

```
ONLINE
↓

OFFLINE
```

atau

```
OFFLINE
↓

ONLINE
```

Selain kondisi tersebut sistem tidak membuat alert baru.

---

# 7. Alert Generation Logic

## Critical Alert

Kondisi

- Status sebelumnya Online
- Status sekarang Offline

Output

Critical Alert dibuat.

---

## Recovery Alert

Kondisi

- Status sebelumnya Offline
- Status sekarang Online

Output

Recovery Alert dibuat.

---

## No Alert

Jika status tetap Online

```
Online

↓

Online
```

Tidak membuat alert.

Jika status tetap Offline

```
Offline

↓

Offline
```

Tidak membuat alert.

---

# 8. Validation Mechanism

Alert tidak langsung dibuat setelah satu kali ping gagal.

Sistem menggunakan Failure Threshold.

Contoh

```
Failure Threshold = 3
```

```
Timeout

↓

Failure 1
```

Belum Alert.

```
Timeout

↓

Failure 2
```

Belum Alert.

```
Timeout

↓

Failure 3
```

Generate Critical Alert.

---

Recovery menggunakan Success Threshold.

Contoh

```
Success Threshold = 2
```

```
Success

↓

1
```

Belum Recovery.

```
Success

↓

2
```

Generate Recovery Alert.

---

# 9. State Transition

```
          Failure >= Threshold

ONLINE ----------------------> OFFLINE
   ▲                              │
   │                              │
   │                              │
   └------------------------------┘
      Success >= Threshold
```

---

# 10. Alert Severity

Versi pertama menggunakan dua severity.

| Severity | Digunakan untuk |
|------------|----------------|
| Critical | Device Down |
| Info | Device Recovery |

---

# 11. Alert Data

Setiap alert menyimpan:

- Alert ID
- Device
- Severity
- Previous Status
- Current Status
- Timestamp
- Message
- Acknowledged

---

# 12. Alert Lifecycle

```
Monitoring

↓

State Changed

↓

Generate Alert

↓

Show Notification

↓

Save Alert

↓

Administrator Open Alert

↓

Administrator Click Acknowledge

↓

Alert Closed
```

---

# 13. Acknowledge

Administrator dapat memberikan status **Acknowledged**.

Tujuan:

- Menandakan alert telah diketahui.
- Menghindari kebingungan antara alert baru dan alert lama.
- Tidak menghapus histori.

---

# 14. Cooldown

Sistem tidak membuat alert baru apabila status perangkat tidak berubah.

Contoh

```
Offline

↓

Offline

↓

Offline
```

Tetap hanya memiliki satu Critical Alert.

---

# 15. Functional Requirements

Sistem harus mampu:

- menghasilkan Critical Alert;
- menghasilkan Recovery Alert;
- menyimpan histori alert;
- menampilkan alert terbaru;
- melakukan acknowledge;
- mengubah status alert menjadi Read;
- mengurutkan alert berdasarkan waktu.

---

# 16. Non Functional Requirements

- Alert muncul kurang dari satu siklus monitoring setelah status tervalidasi berubah.
- Tidak menghasilkan duplicate alert.
- Tidak menghasilkan spam alert.
- Seluruh histori tersimpan.

---

# 17. Future Enhancement

- Telegram Notification
- WhatsApp Notification
- Email Notification
- Browser Push Notification
- Sound Notification
- Escalation Alert
- Multi Severity