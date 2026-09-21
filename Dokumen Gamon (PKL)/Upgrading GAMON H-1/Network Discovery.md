````text
# PRD — Network Discovery (Nmap / Netdiscover)

## 1. Overview

### Nama Fitur
Network Discovery

### Komponen Utama
- Nmap
- Netdiscover

### Status
Planned Upgrade — GAMON

### Tujuan
Menambahkan kemampuan Network Discovery pada GAMON untuk membantu administrator menemukan perangkat yang aktif atau terdeteksi pada jaringan tanpa harus memasukkan IP perangkat satu per satu secara manual.

Fitur ini berfungsi sebagai tahap discovery sebelum perangkat dimasukkan ke dalam sistem monitoring GAMON.

Konsep utamanya:

Network Discovery → Discover Device → Review Result → Select Device → Add to GAMON → Monitoring Engine


## 2. Latar Belakang / Problem

Pada versi GAMON saat ini, administrator harus mengetahui informasi perangkat terlebih dahulu sebelum membuat device monitoring.

Contohnya:

1. Administrator mengetahui IP server `192.168.1.10`
2. Administrator membuka Device Management
3. Administrator memasukkan nama device
4. Administrator memasukkan IP
5. Device disimpan
6. Monitoring Engine mulai melakukan monitoring

Pendekatan tersebut bekerja, tetapi kurang praktis apabila administrator ingin mengetahui perangkat apa saja yang sedang berada atau aktif pada suatu jaringan.

Selain itu, proses pendataan perangkat jaringan secara manual dapat menyebabkan:

- IP perangkat terlewat.
- Perangkat baru tidak langsung diketahui.
- Administrator harus menggunakan tool eksternal terlebih dahulu.
- Proses discovery dan proses monitoring masih terpisah.

Network Discovery dibuat untuk memberikan tahap awal sebelum monitoring.


## 3. Tujuan Fitur

Fitur Network Discovery memiliki beberapa tujuan:

1. Menemukan perangkat yang dapat terdeteksi pada jaringan.
2. Mendapatkan informasi dasar mengenai perangkat hasil discovery.
3. Menampilkan hasil discovery dalam bentuk daftar.
4. Memberikan administrator kemampuan memilih perangkat tertentu.
5. Mempermudah penambahan perangkat hasil discovery ke Device Management.
6. Memisahkan proses discovery dengan proses monitoring.

Network Discovery bukan pengganti Monitoring Engine.

Discovery menjawab:

> "Perangkat apa saja yang dapat ditemukan?"

Monitoring menjawab:

> "Apakah perangkat tersebut masih tersedia dan bagaimana statusnya?"


## 4. Konsep Arsitektur

Network Discovery menjadi modul terpisah dari Monitoring Engine.

```text
                    GAMON
                      │
             ┌────────┴────────┐
             │                 │
     Network Discovery   Monitoring Engine
             │                 │
       ┌─────┴─────┐      ┌────┴────┐
       │           │      │         │
      Nmap     Netdiscover ICMP     TCP
       │           │      │         │
       └─────┬─────┘      └────┬────┘
             │                 │
       Discovery Result    Monitoring Result
             │                 │
             ▼                 ▼
       Discovery UI       Status/History
             │
             ▼
       Select Device
             │
             ▼
      Device Management
             │
             ▼
      Monitoring Engine
````

Dengan struktur tersebut, Nmap dan Netdiscover tidak menjadi bagian dari monitoring cycle.

Contohnya:

```text
Discovery
    ↓
Nmap / Netdiscover
    ↓
192.168.1.1
192.168.1.10
192.168.1.20
192.168.1.25
    ↓
Administrator memilih
192.168.1.10
    ↓
Add to Device Management
    ↓
Monitoring Engine
    ↓
