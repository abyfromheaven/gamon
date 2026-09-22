# ============================================================
# SCRIPT OTOMATIS TAMBAH DEVICE - GAMON
# Untuk persiapan sidang - tambah semua device sekaligus
# ============================================================

$BASE_URL = "http://localhost:8080"

# ==========================================
# LIST DEVICE YANG MAU DITAMBAHKAN
# ==========================================
$devices = @(
    @{
        name     = "Router Core"
        type     = "Router"
        ip       = "192.168.1.1"
        method   = "ICMP Ping"
        interval = 3
        location = "Server Room"
        desc     = "Router utama jaringan"
    },
    @{
        name     = "Switch Core"
        type     = "Switch"
        ip       = "192.168.1.32"
        method   = "ICMP Ping"
        interval = 3
        location = "Server Room"
        desc     = "Switch utama jaringan"
    },
    @{
        name     = "CCTV Server"
        type     = "Server"
        ip       = "192.168.1.164"
        method   = "ICMP Ping"
        interval = 3
        location = "Server Room"
        desc     = "Server monitoring CCTV"
    },
    @{
        name     = "Database Server"
        type     = "Server"
        ip       = "192.168.1.174"
        method   = "ICMP Ping"
        interval = 3
        location = "Server Room"
        desc     = "Server database utama"
    }
)

# ==========================================
# CEK KONEKSI SERVER DULU
# ==========================================
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  GAMON - AUTO ADD DEVICE" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "[*] Cek koneksi ke server GAMON..." -ForegroundColor Yellow

try {
    $health = Invoke-RestMethod -Uri "$BASE_URL/api/health" -Method GET -TimeoutSec 5
    Write-Host "[OK] Server GAMON aktif di $BASE_URL" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] Server GAMON tidak ditemukan di $BASE_URL" -ForegroundColor Red
    Write-Host "[!] Pastikan server GAMON sudah jalan (go run main.go)" -ForegroundColor Red
    Write-Host ""
    exit 1
}

# ==========================================
# TAMBAHKAN DEVICE SATU PER SATU
# ==========================================
Write-Host ""
Write-Host "[*] Menambahkan $($devices.Count) device..." -ForegroundColor Yellow
Write-Host ""

$success = 0
$failed = 0

foreach ($device in $devices) {
    $body = @{
        name           = $device.name
        type           = $device.type
        ip             = $device.ip
        method         = $device.method
        check_interval = $device.interval
        status         = "active"
        location       = $device.location
        description    = $device.desc
    } | ConvertTo-Json

    try {
        $result = Invoke-RestMethod -Uri "$BASE_URL/api/devices" -Method POST -Body $body -ContentType "application/json" -TimeoutSec 10
        Write-Host "[+] $($device.name) ($($device.ip)) - BERHASIL" -ForegroundColor Green
        $success++
    } catch {
        $errMsg = $_.ErrorDetails.Message
        if ($errMsg) {
            Write-Host "[-] $($device.name) ($($device.ip)) - GAGAL: $errMsg" -ForegroundColor Red
        } else {
            Write-Host "[-] $($device.name) ($($device.ip)) - GAGAL: $($_.Exception.Message)" -ForegroundColor Red
        }
        $failed++
    }
}

# ==========================================
# RINGKASAN
# ==========================================
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  RINGKASAN" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Berhasil : $success device" -ForegroundColor Green
if ($failed -gt 0) {
    Write-Host "  Gagal    : $failed device" -ForegroundColor Red
}
Write-Host "  Total    : $($devices.Count) device" -ForegroundColor White
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "[*] Device sudah otomatis dimonitor oleh GAMON" -ForegroundColor Yellow
Write-Host "[*] Buka http://localhost:8080 untuk melihat dashboard" -ForegroundColor Yellow
Write-Host ""
