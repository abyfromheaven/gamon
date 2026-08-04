#!/bin/bash
# Reset Gamon - Hapus semua data untuk testing bersih
set -e

DB_FILE="data/gamon.db"
PID_FILE="data/gamon.pid"

echo "=== Gamon Reset ==="

# Hapus database
if [ -f "$DB_FILE" ]; then
    rm -f "$DB_FILE"
    echo "[OK] Database dihapus: $DB_FILE"
else
    echo "[--] Database tidak ditemukan"
fi

# Hapus WAL/SHM files
rm -f "${DB_FILE}-wal" "${DB_FILE}-shm"
echo "[OK] WAL/SHM files dihapus"

# Hapus log monitoring
rm -f data/monitoring.log 2>/dev/null
echo "[OK] Log monitoring dihapus"

# Hapus binary lama
rm -f gamon 2>/dev/null
echo "[OK] Binary lama dihapus"

echo ""
echo "=== Selesai ==="
echo "Jalankan: go run . untuk memulai ulang dengan data bersih"
