# Alert Center Page

## Overview

Halaman Alert Center merupakan halaman yang digunakan untuk menampilkan seluruh informasi mengenai alert atau gangguan yang terjadi pada perangkat yang sedang dimonitor oleh Gamon.

Berbeda dengan halaman Monitoring yang berfokus pada kondisi perangkat saat ini (current status), Alert Center berfokus pada informasi gangguan (incident) yang pernah maupun sedang terjadi selama proses monitoring berlangsung.

Tujuan utama dari halaman ini adalah membantu pengguna mengetahui kapan gangguan terjadi, perangkat apa yang terdampak, tingkat keparahan gangguan, serta mengetahui apakah gangguan tersebut telah terselesaikan atau masih berlangsung.

Alert Center juga berfungsi sebagai pusat riwayat alert (Alert History) yang dihasilkan secara otomatis oleh Gamon.

## Objectives

* Menampilkan seluruh alert yang dihasilkan oleh Gamon.
* Membantu pengguna mengetahui perangkat yang sedang mengalami gangguan.
* Memberikan informasi mengenai status penyelesaian suatu gangguan.
* Menyediakan riwayat alert yang pernah terjadi selama proses monitoring berlangsung.
* Membantu pengguna melakukan identifikasi gangguan dengan lebih cepat.

## Page Responsibility

Halaman Alert Center hanya bertanggung jawab terhadap informasi mengenai gangguan (incident) yang terjadi selama proses monitoring.

Halaman ini tidak digunakan untuk:

* Melakukan monitoring perangkat secara realtime.
* Mengelola konfigurasi perangkat.
* Menampilkan ringkasan sistem monitoring secara keseluruhan.
* Melakukan analisis performa perangkat.

Dengan demikian, Alert Center hanya berfokus pada informasi alert yang dihasilkan oleh Gamon.

## Alert Mechanism

Alert pada Gamon akan dibuat secara otomatis ketika sistem monitoring mendeteksi adanya perubahan status pada suatu perangkat.

Alur kerjanya sebagai berikut:

1. Gamon melakukan monitoring perangkat secara berkala.
2. Perangkat mengalami gangguan atau kondisi tertentu.
3. Status perangkat berubah menjadi Offline ataupun Warning.
4. Gamon secara otomatis membuat Alert baru.
5. Alert akan ditampilkan pada Dashboard sebagai Latest Alert.
6. Alert disimpan pada halaman Alert Center.
7. Apabila perangkat kembali normal, Gamon akan memperbarui status Alert menjadi Resolved.
8. Riwayat Alert akan tetap disimpan sebagai informasi gangguan yang pernah terjadi.

Alert Center tidak hanya menampilkan gangguan yang sedang berlangsung, namun juga menyimpan riwayat gangguan yang telah selesai ditangani.

## Information Displayed

Informasi yang akan ditampilkan pada halaman Alert Center meliputi:

* Alert Title
* Device Name
* Device Type
* Alert Status
* Severity Level
* Started Time
* Resolved Time
* Alert Description

## Alert Status

Status Alert yang digunakan pada Gamon terdiri dari:

* Ongoing

  * Gangguan masih berlangsung dan belum terselesaikan.

* Resolved

  * Gangguan telah selesai dan perangkat telah kembali normal.

Status Alert berbeda dengan Monitoring Status pada halaman Monitoring.

| Monitoring            | Alert Center             |
| --------------------- | ------------------------ |
| Online                | Tidak menghasilkan Alert |
| Offline               | Menghasilkan Alert       |
| Warning               | Menghasilkan Alert       |
| Unknown               | Tidak menghasilkan Alert |
| Current Device Status | Incident History         |

## Severity Level

Untuk MVP Gamon, Severity Level yang digunakan cukup sederhana dan terdiri dari:

* Low
* Medium
* High
* Critical

Severity Level digunakan untuk membantu pengguna mengetahui tingkat keparahan suatu gangguan.

## Sections

### 1. Alert Summary

Section ini digunakan untuk memberikan gambaran singkat mengenai kondisi Alert yang tercatat pada Gamon.

Informasi yang ditampilkan:

* Total Alert
* Ongoing Alert
* Resolved Alert
* Critical Alert

Tujuan:

* Membantu pengguna mengetahui jumlah alert yang sedang maupun pernah terjadi.
* Memberikan informasi singkat mengenai kondisi gangguan yang sedang berlangsung.

---

### 2. Filter & Search

Section ini digunakan untuk mempermudah pengguna dalam melakukan pencarian maupun filtering Alert.

Fitur yang disediakan:

* Search Alert
* Filter berdasarkan Alert Status:

  * All
  * Ongoing
  * Resolved
* Filter berdasarkan Severity Level:

  * Low
  * Medium
  * High
  * Critical
* Filter berdasarkan Device Type:

  * Server
  * Router
  * Switch
  * Access Point
  * Website
  * Printer
  * Other

Tujuan:

* Mempermudah pengguna menemukan informasi gangguan yang sedang dicari.
* Membantu proses identifikasi gangguan menjadi lebih cepat.

---

### 3. Alert List

Section ini merupakan bagian utama dari halaman Alert Center yang menampilkan seluruh Alert yang pernah tercatat oleh Gamon.

Informasi yang ditampilkan meliputi:

* Alert Title
* Device Name
* Alert Status
* Severity Level
* Started Time

Section ini berfungsi untuk memberikan informasi singkat mengenai gangguan yang sedang maupun pernah terjadi.

---

### 4. Alert Detail

Section ini akan ditampilkan ketika pengguna memilih salah satu Alert pada Alert List.

Informasi yang ditampilkan meliputi:

* Alert Title
* Device Name
* Device Type
* Alert Status
* Severity Level
* Started Time
* Resolved Time (jika tersedia)
* Monitoring Method
* Alert Description

Tujuan:

* Memberikan informasi gangguan secara lebih lengkap.
* Membantu pengguna mengetahui detail dari Alert yang dipilih.

## User Flow

1. Pengguna membuka halaman Alert Center.
2. Pengguna melihat Alert Summary untuk mengetahui jumlah Alert yang tercatat pada Gamon.
3. Pengguna dapat melakukan pencarian ataupun filtering Alert yang ingin dilihat.
4. Pengguna memilih salah satu Alert pada Alert List.
5. Pengguna melihat informasi Alert secara lebih lengkap melalui Alert Detail.
6. Apabila diperlukan, pengguna dapat menuju halaman Monitoring untuk melihat kondisi perangkat secara realtime.

## Scope Limitation (MVP)

Halaman Alert Center pada prototype Gamon hanya berfokus pada manajemen Alert dan riwayat gangguan yang dihasilkan selama proses monitoring berlangsung.

Fitur yang tidak termasuk dalam MVP:

* Alert Escalation System
* Multiple Notification Channel Management
* Incident Assignment
* Ticketing System Integration
* SLA Monitoring
* Alert Analytics
* Predictive Alert System
* AI-based Incident Analysis
* Advanced Reporting System

Prototype ini hanya berfokus pada pembuktian konsep (Proof of Concept) bahwa Gamon mampu menghasilkan Alert secara otomatis ketika terjadi gangguan pada perangkat yang sedang dimonitor serta membantu pengguna mengetahui informasi gangguan yang sedang maupun pernah terjadi secara terpusat.
