## 1. Overview

### Nama Fitur

Network Discovery — OS Detection & Target Information

### Modul

Network Discovery

### Komponen

- Host Discovery
    
- Target Information
    
- OS Detection
    
- Vendor Detection
    
- Port Discovery
    
- Service Detection
    
- Discovery Result
    

### Status

Planned Upgrade — GAMON

### Tujuan

Menambahkan kemampuan untuk memperoleh informasi teknis mengenai host yang ditemukan pada jaringan, khususnya sistem operasi dan informasi pendukung lainnya.

Fitur ini berfokus pada **target yang ditemukan**, bukan sistem operasi atau hardware server yang menjalankan aplikasi GAMON.

Konsep utamanya:

```text
Target Network
      ↓
Network Discovery
      ↓
Host Discovery
      ↓
Target Information
      ├── IP Address
      ├── MAC Address
      ├── Vendor
      ├── OS
      ├── Open Ports
      └── Services
      ↓
Administrator Review
      ↓
Add to Device Management
```

---

## 2. Latar Belakang / Problem

Pada proses monitoring jaringan, administrator membutuhkan informasi dasar mengenai perangkat yang akan dimonitor.

Dalam kondisi manual, administrator mungkin hanya mengetahui:

```text
192.168.1.10
```

Namun informasi tersebut belum menjelaskan perangkat apa yang menggunakan IP tersebut.

Administrator mungkin perlu mengetahui:

- Perangkat tersebut menggunakan OS apa.
    
- Vendor perangkat.
    
- MAC address.
    
- Port yang terbuka.
    
- Service yang terdeteksi.
    
- Informasi teknis lain yang tersedia dari proses discovery.
    

Tanpa discovery information, administrator harus melakukan pemeriksaan menggunakan tool eksternal secara terpisah.

Network Discovery dikembangkan untuk menggabungkan proses tersebut ke dalam workflow GAMON.

---

## 3. Tujuan Fitur

Fitur ini bertujuan untuk:

1. Menemukan host pada target network.
    
2. Mengumpulkan informasi teknis mengenai host yang ditemukan.
    
3. Mendeteksi sistem operasi target apabila informasi yang tersedia memungkinkan.
    
4. Mengidentifikasi vendor perangkat apabila tersedia.
    
5. Menampilkan port dan service yang terdeteksi.
    
6. Menyediakan informasi sebelum administrator memasukkan perangkat ke Device Management.
    
7. Membantu administrator memahami karakteristik target sebelum menentukan konfigurasi monitoring.
    

Fitur ini bukan bertujuan untuk melakukan eksploitasi atau penetration testing terhadap target.

---

## 4. Scope

Fokus fitur dibagi menjadi dua bagian:

### Host Discovery

Menentukan host yang dapat ditemukan pada target network.

Informasi dasar:

```text
IP Address
MAC Address
Host Status
```

### Target Information

Mengumpulkan informasi tambahan mengenai host yang telah ditemukan.

```text
Vendor
OS
Open Ports
Services
```

Struktur:

```text
Network Discovery
│
├── Host Discovery
│   ├── IP Address
│   ├── MAC Address
│   └── Host Status
│
└── Target Information
    ├── Vendor Detection
    ├── OS Detection
    ├── Port Discovery
    └── Service Detection
```

---

## 5. Perbedaan dengan Runtime Environment

Fitur ini harus secara eksplisit membedakan **Target Information** dan **Runtime Environment**.

### Target Information

Informasi mengenai perangkat yang sedang ditemukan:

```text
Target
├── IP
├── MAC
├── Vendor
├── OS
├── Ports
└── Services
```

### Runtime Environment

Informasi mengenai mesin yang menjalankan GAMON:

```text
GAMON Server
├── Operating System
├── CPU
├── RAM
├── Architecture
└── Installed Tools
```

Runtime Environment bukan bagian dari OS Detection target.

Contohnya, apabila GAMON dijalankan pada Ubuntu Server dan melakukan discovery terhadap sebuah Windows PC:

