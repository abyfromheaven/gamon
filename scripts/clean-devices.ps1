# ============================================================
# SCRIPT BERSIHKAN SEMUA DEVICE - GAMON
# Untuk reset data sebelum sidang / demo
# ============================================================

$BASE_URL = "http://localhost:8080"

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  GAMON - CLEAN ALL DEVICES" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# ==========================================
# CEK KONEKSI SERVER
# ==========================================
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
# AMBIL SEMUA DEVICE
# ==========================================
Write-Host ""
Write-Host "[*] Mengambil data semua device..." -ForegroundColor Yellow

try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/api/devices" -Method GET -TimeoutSec 10
    $devices = $response.data
} catch {
    Write-Host "[ERROR] Gagal mengambil data device" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    exit 1
}

if ($null -eq $devices -or $devices.Count -eq 0) {
    Write-Host "[INFO] Tidak ada device yang terdaftar. Database sudah kosong." -ForegroundColor Yellow
    Write-Host ""
    exit 0
}

Write-Host "[INFO] Ditemukan $($devices.Count) device:" -ForegroundColor White
foreach ($d in $devices) {
    Write-Host "  - $($d.name) ($($d.ip)) [ID: $($d.id)]" -ForegroundColor Gray
}

# ==========================================
# KONFIRMASI HAPUS
# ==========================================
Write-Host ""
Write-Host "[WARNING] Semua device akan DIHAPUS permanent!" -ForegroundColor Red
$confirm = Read-Host "Ketik 'YES' untuk konfirmasi"

if ($confirm -ne "YES") {
    Write-Host "[*] Dibatalkan oleh user." -ForegroundColor Yellow
    Write-Host ""
    exit 0
}

# ==========================================
# HAPUS DEVICE SATU PER SATU
# ==========================================
Write-Host ""
Write-Host "[*] Menghapus semua device..." -ForegroundColor Yellow

$deleted = 0
$failed = 0

foreach ($device in $devices) {
    try {
        Invoke-RestMethod -Uri "$BASE_URL/api/devices/$($device.id)" -Method DELETE -TimeoutSec 10 | Out-Null
        Write-Host "[-] $($device.name) ($($device.ip)) - DIHAPUS" -ForegroundColor Red
        $deleted++
    } catch {
        Write-Host "[!] $($device.name) ($($device.ip)) - GAGAL HAPUS: $($_.Exception.Message)" -ForegroundColor Red
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
Write-Host "  Dihapus : $deleted device" -ForegroundColor Red
if ($failed -gt 0) {
    Write-Host "  Gagal   : $failed device" -ForegroundColor Yellow
}
Write-Host "  Sisa    : 0 device" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "[OK] Semua device berhasil dibersihkan!" -ForegroundColor Green
Write-Host "[*] Siap untuk tambah device baru atau demo sidang." -ForegroundColor Yellow
Write-Host ""
