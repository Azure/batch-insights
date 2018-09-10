set -e

apt-get update  
apt-get install -y golang-go git

git clone https://github.com/Azure/batch-insights -b feature/go

cd batch-insights
go run ./..