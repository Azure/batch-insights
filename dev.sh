set -e

apt-get update  
apt-get install -y git

bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
gvm install go1.11

echo GO version $(go --version)

git clone https://github.com/Azure/batch-insights -b feature/go

cd batch-insights
go run ./..
