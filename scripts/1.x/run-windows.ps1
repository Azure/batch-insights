$wd = $env:AZ_BATCH_TASK_WORKING_DIR

$exe = "$wd/batch-insights.exe"

[Net.ServicePointManager]::SecurityProtocol = "tls12, tls11, tls"
Invoke-WebRequest -Uri $env:BATCH_INSIGHTS_DOWNLOAD_URL -OutFile $exe

# Delete if exists
$exists = Get-ScheduledTask | Where-Object {$_.TaskName -like "batchappinsights" };

if($exists)
{
    Write-Output "Scheduled task already exists. Removing it and restarting it";
    Stop-ScheduledTask -TaskName "batchappinsights";
    Unregister-ScheduledTask -Confirm:$false -TaskName "batchappinsights";
}

Write-Output "Starting App insights background process in $wd"

# If using batch insights 1.x you need to have those arguments this way 
$legacyArgs = "$env:AZ_BATCH_POOL_ID $env:AZ_BATCH_NODE_ID $env:APP_INSIGHTS_INSTRUMENTATION_KEY";

$toolArgs =  "$legacyArgs --poolID $env:AZ_BATCH_POOL_ID --nodeID $env:AZ_BATCH_NODE_ID --instKey $env:APP_INSIGHTS_INSTRUMENTATION_KEY $env:AZ_BATCH_INSIGHTS_ARGS"

# TODO-TIM add toolsArgs
$action = New-ScheduledTaskAction -WorkingDirectory $wd -Execute 'cmd.exe' -Argument "/c $exe $toolArgs > $wd\nodestats.log 2>&1"
$principal = New-ScheduledTaskPrincipal -UserID 'NT AUTHORITY\SYSTEM' -LogonType ServiceAccount -RunLevel Highest ; 
$settings = New-ScheduledTaskSettingsSet -RestartCount 255 -RestartInterval ([timespan]::FromMinutes(1)) -ExecutionTimeLimit ([timespan]::FromDays(365))
Register-ScheduledTask -Action $action -Principal $principal -TaskName "batchappinsights" -Settings $settings -Force

Start-ScheduledTask -TaskName "batchappinsights"; 
Get-ScheduledTask -TaskName "batchappinsights";
