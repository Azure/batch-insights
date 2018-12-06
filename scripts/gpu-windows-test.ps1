$ErrorActionPreference = "Stop"

$wd = $env:AZ_BATCH_TASK_WORKING_DIR

Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
choco install -y golang git
$env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path","User")

git clone https://github.com/Azure/batch-insights

Set-Location ./batch-insights

go build

$exe = "$wd/batch-insights/batch-insights.exe"

& $exe