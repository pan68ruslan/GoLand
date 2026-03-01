# Знаходимо процес, який слухає порт 8080
$process = netstat -ano | findstr :27017

if ($process) {
    $procId = ($process -split "\s+")[-1]
    Write-Host "Find 27017 with PID $procId. Terminate..."
    taskkill /PID $procId /F
    Write-Host "Port 27017 free now."
} else {
    Write-Host "Port 27017 free."
}

# Залишаємо PowerShell відкритим
Read-Host "type to close PowerShell"