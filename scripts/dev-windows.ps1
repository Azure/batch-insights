$ErrorActionPreference = "Stop"

$wd = $env:AZ_BATCH_TASK_WORKING_DIR
$branch = $env:BATCH_INSIGHTS_BRANCH

Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
choco install -y golang git
choco install -y -f mingw
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

git clone https://github.com/Azure/batch-insights -b $branch

Set-Location ./batch-insights

go build

$exe = "$wd/batch-insights/batch-insights.exe"

# Delete if exists
$exists = Get-ScheduledTask | Where-Object {$_.TaskName -like "batchappinsights" };

if($exists)
{
    Write-Host "Scheduled task already exists. Removing it and restarting it";
    Stop-ScheduledTask -TaskName "batchappinsights";
    Unregister-ScheduledTask -Confirm:$false -TaskName "batchappinsights";
}


$toolArgs =  "--poolID $env:AZ_BATCH_POOL_ID --nodeID $env:AZ_BATCH_NODE_ID --instKey $env:APP_INSIGHTS_INSTRUMENTATION_KEY $AZ_BATCH_INSIGHTS_ARGS"

Write-Host "Starting App insights background process in $wd"
$action = New-ScheduledTaskAction -WorkingDirectory $wd -Execute 'cmd.exe' -Argument "/c $exe  $toolArgs > .\batch-insights.log 2>&1"  
$principal = New-ScheduledTaskPrincipal -UserID 'NT AUTHORITY\SYSTEM' -LogonType ServiceAccount -RunLevel Highest ; 
Register-ScheduledTask -Action $action -Principal $principal -TaskName "batchappinsights" -Force ; 
Start-ScheduledTask -TaskName "batchappinsights"; 
Get-ScheduledTask -TaskName "batchappinsights";