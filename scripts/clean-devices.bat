@echo off
title GAMON - Clean All Devices
echo.
echo ========================================
echo   GAMON - CLEAN ALL DEVICES
echo ========================================
echo.
echo WARNING: Semua device akan dihapus!
echo.
powershell -ExecutionPolicy Bypass -File "%~dp0clean-devices.ps1"
echo.
pause