```text
GAMON Server:
OS = Ubuntu

Target:
OS = Windows
```

Kedua informasi tersebut memiliki objek yang berbeda dan tidak boleh tercampur.

---

## 6. Konsep Arsitektur

```text
                         GAMON
                           │
                           ▼
                  Network Discovery
                           │
                           ▼
                    Discovery Engine
                           │
                           ▼
                     Target Network
                           │
                           ▼
                    Host Discovery
                           │
                           ▼
                     Target Hosts
                           │
             ┌─────────────┴─────────────┐
             │                           │
             ▼                           ▼
      Basic Information            Target Information
             │                           │
             │                  ┌────────┼─────────┐
             │                  │        │         │
             ▼                  ▼        ▼         ▼
             IP/MAC          Vendor     OS      Ports/Services
             │                  │        │         │
             └──────────────────┴────────┴─────────┘
                                │
                                ▼
                         Discovery Result
                                │
                                ▼
                       Administrator Review
                                │
                                ▼
                       Device Management
```

---

## 7. Discovery Tool

Network Discovery dapat menggunakan tool seperti:

- Nmap
    
- Netdiscover
    

Keduanya dapat digunakan sebagai sumber discovery information, tetapi kemampuan informasi yang tersedia dapat berbeda.

Konsep:

```text
Discovery Engine
       │
       ├── Nmap
       │
       └── Netdiscover
              │
              ▼
       Discovery Result
```

Nmap dapat digunakan ketika diperlukan discovery dengan informasi target yang lebih luas, termasuk informasi mengenai port, service, dan kemungkinan OS.

Netdiscover lebih berfokus pada penemuan host pada jaringan lokal serta informasi seperti IP, MAC, dan vendor apabila tersedia.

---

## 8. OS Detection

OS Detection bertujuan untuk memperoleh **perkiraan sistem operasi yang digunakan oleh target host**.

Contoh hasil:

```text
Target:
192.168.1.10

OS:
Linux
```

atau:

```text
Target:
192.168.1.20

OS:
Windows
```

Apabila informasi tidak cukup:

```text
OS:
Unknown
```

GAMON tidak boleh menganggap hasil deteksi sebagai fakta absolut apabila tool hanya memberikan perkiraan.

---

## 9. OS Detection Result

Hasil OS Detection secara konseptual dapat memiliki:

```text
OS Detection
├── Operating System
├── Version (jika tersedia)
├── Detection Source
└── Confidence / Accuracy Information (jika tersedia)
```

Contoh:

```text
OS:
Linux

Version:
Unknown

Source:
Nmap

Confidence:
Available from discovery result
```

Jika versi OS tidak dapat diketahui, sistem hanya menampilkan informasi yang benar-benar tersedia.

Contoh:

```text
OS: Linux
Version: Unknown
```

bukan memaksakan:

```text
OS: Ubuntu 24.04
```

jika data discovery tidak mendukung informasi tersebut.

---

## 10. Vendor Detection

Vendor Detection memberikan informasi mengenai vendor perangkat apabila dapat diperoleh dari hasil discovery.

Contoh:

```text
MAC Address:
AA:BB:CC:DD:EE:FF

Vendor:
MikroTik
```

atau:

```text
Vendor:
Unknown
```

Vendor merupakan informasi tambahan dan bukan penentu pasti jenis perangkat.

Contohnya:

```text
Vendor: Dell
```

tidak secara otomatis berarti target merupakan server.

Informasi tersebut tetap perlu dikombinasikan dengan informasi lain.

---

## 11. Port Discovery

Port Discovery digunakan untuk mengetahui port yang terdeteksi terbuka pada target.

Contoh:

```text
Target:
192.168.1.10

Detected Ports:
22
80
443
```

Informasi port dapat membantu administrator memahami service yang tersedia pada target.

Namun port yang terdeteksi bukan berarti service tersebut pasti sedang berfungsi secara normal. Status service tetap perlu diverifikasi melalui mekanisme monitoring yang sesuai.

---

## 12. Service Detection

Apabila informasi service tersedia, hasil discovery dapat memberikan informasi seperti:

