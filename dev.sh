set -e

apt-get update  
apt-get install -y git binutils bison build-essential

wget https://dl.google.com/go/go1.11.linux-amd64.tar.gz
tar -xvf go1.11.linux-amd64.tar.gz
mv go /usr/local
export GOROOT=/usr/local/go

echo GO version $(go --version)

git clone https://github.com/Azure/batch-insights -b feature/go

cd batch-insights
go run ./..
