$ErrorActionPreference = "Stop"

$wd = $env:AZ_BATCH_TASK_WORKING_DIR

Invoke-WebRequest -Uri https://dl.google.com/go/go1.11.2.windows-amd64.zip -OutFile go.zip

Expand-Archive go.zip -DestinationPath C:\GO

$env:GOROOT = "C:\GO\"
$env:PATH = "$GOROOT/bin;$env:PATH"

$exe = "$wd/batch-insights.exe"

# Delete if exists
$exists = Get-ScheduledTask | Where-Object {$_.TaskName -like "batchappinsights" };

if($exists)
{
    Write-Host "Scheduled task already exists. Removing it and restarting it";
    Stop-ScheduledTask -TaskName "batchappinsights";
    Unregister-ScheduledTask -Confirm:$false -TaskName "batchappinsights";
}


Write-Host "Starting App insights background process in $wd"
$action = New-ScheduledTaskAction -WorkingDirectory $wd -Execute 'Powershell.exe' -Argument "Start-Process $exe -ArgumentList ('$env:AZ_BATCH_POOL_ID', '$env:AZ_BATCH_NODE_ID', '$env:APP_INSIGHTS_INSTRUMENTATION_KEY')  -RedirectStandardOutput .\node-stats.log -RedirectStandardError .\node-stats.err.log -NoNewWindow"  
$principal = New-ScheduledTaskPrincipal -UserID 'NT AUTHORITY\SYSTEM' -LogonType ServiceAccount -RunLevel Highest ; 
Register-ScheduledTask -Action $action -Principal $principal -TaskName "batchappinsights" -Force ; 
Start-ScheduledTask -TaskName "batchappinsights"; 
Get-ScheduledTask -TaskName "batchappinsights";