```text
Port    Service
22      SSH
80      HTTP
443     HTTPS
```

Informasi tersebut menjadi bagian dari Target Information.

Hubungannya:

```text
Port Discovery
      ↓
Service Detection
      ↓
Target Information
```

---

## 13. Standard Discovery Result

Karena Nmap dan Netdiscover dapat menghasilkan format output yang berbeda, hasil discovery perlu dinormalisasi.

Konsep:

```text
Nmap Output
     │
     ▼
Nmap Parser
     │
     └──────────┐
                │
                ▼
       Standard Discovery Result
                ▲
                │
     ┌──────────┘
     │
Netdiscover Parser
     │
     ▼
Netdiscover Output
```

Standard Discovery Result secara konseptual:

```text
Discovery Result
├── IP Address
├── MAC Address
├── Vendor
├── Host Status
├── OS
├── OS Version
├── Ports
├── Services
├── Source Tool
└── Discovery Time
```

Field yang tidak tersedia tidak boleh diisi menggunakan asumsi.

---

## 14. Discovery Result

Contoh hasil yang ditampilkan kepada administrator:

```text
┌──────────────────────────────────────────────┐
│ Target Information                           │
├──────────────────────────────────────────────┤
│ IP Address : 192.168.1.10                    │
│ MAC        : AA:BB:CC:DD:EE:FF              │
│ Vendor     : Dell                            │
│ OS         : Linux                           │
│ Version    : Unknown                         │
│ Status     : Detected                        │
│ Source     : Nmap                            │
├──────────────────────────────────────────────┤
│ Ports                                      │
│ 22         SSH                               │
│ 80         HTTP                              │
│ 443        HTTPS                             │
└──────────────────────────────────────────────┘
```

Informasi tersebut digunakan untuk membantu administrator melakukan review terhadap target.

---

## 15. Discovery Information sebagai Pre-Monitoring Information

Target Information tidak langsung menjadi konfigurasi monitoring.

Alurnya:

```text
Discovery
   ↓
Target Information
   ↓
Administrator Review
   ↓
Identify Target
   ↓
Add to Device Management
   ↓
Configure Monitoring
   ↓
Monitoring Engine
```

Contoh:

```text
Discovery menemukan:

192.168.1.10
OS: Linux
Port: 22, 80, 443
```

Administrator kemudian dapat menentukan:

```text
Device Name:
Web Server

Monitoring:
ICMP
TCP 22
TCP 80
TCP 443
```

Dengan demikian discovery membantu administrator menentukan konfigurasi monitoring berdasarkan informasi target.

---

## 16. Integrasi dengan Device Management

Setelah administrator memilih target:

```text
Discovery Result
       ↓
[ Add to Device ]
       ↓
Device Management
```

Informasi yang ditemukan dapat digunakan sebagai data awal ketika membuat device.

Contohnya:

```text
IP Address
MAC Address
Vendor
OS
```

Namun Device Management tetap menjadi tempat konfigurasi device yang digunakan oleh Monitoring Engine.

Discovery tidak mengambil alih fungsi Device Management.

---

## 17. Duplicate Detection

Jika target sudah terdaftar di GAMON:

```text
192.168.1.10
Status: Already Registered
```

Administrator tidak perlu membuat device baru.

Contoh:

```text
Discovery Result
       │
       ▼
Check Device Management
       │
       ├── Registered
       │      ↓
       │   Already Registered
       │
       └── Not Registered
              ↓
          Add Device
```

Hal ini mencegah satu target dibuat menjadi beberapa device monitoring.

---

## 18. Discovery Status vs Monitoring Status

Status discovery dan status monitoring harus dipisahkan.

Contoh:

```text
Discovery Status:
Detected

Monitoring Status:
Not Monitored
```

Setelah ditambahkan:

```text
Discovery Status:
Detected

Monitoring Status:
Online
```

Jika monitoring kemudian gagal:

```text
Discovery Status:
Previously Detected

Monitoring Status:
Offline
```