ICMP / TCP Monitoring
```

## 5. Network Discovery vs Monitoring

Kedua fitur harus memiliki tanggung jawab yang berbeda.

|Network Discovery|Monitoring|
|---|---|
|Menemukan perangkat|Memantau perangkat|
|Bersifat on-demand|Berjalan periodik|
|Mencari target|Memeriksa target|
|Menghasilkan discovery result|Menghasilkan monitoring result|
|Tidak menentukan status monitoring|Menentukan status monitoring|
|Tidak membuat alert monitoring|Dapat membuat alert|
|Dapat menggunakan Nmap/Netdiscover|Menggunakan ICMP/TCP/HTTP/SNMP|

Contohnya:

Nmap menemukan:

```text
192.168.1.10
```

Hal tersebut belum berarti perangkat tersebut sudah menjadi device GAMON.

Administrator masih dapat memilih:

```text
[ Add to Monitoring ]
```

Baru setelah itu perangkat masuk ke Device Management.

## 6. Discovery Engine

Network Discovery membutuhkan komponen yang bertugas mengatur proses discovery.

```text
Network Discovery
       │
       ▼
Discovery Engine
       │
       ▼
Tool Selector
       │
 ┌─────┴─────┐
 ▼           ▼
Nmap    Netdiscover
 │           │
 └─────┬─────┘
       ▼
Raw Discovery Result
       │
       ▼
Result Parser / Normalizer
       │
       ▼
Standard Discovery Result
       │
       ▼
Discovery UI
```

Discovery Engine bertanggung jawab terhadap orchestration proses discovery.

Tool seperti Nmap dan Netdiscover bertanggung jawab terhadap proses scanning/discovery yang sebenarnya.

## 7. Nmap

Nmap digunakan sebagai salah satu discovery engine untuk melakukan network scanning.

Dalam konteks GAMON, Nmap dapat digunakan untuk mendapatkan informasi seperti:

- IP address
    
- Host availability
    
- MAC address apabila tersedia
    
- Vendor perangkat apabila tersedia
    
- Port yang terdeteksi
    
- Informasi tambahan hasil scanning yang didukung konfigurasi discovery
    

Nmap lebih cocok ketika GAMON membutuhkan discovery yang lebih luas dan informasi tambahan mengenai host.

## 8. Netdiscover

Netdiscover digunakan sebagai alternatif discovery tool yang berfokus pada discovery perangkat melalui jaringan lokal.

Informasi yang dapat diperoleh dapat berupa:

- IP address
    
- MAC address
    
- Vendor MAC address
    
- Informasi host yang terdeteksi
    

Netdiscover dapat digunakan ketika tujuan utama discovery adalah mengetahui perangkat yang terlihat pada jaringan lokal.

## 9. Tool Selection

Administrator dapat menentukan tool yang digunakan untuk discovery.

Contoh:

```text
Discovery Method

(•) Nmap
( ) Netdiscover
```

Pada tahap awal, pemilihan tool sebaiknya bersifat global untuk satu proses discovery.

Tidak diperlukan konfigurasi tool per-device.

## 10. Discovery Target

Network Discovery harus menentukan target jaringan yang ingin diperiksa.

Contoh:

```text
Target Network:
192.168.1.0/24

[ Start Discovery ]
```

Target dapat berupa network/subnet yang memang menjadi ruang lingkup discovery.

Contoh:

```text
192.168.1.0/24
10.10.10.0/24
172.16.1.0/24
```

Discovery tidak seharusnya berjalan tanpa target yang jelas.

## 11. Discovery Result

Setelah proses selesai, GAMON menampilkan hasil discovery dalam bentuk tabel.

Contoh:

|IP Address|MAC Address|Vendor|Status|Source|
|---|---|---|---|---|
|192.168.1.1|XX:XX:XX:XX|MikroTik|Detected|Nmap|
|192.168.1.10|XX:XX:XX:XX|Dell|Detected|Nmap|
|192.168.1.20|XX:XX:XX:XX|TP-Link|Detected|Netdiscover|
|192.168.1.25|XX:XX:XX:XX|Unknown|Detected|Nmap|

Informasi yang tidak tersedia tidak boleh dibuat-buat.

Contohnya:

```text
Vendor: Unknown
```

lebih baik daripada mengasumsikan jenis perangkat.

## 12. Discovery Result Normalization

Nmap dan Netdiscover dapat menghasilkan format output yang berbeda.

Karena itu hasil keduanya harus dinormalisasi menjadi struktur internal GAMON.

Konsep:

```text
Nmap Output
     │
     ▼
