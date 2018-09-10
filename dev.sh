set -e

add-apt-repository ppa:longsleep/golang-backports
apt-get update  
apt-get install -y golang-go git

git clone https://github.com/Azure/batch-insights -b feature/go

echo GO version $(go --version)
cd batch-insights
go run ./..