Perubahan monitoring status tidak berarti hasil discovery sebelumnya harus dihapus.

---

## 19. Discovery Lifecycle

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
TARGET DETECTION
 │
 ▼
TARGET INFORMATION
 │
 ├── OS Detection
 ├── Vendor Detection
 ├── Port Discovery
 └── Service Detection
 │
 ▼
RESULT NORMALIZATION
 │
 ▼
RESULT READY
 │
 ▼
ADMINISTRATOR REVIEW
 │
 ├── Ignore
 │
 └── Add to Device
          │
          ▼
    Device Management
```

Setelah hasil ditampilkan, proses discovery selesai.

Network Discovery tidak harus terus berjalan seperti Monitoring Engine.

---

## 20. Error Handling

Sistem harus membedakan antara discovery failure dan informasi yang memang tidak tersedia.

Contoh:

```text
OS:
Unknown
```

berarti informasi OS tidak berhasil diperoleh atau tidak tersedia.

Sedangkan:

```text
Discovery Failed
```

berarti proses discovery itu sendiri mengalami masalah.

Beberapa kemungkinan error:

- Discovery tool tidak tersedia.
    
- Target network tidak valid.
    
- Permission tidak mencukupi.
    
- Discovery process gagal.
    
- Parsing output gagal.
    
- Target tidak memberikan informasi yang diperlukan.
    
- Tidak ada host yang berhasil ditemukan.
    

Sistem tidak boleh mengubah error tool menjadi:

```text
No Device Found
```

apabila sebenarnya proses discovery tidak berhasil dijalankan.

---

## 21. Security & Permission

Network Discovery merupakan fitur administratif.

Fitur ini hanya digunakan pada network yang memang berada dalam kewenangan administrator.

Target network harus ditentukan secara eksplisit.

Input target juga harus divalidasi sebelum digunakan oleh discovery engine.

Network Discovery tidak memiliki tujuan:

- Exploitation
    
- Credential attack
    
- Vulnerability exploitation
    
- Password cracking
    
- Unauthorized access
    

Fokusnya adalah asset discovery dan target information.

---

## 22. UI Concept

Halaman dapat memiliki dua bagian utama.

```text
┌──────────────────────────────────────────────┐
│ Network Discovery                            │
├──────────────────────────────────────────────┤
│ Discovery Tool: [ Nmap ▼ ]                   │
│ Target Network: [ 192.168.1.0/24 ]          │
│                                              │
│             [ Start Discovery ]              │
├──────────────────────────────────────────────┤
│ Discovered Hosts                             │
├──────┬──────────────┬──────────┬─────────────┤
│ Select│ IP Address  │ Vendor   │ OS          │
├──────┼──────────────┼──────────┼─────────────┤
│ ☑    │192.168.1.1   │MikroTik  │RouterOS     │
│ ☐    │192.168.1.10  │Dell      │Linux        │
│ ☑    │192.168.1.20  │HP        │Windows      │
└──────┴──────────────┴──────────┴─────────────┘

[ View Information ]   [ Add Selected Devices ]
```

Ketika administrator memilih salah satu host:

```text
┌───────────────────────────────────────┐
│ Target Information                    │
├───────────────────────────────────────┤
│ IP       : 192.168.1.10               │
│ MAC      : AA:BB:CC:DD:EE:FF          │
│ Vendor   : Dell                       │
│ OS       : Linux                      │
│ Version  : Unknown                    │
│ Source   : Nmap                       │
│                                       │
│ Ports                                  │
│ 22       SSH                          │
│ 80       HTTP                         │
│ 443      HTTPS                        │
└───────────────────────────────────────┘
```

---

## 23. Hubungan dengan Modul GAMON

Struktur GAMON setelah penambahan fitur:

```text
GAMON
│
├── Dashboard
│
├── Device Management
│
├── Network Discovery
│   │
│   ├── Host Discovery
│   │
│   └── Target Information
│       ├── Vendor Detection
│       ├── OS Detection
│       ├── Port Discovery
│       └── Service Detection
│
├── Monitoring Engine
│   ├── ICMP Checker
│   │   ├── Ping Executor
│   │   └── Fping Executor
│   │
│   ├── TCP Checker
│   │
│   ├── HTTP/HTTPS Checker
│   └── Future: SNMP
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

