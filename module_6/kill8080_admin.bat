@echo off
:: Запускаємо PowerShell‑скрипт від імені адміністратора
powershell -NoExit -Command "Start-Process PowerShell -Verb RunAs -ArgumentList '-NoExit -File \"%~dp0kill8080.ps1\"'"