Nmap Parser
     │
     ├──────────────┐
     │              │
     ▼              ▼
IP Address       MAC Address
     │              │
     └──────┬───────┘
            ▼
    Standard Discovery Result
```

Begitu pula:

```text
Netdiscover Output
        │
        ▼
Netdiscover Parser
        │
        ▼
Standard Discovery Result
```

Dengan pendekatan tersebut, UI GAMON tidak perlu mengetahui apakah data berasal dari Nmap atau Netdiscover.

## 13. Standard Discovery Result

Secara konseptual, setiap hasil discovery memiliki informasi seperti:

```text
Discovery Result
├── IP Address
├── MAC Address
├── Vendor
├── Host Status
├── Source Tool
└── Discovery Time
```

Informasi tersebut merupakan hasil discovery dan belum otomatis menjadi data permanen Device Management.

## 14. Add Device to GAMON

Salah satu fungsi utama Network Discovery adalah mempermudah administrator memasukkan perangkat hasil discovery ke GAMON.

Contoh:

```text
Discovery Result

☑ 192.168.1.10   Dell
☐ 192.168.1.20   TP-Link
☑ 192.168.1.25   Unknown

[ Add Selected Devices ]
```

Ketika administrator memilih device:

```text
Discovery Result
       ↓
Select Device
       ↓
Add to Device Management
       ↓
Device Configuration
       ↓
Monitoring Engine
```

Administrator tetap dapat mengubah informasi seperti:

- Device Name
    
- Device Type
    
- Monitoring Method
    
- Monitoring Interval
    
- Monitoring Configuration
    

## 15. Duplicate Detection

GAMON harus dapat membedakan perangkat yang sudah terdaftar dengan perangkat baru.

Contoh:

```text
Discovery Result

192.168.1.10
Status: Already Registered
```

Perangkat tersebut tidak boleh dibuat sebagai device baru secara otomatis.

Tujuannya mencegah:

```text
Discovery
    ↓
Add Device
    ↓
Duplicate Device
    ↓
Duplicate Monitoring
    ↓
Duplicate Alert
```

Discovery seharusnya membantu Device Management, bukan membuat data duplikat.

## 16. Discovery Status

Status discovery dan status monitoring harus dibedakan.

Contoh:

```text
Discovery Status:
Detected

Monitoring Status:
Not Monitored
```

Setelah perangkat ditambahkan:

```text
Discovery Status:
Detected

Monitoring Status:
Online
```

Hal ini penting karena perangkat yang ditemukan belum tentu sudah dimonitor oleh GAMON.

## 17. Discovery Process

Alur utama:

```text
Administrator
      │
      ▼
Network Discovery
      │
      ▼
Select Discovery Tool
      │
      ▼
Input Target Network
      │
      ▼
Start Discovery
      │
      ▼
Discovery Engine
      │
      ▼
Nmap / Netdiscover
      │
      ▼
Raw Result
      │
      ▼
Parser / Normalizer
      │
      ▼
Discovery Result
      │
      ▼
Administrator Review
      │
      ├───────────────┐
      │               │
      ▼               ▼
 Ignore          Add Device
                      │
                      ▼
              Device Management
                      │
                      ▼
              Monitoring Engine
```

## 18. Discovery Lifecycle

Network Discovery tidak berjalan terus-menerus seperti Monitoring Engine.

Siklusnya:

```text
IDLE
  │
  ▼
START DISCOVERY
  │
  ▼
SCANNING
  │
  ▼
PROCESSING RESULT
  │
  ▼
RESULT READY
  │
  ▼
REVIEW
  │
  ├── Add Device
  │
  └── Ignore
  │
  ▼
