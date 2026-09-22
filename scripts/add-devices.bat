@echo off
title GAMON - Auto Add Devices
echo.
echo ========================================
echo   GAMON - AUTO ADD DEVICE
echo ========================================
echo.
echo Menambahkan device:
echo   1. Router Core    (192.168.1.1)
echo   2. Switch Core    (192.168.1.32)
echo   3. CCTV Server    (192.168.1.164)
echo   4. Database Server(192.168.1.174)
echo.
echo Pastikan server GAMON sudah jalan!
echo.
powershell -ExecutionPolicy Bypass -File "%~dp0add-devices.ps1"
echo.
pause
