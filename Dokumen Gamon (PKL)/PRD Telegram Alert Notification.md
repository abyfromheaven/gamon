Berikut adalah PRD yang saya sarankan untuk **fitur Telegram Bot Alert Notification** pada GAMON. Dokumen ini **khusus membahas integrasi Telegram sebagai media notifikasi**, bukan sistem monitoring ataupun Alert Notification inti yang sudah dibuat sebelumnya.

````md
# Product Requirements Document (PRD)
# GAMON (Garda Monitoring)

## Feature: Telegram Bot Alert Notification

---

# Document Information

| Item | Value |
|------|------|
| Feature | Telegram Bot Alert Notification |
| Version | 1.0 |
| Status | Draft |
| Priority | High |
| Platform | Web Application |
| Dependencies | Alert Notification Module |

---

# 1. Overview

## Deskripsi

Telegram Bot Alert Notification merupakan fitur yang berfungsi mengirimkan informasi gangguan perangkat secara otomatis kepada administrator melalui aplikasi Telegram.

Fitur ini bekerja sebagai **Notification Delivery Channel**, yaitu menerima alert yang telah dihasilkan oleh Alert Notification Engine, kemudian mengirimkan informasi tersebut ke Telegram menggunakan Telegram Bot API.

Telegram Notification bukan bagian dari proses monitoring maupun pendeteksian perangkat. Seluruh proses deteksi tetap dilakukan oleh Alert Engine, sedangkan Telegram hanya bertugas sebagai media penyampaian informasi.

---

# 2. Background

Administrator tidak selalu berada di depan dashboard monitoring selama jam operasional. Walaupun aplikasi GAMON telah menyediakan Alert Notification pada website, administrator tetap harus membuka aplikasi untuk mengetahui adanya gangguan.

Dalam kondisi tertentu, keterlambatan membuka dashboard dapat menyebabkan keterlambatan penanganan terhadap perangkat yang mengalami gangguan.

Oleh karena itu diperlukan media notifikasi eksternal yang mampu memberikan informasi secara real-time langsung ke perangkat administrator.

Telegram dipilih karena:

- Gratis.
- Ringan.
- Mudah diintegrasikan.
- Memiliki Bot API yang stabil.
- Mendukung Group maupun Personal Chat.
- Tidak memerlukan infrastruktur tambahan.

---

# 3. Objectives

Fitur Telegram Notification bertujuan untuk:

- Mengirimkan alert secara real-time.
- Mengurangi waktu respon administrator.
- Memberikan pemberitahuan tanpa harus membuka dashboard.
- Menjadi media distribusi alert eksternal.
- Mendukung monitoring jarak jauh.

---

# 4. Scope

## Included

- Critical Alert
- Recovery Alert
- Telegram Bot API
- Group Chat
- Personal Chat
- Retry Mechanism
- Delivery Status

## Excluded

- WhatsApp
- Discord
- Email
- SMS
- Push Browser Notification
- Escalation Notification

---

# 5. Actors

Primary User

Administrator

Secondary System

Telegram Bot API

---

# 6. Dependencies

Fitur ini bergantung pada:

- Monitoring Engine
- Alert Notification Module
- Telegram Bot API

Telegram Notification tidak dapat berjalan apabila Alert Notification belum menghasilkan Alert Event.

---

# 7. System Architecture

```
Monitoring Engine
        │
        ▼
Alert Notification Engine
        │
        ▼
Generate Alert Event
        │
        ▼
Save Alert Database
        │
        ▼
Notification Service
        │
        ▼
Telegram Bot API
        │
        ▼
Administrator
```

---

# 8. Workflow

```
Device Monitoring

↓

State Transition

↓

Generate Alert

↓

Save Alert

↓

Notification Queue

↓

Telegram Service

↓

Send Message

↓

Update Delivery Status
```

---

# 9. Alert Trigger

Telegram hanya dikirim ketika terjadi perubahan status perangkat.

## Critical

```
ONLINE

↓

OFFLINE
```

Generate Telegram Alert.

---

## Recovery

```
OFFLINE

↓

ONLINE
```

Generate Recovery Message.

---

## Tidak Mengirim

```
ONLINE

↓

ONLINE
```

Tidak mengirim Telegram.

```
OFFLINE

↓

OFFLINE
```

Tidak mengirim Telegram.

---

# 10. Notification Flow

```
Monitoring

↓

Threshold Validation

↓

Status Changed

↓

Generate Alert

↓

Database

↓

Telegram Queue

↓

Telegram API

↓

Administrator
```

---

# 11. Delivery Mechanism

Telegram dikirim secara asynchronous.

Artinya proses monitoring tidak menunggu Telegram selesai dikirim.

Keuntungan:

- Monitoring tetap berjalan.
- Telegram gagal tidak mempengaruhi monitoring.
- Alert tetap tersimpan.

---

# 12. Retry Mechanism

Apabila Telegram API gagal merespon.

```
Attempt 1

↓

Failed

↓

Wait 5 Seconds

↓

Attempt 2

↓

Failed

↓

Wait 10 Seconds

↓

Attempt 3

↓

Success
```

Apabila seluruh percobaan gagal.

Status Delivery menjadi:

Failed

Alert tetap tersedia pada website.

---

# 13. Delivery Status