IDLE
```

Discovery dapat dijalankan kembali kapan pun administrator membutuhkan informasi terbaru.

## 19. Error Handling

Network Discovery harus membedakan error tool dengan hasil discovery.

Contoh:

```text
Nmap tidak tersedia
```

bukan:

```text
Semua device offline
```

Begitu juga:

```text
Netdiscover tidak tersedia
```

harus menghasilkan error konfigurasi/environment, bukan discovery result kosong yang dianggap sebagai tidak ada perangkat.

Jenis error yang perlu diperhatikan:

- Tool tidak terinstall.
    
- Tool tidak dapat dieksekusi.
    
- Target network tidak valid.
    
- Permission tidak mencukupi.
    
- Proses discovery gagal.
    
- Parsing output gagal.
    
- Discovery timeout.
    
- Tidak ditemukan perangkat.
    

## 20. Security & Permission

Network Discovery merupakan fitur administratif dan harus digunakan hanya pada jaringan yang memang memiliki izin untuk dikelola.

Akses fitur dapat dibatasi kepada administrator.

GAMON juga harus memperlakukan command execution sebagai komponen yang terkontrol.

Input target network tidak boleh langsung diperlakukan sebagai command mentah tanpa validasi.

## 21. Integrasi dengan Device Management

Network Discovery tidak menggantikan Device Management.

Hubungannya:

```text
Network Discovery
       │
       │ discovered device
       ▼
Device Management
       │
       │ configured device
       ▼
Monitoring Engine
```

Device Management tetap menjadi sumber konfigurasi device yang digunakan Monitoring Engine.

## 22. Integrasi dengan Monitoring Engine

Setelah device berhasil ditambahkan:

```text
Discovery
   ↓
Device Management
   ↓
Monitoring Configuration
   ↓
Monitoring Engine
   ↓
ICMP Checker
   ↓
TCP Checker
   ↓
Result Processing
   ↓