Network Discovery menjadi modul yang bertugas mengenali target sebelum target masuk ke proses monitoring.

---

## 24. Scope MVP

Versi awal fitur cukup mencakup:

1. Menentukan target network.
    
2. Menjalankan discovery menggunakan Nmap atau Netdiscover.
    
3. Menemukan host.
    
4. Mengambil IP address.
    
5. Mengambil MAC address jika tersedia.
    
6. Mengambil vendor jika tersedia.
    
7. Melakukan OS Detection jika didukung hasil discovery.
    
8. Menampilkan OS sebagai informasi target.
    
9. Menampilkan port/service yang terdeteksi apabila tersedia.
    
10. Menampilkan hasil dalam Discovery Result.
    
11. Mendeteksi apakah target sudah terdaftar.
    
12. Memberikan opsi Add to Device Management.
    

---

## 25. Out of Scope

Untuk versi awal, fitur berikut tidak termasuk:

- Vulnerability scanning.
    
- Exploitation.
    
- Credential discovery.
    
- Password attack.
    
- Automatic penetration testing.
    
- Automatic vulnerability remediation.
    
- Continuous network scanning.
    
- Automatic addition seluruh hasil discovery.
    
- Network topology mapping.
    
- Full asset management.
    
- Runtime monitoring terhadap server GAMON.
    

Runtime monitoring merupakan fitur yang berbeda dan tidak termasuk dalam Target Information.

---

## 26. Future Development

Target Information dapat dikembangkan menjadi lebih lengkap:

```text
Target Information
│
├── Identity
│   ├── IP
│   ├── MAC
│   └── Vendor
│
├── OS Information
│   ├── OS Family
│   ├── OS Version
│   └── Detection Confidence
│
├── Network Information
│   ├── Ports
│   └── Services
│
└── Future
    ├── Hostname
    ├── Device Type
    ├── Service Version
    └── Additional Network Metadata
```

Informasi tambahan tetap harus ditampilkan berdasarkan data yang benar-benar berhasil diperoleh dari discovery.

---

## 27. Prinsip Arsitektur

Prinsip utama:

```text
Target Information ≠ Runtime Environment

Discovery ≠ Monitoring

Discovery Result ≠ Device

Detected ≠ Online

OS Detection ≠ OS Verification
```

Artinya:

- OS Detection memberikan informasi mengenai target.
    
- Runtime Environment memberikan informasi mengenai server GAMON.
    
- Discovery menemukan dan mengidentifikasi target.
    
- Monitoring memantau target secara periodik.
    
- Target yang ditemukan belum otomatis menjadi device.
    
- Hasil OS Detection merupakan hasil identifikasi berdasarkan informasi yang tersedia, bukan jaminan absolut mengenai OS target.
    

---

## 28. Expected Result

Sebelum fitur ini:

```text
Administrator
      ↓
Mengetahui IP secara manual
      ↓
Input Device
      ↓
Monitoring
```

Setelah fitur ini:

```text
Administrator
      ↓
Menentukan Target Network
      ↓
Network Discovery
      ↓
Host Discovery
      ↓
Target Information
      ├── IP
      ├── MAC
      ├── Vendor
      ├── OS
      ├── Ports
      └── Services
      ↓
Administrator Review
      ↓
Add Device
      ↓
Device Management
      ↓
Monitoring Engine
      ↓
ICMP / TCP / HTTP / SNMP
      ↓
History / Alert / Notification
```

Dengan demikian, Network Discovery menjadi tahap **asset identification dan target information gathering** sebelum perangkat dimasukkan ke dalam monitoring GAMON.

Fitur ini memperluas GAMON dari sekadar sistem yang menerima IP secara manual menjadi sistem yang dapat membantu administrator mengenali perangkat jaringan terlebih dahulu, kemudian menentukan perangkat mana yang akan dimonitor.