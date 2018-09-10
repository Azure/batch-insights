set -e

apt-get update  
apt-get install -y git binutils bison build-essential

bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source /mnt/batch/tasks/startup/wd/.gvm/scripts/gvm
gvm install go1.4.3
gvm use go1.4.3
gvm install go1.11
gvm use go1.11

echo GO version $(go --version)

git clone https://github.com/Azure/batch-insights -b feature/go

cd batch-insights
go run ./..