History / Alert
```

Network Discovery tidak menentukan apakah device harus menggunakan ICMP, TCP, HTTP, atau metode monitoring lainnya.

Pemilihan monitoring method tetap menjadi tanggung jawab konfigurasi monitoring.

## 23. Integrasi dengan Alert System

Network Discovery tidak menghasilkan alert monitoring.

Contoh:

```text
Nmap menemukan 192.168.1.10
```

tidak otomatis berarti:

```text
Alert: Device Online
```

Alert baru menjadi tanggung jawab Monitoring Engine setelah device masuk ke monitoring.

Dengan demikian:

```text
Discovery Event ≠ Monitoring Event
```

Hal ini menjaga separation of concern antara discovery dan monitoring.

## 24. UI Concept

Halaman Network Discovery dapat memiliki struktur:

```text
┌─────────────────────────────────────────────┐
│ Network Discovery                           │
├─────────────────────────────────────────────┤
│                                             │
│ Discovery Method: [ Nmap ▼ ]                │
│ Target Network:   [ 192.168.1.0/24 ]       │
│                                             │
│              [ Start Discovery ]            │
│                                             │
├─────────────────────────────────────────────┤
│ Discovery Result                            │
├───────┬──────────────┬─────────┬────────────┤
│ Select│ IP Address   │ Vendor  │ Status     │
├───────┼──────────────┼─────────┼────────────┤
│  ☑    │192.168.1.1   │MikroTik │Detected    │
│  ☑    │192.168.1.10  │Dell     │Registered  │
│  ☐    │192.168.1.20  │TP-Link  │Detected    │
├───────┴──────────────┴─────────┴────────────┤
│                                             │
│        [ Add Selected Devices ]             │
└─────────────────────────────────────────────┘
```

## 25. Hubungan dengan Fitur GAMON Lain

Struktur fitur GAMON menjadi:

```text
GAMON
│
├── Dashboard
│
├── Device Management
│
├── Network Discovery
│   ├── Nmap
│   └── Netdiscover
│
├── Monitoring Engine
│   ├── ICMP Checker
│   │   ├── Ping Executor
│   │   └── Fping Executor
│   │
│   ├── TCP Checker
│   │
│   └── Future
│       ├── HTTP/HTTPS
│       └── SNMP
│
├── Alert Engine
│
├── Alert Center
│
├── Notification
│   ├── Telegram
│   └── Sound
│
└── History
```

Dengan struktur ini, Network Discovery menjadi modul yang berdiri sendiri dan tidak mengacaukan arsitektur Monitoring Engine.

## 26. Scope MVP

Untuk tahap awal, Network Discovery cukup mencakup:

1. Memilih Nmap atau Netdiscover.
    
2. Memasukkan target network.
    
3. Menjalankan discovery.
    
4. Mengambil hasil discovery.
    
5. Normalisasi hasil.
    
6. Menampilkan IP/MAC/Vendor jika tersedia.
    
7. Menandai device yang sudah terdaftar.
    
8. Memilih device hasil discovery.
    
9. Menambahkan device terpilih ke Device Management.
    

Tidak perlu langsung memasukkan seluruh kemampuan Nmap ke GAMON.

## 27. Out of Scope

Untuk versi awal, fitur berikut tidak menjadi bagian Network Discovery:

- Vulnerability scanning.
    
- Exploitation.
    
- Password attack.
    
- Credential discovery.
    
- Automatic penetration testing.
    
- Automatic OS exploitation.
    
- Continuous network scanning.
    
- Automatic deletion of devices.
    
- Automatic addition seluruh hasil discovery tanpa review administrator.
    
- Integrasi SNMP discovery kompleks.
    
- Network topology mapping.
    

Network Discovery harus tetap fokus pada asset discovery.

## 28. Future Development

Setelah MVP stabil, Network Discovery dapat dikembangkan menjadi:

```text
Network Discovery
│
├── Nmap
├── Netdiscover
├── ARP Discovery
├── Host Discovery
├── Port Discovery
├── OS Detection
├── Service Detection
├── Vendor Detection
├── Network Range Detection
└── Network Topology
```

Namun fitur tersebut sebaiknya ditambahkan secara bertahap.

## 29. Prinsip Arsitektur

Prinsip utama desain:

```text
Discovery ≠ Monitoring
Tool ≠ Discovery Engine
Discovery Result ≠ Device
Device ≠ Monitoring Status
```

Nmap dan Netdiscover hanya merupakan tool/executor yang digunakan untuk memperoleh informasi.

Discovery Engine mengatur proses.

Result Normalizer menyamakan output.

Device Management menyimpan device yang dipilih.

Monitoring Engine melakukan monitoring terhadap device tersebut.

Alert Engine menangani perubahan status monitoring.

## 30. Expected Result

Setelah fitur diterapkan, alur penggunaan GAMON menjadi lebih lengkap:

Sebelum Network Discovery:

```text
Administrator tahu IP
        ↓
Input Device Manual
        ↓
Monitoring
```

Setelah Network Discovery:

```text
Administrator menentukan Network
        ↓
Network Discovery
        ↓
Nmap / Netdiscover
        ↓
Daftar perangkat
        ↓
Administrator memilih perangkat
        ↓
Device Management
        ↓
Monitoring Engine
        ↓
ICMP / TCP / HTTP / SNMP
        ↓
History
        ↓
Alert
        ↓
Telegram / Sound / Dashboard
```

Dengan demikian Network Discovery berfungsi sebagai pintu masuk dari kondisi "belum mengetahui asset jaringan" menuju kondisi "asset sudah terdaftar dan dapat dimonitor oleh GAMON".

## 31. Batasan Implementasi

Fitur discovery tidak boleh dianggap sebagai pengganti sistem asset management secara penuh.

Hasil discovery hanya menunjukkan perangkat yang berhasil terdeteksi pada saat proses dijalankan. Perangkat yang tidak terlihat pada saat scanning tidak otomatis berarti perangkat tersebut tidak ada secara permanen.

Selain itu, informasi seperti vendor, MAC address, OS, atau service bergantung pada metode discovery, jaringan, permission, konfigurasi target, dan kemampuan tool yang digunakan.

Karena itu hasil discovery harus dianggap sebagai informasi hasil scanning yang perlu diverifikasi administrator sebelum digunakan sebagai konfigurasi monitoring.