set -e

apt-get update  
apt-get install -y golang-1.11-go git

git clone https://github.com/Azure/batch-insights -b feature/go

echo GO version $(go --version)
cd batch-insights
go run ./..
