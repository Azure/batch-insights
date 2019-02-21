# Script to be used as one liner

### For linux

* `run-linux.sh`: Run a published version for linux
* `dev.sh`: Will install go, git then build and run on the fly


### For windows
* `run-windows.ps1`: Run a publush version for windows
* `dev-windows.ps1`: Will install go, git then build and run on the fly


## Development
There is some dev script that will install go and other needed dependencies to build and run this project on the fly.
Set `BATCH_INSIGHTS_BRANCH` environment variable to the branch you are testing

On linux 
```bash
/bin/bash -c 'wget  -O - https://raw.githubusercontent.com/Azure/batch-insights/$BATCH_INSIGHTS_BRANCH/scripts/dev.sh | bash'
```

On windows

```powershell
cmd /c @"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -InputFormat None -ExecutionPolicy Bypass -Command "iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/Azure/batch-insights/$env:BATCH_INSIGHTS_BRANCH/dev-windows.ps1'))"
```
