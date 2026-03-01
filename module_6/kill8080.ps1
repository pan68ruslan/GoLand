# Знаходимо процес, який слухає порт 8080
$process = netstat -ano | findstr :8080

if ($process) {
    $procId = ($process -split "\s+")[-1]
    Write-Host "Find 8080 with PID $procId. Terminate..."
    taskkill /PID $procId /F
    Write-Host "Port 8080 free now."
} else {
    Write-Host "Port 8080 free."
}

# Залишаємо PowerShell відкритим
Read-Host "type to close PowerShell"