Setiap pengiriman memiliki status.

| Status | Keterangan |
|---------|------------|
| Pending | Menunggu dikirim |
| Sending | Sedang dikirim |
| Sent | Berhasil dikirim |
| Failed | Gagal dikirim |

---

# 14. Message Structure

Setiap pesan Telegram minimal berisi informasi berikut.

- Severity
- Device Name
- IP Address
- Current Status
- Timestamp
- Alert Message

---

# 15. Critical Alert Template

```
🚨 GAMON ALERT

Severity :
CRITICAL

Device :
Core Router

IP Address :
192.168.10.1

Status :
OFFLINE

Time :
21 Juli 2026
13:20 WIB

Message :
Device is unreachable.
```

---

# 16. Recovery Alert Template

```
✅ GAMON RECOVERY

Device :
Core Router

IP Address :
192.168.10.1

Status :
ONLINE

Recovery Time :
13:24 WIB

Message :
Device is reachable again.
```

---

# 17. Spam Prevention

Telegram tidak dikirim berulang selama status perangkat tidak berubah.

Contoh

```
ONLINE

↓

OFFLINE
```

Telegram dikirim.

```
OFFLINE

↓

OFFLINE

↓

OFFLINE

↓

OFFLINE
```

Tidak ada Telegram tambahan.

Telegram baru dikirim kembali apabila status berubah menjadi Online.

---

# 18. Business Rules

BR-01

Telegram hanya dikirim apabila Alert berhasil dibuat.

---

BR-02

Telegram tidak boleh dikirim sebelum Alert disimpan ke database.

---

BR-03

Telegram tidak boleh dikirim apabila status perangkat tidak berubah.

---

BR-04

Telegram mengikuti hasil validasi Failure Threshold dan Success Threshold.

---

BR-05

Telegram hanya dikirim satu kali untuk setiap perubahan status.

---

# 19. Functional Requirements

FR-01

Sistem mampu mengirim Critical Alert.

---

FR-02

Sistem mampu mengirim Recovery Alert.

---

FR-03

Sistem mampu mengirim ke Group Chat.

---

FR-04

Sistem mampu mengirim ke Personal Chat.

---

FR-05

Sistem menyimpan status pengiriman.

---

FR-06

Sistem melakukan Retry apabila pengiriman gagal.

---

FR-07

Sistem tidak mengirim duplicate notification.

---

# 20. Non Functional Requirements

- Delivery cepat.
- Reliable.
- Tidak duplicate.
- Tidak blocking monitoring.
- Mudah dikonfigurasi.
- Mudah dikembangkan ke media lain.

---

# 21. Configuration

Administrator dapat mengatur:

- Bot Token
- Chat ID
- Enable Telegram Notification
- Enable Critical Alert
- Enable Recovery Alert

---

# 22. Future Enhancement

- Multiple Telegram Group
- Multiple Chat ID
- Discord Notification
- WhatsApp Notification
- Email Notification
- Escalation Alert
- Scheduled Summary Report
- Daily Report
- Weekly Report
- Monthly Report
- Maintenance Mode
- Notification Filter
- Mention Administrator
- Rich Markdown Message
- Inline Button
- Open Alert Dashboard Button

---

# 23. Acceptance Criteria

| ID | Acceptance Criteria |
|----|---------------------|
| AC-01 | Telegram terkirim ketika perangkat berubah dari Online menjadi Offline. |
| AC-02 | Telegram terkirim ketika perangkat berubah dari Offline menjadi Online. |
| AC-03 | Telegram tidak terkirim apabila status perangkat tidak berubah. |
| AC-04 | Alert tetap tersimpan walaupun Telegram gagal dikirim. |
| AC-05 | Retry dijalankan ketika Telegram API gagal. |
| AC-06 | Delivery Status diperbarui setelah proses pengiriman selesai. |

---

# 24. Conclusion

Telegram Bot Alert Notification merupakan fitur pendukung yang berfungsi sebagai media distribusi alert eksternal pada aplikasi GAMON. Fitur ini menerima Alert Event dari Alert Notification Engine, kemudian mengirimkan informasi tersebut kepada administrator melalui Telegram secara asynchronous.

Dengan menerapkan mekanisme Trigger Validation, Retry Mechanism, Spam Prevention, dan Delivery Status Tracking, sistem mampu memberikan notifikasi secara cepat tanpa mengganggu proses monitoring utama. Pendekatan ini menjaga pemisahan tanggung jawab antara proses deteksi gangguan dan proses distribusi notifikasi sehingga arsitektur sistem menjadi lebih modular, andal, dan mudah dikembangkan pada versi berikutnya.
````

## Catatan Arsitektur

Saya justru menyarankan agar Anda menganggap **Telegram Bot Notification sebagai "plugin" atau "integration module"**, bukan fitur inti GAMON. Dengan desain seperti ini, nanti jika ingin menambah **WhatsApp, Discord, Email, atau Microsoft Teams**, Anda tidak perlu mengubah logika monitoring maupun alert. Semua media notifikasi cukup mengonsumsi **Alert Event** yang sama dari modul inti. Pendekatan ini lebih sesuai dengan prinsip _separation of concerns_ dan akan menjadi nilai tambah ketika menjelaskan desain sistem pada BAB IV maupun saat sidang PKL.